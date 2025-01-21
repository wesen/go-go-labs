I'll create a comprehensive specification document that merges the CSV format understanding with the bank statement processing requirements.

# Bank Statement Processor Specification

## 1. Input Format

### 1.1 CSV Structure
The input is a CSV file with three columns:
```
page,lineno,text
```
- `page`: Integer, 1-based page number in the statement
- `lineno`: Integer, 1-based line number within the page
- `text`: String, the actual content of the line

Example:
```csv
1,1,BANK OF AMERICA
1,2,Preferred Rewards
3,2,Account number: 0012 2984 9384 8383
```

### 1.2 Line Ordering
- Lines are ordered by page number first, then by line number
- Line numbers restart at 1 for each new page
- Lines must be processed sequentially to maintain context
- Multi-line entries (like transactions) span consecutive line numbers

## 2. Data Structures

### 2.1 Input Line Structure
```go
type Line struct {
    Page       int    // Page number (1-based)
    LineNumber int    // Line number within page (1-based)
    Text       string // Raw text content
}
```

### 2.2 Account Structure
```go
type Account struct {
    Number      string  // Normalized account number (spaces/dashes removed)
    Type        string  // e.g. "Adv Plus Banking", "Regular Savings"
    Balance     float64 // Current balance
    DetailsPage int     // Page number where details start
}
```

### 2.3 Transaction Structure
```go
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

type TransactionType string
const (
    TransactionDeposit    TransactionType = "deposit"
    TransactionWithdrawal TransactionType = "withdrawal"
    TransactionServiceFee TransactionType = "service_fee"
)
```

### 2.4 Account Summary Structure
```go
type AccountSummary struct {
    AccountNumber    string
    Period          string    // Statement period
    BeginBalance    float64
    EndBalance      float64
    DepositsTotal   float64
    WithdrawalsTotal float64
    ServiceFeesTotal float64
    
    // New fields to track parsed vs calculated totals
    ParsedDepositsTotal    float64
    ParsedWithdrawalsTotal float64
    ParsedServiceFeesTotal float64
}
```

## 3. Line Pattern Recognition

### 3.1 Account Identification
```
Pattern 1: "Account number: XXXX XXXX XXXX"
Pattern 2: "Account # XXXX XXXX XXXX"
Example: 3,2,Account number: 1234 5678 9012 2345
```

### 3.2 Statement Period
```
Pattern: "for [Month DD, YYYY] to [Month DD, YYYY]"
Example: 1,18,"for December 28, 2023 to January 29, 2024"
```

### 3.3 Section Headers
```
Pattern 1: "Deposits and other additions"
Pattern 2: "Withdrawals and other subtractions"
Pattern 3: "Service fees"
```

### 3.4 Transaction Formats
```
Single-line:
Page,LineNo,MM/DD/YY
Page,LineNo,Description
Page,LineNo,Amount

Multi-line:
Page,LineNo,MM/DD/YY
Page,LineNo+1,Description Line 1
Page,LineNo+2,Description Line 2
...
Page,LineNo+n,Amount
```

### 3.5 Amount Formats
```
Positive: "$1,234.56" or "1,234.56"
Negative: "-1,234.56" or "($1,234.56)"
```

## 4. Processing Rules

### 4.1 State Management
```go
type ProcessorState struct {
    CurrentAccount     string
    CurrentSection     string // "deposits", "withdrawals", "fees"
    CurrentTransaction *Transaction
    PageContext       int
    LastLineNumber    int
}
```

### 4.2 Transaction Building Rules
1. Start new transaction when:
   - Date pattern (MM/DD/YY) is detected
   - Previous transaction is complete (has amount)

2. Continue transaction when:
   - Line is part of description (no date/amount pattern)
   - Line contains amount for current transaction

3. Complete transaction when:
   - Amount is found
   - New date pattern is detected
   - Section changes

### 4.3 Multi-line Entry Handling
1. Track consecutive lines by page and line number
2. Combine description lines until amount is found
3. Handle page breaks in multi-line entries
4. Maintain context across "continued on next page" markers

## 5. Public Interface

```go
type StatementProcessor interface {
    // Process a single line from the CSV
    ProcessLine(line Line) error
    
    // Get all detected accounts
    GetAccounts() map[string]*Account
    
    // Get transactions for an account
    GetTransactions(accountNumber string) []Transaction
    
    // Get account summary
    GetAccountSummary(accountNumber string) *AccountSummary
    
    // Reset processor state
    Reset()
}
```

## 6. Processing Steps

### 6.1 Line Processing
1. Parse CSV line into Line struct
2. Validate page and line numbers
3. Check for section headers
4. Process line based on current state
5. Update state based on line content

### 6.2 Transaction Processing
1. Detect transaction start (date pattern)
2. Accumulate description lines
3. Parse amount when found
4. Validate transaction completeness
5. Store in appropriate collection

### 6.3 Summary Processing
1. Track running totals for each section
2. Validate against stated totals
3. Update account summaries
4. Handle balance calculations

## 7. Error Handling

### 7.1 Required Error Checks
1. Invalid page/line numbers
2. Malformed amounts
3. Invalid dates
4. Missing required fields
5. Section total mismatches
6. Incomplete transactions
7. Invalid account numbers

### 7.2 Error Types
```go
type ProcessingError struct {
    Page       int
    LineNumber int
    ErrorType  string
    Message    string
    Raw        string
}
```

## 8. Validation Requirements

### 8.1 Data Validation
1. All amounts must parse to valid float64
2. All dates must parse to valid time.Time
3. Account numbers must be properly normalized
4. Transaction totals must match summary totals:
   - Calculated totals from individual transactions
   - Parsed totals from summary section
   - Both should match and be validated
5. Beginning + Credits - Debits = Ending balance

### 8.2 Summary Validation Rules
1. Store both parsed and calculated totals separately
2. When processing summary section:
   - Parse "Deposits and other additions" total into ParsedDepositsTotal
   - Parse "Withdrawals and other subtractions" total into ParsedWithdrawalsTotal
   - Parse "Service fees" total into ParsedServiceFeesTotal
3. When processing transactions:
   - Accumulate into DepositsTotal, WithdrawalsTotal, ServiceFeesTotal
4. Validate that parsed totals match calculated totals within epsilon (0.01)
5. Report discrepancies with detailed breakdown of differences

### 8.3 Summary Section Format
```
Account summary
Beginning balance on [Date]              $X,XXX.XX
Deposits and other additions             $X,XXX.XX
Withdrawals and other subtractions      -$X,XXX.XX
Checks                                  -$X,XXX.XX
Service fees                            -$X,XXX.XX
Ending balance on [Date]                 $X,XXX.XX
```

### 8.4 Validation Methods
```go
// ValidateSummaryTotals checks if transaction totals match summary amounts
// Returns nil if totals match, otherwise returns error with details
func (p *StatementProcessor) ValidateSummaryTotals(accountNumber string) error

// ValidateParsedTotals checks if parsed summary totals match calculated totals
// Returns nil if totals match, otherwise returns error with details
func (p *StatementProcessor) ValidateParsedTotals(accountNumber string) error
```

This specification provides a complete framework for processing bank statement data from the CSV format while maintaining accuracy and proper organization of financial information.
