package cmd

import (
	"context"
	"fmt"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"

	"github.com/go-go-golems/go-go-labs/cmd/apps/kitty-control/pkg"
)

type SetTabTitleCommand struct {
	*cmds.CommandDescription
}

type SetTabTitleSettings struct {
	SocketPath string `glazed.parameter:"socket-path"`
	Title      string `glazed.parameter:"title"`
	Match      string `glazed.parameter:"match"`
}

var _ cmds.BareCommand = &SetTabTitleCommand{}

func (c *SetTabTitleCommand) Run(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
) error {
	s := &SetTabTitleSettings{}
	if err := parsedLayers.InitializeStruct(layers.DefaultSlug, s); err != nil {
		return fmt.Errorf("failed to initialize settings: %w", err)
	}

	if s.Title == "" {
		return fmt.Errorf("title is required")
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
	payload := &pkg.SetTabTitlePayload{
		Title: s.Title,
	}

	if s.Match != "" {
		payload.Match = &s.Match
	}

	// Send command
	err = client.SetTabTitle(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to set tab title: %w", err)
	}

	if s.Match != "" {
		fmt.Printf("Set title for tab(s) matching: %s\n", s.Match)
	} else {
		fmt.Printf("Set tab title to: %s\n", s.Title)
	}

	return nil
}

func NewSetTabTitleCommand() (*SetTabTitleCommand, error) {
	cmdDesc := cmds.NewCommandDescription(
		"set-tab-title",
		cmds.WithShort("Set the title of kitty tabs"),
		cmds.WithLong(`Set the title of one or more kitty tabs.

This command allows you to set the title of specific tabs using match patterns
or set the title of the current tab. The title is required for this command.

Match patterns examples:
  title:*work*    - Tabs with "work" in the title
  id:1           - Tab with ID 1
  index:0        - First tab (zero-indexed)

Examples:
  kitty-control set-tab-title --title="Development" --match="title:*dev*"
  kitty-control set-tab-title --title="Main Terminal"
  kitty-control set-tab-title --title="Database" --match="id:2"`),

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
				parameters.WithHelp("Title to set for the tab(s) (required)"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"match",
				parameters.ParameterTypeString,
				parameters.WithHelp("Pattern to match tabs to set title for"),
				parameters.WithDefault(""),
			),
		),
	)

	return &SetTabTitleCommand{
		CommandDescription: cmdDesc,
	}, nil
}
