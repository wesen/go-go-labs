package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func initLogger(logLevel string) error {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("invalid log level '%s': %w", logLevel, err)
	}

	output := zerolog.ConsoleWriter{Out: os.Stderr}
	log.Logger = zerolog.New(output).Level(level).With().Timestamp().Caller().Logger()
	log.Info().Str("level", level.String()).Msg("Logger initialized")
	return nil
}

func init() {
	// Default logger until properly configured
	log.Logger = zerolog.New(os.Stderr).Level(zerolog.InfoLevel).With().Timestamp().Logger()
}

func Execute() {

	rootCmd := &cobra.Command{
		Use:   "remote-js-executor",
		Short: "Watch JavaScript files and execute them in Chrome",
		Long:  "A tool that watches a directory for JavaScript file changes and executes them in Chrome using the Chrome DevTools Protocol.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logLevel, err := cmd.Flags().GetString("log-level")
			if err != nil {
				return fmt.Errorf("failed to get log-level flag: %w", err)
			}
			return initLogger(logLevel)
		},
	}

	// Add global flags
	rootCmd.PersistentFlags().String("log-level", "info", "Log level (debug, info, warn, error)")

	// Add commands
	execCmd, err := newWatchCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating watch command: %s\n", err)
		os.Exit(1)
	}
	rootCmd.AddCommand(execCmd)

	startCmd, err := newStartChromeCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating start-chrome command: %s\n", err)
		os.Exit(1)
	}
	rootCmd.AddCommand(startCmd)

	execOnceCmd, err := newExecuteOnceCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating execute-once command: %s\n", err)
		os.Exit(1)
	}
	rootCmd.AddCommand(execOnceCmd)

	baseCtx := context.Background()
	ctx, cancel := signal.NotifyContext(baseCtx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}