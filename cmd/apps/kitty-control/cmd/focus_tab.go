package cmd

import (
	"context"
	"fmt"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"

	"github.com/go-go-golems/go-go-labs/cmd/apps/kitty-control/pkg"
)

type FocusTabCommand struct {
	*cmds.CommandDescription
}

type FocusTabSettings struct {
	SocketPath string `glazed.parameter:"socket-path"`
	Match      string `glazed.parameter:"match"`
}

var _ cmds.BareCommand = &FocusTabCommand{}

func (c *FocusTabCommand) Run(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
) error {
	s := &FocusTabSettings{}
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
	payload := &pkg.FocusTabPayload{}
	if s.Match != "" {
		payload.Match = &s.Match
	}

	// Send command
	err = client.FocusTab(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to focus tab: %w", err)
	}

	if s.Match != "" {
		fmt.Printf("Focused tab matching: %s\n", s.Match)
	} else {
		fmt.Printf("Focused tab\n")
	}

	return nil
}

func NewFocusTabCommand() (*FocusTabCommand, error) {
	cmdDesc := cmds.NewCommandDescription(
		"focus-tab",
		cmds.WithShort("Focus a kitty tab"),
		cmds.WithLong(`Focus a specific kitty tab using a match pattern.

This command allows you to bring a specific tab into focus. You can specify
which tab to focus using various match patterns based on tab title, window
properties, and more.

Match patterns examples:
  title:*work*    - Tabs with "work" in the title
  index:0         - First tab (0-indexed)
  index:-1        - Last tab
  recent:1        - Previously active tab

Examples:
  kitty-control focus-tab --match="title:*work*"    # Focus tab with "work" in title
  kitty-control focus-tab --match="index:0"        # Focus first tab
  kitty-control focus-tab --match="recent:1"       # Focus previously active tab
  kitty-control focus-tab --match="index:-1"       # Focus last tab`),

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
				parameters.WithHelp("Pattern to match the tab to focus"),
				parameters.WithRequired(true),
			),
		),
	)

	return &FocusTabCommand{
		CommandDescription: cmdDesc,
	}, nil
}
