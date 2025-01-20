package debug

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-go-golems/go-go-labs/cmd/apps/textractor/pkg/textract/parser"
	"github.com/spf13/cobra"
)

func newCSVCommand() *cobra.Command {
	var outputFile string

	cmd := &cobra.Command{
		Use:   "csv [flags] file...",
		Short: "Export tables and lines from Textract JSON files to CSV",
		Long:  "Parses Textract JSON files and exports all tables and lines to CSV format, merging tables with matching headers",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("at least one input file is required")
			}

			// Generate lines output filename
			ext := filepath.Ext(outputFile)
			base := strings.TrimSuffix(outputFile, ext)
			linesOutputFile := fmt.Sprintf("%s-lines%s", base, ext)

			// Create lines output file
			linesFile, err := os.Create(linesOutputFile)
			if err != nil {
				return fmt.Errorf("creating lines output file %s: %w", linesOutputFile, err)
			}
			defer linesFile.Close()

			linesWriter := csv.NewWriter(linesFile)
			// Write header
			if err := linesWriter.Write([]string{"page", "lineno", "text"}); err != nil {
				return fmt.Errorf("writing lines CSV header: %w", err)
			}

			// Replace the tablesByHeader map with a slice to maintain order
			type TableData struct {
				headers    []string
				rows       [][]string
				sourceInfo []string // Store file and page information
			}
			var tables []TableData

			// Process all input files
			for _, file := range args {
				docs, err := parser.LoadFromJSON(file)
				if err != nil {
					return fmt.Errorf("loading JSON from %s: %w", file, err)
				}

				// Process each document
				for _, doc := range docs {
					for _, page := range doc.Pages() {
						// Process lines
						for i, line := range page.Lines() {
							record := []string{
								fmt.Sprintf("%d", page.Number()),
								fmt.Sprintf("%d", i+1),
								strings.TrimSpace(line.Text()),
							}
							if err := linesWriter.Write(record); err != nil {
								return fmt.Errorf("writing line to CSV: %w", err)
							}
						}

						for _, table := range page.Tables() {
							if table.RowCount() == 0 {
								continue
							}

							headerCells := table.GetHeaders()
							var headerRow []string
							var dataRows []parser.TableRow

							if headerCells != nil {
								headerRow = make([]string, len(headerCells))
								for i, cell := range headerCells {
									headerRow[i] = strings.TrimSpace(cell.Text())
								}
								dataRows = table.Rows()[1:] // Skip header row
							} else {
								headerRow = make([]string, table.ColumnCount())
								for i := range headerRow {
									headerRow[i] = fmt.Sprintf("Column%d", i+1)
								}
								dataRows = table.Rows() // Use all rows
							}

							// Convert rows to string matrix
							var rows [][]string
							for _, row := range dataRows {
								stringRow := make([]string, len(row.Cells()))
								for i, cell := range row.Cells() {
									stringRow[i] = strings.TrimSpace(cell.Text())
								}
								rows = append(rows, stringRow)
							}

							sourceInfo := fmt.Sprintf("# Source: %s, Page: %d", file, page.Number())

							// Check if we can append to the last table
							if len(tables) > 0 {
								lastTable := &tables[len(tables)-1]
								if len(lastTable.headers) == len(headerRow) {
									match := true
									for i := range headerRow {
										if lastTable.headers[i] != headerRow[i] {
											match = false
											break
										}
									}
									if match {
										lastTable.rows = append(lastTable.rows, rows...)
										lastTable.sourceInfo = append(lastTable.sourceInfo, sourceInfo)
										continue
									}
								}
							}

							// Create new table entry
							tables = append(tables, TableData{
								headers:    headerRow,
								rows:       rows,
								sourceInfo: []string{sourceInfo},
							})
						}
					}
				}
			}

			linesWriter.Flush()
			if err := linesWriter.Error(); err != nil {
				return fmt.Errorf("flushing lines CSV writer: %w", err)
			}

			// Write tables to CSV files
			baseFile := outputFile
			if baseFile == "" {
				baseFile = "output.csv"
			}

			for i, table := range tables {
				filename := baseFile
				if len(tables) > 1 {
					filename = fmt.Sprintf("%s-%d%s", base, i+1, ext)
				}

				f, err := os.Create(filename)
				if err != nil {
					return fmt.Errorf("creating output file %s: %w", filename, err)
				}
				defer f.Close()

				// Write source information as comments
				for _, info := range table.sourceInfo {
					if _, err := fmt.Fprintln(f, info); err != nil {
						return fmt.Errorf("writing source info to %s: %w", filename, err)
					}
				}

				w := csv.NewWriter(f)
				if err := w.Write(table.headers); err != nil {
					return fmt.Errorf("writing headers to %s: %w", filename, err)
				}
				if err := w.WriteAll(table.rows); err != nil {
					return fmt.Errorf("writing rows to %s: %w", filename, err)
				}
				w.Flush()

				if err := w.Error(); err != nil {
					return fmt.Errorf("flushing CSV writer for %s: %w", filename, err)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "output.csv", "Output file path for tables (lines will be written to {output}-lines.csv)")
	return cmd
}
