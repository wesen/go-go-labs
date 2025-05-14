package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-go-golems/clay/pkg/watcher"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type WatchSettings struct {
	Directory    string `glazed.parameter:"directory"`
	Pattern      string `glazed.parameter:"pattern"`
	ChromeURL    string `glazed.parameter:"chrome-url"`
	NavigateToURL string `glazed.parameter:"navigate-to"`
	ThrottleMs   int    `glazed.parameter:"throttle-ms"`
}

type WatchCommand struct {
	*cmds.CommandDescription
	settings WatchSettings
}

func newWatchCommand() (*cobra.Command, error) {
	glazeLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, fmt.Errorf("could not create Glazed parameter layer: %w", err)
	}

	cmd := &WatchCommand{
		CommandDescription: cmds.NewCommandDescription(
			"watch",
			cmds.WithShort("Watch JavaScript files and execute them in Chrome"),
			cmds.WithLong("Watch a directory for JavaScript file changes and execute them in Chrome using the Chrome DevTools Protocol."),
			cmds.WithFlags(
				parameters.NewParameterDefinition(
					"directory",
					parameters.ParameterTypeString,
					parameters.WithHelp("Directory to watch for JavaScript files"),
					parameters.WithRequired(true),
				),
				parameters.NewParameterDefinition(
					"pattern",
					parameters.ParameterTypeString,
					parameters.WithHelp("File pattern to watch (e.g., '**/*.js')"),
					parameters.WithDefault("**/*.js"),
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
				parameters.NewParameterDefinition(
					"throttle-ms",
					parameters.ParameterTypeInteger,
					parameters.WithHelp("Throttle execution to prevent multiple runs on rapid file changes (milliseconds)"),
					parameters.WithDefault(500),
				),
			),
			cmds.WithLayersList(glazeLayer),
		),
	}

	return cli.BuildCobraCommandFromGlazeCommand(cmd)
}

func (c *WatchCommand) RunIntoGlazeProcessor(ctx context.Context, parsedLayers *layers.ParsedLayers, gp middlewares.Processor) error {
	if err := parsedLayers.InitializeStruct(layers.DefaultSlug, &c.settings); err != nil {
		return err
	}

	// Validate directory exists
	dirInfo, err := os.Stat(c.settings.Directory)
	if err != nil {
		return fmt.Errorf("error accessing directory %s: %w", c.settings.Directory, err)
	}
	if !dirInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", c.settings.Directory)
	}

	// Log settings
	log.Info().Str("directory", c.settings.Directory).
		Str("pattern", c.settings.Pattern).
		Str("chromeURL", c.settings.ChromeURL).
		Str("navigateToURL", c.settings.NavigateToURL).
		Int("throttleMs", c.settings.ThrottleMs).
		Msg("Starting file watcher")

	// Create Chrome executor
	executor, err := NewChromeExecutor(c.settings.ChromeURL, c.settings.NavigateToURL)
	if err != nil {
		return fmt.Errorf("failed to create Chrome executor: %w", err)
	}
	defer executor.Close()

	// Create a throttled execution function
	lastExecution := time.Now().Add(-24 * time.Hour) // Initialize to a day ago
	mutex := &sync.Mutex{}

	handleFileChange := func(path string) error {
		mutex.Lock()
		now := time.Now()
		elapsedMs := now.Sub(lastExecution).Milliseconds()
		if elapsedMs < int64(c.settings.ThrottleMs) {
			mutex.Unlock()
			log.Debug().Str("path", path).Int64("elapsedMs", elapsedMs).Msg("Skipping execution due to throttling")
			return nil
		}
		lastExecution = now
		mutex.Unlock()

		log.Info().Str("path", path).Msg("Detected JavaScript file change")

		// Read the file
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Execute the JavaScript
		result, err := executor.ExecuteJavaScript(ctx, string(content))
		if err != nil {
			return fmt.Errorf("failed to execute JavaScript from file %s: %w", path, err)
		}

		filename := filepath.Base(path)
		log.Info().Str("file", filename).Str("result", result).Msg("JavaScript execution completed")
		return nil
	}

	// Configure watcher
	w := watcher.NewWatcher(
		watcher.WithPaths(c.settings.Directory),
		watcher.WithMask(c.settings.Pattern),
		watcher.WithWriteCallback(handleFileChange),
		watcher.WithBreakOnError(false),
	)

	log.Info().Msg("Watcher initialized. Waiting for JavaScript file changes...")

	// Run the watcher until context is canceled
	return w.Run(ctx)
}