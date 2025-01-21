package boa

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

// Test account numbers
const (
	testAccount1 = "1111 2222 3333" // Will normalize to "111122223333"
	testAccount2 = "4444 5555 6666" // Will normalize to "444455556666"
)

func init() {
	// Configure zerolog to use console writer for tests
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "15:04:05.000",
	}
	log.Logger = log.Output(consoleWriter)

	// Set log level from TEST_LOG_LEVEL environment variable
	// Valid values: trace, debug, info, warn, error (default: info)
	logLevel := strings.ToLower(os.Getenv("TEST_LOG_LEVEL"))
	switch logLevel {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

// Helper function to create a logger for a test
func testLogger(t *testing.T) zerolog.Logger {
	return log.With().
		Str("test", t.Name()).
		Logger()
}

func TestProcessorBasics(t *testing.T) {
	p := NewProcessor()
	assert.NotNil(t, p)
	assert.Empty(t, p.GetAccounts())
	assert.Empty(t, p.GetTransactions("any"))
	assert.Nil(t, p.GetAccountSummary("any"))
}

func TestAccountDetection(t *testing.T) {
	p := NewProcessor()

	// Test both account number formats
	testCases := []struct {
		text     string
		expected string
	}{
		{"Account number: " + testAccount1, "111122223333"},
		{"Account # " + testAccount2, "444455556666"},
	}

	for _, tc := range testCases {
		err := p.ProcessLine(Line{
			Page:       1,
			LineNumber: 1,
			Text:       tc.text,
		})
		assert.NoError(t, err)

		acc := p.GetAccounts()[tc.expected]
		assert.NotNil(t, acc)
		assert.Equal(t, tc.expected, acc.Number)
	}
}

func TestStatementPeriod(t *testing.T) {
	p := NewProcessor()

	// First set account
	p.ProcessLine(Line{
		Page:       1,
		LineNumber: 1,
		Text:       "Account number: " + testAccount1,
	})

	// Then process period
	err := p.ProcessLine(Line{
		Page:       1,
		LineNumber: 2,
		Text:       "for December 28, 2023 to January 29, 2024",
	})
	assert.NoError(t, err)

	summary := p.GetAccountSummary("111122223333")
	assert.NotNil(t, summary)
	assert.Equal(t, "December 28, 2023 to January 29, 2024", summary.Period)
}

func TestTransactionProcessing(t *testing.T) {
	p := NewProcessor()

	// Setup account and section
	lines := []Line{
		{Page: 1, LineNumber: 1, Text: "Account number: " + testAccount1},
		{Page: 1, LineNumber: 2, Text: "Deposits and other additions"},
		{Page: 1, LineNumber: 3, Text: "12/29/23"},
		{Page: 1, LineNumber: 4, Text: "THE TREE CENTER DES:PAYROLL"},
		{Page: 1, LineNumber: 5, Text: "2,774.40"},
	}

	for _, line := range lines {
		err := p.ProcessLine(line)
		assert.NoError(t, err)
	}

	txns := p.GetTransactions("111122223333")
	assert.Len(t, txns, 1)

	tx := txns[0]
	assert.Equal(t, time.Date(2023, 12, 29, 0, 0, 0, 0, time.UTC), tx.Date)
	assert.Equal(t, "THE TREE CENTER DES:PAYROLL", tx.Description)
	assert.Equal(t, 2774.40, tx.Amount)
	assert.Equal(t, TransactionDeposit, tx.Type)
}

func TestMultiLineTransaction(t *testing.T) {
	p := NewProcessor()

	// Setup account and section
	lines := []Line{
		{Page: 1, LineNumber: 1, Text: "Account number: " + testAccount1},
		{Page: 1, LineNumber: 2, Text: "Withdrawals and other subtractions"},
		{Page: 1, LineNumber: 3, Text: "01/02/24"},
		{Page: 1, LineNumber: 4, Text: "CHECKCARD 1225 MOXY KARLSRUHE"},
		{Page: 1, LineNumber: 5, Text: "KARLSRUHE 24463684001520020561103"},
		{Page: 1, LineNumber: 6, Text: "-798.43"},
	}

	for _, line := range lines {
		err := p.ProcessLine(line)
		assert.NoError(t, err)
	}

	txns := p.GetTransactions("111122223333")
	assert.Len(t, txns, 1)

	tx := txns[0]
	assert.Equal(t, time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), tx.Date)
	assert.Equal(t, "CHECKCARD 1225 MOXY KARLSRUHE KARLSRUHE 24463684001520020561103", tx.Description)
	assert.Equal(t, -798.43, tx.Amount)
	assert.Equal(t, TransactionWithdrawal, tx.Type)
}

func TestServiceFees(t *testing.T) {
	p := NewProcessor()

	// Setup account and section
	lines := []Line{
		{Page: 1, LineNumber: 1, Text: "Account number: " + testAccount1},
		{Page: 1, LineNumber: 2, Text: "Service fees"},
		{Page: 1, LineNumber: 3, Text: "12/29/23"},
		{Page: 1, LineNumber: 4, Text: "INTERNATIONAL TRANSACTION FEE"},
		{Page: 1, LineNumber: 5, Text: "-5.00"},
	}

	for _, line := range lines {
		err := p.ProcessLine(line)
		assert.NoError(t, err)
	}

	txns := p.GetTransactions("111122223333")
	assert.Len(t, txns, 1)

	tx := txns[0]
	assert.Equal(t, time.Date(2023, 12, 29, 0, 0, 0, 0, time.UTC), tx.Date)
	assert.Equal(t, "INTERNATIONAL TRANSACTION FEE", tx.Description)
	assert.Equal(t, -5.00, tx.Amount)
	assert.Equal(t, TransactionServiceFee, tx.Type)
}

func TestAmountParsing(t *testing.T) {
	p := NewProcessor()

	testCases := []struct {
		text     string
		expected float64
		valid    bool
	}{
		{"$1,234.56", 1234.56, true},
		{"1,234.56", 1234.56, true},
		{"-1,234.56", -1234.56, true},
		{"($1,234.56)", -1234.56, true},
		{"not an amount", 0, false},
	}

	for _, tc := range testCases {
		amount, ok := p.tryParseAmount(tc.text)
		assert.Equal(t, tc.valid, ok)
		if tc.valid {
			assert.Equal(t, tc.expected, amount)
		}
	}
}

func TestReset(t *testing.T) {
	p := NewProcessor()

	// Add some data
	p.ProcessLine(Line{Page: 1, LineNumber: 1, Text: "Account number: " + testAccount1})
	p.ProcessLine(Line{Page: 1, LineNumber: 2, Text: "Deposits and other additions"})
	p.ProcessLine(Line{Page: 1, LineNumber: 3, Text: "12/29/23"})
	p.ProcessLine(Line{Page: 1, LineNumber: 4, Text: "DEPOSIT"})
	p.ProcessLine(Line{Page: 1, LineNumber: 5, Text: "100.00"})

	// Verify data exists
	assert.NotEmpty(t, p.GetAccounts())
	assert.NotEmpty(t, p.GetTransactions("111122223333"))

	// Reset
	p.Reset()

	// Verify everything is cleared
	assert.Empty(t, p.GetAccounts())
	assert.Empty(t, p.GetTransactions("111122223333"))
	assert.Nil(t, p.GetAccountSummary("111122223333"))
}

func TestMultipleSections(t *testing.T) {
	p := NewProcessor()

	// Setup test data with multiple sections
	lines := []Line{
		{Page: 1, LineNumber: 1, Text: "Account number: " + testAccount1},
		// First section - deposits
		{Page: 1, LineNumber: 2, Text: "Deposits and other additions"},
		{Page: 1, LineNumber: 3, Text: "12/29/23"},
		{Page: 1, LineNumber: 4, Text: "DIRECT DEPOSIT FROM EMPLOYER"},
		{Page: 1, LineNumber: 5, Text: "2,500.00"},
		// Second section - withdrawals
		{Page: 1, LineNumber: 6, Text: "Withdrawals and other subtractions"},
		{Page: 1, LineNumber: 7, Text: "12/30/23"},
		{Page: 1, LineNumber: 8, Text: "ATM WITHDRAWAL"},
		{Page: 1, LineNumber: 9, Text: "-100.00"},
		// Third section - service fees
		{Page: 1, LineNumber: 10, Text: "Service fees"},
		{Page: 1, LineNumber: 11, Text: "12/31/23"},
		{Page: 1, LineNumber: 12, Text: "MONTHLY MAINTENANCE FEE"},
		{Page: 1, LineNumber: 13, Text: "-12.00"},
		// Back to deposits section
		{Page: 1, LineNumber: 14, Text: "Deposits and other additions"},
		{Page: 1, LineNumber: 15, Text: "01/02/24"},
		{Page: 1, LineNumber: 16, Text: "MOBILE DEPOSIT"},
		{Page: 1, LineNumber: 17, Text: "1,000.00"},
	}

	for _, line := range lines {
		err := p.ProcessLine(line)
		assert.NoError(t, err)
	}

	// Get all transactions
	txns := p.GetTransactions("111122223333")
	assert.Len(t, txns, 4)

	// Verify transactions are in the correct order and have correct types
	assert.Equal(t, TransactionDeposit, txns[0].Type)
	assert.Equal(t, 2500.00, txns[0].Amount)
	assert.Equal(t, "DIRECT DEPOSIT FROM EMPLOYER", txns[0].Description)

	assert.Equal(t, TransactionWithdrawal, txns[1].Type)
	assert.Equal(t, -100.00, txns[1].Amount)
	assert.Equal(t, "ATM WITHDRAWAL", txns[1].Description)

	assert.Equal(t, TransactionServiceFee, txns[2].Type)
	assert.Equal(t, -12.00, txns[2].Amount)
	assert.Equal(t, "MONTHLY MAINTENANCE FEE", txns[2].Description)

	assert.Equal(t, TransactionDeposit, txns[3].Type)
	assert.Equal(t, 1000.00, txns[3].Amount)
	assert.Equal(t, "MOBILE DEPOSIT", txns[3].Description)

	// Verify summary totals
	summary := p.GetAccountSummary("111122223333")
	assert.NotNil(t, summary)
	assert.Equal(t, 3500.00, summary.DepositsTotal)
	assert.Equal(t, -100.00, summary.WithdrawalsTotal)
	assert.Equal(t, -12.00, summary.ServiceFeesTotal)
}

func TestAnonymization(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"1234 5678 9012", "1234 XXXX 9012"},
		{"1234567890123456", "1234 XXXXXXXX 3456"},
		{"12345", "12345"}, // Too short to anonymize
		{"1234-5678-9012", "1234 XXXX 9012"},
	}

	for _, tc := range testCases {
		result := anonymizeAccountNumber(tc.input)
		assert.Equal(t, tc.expected, result)
	}
}

func TestAccountSummaryParsing(t *testing.T) {
	p := NewProcessor()

	// Setup test data with anonymized account summary
	lines := []Line{
		{Page: 1, LineNumber: 1, Text: "Account number: 1234 5678 9012"},
		{Page: 1, LineNumber: 2, Text: "Your Adv Plus Banking"},
		{Page: 1, LineNumber: 3, Text: "Preferred Rewards Platinum"},
		{Page: 1, LineNumber: 4, Text: "Account summary"},
		{Page: 1, LineNumber: 5, Text: "Beginning balance on December 28, 2023"},
		{Page: 1, LineNumber: 6, Text: "$33,192.29"},
		{Page: 1, LineNumber: 7, Text: "Deposits and other additions"},
		{Page: 1, LineNumber: 8, Text: "13,067.36"},
		{Page: 1, LineNumber: 9, Text: "Withdrawals and other subtractions"},
		{Page: 1, LineNumber: 10, Text: "-35,749.10"},
		{Page: 1, LineNumber: 11, Text: "Checks"},
		{Page: 1, LineNumber: 12, Text: "-0.00"},
		{Page: 1, LineNumber: 13, Text: "Service fees"},
		{Page: 1, LineNumber: 14, Text: "-173.07"},
		{Page: 1, LineNumber: 15, Text: "Ending balance on January 29, 2024"},
		{Page: 1, LineNumber: 16, Text: "$10,337.48"},
	}

	for _, line := range lines {
		err := p.ProcessLine(line)
		assert.NoError(t, err)
	}

	// Get account summary
	summary := p.GetAccountSummary("123456789012")
	assert.NotNil(t, summary)

	// Verify summary details
	assert.Equal(t, 33192.29, summary.BeginBalance)
	assert.Equal(t, 10337.48, summary.EndBalance)

	// Verify parsed totals
	assert.Equal(t, 13067.36, summary.ParsedDepositsTotal)
	assert.Equal(t, -35749.10, summary.ParsedWithdrawalsTotal)
	assert.Equal(t, -173.07, summary.ParsedServiceFeesTotal)

	// Verify account details
	account := p.GetAccounts()["123456789012"]
	assert.NotNil(t, account)
	assert.Equal(t, "Adv Plus Banking", account.Type)
	assert.Equal(t, 10337.48, account.Balance)

	// Verify balance equation
	expectedEndBalance := summary.BeginBalance +
		summary.ParsedDepositsTotal +
		summary.ParsedWithdrawalsTotal +
		summary.ParsedServiceFeesTotal
	assert.InDelta(t, expectedEndBalance, summary.EndBalance, 0.01)
}

func TestSummaryValidation(t *testing.T) {
	p := NewProcessor()

	// Setup test data with both transactions and summary section
	lines := []Line{
		{Page: 1, LineNumber: 1, Text: "Account number: " + testAccount1},

		// Account summary section
		{Page: 1, LineNumber: 2, Text: "Beginning balance on December 28, 2023"},
		{Page: 1, LineNumber: 3, Text: "$1,000.00"},
		{Page: 1, LineNumber: 4, Text: "Deposits and other additions"},
		{Page: 1, LineNumber: 5, Text: "2,500.00"},
		{Page: 1, LineNumber: 6, Text: "Withdrawals and other subtractions"},
		{Page: 1, LineNumber: 7, Text: "-1,200.00"},
		{Page: 1, LineNumber: 8, Text: "Service fees"},
		{Page: 1, LineNumber: 9, Text: "-12.00"},
		{Page: 1, LineNumber: 10, Text: "Ending balance on January 29, 2024"},
		{Page: 1, LineNumber: 11, Text: "$2,288.00"},

		// Actual transactions
		{Page: 1, LineNumber: 12, Text: "Deposits and other additions"},
		{Page: 1, LineNumber: 13, Text: "01/02/24"},
		{Page: 1, LineNumber: 14, Text: "DIRECT DEPOSIT"},
		{Page: 1, LineNumber: 15, Text: "2,500.00"},

		{Page: 1, LineNumber: 16, Text: "Withdrawals and other subtractions"},
		{Page: 1, LineNumber: 17, Text: "01/03/24"},
		{Page: 1, LineNumber: 18, Text: "ATM WITHDRAWAL"},
		{Page: 1, LineNumber: 19, Text: "-1,200.00"},

		{Page: 1, LineNumber: 20, Text: "Service fees"},
		{Page: 1, LineNumber: 21, Text: "01/04/24"},
		{Page: 1, LineNumber: 22, Text: "MONTHLY FEE"},
		{Page: 1, LineNumber: 23, Text: "-12.00"},
	}

	for _, line := range lines {
		err := p.ProcessLine(line)
		assert.NoError(t, err)
	}

	// First verify individual transactions
	txns := p.GetTransactions("111122223333")
	assert.Len(t, txns, 3, "Should have exactly 3 transactions")

	// Verify deposit transaction
	assert.Equal(t, TransactionDeposit, txns[0].Type)
	assert.Equal(t, time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), txns[0].Date)
	assert.Equal(t, "DIRECT DEPOSIT", txns[0].Description)
	assert.Equal(t, 2500.00, txns[0].Amount)

	// Verify withdrawal transaction
	assert.Equal(t, TransactionWithdrawal, txns[1].Type)
	assert.Equal(t, time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC), txns[1].Date)
	assert.Equal(t, "ATM WITHDRAWAL", txns[1].Description)
	assert.Equal(t, -1200.00, txns[1].Amount)

	// Verify service fee transaction
	assert.Equal(t, TransactionServiceFee, txns[2].Type)
	assert.Equal(t, time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC), txns[2].Date)
	assert.Equal(t, "MONTHLY FEE", txns[2].Description)
	assert.Equal(t, -12.00, txns[2].Amount)

	// Get account summary
	summary := p.GetAccountSummary("111122223333")
	assert.NotNil(t, summary)

	// Verify parsed totals
	assert.Equal(t, 1000.00, summary.BeginBalance)
	assert.Equal(t, 2500.00, summary.ParsedDepositsTotal)
	assert.Equal(t, -1200.00, summary.ParsedWithdrawalsTotal)
	assert.Equal(t, -12.00, summary.ParsedServiceFeesTotal)
	assert.Equal(t, 2288.00, summary.EndBalance)

	// Verify calculated totals match parsed totals
	assert.Equal(t, summary.ParsedDepositsTotal, summary.DepositsTotal)
	assert.Equal(t, summary.ParsedWithdrawalsTotal, summary.WithdrawalsTotal)
	assert.Equal(t, summary.ParsedServiceFeesTotal, summary.ServiceFeesTotal)

	// Validate totals using the validation methods
	assert.NoError(t, p.ValidateSummaryTotals("111122223333"))
	assert.NoError(t, p.ValidateParsedTotals("111122223333"))

	// Test validation with mismatched data
	p.summaries["111122223333"].ParsedDepositsTotal += 0.01
	assert.Error(t, p.ValidateParsedTotals("111122223333"))
}
