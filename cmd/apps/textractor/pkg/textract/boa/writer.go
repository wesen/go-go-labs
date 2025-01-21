package boa

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// WriteOutput writes all processed data to files in the specified directory
func (p *StatementProcessor) WriteOutput(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Write data for each account
	for accNum := range p.accounts {
		// Write summary JSON
		if err := p.writeSummaryJSON(outputDir, accNum); err != nil {
			return fmt.Errorf("writing summary for account %s: %w", accNum, err)
		}

		// Get all transactions for this account
		txns := p.transactions[accNum]

		// Split transactions by type
		var deposits, withdrawals, serviceFees []Transaction
		for _, tx := range txns {
			switch tx.Type {
			case TransactionDeposit:
				deposits = append(deposits, tx)
			case TransactionWithdrawal:
				withdrawals = append(withdrawals, tx)
			case TransactionServiceFee:
				serviceFees = append(serviceFees, tx)
			}
		}

		// Write deposits CSV
		if err := p.writeTransactionsCSV(outputDir, accNum, "deposits", deposits); err != nil {
			return fmt.Errorf("writing deposits for account %s: %w", accNum, err)
		}

		// Write withdrawals CSV
		if err := p.writeTransactionsCSV(outputDir, accNum, "withdrawals", withdrawals); err != nil {
			return fmt.Errorf("writing withdrawals for account %s: %w", accNum, err)
		}

		// Write service fees CSV
		if err := p.writeTransactionsCSV(outputDir, accNum, "service-fees", serviceFees); err != nil {
			return fmt.Errorf("writing service fees for account %s: %w", accNum, err)
		}
	}

	return nil
}

func (p *StatementProcessor) writeSummaryJSON(outputDir, accNum string) error {
	filename := filepath.Join(outputDir, fmt.Sprintf("summary-%s.json", accNum))

	summary := p.summaries[accNum]
	if summary == nil {
		summary = &AccountSummary{
			AccountNumber: accNum,
		}
	}

	// Add account info
	if acc := p.accounts[accNum]; acc != nil {
		summary.EndBalance = acc.Balance
	}

	// Calculate totals if not set
	if summary.DepositsTotal == 0 || summary.WithdrawalsTotal == 0 || summary.ServiceFeesTotal == 0 {
		for _, tx := range p.transactions[accNum] {
			switch tx.Type {
			case TransactionDeposit:
				summary.DepositsTotal += tx.Amount
			case TransactionWithdrawal:
				summary.WithdrawalsTotal += tx.Amount
			case TransactionServiceFee:
				summary.ServiceFeesTotal += tx.Amount
			}
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(summary)
}

func (p *StatementProcessor) writeTransactionsCSV(outputDir, accNum, txType string, transactions []Transaction) error {
	if len(transactions) == 0 {
		return nil // Skip empty transaction lists
	}

	filename := filepath.Join(outputDir, fmt.Sprintf("%s-%s.csv", txType, accNum))
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// Write header
	if err := w.Write([]string{"Date", "Page", "Line", "Description", "Amount"}); err != nil {
		return err
	}

	// Sort transactions by date
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Date.Before(transactions[j].Date)
	})

	// Write transactions
	for _, tx := range transactions {
		record := []string{
			tx.Date.Format("2006-01-02"),
			fmt.Sprintf("%d", tx.Page),
			fmt.Sprintf("%d", tx.LineNumber),
			tx.Description,
			fmt.Sprintf("%.2f", tx.Amount),
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}

	return nil
}
