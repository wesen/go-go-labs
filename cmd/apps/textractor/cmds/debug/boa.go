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

			processor := boa.NewProcessor()

			// Process all input files
			for _, file := range args {
				ext := filepath.Ext(file)
				var err error

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

			// Write output files
			if err := processor.WriteOutput(outputDir); err != nil {
				return fmt.Errorf("writing output: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output-dir", "o", ".", "Directory to write output files")
	return cmd
}
