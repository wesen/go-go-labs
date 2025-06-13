package cmd

import (
	"fmt"
	"os"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds/logging"
	"github.com/go-go-golems/glazed/pkg/help"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "kitty-control",
	Short: "A CLI tool for controlling kitty terminal emulator",
	Long: `kitty-control is a command-line interface for interacting with kitty terminal emulator
using its remote control protocol. It allows you to list windows, send text, focus windows,
launch new windows, get text content, and perform various other operations.

Logging:
  Use --log-level to control logging verbosity (trace, debug, info, warn, error, fatal, panic).
  Use --log-format to control log format (json, text).

Examples:
  kitty-control list                           # List all windows and tabs
  kitty-control send-text "Hello World"       # Send text to active window
  kitty-control get-text --extent=all         # Get all scrollback content
  kitty-control focus-window --match="title:*vim*"  # Focus window with vim in title
  kitty-control close-window --match="title:*old*"  # Close windows with "old" in title
  kitty-control launch "htop"                 # Launch htop in new window
  kitty-control --log-level=debug list        # List with debug logging enabled`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := logging.InitLoggerFromViper()
		cobra.CheckErr(err)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add logging support
	err := logging.AddLoggingLayerToRootCommand(rootCmd, "kitty-control")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding logging layer: %v\n", err)
		os.Exit(1)
	}

	// Bind viper to flags
	err = viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error binding flags: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	err = logging.InitLoggerFromViper()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	// Initialize help system
	helpSystem := help.NewHelpSystem()
	helpSystem.SetupCobraRootCommand(rootCmd)

	// Add commands
	addListCommand()
	addSendTextCommand()
	addFocusWindowCommand()
	addLaunchCommand()
	addCloseWindowCommand()
	addGetTextCommand()
	addResizeWindowCommand()
	addFocusTabCommand()
	addNewWindowCommand()
	addScrollWindowCommand()
	addSetWindowTitleCommand()
	addSetTabTitleCommand()
	addCreateMarkerCommand()
	addRemoveMarkerCommand()
	addEnvCommand()
	// More commands will be added by subagents later
}

func addListCommand() {
	listCmd, err := NewListCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating list command: %v\n", err)
		return
	}

	cobraCmd, err := cli.BuildCobraCommandFromCommand(listCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building list command: %v\n", err)
		return
	}

	rootCmd.AddCommand(cobraCmd)
}

func addSendTextCommand() {
	sendTextCmd, err := NewSendTextCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating send-text command: %v\n", err)
		return
	}

	cobraCmd, err := cli.BuildCobraCommandFromCommand(sendTextCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building send-text command: %v\n", err)
		return
	}

	rootCmd.AddCommand(cobraCmd)
}

func addFocusWindowCommand() {
	focusCmd, err := NewFocusWindowCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating focus-window command: %v\n", err)
		return
	}

	cobraCmd, err := cli.BuildCobraCommandFromCommand(focusCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building focus-window command: %v\n", err)
		return
	}

	rootCmd.AddCommand(cobraCmd)
}

func addLaunchCommand() {
	launchCmd, err := NewLaunchCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating launch command: %v\n", err)
		return
	}

	cobraCmd, err := cli.BuildCobraCommandFromCommand(launchCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building launch command: %v\n", err)
		return
	}

	rootCmd.AddCommand(cobraCmd)
}

func addCloseWindowCommand() {
	closeCmd, err := NewCloseWindowCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating close-window command: %v\n", err)
		return
	}

	cobraCmd, err := cli.BuildCobraCommandFromCommand(closeCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building close-window command: %v\n", err)
		return
	}

	rootCmd.AddCommand(cobraCmd)
}

func addGetTextCommand() {
	getTextCmd, err := NewGetTextCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating get-text command: %v\n", err)
		return
	}

	cobraCmd, err := cli.BuildCobraCommandFromWriterCommand(getTextCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building get-text command: %v\n", err)
		return
	}

	rootCmd.AddCommand(cobraCmd)
}

func addResizeWindowCommand() {
	resizeCmd, err := NewResizeWindowCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating resize-window command: %v\n", err)
		return
	}

	cobraCmd, err := cli.BuildCobraCommandFromCommand(resizeCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building resize-window command: %v\n", err)
		return
	}

	rootCmd.AddCommand(cobraCmd)
}

func addFocusTabCommand() {
	focusTabCmd, err := NewFocusTabCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating focus-tab command: %v\n", err)
		return
	}

	cobraCmd, err := cli.BuildCobraCommandFromCommand(focusTabCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building focus-tab command: %v\n", err)
		return
	}

	rootCmd.AddCommand(cobraCmd)
}

func addNewWindowCommand() {
	newWindowCmd, err := NewNewWindowCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating new-window command: %v\n", err)
		return
	}

	cobraCmd, err := cli.BuildCobraCommandFromCommand(newWindowCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building new-window command: %v\n", err)
		return
	}

	rootCmd.AddCommand(cobraCmd)
}

func addScrollWindowCommand() {
	scrollCmd, err := NewScrollWindowCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating scroll-window command: %v\n", err)
		return
	}

	cobraCmd, err := cli.BuildCobraCommandFromCommand(scrollCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building scroll-window command: %v\n", err)
		return
	}

	rootCmd.AddCommand(cobraCmd)
}

func addSetWindowTitleCommand() {
	setWindowTitleCmd, err := NewSetWindowTitleCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating set-window-title command: %v\n", err)
		return
	}

	cobraCmd, err := cli.BuildCobraCommandFromCommand(setWindowTitleCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building set-window-title command: %v\n", err)
		return
	}

	rootCmd.AddCommand(cobraCmd)
}

func addSetTabTitleCommand() {
	setTabTitleCmd, err := NewSetTabTitleCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating set-tab-title command: %v\n", err)
		return
	}

	cobraCmd, err := cli.BuildCobraCommandFromCommand(setTabTitleCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building set-tab-title command: %v\n", err)
		return
	}

	rootCmd.AddCommand(cobraCmd)
}

func addCreateMarkerCommand() {
	createMarkerCmd, err := NewCreateMarkerCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating create-marker command: %v\n", err)
		return
	}

	cobraCmd, err := cli.BuildCobraCommandFromCommand(createMarkerCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building create-marker command: %v\n", err)
		return
	}

	rootCmd.AddCommand(cobraCmd)
}

func addRemoveMarkerCommand() {
	removeMarkerCmd, err := NewRemoveMarkerCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating remove-marker command: %v\n", err)
		return
	}

	cobraCmd, err := cli.BuildCobraCommandFromCommand(removeMarkerCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building remove-marker command: %v\n", err)
		return
	}

	rootCmd.AddCommand(cobraCmd)
}

func addEnvCommand() {
	envCmd, err := NewEnvCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating env command: %v\n", err)
		return
	}

	cobraCmd, err := cli.BuildCobraCommandFromCommand(envCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building env command: %v\n", err)
		return
	}

	rootCmd.AddCommand(cobraCmd)
}
