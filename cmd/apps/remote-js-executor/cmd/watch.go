package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-go-golems/clay/pkg/watcher"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/logging"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type WatchSettings struct {
	Path           string `glazed.parameter:"path"`
	Pattern        string `glazed.parameter:"pattern"`
	ChromeURL      string `glazed.parameter:"chrome-url"`
	NavigateToURL  string `glazed.parameter:"navigate-to"`
	ThrottleMs     int    `glazed.parameter:"throttle-ms"`
	StartChrome    bool   `glazed.parameter:"start-chrome"`
	ChromePort     int    `glazed.parameter:"chrome-port"`
	Headless       bool   `glazed.parameter:"headless"`
	LoadAllOnStart bool   `glazed.parameter:"load-all-on-start"`
}

type WatchCommand struct {
	*cmds.CommandDescription
	settings WatchSettings
}

func newWatchCommand() (*cobra.Command, error) {
	cmd := &WatchCommand{
		CommandDescription: cmds.NewCommandDescription(
			"watch",
			cmds.WithShort("Watch JavaScript files and execute them in Chrome"),
			cmds.WithLong("Watch a directory or file for JavaScript changes and execute them in Chrome using the Chrome DevTools Protocol."),
			cmds.WithFlags(
				parameters.NewParameterDefinition(
					"path",
					parameters.ParameterTypeString,
					parameters.WithHelp("Path to a directory or specific JS file"),
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
				parameters.NewParameterDefinition(
					"start-chrome",
					parameters.ParameterTypeBool,
					parameters.WithHelp("Start Chrome automatically before watching"),
					parameters.WithDefault(false),
				),
				parameters.NewParameterDefinition(
					"chrome-port",
					parameters.ParameterTypeInteger,
					parameters.WithHelp("Port to use when starting Chrome"),
					parameters.WithDefault(9222),
				),
				parameters.NewParameterDefinition(
					"headless",
					parameters.ParameterTypeBool,
					parameters.WithHelp("Run Chrome in headless mode when starting it"),
					parameters.WithDefault(false),
				),
				parameters.NewParameterDefinition(
					"load-all-on-start",
					parameters.ParameterTypeBool,
					parameters.WithHelp("Load and execute all JavaScript files in the directory on startup"),
					parameters.WithDefault(false),
				),
			),
		),
	}
	_, err := logging.AddLoggingLayerToCommand(cmd)
	if err != nil {
		return nil, err
	}

	return cli.BuildCobraCommandFromBareCommand(cmd)
}

// extractVisitURL extracts the URL from a // visitUrl: comment at the beginning of JavaScript code
func extractVisitURL(content string) string {
	// Look for a visitUrl comment at the beginning of the file
	// The format is: // visitUrl: https://example.com
	// We'll use a simple string check for now
	const visitUrlPrefix = "// visitUrl:"

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, visitUrlPrefix) {
			return strings.TrimSpace(line[len(visitUrlPrefix):])
		}
		// Only check the first few lines
		if i > 5 {
			break
		}
	}

	return ""
}

func (c *WatchCommand) Run(ctx context.Context, parsedLayers *layers.ParsedLayers) error {
	if err := parsedLayers.InitializeStruct(layers.DefaultSlug, &c.settings); err != nil {
		return err
	}

	// Validate path exists
	pathInfo, err := os.Stat(c.settings.Path)
	if err != nil {
		return fmt.Errorf("error accessing path %s: %w", c.settings.Path, err)
	}

	// Track if we're working with a single file or directory
	isSingleFile := !pathInfo.IsDir()

	// Log settings
	log.Info().Str("path", c.settings.Path).
		Str("pattern", c.settings.Pattern).
		Str("chromeURL", c.settings.ChromeURL).
		Str("navigateToURL", c.settings.NavigateToURL).
		Int("throttleMs", c.settings.ThrottleMs).
		Bool("startChrome", c.settings.StartChrome).
		Int("chromePort", c.settings.ChromePort).
		Bool("headless", c.settings.Headless).
		Bool("loadAllOnStart", c.settings.LoadAllOnStart).
		Bool("isSingleFile", isSingleFile).
		Msg("Starting file watcher")

	// Start Chrome if requested
	var chromeCmd *exec.Cmd
	if c.settings.StartChrome {
		log.Info().Msg("Starting Chrome before watching")

		// Find Chrome executable
		browserPath, err := detectChromePath()
		if err != nil {
			return fmt.Errorf("failed to auto-detect Chrome/Chromium path: %w", err)
		}

		// Create temporary user data directory
		tempDir, err := os.MkdirTemp("", "chrome-remote-js-*")
		if err != nil {
			return fmt.Errorf("failed to create temporary user data directory: %w", err)
		}
		log.Info().Str("tempDir", tempDir).Msg("Created temporary Chrome user data directory")

		// Build arguments
		args := []string{
			"--remote-debugging-port=" + fmt.Sprintf("%d", c.settings.ChromePort),
			"--user-data-dir=" + tempDir,
			"--no-first-run",
			"--no-default-browser-check",
		}

		if c.settings.Headless {
			args = append(args, "--headless=new")
		}

		// Create and start the command
		chromeCmd = exec.CommandContext(ctx, browserPath, args...)
		chromeCmd.Stdout = newLogWriter("chrome-stdout")
		chromeCmd.Stderr = newLogWriter("chrome-stderr")

		log.Info().Strs("args", args).Msg("Starting Chrome with arguments")
		if err := chromeCmd.Start(); err != nil {
			return fmt.Errorf("failed to start Chrome: %w", err)
		}

		// Give Chrome time to initialize
		log.Info().Msg("Waiting for Chrome to initialize")
		time.Sleep(1 * time.Second)

		// Update Chrome URL if using the default
		if c.settings.ChromeURL == "http://localhost:9222" && c.settings.ChromePort != 9222 {
			c.settings.ChromeURL = fmt.Sprintf("http://localhost:%d", c.settings.ChromePort)
			log.Info().Str("chromeURL", c.settings.ChromeURL).Msg("Updated Chrome URL based on port")
		}

		fmt.Printf("Chrome started successfully on port %d\n", c.settings.ChromePort)
	}

	// Create Chrome executor
	log.Debug().
		Str("chromeURL", c.settings.ChromeURL).
		Str("navigateToURL", c.settings.NavigateToURL).
		Msg("WatchCommand: Attempting to create ChromeExecutor")
	executor, err := NewChromeExecutor(c.settings.ChromeURL, c.settings.NavigateToURL)
	if err != nil {
		log.Error().Err(err).Msg("WatchCommand: Failed to create Chrome executor during setup")
		return fmt.Errorf("failed to create Chrome executor: %w", err)
	}
	if executor == nil {
		log.Error().Msg("WatchCommand: Chrome executor is nil after creation, even without explicit error")
		return fmt.Errorf("failed to create Chrome executor: instance is nil")
	}
	log.Info().Msg("WatchCommand: Chrome executor created successfully")
	defer func() {
		log.Debug().Msg("WatchCommand: Closing Chrome executor")
		executor.Close()
	}()

	// Create a throttled execution function
	lastExecution := time.Now().Add(-24 * time.Hour) // Initialize to a day ago
	mutex := &sync.Mutex{}

	executeJSFile := func(path string, applyThrottle bool) error {
		if applyThrottle {
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
		}

		log.Info().Str("path", path).Msg("Processing JavaScript file")

		// Log the context being passed to ExecuteJavaScript
		if ctx.Err() != nil {
			log.Warn().Err(ctx.Err()).Str("path", path).Msg("WatchCommand: Context for ExecuteJavaScript is already canceled before call")
		} else {
			log.Debug().Str("path", path).Msg("WatchCommand: Context for ExecuteJavaScript appears valid before call")
		}

		// Read the file
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Check for visitUrl comment
		visitUrl := extractVisitURL(string(content))
		if visitUrl != "" {
			log.Info().Str("path", path).Str("visitUrl", visitUrl).Msg("Found visitUrl comment, navigating before execution")
			err := executor.NavigateToPage(ctx, visitUrl)
			if err != nil {
				log.Error().Err(err).Str("path", path).Str("visitUrl", visitUrl).Msg("Failed to navigate to URL specified in visitUrl comment")
				return fmt.Errorf("failed to navigate to URL %s specified in file %s: %w", visitUrl, path, err)
			}
		}

		// Execute the JavaScript
		log.Debug().Str("path", path).Msg("WatchCommand: Calling executor.ExecuteJavaScript")
		result, err := executor.ExecuteJavaScript(ctx, string(content))
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("WatchCommand: Error returned from executor.ExecuteJavaScript")
			return fmt.Errorf("failed to execute JavaScript from file %s: %w", path, err)
		}

		filename := filepath.Base(path)
		log.Info().Str("file", filename).Str("result", result).Msg("JavaScript execution completed")
		return nil
	}

	handleFileChange := func(path string) error {
		return executeJSFile(path, true)
	}

	// Configure watcher
	// Use appropriate mask pattern depending on if it's a file or directory
	mask := c.settings.Pattern
	if isSingleFile {
		mask = "*" // For single files, use a catch-all pattern
	}

	w := watcher.NewWatcher(
		watcher.WithPaths(c.settings.Path),
		watcher.WithMask(mask),
		watcher.WithWriteCallback(handleFileChange),
		watcher.WithBreakOnError(false),
	)

	log.Info().Msg("Watcher initialized. Waiting for JavaScript file changes...")

	// Load all JavaScript files on startup if requested
	if c.settings.LoadAllOnStart {
		log.Info().Msg("Loading all JavaScript files in directory on startup")

		if isSingleFile {
			// If path is a single file, just execute it
			if filepath.Ext(c.settings.Path) == ".js" {
				err = executeJSFile(c.settings.Path, false) // Don't apply throttling for initial load
				if err != nil {
					log.Error().Err(err).Str("path", c.settings.Path).Msg("Failed to execute JavaScript file during initial load")
				}
			}
		} else {
			// Use filepath.Walk to find all JS files in the directory
			err := filepath.Walk(c.settings.Path, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// Skip directories
				if info.IsDir() {
					return nil
				}

				// Check if file matches pattern (*.js)
				if filepath.Ext(path) == ".js" {
					err = executeJSFile(path, false) // Don't apply throttling for initial load
					if err != nil {
						log.Error().Err(err).Str("path", path).Msg("Failed to execute JavaScript file during initial load")
						// Continue with other files even if one fails
					}
				}

				return nil
			})

			if err != nil {
				log.Error().Err(err).Msg("Error during initial load of JavaScript files")
			}
		}
	}

	// Run the watcher until context is canceled
	return w.Run(ctx)
}
