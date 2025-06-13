package cmd

import (
	"context"
	"fmt"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"

	"github.com/go-go-golems/go-go-labs/cmd/apps/kitty-control/pkg"
)

type CreateMarkerCommand struct {
	*cmds.CommandDescription
}

type CreateMarkerSettings struct {
	SocketPath string   `glazed.parameter:"socket-path"`
	Match      string   `glazed.parameter:"match"`
	Self       bool     `glazed.parameter:"self"`
	MarkerSpec []string `glazed.parameter:"marker-spec"`
}

var _ cmds.BareCommand = &CreateMarkerCommand{}

func (c *CreateMarkerCommand) Run(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
) error {
	s := &CreateMarkerSettings{}
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
	payload := &pkg.CreateMarkerPayload{}

	if s.Match != "" {
		payload.Match = &s.Match
	}
	if s.Self {
		payload.Self = &s.Self
	}
	if len(s.MarkerSpec) > 0 {
		payload.MarkerSpec = s.MarkerSpec
	}

	// Send command
	err = client.CreateMarker(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to create marker: %w", err)
	}

	if s.Self {
		fmt.Printf("Created marker in current window\n")
	} else if s.Match != "" {
		fmt.Printf("Created marker in window(s) matching: %s\n", s.Match)
	} else {
		fmt.Printf("Created marker\n")
	}

	return nil
}

func NewCreateMarkerCommand() (*CreateMarkerCommand, error) {
	cmdDesc := cmds.NewCommandDescription(
		"create-marker",
		cmds.WithShort("Create text markers in kitty windows"),
		cmds.WithLong(`Create text markers in one or more kitty windows.

Text markers are used to mark specific locations in terminal scrollback that can
be jumped to later. This is useful for marking important output, errors, or
specific points in long-running processes.

Match patterns examples:
  title:*vim*     - Windows with "vim" in the title
  id:1           - Window with ID 1
  pid:12345      - Window with process ID 12345
  cwd:/home/user - Windows in the /home/user directory

Marker specifications can include:
  - Text patterns to match
  - Colors and styles for the markers
  - Custom marker names

Examples:
  kitty-control create-marker --self --marker-spec="ERROR" --marker-spec="WARN"
  kitty-control create-marker --match="title:*log*" --marker-spec="[ERROR]"
  kitty-control create-marker --match="id:1" --marker-spec="Build complete"`),

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
				parameters.WithHelp("Pattern to match windows to create markers in"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"self",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Create marker in the window this command is run in"),
				parameters.WithDefault(false),
			),
			parameters.NewParameterDefinition(
				"marker-spec",
				parameters.ParameterTypeStringList,
				parameters.WithHelp("Marker specifications (can be used multiple times)"),
			),
		),
	)

	return &CreateMarkerCommand{
		CommandDescription: cmdDesc,
	}, nil
}
