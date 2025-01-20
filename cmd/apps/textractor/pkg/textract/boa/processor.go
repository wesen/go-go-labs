package boa

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// Constants for section headers
const (
	SectionDeposits    = "deposits"
	SectionWithdrawals = "withdrawals"
	SectionServiceFees = "fees"
)

type TransactionType string

const (
	TransactionDeposit    TransactionType = "deposit"
	TransactionWithdrawal TransactionType = "withdrawal"
	TransactionServiceFee TransactionType = "service_fee"
)

// Account represents a bank account
type Account struct {
	Number      string  // Normalized account number (spaces/dashes removed)
	Type        string  // e.g. "Adv Plus Banking", "Regular Savings"
	Balance     float64 // Current balance
	DetailsPage int     // Page number where details start
}

// Transaction represents a financial transaction
type Transaction struct {
	Date          time.Time
	Description   string
	Amount        float64
	Type          TransactionType
	AccountNumber string
	Page          int    // Page where transaction starts
	LineNumber    int    // Line number where transaction starts
	RawLines      []Line // All lines that make up this transaction
}

// AccountSummary represents the summary of an account
type AccountSummary struct {
	AccountNumber    string
	Period           string // Statement period
	BeginBalance     float64
	EndBalance       float64
	DepositsTotal    float64
	WithdrawalsTotal float64
	ServiceFeesTotal float64
}

// ProcessorState tracks the current state during processing
type ProcessorState struct {
	CurrentAccount     string
	CurrentSection     string // "deposits", "withdrawals", "fees"
	CurrentTransaction *Transaction
	PageContext        int
	LastLineNumber     int
}

// StatementProcessor processes bank statements
type StatementProcessor struct {
	state        ProcessorState
	accounts     map[string]*Account
	summaries    map[string]*AccountSummary
	transactions map[string][]Transaction
}

// Regular expressions for pattern matching
var (
	accountNumberPattern = regexp.MustCompile(`Account (?:number:|#) (\d{4} \d{4} \d{4})`)
	datePattern          = regexp.MustCompile(`^(\d{2}/\d{2}/\d{2})`) // Must be at start of line
	amountPattern        = regexp.MustCompile(`^[-$(),.0-9]+$`)       // Must be entire line
	periodPattern        = regexp.MustCompile(`for ([A-Za-z]+ \d{1,2}, \d{4}) to ([A-Za-z]+ \d{1,2}, \d{4})`)
)

// NewProcessor creates a new StatementProcessor
func NewProcessor() *StatementProcessor {
	return &StatementProcessor{
		accounts:     make(map[string]*Account),
		summaries:    make(map[string]*AccountSummary),
		transactions: make(map[string][]Transaction),
	}
}

// anonymizeAccountNumber replaces middle digits with X's, keeping first and last 4 digits
func anonymizeAccountNumber(accNum string) string {
	// Remove spaces and dashes first
	clean := strings.ReplaceAll(strings.ReplaceAll(accNum, " ", ""), "-", "")
	if len(clean) < 8 { // Need at least 8 digits to anonymize
		return clean
	}

	// Keep first 4 and last 4 digits, replace rest with X
	first4 := clean[:4]
	last4 := clean[len(clean)-4:]
	middle := strings.Repeat("X", len(clean)-8)

	// Format with spaces for readability
	return fmt.Sprintf("%s %s %s", first4, middle, last4)
}

// normalizeAccountNumber removes spaces and dashes from account numbers
func normalizeAccountNumber(accNum string) string {
	return strings.ReplaceAll(strings.ReplaceAll(accNum, " ", ""), "-", "")
}

// ProcessLine processes a single line from the CSV
func (p *StatementProcessor) ProcessLine(line Line) error {
	// Anonymize any account numbers in the line text for logging
	logText := accountNumberPattern.ReplaceAllStringFunc(line.Text, func(match string) string {
		if groups := accountNumberPattern.FindStringSubmatch(match); len(groups) > 1 {
			return "Account " + anonymizeAccountNumber(groups[1])
		}
		return match
	})

	logger := log.With().
		Int("page", line.Page).
		Int("line", line.LineNumber).
		Str("text", logText).
		Logger()

	// Update state
	p.state.PageContext = line.Page
	p.state.LastLineNumber = line.LineNumber

	// Check for account number
	if match := accountNumberPattern.FindStringSubmatch(line.Text); match != nil {
		accNum := normalizeAccountNumber(match[1])
		logger.Debug().
			Str("account", anonymizeAccountNumber(accNum)).
			Msg("Found account number")
		p.state.CurrentAccount = accNum
		if _, exists := p.accounts[accNum]; !exists {
			p.accounts[accNum] = &Account{
				Number:      accNum,
				DetailsPage: line.Page,
			}
		}
		return nil
	}

	// Check for statement period
	if match := periodPattern.FindStringSubmatch(line.Text); match != nil {
		if p.state.CurrentAccount != "" {
			logger.Debug().
				Str("account", p.state.CurrentAccount).
				Str("period", fmt.Sprintf("%s to %s", match[1], match[2])).
				Msg("Found statement period")
			if summary, exists := p.summaries[p.state.CurrentAccount]; exists {
				summary.Period = fmt.Sprintf("%s to %s", match[1], match[2])
			} else {
				p.summaries[p.state.CurrentAccount] = &AccountSummary{
					AccountNumber: p.state.CurrentAccount,
					Period:        fmt.Sprintf("%s to %s", match[1], match[2]),
				}
			}
		}
		return nil
	}

	// Check for section headers
	switch {
	case strings.Contains(line.Text, "Deposits and other additions"):
		logger.Debug().Msg("Entering deposits section")
		p.state.CurrentSection = SectionDeposits
		return nil
	case strings.Contains(line.Text, "Withdrawals and other subtractions"):
		logger.Debug().Msg("Entering withdrawals section")
		p.state.CurrentSection = SectionWithdrawals
		return nil
	case strings.Contains(line.Text, "Service fees"):
		logger.Debug().Msg("Entering service fees section")
		p.state.CurrentSection = SectionServiceFees
		return nil
	}

	// Process transactions
	if p.state.CurrentSection != "" && p.state.CurrentAccount != "" {
		return p.processTransactionLine(line)
	}

	return nil
}

// processTransactionLine processes a line that might be part of a transaction
func (p *StatementProcessor) processTransactionLine(line Line) error {
	logger := log.With().
		Int("page", line.Page).
		Int("line", line.LineNumber).
		Str("text", line.Text).
		Str("section", p.state.CurrentSection).
		Logger()

	// Start new transaction if we see a date at the start of a line
	if match := datePattern.FindStringSubmatch(line.Text); match != nil {
		if p.state.CurrentTransaction != nil {
			logger.Debug().
				Interface("current_tx", p.state.CurrentTransaction).
				Msg("Completing previous transaction before starting new one")
			p.completeTransaction()
		}
		date, err := time.Parse("01/02/06", match[1])
		if err != nil {
			return fmt.Errorf("parsing date: %w", err)
		}
		logger.Debug().
			Time("date", date).
			Msg("Starting new transaction")
		p.state.CurrentTransaction = &Transaction{
			Date:          date,
			AccountNumber: p.state.CurrentAccount,
			Page:          line.Page,
			LineNumber:    line.LineNumber,
			Type:          p.transactionTypeFromSection(),
			RawLines:      []Line{line},
		}
		return nil
	}

	// If we have a current transaction, add to description or try to parse amount
	if p.state.CurrentTransaction != nil {
		trimmedText := strings.TrimSpace(line.Text)
		logger := logger.With().
			Interface("current_tx", p.state.CurrentTransaction).
			Str("trimmed_text", trimmedText).
			Logger()

		// Try to parse amount only if the line contains just the amount
		if trimmedText != "" && amountPattern.MatchString(trimmedText) {
			if amount, ok := p.tryParseAmount(trimmedText); ok {
				logger.Debug().
					Float64("amount", amount).
					Msg("Found transaction amount")
				p.state.CurrentTransaction.Amount = amount
				p.state.CurrentTransaction.RawLines = append(p.state.CurrentTransaction.RawLines, line)
				p.completeTransaction()
				return nil
			}
		}

		// Add to description if not an amount
		if p.state.CurrentTransaction.Description != "" {
			p.state.CurrentTransaction.Description += " "
		}
		p.state.CurrentTransaction.Description += trimmedText
		p.state.CurrentTransaction.RawLines = append(p.state.CurrentTransaction.RawLines, line)
		logger.Debug().
			Str("description", p.state.CurrentTransaction.Description).
			Msg("Updated transaction description")
	}

	return nil
}

// completeTransaction finalizes and stores the current transaction
func (p *StatementProcessor) completeTransaction() {
	if p.state.CurrentTransaction == nil {
		return
	}

	logger := log.With().
		Interface("transaction", p.state.CurrentTransaction).
		Str("account", p.state.CurrentTransaction.AccountNumber).
		Logger()

	p.transactions[p.state.CurrentAccount] = append(
		p.transactions[p.state.CurrentAccount],
		*p.state.CurrentTransaction,
	)

	// Update summary
	summary, exists := p.summaries[p.state.CurrentAccount]
	if !exists {
		summary = &AccountSummary{AccountNumber: p.state.CurrentAccount}
		p.summaries[p.state.CurrentAccount] = summary
	}

	switch p.state.CurrentTransaction.Type {
	case TransactionDeposit:
		summary.DepositsTotal += p.state.CurrentTransaction.Amount
	case TransactionWithdrawal:
		summary.WithdrawalsTotal += p.state.CurrentTransaction.Amount
	case TransactionServiceFee:
		summary.ServiceFeesTotal += p.state.CurrentTransaction.Amount
	}

	logger.Debug().
		Float64("deposits_total", summary.DepositsTotal).
		Float64("withdrawals_total", summary.WithdrawalsTotal).
		Float64("fees_total", summary.ServiceFeesTotal).
		Msg("Completed transaction and updated summary")

	p.state.CurrentTransaction = nil
}

// transactionTypeFromSection determines transaction type based on current section
func (p *StatementProcessor) transactionTypeFromSection() TransactionType {
	switch p.state.CurrentSection {
	case SectionDeposits:
		return TransactionDeposit
	case SectionWithdrawals:
		return TransactionWithdrawal
	case SectionServiceFees:
		return TransactionServiceFee
	default:
		return ""
	}
}

// tryParseAmount attempts to parse an amount from text
func (p *StatementProcessor) tryParseAmount(text string) (float64, bool) {
	if match := amountPattern.FindString(text); match != "" {
		// Clean up the amount string
		clean := strings.ReplaceAll(match, "$", "")
		clean = strings.ReplaceAll(clean, ",", "")
		clean = strings.ReplaceAll(clean, "(", "-")
		clean = strings.ReplaceAll(clean, ")", "")

		amount, err := strconv.ParseFloat(clean, 64)
		if err == nil {
			return amount, true
		}
	}
	return 0, false
}

// GetAccounts returns all detected accounts
func (p *StatementProcessor) GetAccounts() map[string]*Account {
	return p.accounts
}

// GetTransactions returns transactions for an account
func (p *StatementProcessor) GetTransactions(accountNumber string) []Transaction {
	return p.transactions[accountNumber]
}

// GetAccountSummary returns the summary for an account
func (p *StatementProcessor) GetAccountSummary(accountNumber string) *AccountSummary {
	return p.summaries[accountNumber]
}

// Reset resets the processor state
func (p *StatementProcessor) Reset() {
	p.state = ProcessorState{}
	p.accounts = make(map[string]*Account)
	p.summaries = make(map[string]*AccountSummary)
	p.transactions = make(map[string][]Transaction)
}
