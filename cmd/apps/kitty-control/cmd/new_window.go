package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"

	"github.com/go-go-golems/go-go-labs/cmd/apps/kitty-control/pkg"
)

type NewWindowCommand struct {
	*cmds.CommandDescription
}

type NewWindowSettings struct {
	SocketPath string   `glazed.parameter:"socket-path"`
	Args       []string `glazed.parameter:"args"`
	Match      string   `glazed.parameter:"match"`
	Title      string   `glazed.parameter:"title"`
	Cwd        string   `glazed.parameter:"cwd"`
	KeepFocus  bool     `glazed.parameter:"keep-focus"`
	WindowType string   `glazed.parameter:"window-type"`
	NewTab     bool     `glazed.parameter:"new-tab"`
	TabTitle   string   `glazed.parameter:"tab-title"`
}

var _ cmds.BareCommand = &NewWindowCommand{}

func (c *NewWindowCommand) Run(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
) error {
	s := &NewWindowSettings{}
	if err := parsedLayers.InitializeStruct(layers.DefaultSlug, s); err != nil {
		return fmt.Errorf("failed to initialize settings: %w", err)
	}

	// Create client
	var client *pkg.Client
	var err error

	if s.SocketPath != "" {
		client = pkg.NewClient(s.SocketPath)
	} else {
		client, err = pkg.NewClientFromEnv()
		if err != nil {
			return fmt.Errorf("failed to create client from environment: %w", err)
		}
	}

	// Create payload
	payload := &pkg.NewWindowPayload{}

	if len(s.Args) > 0 {
		payload.Args = s.Args
	}
	if s.Match != "" {
		payload.Match = &s.Match
	}
	if s.Title != "" {
		payload.Title = &s.Title
	}
	if s.Cwd != "" {
		payload.Cwd = &s.Cwd
	}
	if s.KeepFocus {
		payload.KeepFocus = &s.KeepFocus
	}
	if s.WindowType != "" {
		payload.WindowType = &s.WindowType
	}
	if s.NewTab {
		payload.NewTab = &s.NewTab
	}
	if s.TabTitle != "" {
		payload.TabTitle = &s.TabTitle
	}

	// Send command
	err = client.NewWindow(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to create new window: %w", err)
	}

	// Build output message
	var msg strings.Builder
	if s.NewTab {
		msg.WriteString("Created new tab")
		if s.TabTitle != "" {
			msg.WriteString(fmt.Sprintf(" with title '%s'", s.TabTitle))
		}
	} else {
		msg.WriteString("Created new window")
		if s.Title != "" {
			msg.WriteString(fmt.Sprintf(" with title '%s'", s.Title))
		}
	}

	if len(s.Args) > 0 {
		msg.WriteString(fmt.Sprintf(" running: %s", strings.Join(s.Args, " ")))
	}

	if s.Cwd != "" {
		msg.WriteString(fmt.Sprintf(" in directory: %s", s.Cwd))
	}

	fmt.Println(msg.String())
	return nil
}

func NewNewWindowCommand() (*NewWindowCommand, error) {
	cmdDesc := cmds.NewCommandDescription(
		"new-window",
		cmds.WithShort("Create a new kitty window"),
		cmds.WithLong(`Create a new kitty window with optional command and configuration.

This command creates a new window or tab in kitty with the specified program,
title, working directory and other options. It's a simpler version of the
launch command focused on basic window creation.

Window types:
  window  - Create a new window (default)
  tab     - Create in a new tab
  overlay - Create as overlay

Examples:
  kitty-control new-window                                    # Create empty window
  kitty-control new-window --args="htop"                     # Create window running htop
  kitty-control new-window --args="vim,file.txt" --title="Editor"  # Create vim window with title
  kitty-control new-window --new-tab --tab-title="Development"     # Create new tab
  kitty-control new-window --cwd="/tmp" --keep-focus              # Create in /tmp, keep current focus
  kitty-control new-window --window-type="overlay" --args="top"   # Create overlay window`),

		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"socket-path",
				parameters.ParameterTypeString,
				parameters.WithHelp("Path to kitty socket (defaults to KITTY_LISTEN_ON)"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"args",
				parameters.ParameterTypeStringList,
				parameters.WithHelp("Command and arguments to run in the new window"),
			),
			parameters.NewParameterDefinition(
				"match",
				parameters.ParameterTypeString,
				parameters.WithHelp("Pattern to match where to create the new window"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"title",
				parameters.ParameterTypeString,
				parameters.WithHelp("Title for the new window"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"cwd",
				parameters.ParameterTypeString,
				parameters.WithHelp("Working directory for the new window"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"keep-focus",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Do not focus the new window"),
				parameters.WithDefault(false),
			),
			parameters.NewParameterDefinition(
				"window-type",
				parameters.ParameterTypeChoice,
				parameters.WithHelp("Type of window to create"),
				parameters.WithChoices("window", "tab", "overlay"),
				parameters.WithDefault("window"),
			),
			parameters.NewParameterDefinition(
				"new-tab",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Create the window in a new tab"),
				parameters.WithDefault(false),
			),
			parameters.NewParameterDefinition(
				"tab-title",
				parameters.ParameterTypeString,
				parameters.WithHelp("Title for the new tab (when using --new-tab)"),
				parameters.WithDefault(""),
			),
		),
	)

	return &NewWindowCommand{
		CommandDescription: cmdDesc,
	}, nil
}
