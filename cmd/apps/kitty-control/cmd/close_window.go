package cmd

import (
	"context"
	"fmt"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"

	"github.com/go-go-golems/go-go-labs/cmd/apps/kitty-control/pkg"
)

type CloseWindowCommand struct {
	*cmds.CommandDescription
}

type CloseWindowSettings struct {
	SocketPath    string `glazed.parameter:"socket-path"`
	Match         string `glazed.parameter:"match"`
	Self          bool   `glazed.parameter:"self"`
	IgnoreNoMatch bool   `glazed.parameter:"ignore-no-match"`
}

var _ cmds.BareCommand = &CloseWindowCommand{}

func (c *CloseWindowCommand) Run(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
) error {
	s := &CloseWindowSettings{}
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
	payload := &pkg.CloseWindowPayload{}

	if s.Match != "" {
		payload.Match = &s.Match
	}
	if s.Self {
		payload.Self = &s.Self
	}
	if s.IgnoreNoMatch {
		payload.IgnoreNoMatch = &s.IgnoreNoMatch
	}

	// Send command
	err = client.CloseWindowWithPayload(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to close window: %w", err)
	}

	if s.Self {
		fmt.Printf("Closed current window\n")
	} else if s.Match != "" {
		fmt.Printf("Closed window(s) matching: %s\n", s.Match)
	} else {
		fmt.Printf("Closed window\n")
	}

	return nil
}

func NewCloseWindowCommand() (*CloseWindowCommand, error) {
	cmdDesc := cmds.NewCommandDescription(
		"close-window",
		cmds.WithShort("Close a kitty window"),
		cmds.WithLong(`Close one or more kitty windows.

This command allows you to close specific windows using match patterns, 
close the current window, or control error handling when no matches are found.

Match patterns examples:
  title:*vim*     - Windows with "vim" in the title
  id:1           - Window with ID 1
  pid:12345      - Window with process ID 12345
  cwd:/home/user - Windows in the /home/user directory

Examples:
  kitty-control close-window --match="title:*vim*"     # Close vim windows
  kitty-control close-window --match="id:1"           # Close window ID 1
  kitty-control close-window --self                   # Close current window
  kitty-control close-window --match="title:*old*" --ignore-no-match  # Ignore if no matches`),

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
				parameters.WithHelp("Pattern to match windows to close"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"self",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Close the window this command is run in"),
				parameters.WithDefault(false),
			),
			parameters.NewParameterDefinition(
				"ignore-no-match",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Do not report an error if no windows match"),
				parameters.WithDefault(false),
			),
		),
	)

	return &CloseWindowCommand{
		CommandDescription: cmdDesc,
	}, nil
}
