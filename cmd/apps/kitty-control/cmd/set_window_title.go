package cmd

import (
	"context"
	"fmt"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"

	"github.com/go-go-golems/go-go-labs/cmd/apps/kitty-control/pkg"
)

type SetWindowTitleCommand struct {
	*cmds.CommandDescription
}

type SetWindowTitleSettings struct {
	SocketPath string `glazed.parameter:"socket-path"`
	Title      string `glazed.parameter:"title"`
	Match      string `glazed.parameter:"match"`
	Temporary  bool   `glazed.parameter:"temporary"`
}

var _ cmds.BareCommand = &SetWindowTitleCommand{}

func (c *SetWindowTitleCommand) Run(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
) error {
	s := &SetWindowTitleSettings{}
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
	payload := &pkg.SetWindowTitlePayload{}

	if s.Title != "" {
		payload.Title = &s.Title
	}
	if s.Match != "" {
		payload.Match = &s.Match
	}
	if s.Temporary {
		payload.Temporary = &s.Temporary
	}

	// Send command
	err = client.SetWindowTitle(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to set window title: %w", err)
	}

	if s.Match != "" {
		fmt.Printf("Set title for window(s) matching: %s\n", s.Match)
	} else {
		fmt.Printf("Set window title\n")
	}

	return nil
}

func NewSetWindowTitleCommand() (*SetWindowTitleCommand, error) {
	cmdDesc := cmds.NewCommandDescription(
		"set-window-title",
		cmds.WithShort("Set the title of kitty windows"),
		cmds.WithLong(`Set the title of one or more kitty windows.

This command allows you to set the title of specific windows using match patterns
or set the title of the current window. The title can be set temporarily (until 
the next title change by the application) or permanently.

Match patterns examples:
  title:*vim*     - Windows with "vim" in the title
  id:1           - Window with ID 1
  pid:12345      - Window with process ID 12345
  cwd:/home/user - Windows in the /home/user directory

Examples:
  kitty-control set-window-title --title="My Editor" --match="title:*vim*"
  kitty-control set-window-title --title="Terminal" --temporary
  kitty-control set-window-title --title="Debug Session" --match="id:1"`),

		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"socket-path",
				parameters.ParameterTypeString,
				parameters.WithHelp("Path to kitty socket (defaults to KITTY_LISTEN_ON)"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"title",
				parameters.ParameterTypeString,
				parameters.WithHelp("Title to set for the window(s)"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"match",
				parameters.ParameterTypeString,
				parameters.WithHelp("Pattern to match windows to set title for"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"temporary",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Make the title change temporary (until next app title change)"),
				parameters.WithDefault(false),
			),
		),
	)

	return &SetWindowTitleCommand{
		CommandDescription: cmdDesc,
	}, nil
}
