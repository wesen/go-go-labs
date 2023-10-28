package main

import "github.com/spf13/cobra"

// https://chat.openai.com/share/523b6523-8984-4a06-ba40-ddeb98cc019b

import (
	"bufio"
	"fmt"
	"github.com/go-go-golems/clay/pkg"
	"github.com/spf13/viper"
	"os"
	"strings"
	"text/template"
	"time"
)

const mdTemplate = `
## {{.Title}} 
Date: {{.CurrentDateTime}}

{{- if .Descriptions -}}
### Description

{{range .Descriptions}}
- {{.}}
{{- end}}
{{- end}}

{{- if .Tags -}}
### Tags

{{range .Tags}}
- {{.}}
{{- end}}
{{- end}}

{{if .FileContent -}}
{{.FileContent}}
{{- end -}}
`

type EntryData struct {
	Title           string
	Descriptions    []string
	Tags            []string
	FileContent     string
	CurrentDateTime string
}

func GetOutputFilename(cmd *cobra.Command) string {
	if output, _ := cmd.Flags().GetString("output"); output != "" {
		return output
	}

	if envOutput := os.Getenv("CAPTURE_OUTPUT"); envOutput != "" {
		return envOutput
	}

	err := viper.ReadInConfig() // Find and read the config file
	if err == nil {             // Handle errors reading the config file
		if configFileOutput := viper.GetString("output"); configFileOutput != "" {
			return configFileOutput
		}
	}

	return "log-{{DATE}}.md"
}

func main() {
	var (
		title       string
		description []string
		tags        []string
	)

	var rootCmd = &cobra.Command{
		Use:   "capture",
		Short: "CLI tool to append file contents and metadata to foobar.md",
		Long: `A simple CLI tool that takes a list of files along with title, description,
and tags and appends them to a markdown file called foobar.md.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Aggregate content from files
			var contentBuilder strings.Builder

			for i, file := range args {
				if file != "-" {
					contentBuilder.WriteString(fmt.Sprintf("### File: %s\n\n", file))
				}

				if file == "-" {
					scanner := bufio.NewScanner(os.Stdin)
					for scanner.Scan() {
						contentBuilder.WriteString(scanner.Text() + "\n")
					}
				} else {
					data, err := os.ReadFile(file)
					if err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
						continue
					}
					contentBuilder.WriteString(string(data) + "\n")
				}

				if i < len(args)-1 {
					contentBuilder.WriteString("\n---\n\n")
				}
			}

			currentDate := time.Now().Format("2006-01-02")
			currentDateTime := time.Now().Format("2006-01-02 15:04:05")

			data := EntryData{
				Title:           title,
				Descriptions:    description,
				Tags:            tags,
				FileContent:     contentBuilder.String(),
				CurrentDateTime: currentDateTime,
			}

			// Parse and execute the template
			tmpl, err := template.New("markdown").Parse(mdTemplate)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error parsing template: %s\n", err)
				return
			}

			var renderedContent strings.Builder
			if err := tmpl.Execute(&renderedContent, data); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error executing template: %s\n", err)
				return
			}

			// Determine the output file
			outputFilename := GetOutputFilename(cmd)

			outputFilename = strings.Replace(outputFilename, "{{DATE}}", currentDate, -1)

			// Append to the determined output file or create it if it doesn't exist
			f, err := os.OpenFile(outputFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error opening %s: %s\n", outputFilename, err)
				return
			}
			defer func(f *os.File) {
				_ = f.Close()
			}(f)

			if _, err := f.WriteString(renderedContent.String()); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error writing to %s: %s\n", outputFilename, err)
				return
			}

			fmt.Printf("Content appended to %s successfully.\n", outputFilename)
		},
	}

	rootCmd.Flags().StringVarP(&title, "title", "t", "", "Title for the entry")
	rootCmd.Flags().StringArrayVarP(&description, "description", "d", []string{}, "Description for the entry (repeatable)")
	rootCmd.Flags().StringArrayVarP(&tags, "tags", "g", []string{}, "Tags for the entry (repeatable)")
	rootCmd.Flags().StringP("output", "o", "", "Specify an output file (overrides the CAPTURE_OUTPUT environment variable and config file)")

	_ = rootCmd.MarkFlagRequired("title")

	err := pkg.InitViper("capture", rootCmd)
	cobra.CheckErr(err)

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
