package boa

// Line represents a single line of text from a statement
type Line struct {
	Text       string
	Page       int
	LineNumber int
}

// Table represents a table from a statement
type Table struct {
	Headers []string
	Rows    [][]string
	Page    int
}

// StatementPage represents a page from a bank statement
type StatementPage interface {
	Lines() []Line
	Tables() []Table
	Number() int
}
