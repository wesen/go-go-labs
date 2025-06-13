package cmd

import (
	"context"
	"fmt"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"

	"github.com/go-go-golems/go-go-labs/cmd/apps/kitty-control/pkg"
)

type RemoveMarkerCommand struct {
	*cmds.CommandDescription
}

type RemoveMarkerSettings struct {
	SocketPath string `glazed.parameter:"socket-path"`
	Match      string `glazed.parameter:"match"`
	Self       bool   `glazed.parameter:"self"`
}

var _ cmds.BareCommand = &RemoveMarkerCommand{}

func (c *RemoveMarkerCommand) Run(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
) error {
	s := &RemoveMarkerSettings{}
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
	payload := &pkg.RemoveMarkerPayload{}

	if s.Match != "" {
		payload.Match = &s.Match
	}
	if s.Self {
		payload.Self = &s.Self
	}

	// Send command
	err = client.RemoveMarker(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to remove marker: %w", err)
	}

	if s.Self {
		fmt.Printf("Removed markers from current window\n")
	} else if s.Match != "" {
		fmt.Printf("Removed markers from window(s) matching: %s\n", s.Match)
	} else {
		fmt.Printf("Removed markers\n")
	}

	return nil
}

func NewRemoveMarkerCommand() (*RemoveMarkerCommand, error) {
	cmdDesc := cmds.NewCommandDescription(
		"remove-marker",
		cmds.WithShort("Remove text markers from kitty windows"),
		cmds.WithLong(`Remove text markers from one or more kitty windows.

This command removes previously created text markers from the terminal scrollback.
You can remove markers from specific windows using match patterns or from the
current window.

Match patterns examples:
  title:*vim*     - Windows with "vim" in the title
  id:1           - Window with ID 1
  pid:12345      - Window with process ID 12345
  cwd:/home/user - Windows in the /home/user directory

Examples:
  kitty-control remove-marker --self                    # Remove markers from current window
  kitty-control remove-marker --match="title:*log*"     # Remove markers from log windows
  kitty-control remove-marker --match="id:1"            # Remove markers from window ID 1`),

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
				parameters.WithHelp("Pattern to match windows to remove markers from"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"self",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Remove markers from the window this command is run in"),
				parameters.WithDefault(false),
			),
		),
	)

	return &RemoveMarkerCommand{
		CommandDescription: cmdDesc,
	}, nil
}
