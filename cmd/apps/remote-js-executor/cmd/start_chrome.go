package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/go-go-golems/glazed/pkg/types"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type StartChromeSettings struct {
	Port       int    `glazed.parameter:"port"`
	BrowserPath string `glazed.parameter:"browser-path"`
	UserDataDir string `glazed.parameter:"user-data-dir"`
	Headless   bool   `glazed.parameter:"headless"`
	WaitMs     int    `glazed.parameter:"wait-ms"`
}

type StartChromeCommand struct {
	*cmds.CommandDescription
	settings StartChromeSettings
}

func newStartChromeCommand() (*cobra.Command, error) {
	glazeLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, fmt.Errorf("could not create Glazed parameter layer: %w", err)
	}

	cmd := &StartChromeCommand{
		CommandDescription: cmds.NewCommandDescription(
			"start-chrome",
			cmds.WithShort("Start Chrome with remote debugging enabled"),
			cmds.WithLong("Start Chrome or Chromium browser with remote debugging enabled on the specified port."),
			cmds.WithFlags(
				parameters.NewParameterDefinition(
					"port",
					parameters.ParameterTypeInteger,
					parameters.WithHelp("Port for Chrome remote debugging"),
					parameters.WithDefault(9222),
				),
				parameters.NewParameterDefinition(
					"browser-path",
					parameters.ParameterTypeString,
					parameters.WithHelp("Path to Chrome/Chromium executable (leave empty for auto-detection)"),
					parameters.WithDefault(""),
				),
				parameters.NewParameterDefinition(
					"user-data-dir",
					parameters.ParameterTypeString,
					parameters.WithHelp("User data directory for Chrome (leave empty for temporary profile)"),
					parameters.WithDefault(""),
				),
				parameters.NewParameterDefinition(
					"headless",
					parameters.ParameterTypeBool,
					parameters.WithHelp("Run Chrome in headless mode"),
					parameters.WithDefault(false),
				),
				parameters.NewParameterDefinition(
					"wait-ms",
					parameters.ParameterTypeInteger,
					parameters.WithHelp("Time to wait after starting Chrome (milliseconds)"),
					parameters.WithDefault(1000),
				),
			),
			cmds.WithLayersList(glazeLayer),
		),
	}

	return cli.BuildCobraCommandFromGlazeCommand(cmd)
}

func (c *StartChromeCommand) RunIntoGlazeProcessor(ctx context.Context, parsedLayers *layers.ParsedLayers, gp middlewares.Processor) error {
	if err := parsedLayers.InitializeStruct(layers.DefaultSlug, &c.settings); err != nil {
		return err
	}

	// Find browser executable if not specified
	browserPath := c.settings.BrowserPath
	if browserPath == "" {
		detectedPath, err := detectChromePath()
		if err != nil {
			return fmt.Errorf("failed to auto-detect Chrome/Chromium path: %w", err)
		}
		browserPath = detectedPath
	}

	// Verify browser path
	if _, err := os.Stat(browserPath); err != nil {
		return fmt.Errorf("browser executable not found at '%s': %w", browserPath, err)
	}

	// Prepare user data directory
	userDataDir := c.settings.UserDataDir
	if userDataDir == "" {
		// Create temporary directory
		tempDir, err := os.MkdirTemp("", "chrome-remote-js-*")
		if err != nil {
			return fmt.Errorf("failed to create temporary user data directory: %w", err)
		}
		userDataDir = tempDir
		log.Info().Str("tempDir", tempDir).Msg("Created temporary Chrome user data directory")
	} else {
		// Ensure the directory exists
		if err := os.MkdirAll(userDataDir, 0755); err != nil {
			return fmt.Errorf("failed to create user data directory at '%s': %w", userDataDir, err)
		}
	}

	// Log the settings
	log.Info().Str("browserPath", browserPath).
		Int("port", c.settings.Port).
		Str("userDataDir", userDataDir).
		Bool("headless", c.settings.Headless).
		Msg("Starting Chrome")

	// Build command arguments
	args := []string{
		"--remote-debugging-port=" + fmt.Sprintf("%d", c.settings.Port),
		"--user-data-dir=" + userDataDir,
		"--no-first-run",
		"--no-default-browser-check",
	}

	// Add headless mode if requested
	if c.settings.Headless {
		// Chrome 112 and newer use --headless=new instead of just --headless
		args = append(args, "--headless=new")
	}

	// Create the command
	cmd := exec.CommandContext(ctx, browserPath, args...)

	// Set the command output to our logger
	cmd.Stdout = newLogWriter("chrome-stdout")
	cmd.Stderr = newLogWriter("chrome-stderr")

	// Start Chrome
	log.Info().Strs("args", args).Msg("Starting Chrome with arguments")
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start Chrome: %w", err)
	}

	// Give Chrome some time to start
	log.Info().Int("waitMs", c.settings.WaitMs).Msg("Waiting for Chrome to initialize")
	time.Sleep(time.Duration(c.settings.WaitMs) * time.Millisecond)

	// Output success message
	log.Info().Int("pid", cmd.Process.Pid).Str("debugURL", fmt.Sprintf("http://localhost:%d", c.settings.Port)).Msg("Chrome started successfully")

	// Wait for Chrome to exit or context to be canceled
	go func() {
		<-ctx.Done()
		log.Info().Int("pid", cmd.Process.Pid).Msg("Stopping Chrome due to context cancellation")
		cmd.Process.Kill()
	}()

	// Add output data to the processor
	row := types.NewRow(
		types.MRP("status", "success"),
		types.MRP("browser", filepath.Base(browserPath)),
		types.MRP("browser_path", browserPath),
		types.MRP("pid", cmd.Process.Pid),
		types.MRP("debug_url", fmt.Sprintf("http://localhost:%d", c.settings.Port)),
		types.MRP("user_data_dir", userDataDir),
		types.MRP("headless", c.settings.Headless),
	)

	// Output information for the user
	if err := gp.AddRow(ctx, row); err != nil {
		return fmt.Errorf("failed to add row to processor: %w", err)
	}

	// Block until Chrome exits if we're not waiting interactively
	if err := cmd.Wait(); err != nil {
		log.Warn().Err(err).Msg("Chrome exited with an error")
		return fmt.Errorf("Chrome exited with an error: %w", err)
	}

	log.Info().Msg("Chrome exited successfully")
	return nil
}

// logWriter is a simple writer that logs output to zerolog
type logWriter struct {
	prefix string
}

func newLogWriter(prefix string) *logWriter {
	return &logWriter{prefix: prefix}
}

func (l *logWriter) Write(p []byte) (n int, err error) {
	// Trim trailing whitespace and only log non-empty lines
	text := strings.TrimSpace(string(p))
	if text != "" {
		log.Debug().Str("source", l.prefix).Msg(text)
	}
	return len(p), nil
}

// detectChromePath attempts to find Chrome or Chromium in standard locations
func detectChromePath() (string, error) {
	log.Debug().Msg("Attempting to detect Chrome/Chromium location")

	paths := []string{}

	switch runtime.GOOS {
	case "linux":
		// Common Linux locations
		paths = []string{
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
			"/usr/bin/google-chrome",
			"/usr/bin/google-chrome-stable",
			"/snap/bin/chromium",
		}
	case "darwin":
		// macOS locations
		paths = []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
		}
	case "windows":
		// Windows locations
		progFiles := os.Getenv("PROGRAMFILES")
		progFilesX86 := os.Getenv("PROGRAMFILES(X86)")
		paths = []string{
			progFiles + "\\Google\\Chrome\\Application\\chrome.exe",
			progFilesX86 + "\\Google\\Chrome\\Application\\chrome.exe",
			progFiles + "\\Chromium\\Application\\chrome.exe",
			progFilesX86 + "\\Chromium\\Application\\chrome.exe",
		}
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Try each path
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			log.Debug().Str("path", path).Msg("Found Chrome/Chromium at path")
			return path, nil
		}
	}

	// Try to find chrome in PATH
	for _, browser := range []string{"chromium", "chromium-browser", "google-chrome", "chrome"} {
		if path, err := exec.LookPath(browser); err == nil {
			log.Debug().Str("path", path).Str("browser", browser).Msg("Found Chrome/Chromium in PATH")
			return path, nil
		}
	}

	return "", fmt.Errorf("could not find Chrome or Chromium browser in standard locations")
}