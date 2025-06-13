package cmd

import (
	"context"
	"fmt"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"

	"github.com/go-go-golems/go-go-labs/cmd/apps/kitty-control/pkg"
)

type ResizeWindowCommand struct {
	*cmds.CommandDescription
}

type ResizeWindowSettings struct {
	SocketPath string `glazed.parameter:"socket-path"`
	Match      string `glazed.parameter:"match"`
	Self       bool   `glazed.parameter:"self"`
	Increment  bool   `glazed.parameter:"increment"`
	Axis       string `glazed.parameter:"axis"`
}

var _ cmds.BareCommand = &ResizeWindowCommand{}

func (c *ResizeWindowCommand) Run(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
) error {
	s := &ResizeWindowSettings{}
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
	payload := &pkg.ResizeWindowPayload{}

	if s.Match != "" {
		payload.Match = &s.Match
	}
	if s.Self {
		payload.Self = &s.Self
	}
	if s.Increment {
		payload.Increment = &s.Increment
	}
	if s.Axis != "" {
		payload.Axis = &s.Axis
	}

	// Send command
	err = client.ResizeWindow(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to resize window: %w", err)
	}

	if s.Self {
		fmt.Printf("Resized current window (axis: %s)\n", s.Axis)
	} else if s.Match != "" {
		fmt.Printf("Resized window(s) matching: %s (axis: %s)\n", s.Match, s.Axis)
	} else {
		fmt.Printf("Resized window (axis: %s)\n", s.Axis)
	}

	return nil
}

func NewResizeWindowCommand() (*ResizeWindowCommand, error) {
	cmdDesc := cmds.NewCommandDescription(
		"resize-window",
		cmds.WithShort("Resize a kitty window"),
		cmds.WithLong(`Resize one or more kitty windows along specified axis.

This command allows you to resize windows horizontally, vertically, or reset to default size.
You can target specific windows using match patterns or resize the current window.

Axis options:
  horizontal - Resize window width
  vertical   - Resize window height
  reset      - Reset window to default size

Match patterns examples:
  title:*vim*     - Windows with "vim" in the title
  id:1           - Window with ID 1
  pid:12345      - Window with process ID 12345
  cwd:/home/user - Windows in the /home/user directory

Examples:
  kitty-control resize-window --axis=horizontal --increment          # Grow current window horizontally
  kitty-control resize-window --axis=vertical --match="title:*vim*"  # Resize vim window vertically
  kitty-control resize-window --axis=reset --self                   # Reset current window size`),

		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"socket-path",
				parameters.ParameterTypeString,
				parameters.WithHelp("Path to kitty socket (defaults to KITTY_LISTEN_ON)"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"match",
				parameters.ParameterTypeString,
				parameters.WithHelp("Pattern to match windows to resize"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"self",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Resize the window this command is run in"),
				parameters.WithDefault(false),
			),
			parameters.NewParameterDefinition(
				"increment",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Increment the window size rather than decrement"),
				parameters.WithDefault(false),
			),
			parameters.NewParameterDefinition(
				"axis",
				parameters.ParameterTypeChoice,
				parameters.WithHelp("Axis along which to resize the window"),
				parameters.WithChoices("horizontal", "vertical", "reset"),
				parameters.WithDefault("horizontal"),
			),
		),
	)

	return &ResizeWindowCommand{
		CommandDescription: cmdDesc,
	}, nil
}
