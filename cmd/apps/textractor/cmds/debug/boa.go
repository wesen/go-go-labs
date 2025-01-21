package debug

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-go-golems/go-go-labs/cmd/apps/textractor/pkg/textract/boa"
	"github.com/go-go-golems/go-go-labs/cmd/apps/textractor/pkg/textract/parser"
	"github.com/spf13/cobra"
)

type textractPage struct {
	page parser.Page
}

func (t *textractPage) Lines() []boa.Line {
	lines := t.page.Lines()
	result := make([]boa.Line, len(lines))
	for i, line := range lines {
		result[i] = boa.Line{
			Text:       line.Text(),
			Page:       t.page.Number(),
			LineNumber: i + 1,
		}
	}
	return result
}

func processCSVFile(processor *boa.StatementProcessor, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("opening CSV file %s: %w", filename, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("reading CSV file %s: %w", filename, err)
	}

	// Skip header row
	if len(records) > 0 {
		records = records[1:]
	}

	for _, record := range records {
		if len(record) != 3 {
			continue // Skip malformed rows
		}

		page, err := strconv.Atoi(record[0])
		if err != nil {
			continue // Skip rows with invalid page numbers
		}

		lineNumber, err := strconv.Atoi(record[1])
		if err != nil {
			continue // Skip rows with invalid line numbers
		}

		line := boa.Line{
			Page:       page,
			LineNumber: lineNumber,
			Text:       record[2],
		}

		if err := processor.ProcessLine(line); err != nil {
			return fmt.Errorf("processing line %d on page %d: %w", lineNumber, page, err)
		}
	}

	return nil
}

func processJSONFile(processor *boa.StatementProcessor, filename string) error {
	docs, err := parser.LoadFromJSON(filename)
	if err != nil {
		return fmt.Errorf("loading JSON from %s: %w", filename, err)
	}

	for _, doc := range docs {
		for _, page := range doc.Pages() {
			tPage := &textractPage{page: page}
			for _, line := range tPage.Lines() {
				if err := processor.ProcessLine(line); err != nil {
					return fmt.Errorf("processing line on page %d: %w", page.Number(), err)
				}
			}
		}
	}

	return nil
}

func newBoACommand() *cobra.Command {
	var outputDir string

	cmd := &cobra.Command{
		Use:   "boa [flags] file...",
		Short: "Process Bank of America statements",
		Long:  "Parse Bank of America statements from Textract JSON files or CSV files and extract transactions and summaries",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("at least one input file is required")
			}

			// If no output directory specified, use first input file's directory + "boa-output/"
			if outputDir == "" {
				outputDir = filepath.Join(filepath.Dir(args[0]), "boa-output")
			}

			// Create output directory if it doesn't exist
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("creating output directory: %w", err)
			}

			processor := boa.NewProcessor()

			// Process all input files
			for _, file := range args {
				ext := filepath.Ext(file)
				var err error

				fmt.Printf("Processing %s...\n", file)
				switch ext {
				case ".csv":
					err = processCSVFile(processor, file)
				case ".json":
					err = processJSONFile(processor, file)
				default:
					return fmt.Errorf("unsupported file type for %s: %s", file, ext)
				}

				if err != nil {
					return err
				}
			}

			// Print summary for each account
			for accNum, acc := range processor.GetAccounts() {
				fmt.Printf("\nAccount: %s (found on page %d)\n", accNum, acc.DetailsPage)

				summary := processor.GetAccountSummary(accNum)
				if summary != nil {
					fmt.Printf("Statement period: %s\n", summary.Period)
					fmt.Printf("Beginning balance: $%.2f\n", summary.BeginBalance)
					fmt.Printf("Ending balance: $%.2f\n", summary.EndBalance)

					txns := processor.GetTransactions(accNum)
					var deposits, withdrawals, fees int
					for _, tx := range txns {
						switch tx.Type {
						case boa.TransactionDeposit:
							deposits++
						case boa.TransactionWithdrawal:
							withdrawals++
						case boa.TransactionServiceFee:
							fees++
						}
					}

					fmt.Printf("\nTransactions:\n")
					fmt.Printf("- Deposits: %d transactions totaling $%.2f (parsed from summary: $%.2f)\n",
						deposits, summary.DepositsTotal, summary.ParsedDepositsTotal)
					fmt.Printf("- Withdrawals: %d transactions totaling $%.2f (parsed from summary: $%.2f)\n",
						withdrawals, summary.WithdrawalsTotal, summary.ParsedWithdrawalsTotal)
					fmt.Printf("- Service fees: %d transactions totaling $%.2f (parsed from summary: $%.2f)\n",
						fees, summary.ServiceFeesTotal, summary.ParsedServiceFeesTotal)

					// Validate totals
					if err := processor.ValidateSummaryTotals(accNum); err != nil {
						fmt.Printf("⚠️  Warning: %v\n", err)
					} else {
						fmt.Printf("✅ All transaction totals validated successfully\n")
					}

					// Validate parsed totals
					if err := processor.ValidateParsedTotals(accNum); err != nil {
						fmt.Printf("⚠️  Warning: %v\n", err)
					} else {
						fmt.Printf("✅ All parsed totals validated successfully\n")
					}
				}
			}

			// Write output files and print their locations
			if err := processor.WriteOutput(outputDir); err != nil {
				return fmt.Errorf("writing output: %w", err)
			}

			fmt.Printf("\nOutput files written to:\n")
			for accNum := range processor.GetAccounts() {
				fmt.Printf("- Summary: %s\n", filepath.Join(outputDir, fmt.Sprintf("summary-%s.json", accNum)))
				fmt.Printf("- Deposits: %s\n", filepath.Join(outputDir, fmt.Sprintf("deposits-%s.csv", accNum)))
				fmt.Printf("- Withdrawals: %s\n", filepath.Join(outputDir, fmt.Sprintf("withdrawals-%s.csv", accNum)))
				fmt.Printf("- Service fees: %s\n", filepath.Join(outputDir, fmt.Sprintf("service-fees-%s.csv", accNum)))
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "Directory to write output files (default: input-file-dir/boa-output)")
	return cmd
}
