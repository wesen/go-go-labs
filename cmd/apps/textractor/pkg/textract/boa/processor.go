package boa

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
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

	// New fields for parsed totals
	ParsedDepositsTotal    float64
	ParsedWithdrawalsTotal float64
	ParsedServiceFeesTotal float64
}

// Section represents different parts of the statement
type Section string

const (
	SectionUnknown     Section = ""
	SectionSummary     Section = "summary"
	SectionDeposits    Section = "deposits"
	SectionWithdrawals Section = "withdrawals"
	SectionServiceFees Section = "fees"
)

// ProcessorState tracks the current state during processing
type ProcessorState struct {
	CurrentAccount     string
	CurrentSection     Section // Use the Section type
	CurrentTransaction *Transaction
	PageContext        int
	LastLineNumber     int
	expectingAmount    string // New field to track expected amount type
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
	datePattern          = regexp.MustCompile(`^(\d{2}/\d{2}/\d{2})`)   // Must be at start of line
	amountPattern        = regexp.MustCompile(`^\(?[-$]?\d+(?:,\d{3})*\.\d{2}\)?$`) // Matches ($1,234.56), -$1,234.56, 1,234.56 etc
	periodPattern        = regexp.MustCompile(`for ([A-Za-z]+ \d{1,2}, \d{4}) to ([A-Za-z]+ \d{1,2}, \d{4})`)
	beginBalancePattern  = regexp.MustCompile(`Beginning balance on (.+)`)
	endBalancePattern    = regexp.MustCompile(`Ending balance on (.+)`)
	accountTypePattern   = regexp.MustCompile(`Your (.*Banking)`)
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
	logger := log.With().
		Int("page", line.Page).
		Int("line", line.LineNumber).
		Str("text", line.Text).
		Logger()

	// Check for beginning of summary section
	if match := beginBalancePattern.FindStringSubmatch(line.Text); match != nil {
		logger.Debug().Msg("Entering summary section")
		p.state.CurrentSection = SectionSummary
		p.state.expectingAmount = "begin_balance"
		return nil
	}

	// Check for end of summary section
	if match := endBalancePattern.FindStringSubmatch(line.Text); match != nil {
		logger.Debug().Msg("Processing end balance")
		p.state.expectingAmount = "end_balance"
		return nil
	}

	// Process section headers based on current section
	if p.state.CurrentSection == SectionSummary {
		switch {
		case strings.Contains(line.Text, "Deposits and other additions"):
			logger.Debug().Msg("Processing summary deposits section")
			p.state.expectingAmount = "deposits_total"
			return nil
		case strings.Contains(line.Text, "Withdrawals and other subtractions"):
			logger.Debug().Msg("Processing summary withdrawals section")
			p.state.expectingAmount = "withdrawals_total"
			return nil
		case strings.Contains(line.Text, "Service fees"):
			logger.Debug().Msg("Processing summary service fees section")
			p.state.expectingAmount = "fees_total"
			return nil
		}

		// Handle expected amounts
		if p.state.expectingAmount != "" && p.state.CurrentAccount != "" {
			if amount, ok := p.tryParseAmount(line.Text); ok {
				summary := p.summaries[p.state.CurrentAccount]
				if summary == nil {
					summary = &AccountSummary{AccountNumber: p.state.CurrentAccount}
					p.summaries[p.state.CurrentAccount] = summary
				}

				switch p.state.expectingAmount {
				case "begin_balance":
					summary.BeginBalance = amount
				case "end_balance":
					summary.EndBalance = amount
					if acc := p.accounts[p.state.CurrentAccount]; acc != nil {
						acc.Balance = amount
					}
					// Exit summary section after processing end balance
					p.state.CurrentSection = SectionUnknown
				case "deposits_total":
					summary.ParsedDepositsTotal = amount
				case "withdrawals_total":
					summary.ParsedWithdrawalsTotal = amount
				case "fees_total":
					summary.ParsedServiceFeesTotal = amount
				}
				p.state.expectingAmount = ""
				return nil
			}
		}
	}

	// Only process transaction sections when not in summary
	if p.state.CurrentSection != SectionSummary {
		switch {
		case strings.Contains(line.Text, "Deposits and other additions"):
			logger.Debug().Msg("Entering deposits transaction section")
			p.state.CurrentSection = SectionDeposits
			return nil
		case strings.Contains(line.Text, "Withdrawals and other subtractions"):
			logger.Debug().Msg("Entering withdrawals transaction section")
			p.state.CurrentSection = SectionWithdrawals
			return nil
		case strings.Contains(line.Text, "Service fees"):
			logger.Debug().Msg("Entering service fees transaction section")
			p.state.CurrentSection = SectionServiceFees
			return nil
		}
	}

	// Anonymize any account numbers in the line text for logging
	logText := accountNumberPattern.ReplaceAllStringFunc(line.Text, func(match string) string {
		if groups := accountNumberPattern.FindStringSubmatch(match); len(groups) > 1 {
			return "Account " + anonymizeAccountNumber(groups[1])
		}
		return match
	})

	logger.Debug().
		Int("page", line.Page).
		Int("line", line.LineNumber).
		Str("text", logText).
		Msg("Processing line")

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

	// Check for account type
	if match := accountTypePattern.FindStringSubmatch(line.Text); match != nil && p.state.CurrentAccount != "" {
		if acc := p.accounts[p.state.CurrentAccount]; acc != nil {
			acc.Type = match[1]
		}
		return nil
	}

	// Process transactions
	if p.state.CurrentSection != SectionUnknown && p.state.CurrentSection != SectionSummary && p.state.CurrentAccount != "" {
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
		Str("section", string(p.state.CurrentSection)).
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
		Str("type", string(p.state.CurrentTransaction.Type)).
		Float64("amount", p.state.CurrentTransaction.Amount).
		Str("description", p.state.CurrentTransaction.Description).
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
	p.state = ProcessorState{
		CurrentSection: SectionUnknown,
	}
	p.accounts = make(map[string]*Account)
	p.summaries = make(map[string]*AccountSummary)
	p.transactions = make(map[string][]Transaction)
}

// ValidateSummaryTotals checks if transaction totals match summary amounts
// Returns nil if totals match, otherwise returns an error with details
func (p *StatementProcessor) ValidateSummaryTotals(accountNumber string) error {
	summary := p.summaries[accountNumber]
	if summary == nil {
		return fmt.Errorf("no summary found for account %s", accountNumber)
	}

	var calcDeposits, calcWithdrawals, calcFees float64
	for _, tx := range p.transactions[accountNumber] {
		switch tx.Type {
		case TransactionDeposit:
			calcDeposits += tx.Amount
		case TransactionWithdrawal:
			calcWithdrawals += tx.Amount
		case TransactionServiceFee:
			calcFees += tx.Amount
		}
	}

	// Use small epsilon for float comparison
	const epsilon = 0.01
	if math.Abs(calcDeposits-summary.DepositsTotal) > epsilon {
		return fmt.Errorf("deposits total mismatch for account %s: calculated %.2f vs summary %.2f",
			accountNumber, calcDeposits, summary.DepositsTotal)
	}
	if math.Abs(calcWithdrawals-summary.WithdrawalsTotal) > epsilon {
		return fmt.Errorf("withdrawals total mismatch for account %s: calculated %.2f vs summary %.2f",
			accountNumber, calcWithdrawals, summary.WithdrawalsTotal)
	}
	if math.Abs(calcFees-summary.ServiceFeesTotal) > epsilon {
		return fmt.Errorf("service fees total mismatch for account %s: calculated %.2f vs summary %.2f",
			accountNumber, calcFees, summary.ServiceFeesTotal)
	}

	return nil
}

// Add new validation method
func (p *StatementProcessor) ValidateParsedTotals(accountNumber string) error {
	summary := p.summaries[accountNumber]
	if summary == nil {
		return fmt.Errorf("no summary found for account %s", accountNumber)
	}

	// Use small epsilon for float comparison
	const epsilon = 0.01

	// Check each total type
	if math.Abs(summary.DepositsTotal-summary.ParsedDepositsTotal) > epsilon {
		return fmt.Errorf("deposits total mismatch: calculated %.2f vs parsed %.2f",
			summary.DepositsTotal, summary.ParsedDepositsTotal)
	}

	if math.Abs(summary.WithdrawalsTotal-summary.ParsedWithdrawalsTotal) > epsilon {
		return fmt.Errorf("withdrawals total mismatch: calculated %.2f vs parsed %.2f",
			summary.WithdrawalsTotal, summary.ParsedWithdrawalsTotal)
	}

	if math.Abs(summary.ServiceFeesTotal-summary.ParsedServiceFeesTotal) > epsilon {
		return fmt.Errorf("service fees total mismatch: calculated %.2f vs parsed %.2f",
			summary.ServiceFeesTotal, summary.ParsedServiceFeesTotal)
	}

	// Validate that beginning balance + net changes = ending balance
	netChange := summary.DepositsTotal + summary.WithdrawalsTotal + summary.ServiceFeesTotal
	expectedEndBalance := summary.BeginBalance + netChange
	if math.Abs(expectedEndBalance-summary.EndBalance) > epsilon {
		return fmt.Errorf("balance equation mismatch: begin %.2f + net change %.2f != end %.2f",
			summary.BeginBalance, netChange, summary.EndBalance)
	}

	return nil
}
