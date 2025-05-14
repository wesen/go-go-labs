package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type ExecuteOnceSettings struct {
	FilePath      string `glazed.parameter:"file"`
	Script        string `glazed.parameter:"script"`
	ChromeURL     string `glazed.parameter:"chrome-url"`
	NavigateToURL string `glazed.parameter:"navigate-to"`
}

type ExecuteOnceCommand struct {
	*cmds.CommandDescription
	settings ExecuteOnceSettings
}

func newExecuteOnceCommand() (*cobra.Command, error) {
	cmd := &ExecuteOnceCommand{
		CommandDescription: cmds.NewCommandDescription(
			"execute",
			cmds.WithShort("Execute a JavaScript file once in Chrome"),
			cmds.WithLong("Execute a JavaScript file or provided script once in Chrome using the Chrome DevTools Protocol."),
			cmds.WithFlags(
				parameters.NewParameterDefinition(
					"file",
					parameters.ParameterTypeString,
					parameters.WithHelp("Path to JavaScript file to execute"),
					parameters.WithDefault(""),
				),
				parameters.NewParameterDefinition(
					"script",
					parameters.ParameterTypeString,
					parameters.WithHelp("JavaScript script to execute (alternative to file)"),
					parameters.WithDefault(""),
				),
				parameters.NewParameterDefinition(
					"chrome-url",
					parameters.ParameterTypeString,
					parameters.WithHelp("Chrome DevTools Protocol URL"),
					parameters.WithDefault("http://localhost:9222"),
				),
				parameters.NewParameterDefinition(
					"navigate-to",
					parameters.ParameterTypeString,
					parameters.WithHelp("URL to navigate to before executing JavaScript (optional)"),
					parameters.WithDefault(""),
				),
			),
		),
	}

	return cli.BuildCobraCommandFromBareCommand(cmd)
}

func (c *ExecuteOnceCommand) Run(ctx context.Context, parsedLayers *layers.ParsedLayers) error {
	if err := parsedLayers.InitializeStruct(layers.DefaultSlug, &c.settings); err != nil {
		return err
	}

	// Validate input - need either file or script
	if c.settings.FilePath == "" && c.settings.Script == "" {
		return fmt.Errorf("either --file or --script must be provided")
	}

	// Get JavaScript to execute
	javascript := c.settings.Script
	if c.settings.FilePath != "" {
		// Read from file
		fileInfo, err := os.Stat(c.settings.FilePath)
		if err != nil {
			return fmt.Errorf("error accessing file %s: %w", c.settings.FilePath, err)
		}
		if fileInfo.IsDir() {
			return fmt.Errorf("%s is a directory, not a file", c.settings.FilePath)
		}

		content, err := os.ReadFile(c.settings.FilePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", c.settings.FilePath, err)
		}
		javascript = string(content)
	}

	// Log settings
	sourceType := "script"
	sourceName := "[inline-script]"
	if c.settings.FilePath != "" {
		sourceType = "file"
		sourceName = c.settings.FilePath
	}

	log.Info().Str("sourceType", sourceType).
		Str("source", sourceName).
		Str("chromeURL", c.settings.ChromeURL).
		Str("navigateToURL", c.settings.NavigateToURL).
		Msg("Executing JavaScript")

	// Create Chrome executor
	executor, err := NewChromeExecutor(c.settings.ChromeURL, c.settings.NavigateToURL)
	if err != nil {
		return fmt.Errorf("failed to create Chrome executor: %w", err)
	}
	defer executor.Close()

	// Execute the JavaScript
	result, err := executor.ExecuteJavaScript(ctx, javascript)
	if err != nil {
		return fmt.Errorf("failed to execute JavaScript: %w", err)
	}

	log.Info().Str("source", sourceName).Msg("JavaScript execution completed")

	// Print the results
	fmt.Printf("JavaScript Execution Results:\n")
	fmt.Printf("  Source: %s\n", sourceName)
	fmt.Printf("  Result: %s\n", result)

	return nil
}
