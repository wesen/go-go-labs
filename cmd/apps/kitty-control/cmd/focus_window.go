package cmd

import (
	"context"
	"fmt"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"

	"github.com/go-go-golems/go-go-labs/cmd/apps/kitty-control/pkg"
)

type FocusWindowCommand struct {
	*cmds.CommandDescription
}

type FocusWindowSettings struct {
	SocketPath string `glazed.parameter:"socket-path"`
	Match      string `glazed.parameter:"match"`
}

var _ cmds.BareCommand = &FocusWindowCommand{}

func (c *FocusWindowCommand) Run(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
) error {
	s := &FocusWindowSettings{}
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

	// Send command
	err = client.FocusWindow(ctx, s.Match)
	if err != nil {
		return fmt.Errorf("failed to focus window: %w", err)
	}

	if s.Match != "" {
		fmt.Printf("Focused window matching: %s\n", s.Match)
	} else {
		fmt.Printf("Focused window\n")
	}

	return nil
}

func NewFocusWindowCommand() (*FocusWindowCommand, error) {
	cmdDesc := cmds.NewCommandDescription(
		"focus-window",
		cmds.WithShort("Focus a kitty window"),
		cmds.WithLong(`Focus a specific kitty window using a match pattern.

This command allows you to bring a specific window into focus. You can specify
which window to focus using various match patterns based on window title,
process name, working directory, and more.

Match patterns examples:
  title:*vim*     - Windows with "vim" in the title
  id:1           - Window with ID 1
  pid:12345      - Window with process ID 12345
  cwd:/home/user - Windows in the /home/user directory

Examples:
  kitty-control focus-window --match="title:*vim*"     # Focus vim window
  kitty-control focus-window --match="id:1"           # Focus window ID 1
  kitty-control focus-window --match="cwd:/tmp"       # Focus window in /tmp`),

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
				parameters.WithHelp("Pattern to match the window to focus"),
				parameters.WithRequired(true),
			),
		),
	)

	return &FocusWindowCommand{
		CommandDescription: cmdDesc,
	}, nil
}
