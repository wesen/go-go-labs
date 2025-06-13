package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"

	"github.com/go-go-golems/go-go-labs/cmd/apps/kitty-control/pkg"
)

type ScrollWindowCommand struct {
	*cmds.CommandDescription
}

type ScrollWindowSettings struct {
	SocketPath string   `glazed.parameter:"socket-path"`
	Amount     []string `glazed.parameter:"amount"`
	Match      string   `glazed.parameter:"match"`
}

var _ cmds.BareCommand = &ScrollWindowCommand{}

func (c *ScrollWindowCommand) Run(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
) error {
	s := &ScrollWindowSettings{}
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
	payload := &pkg.ScrollWindowPayload{}

	if s.Match != "" {
		payload.Match = &s.Match
	}

	// Parse amount - should be a 2-item list
	if len(s.Amount) > 0 {
		payload.Amount = make([]interface{}, 0, len(s.Amount))
		for _, amountStr := range s.Amount {
			// Try to parse as integer first
			if intVal, err := strconv.Atoi(amountStr); err == nil {
				payload.Amount = append(payload.Amount, intVal)
			} else {
				// Use as string if not a number
				payload.Amount = append(payload.Amount, amountStr)
			}
		}
	}

	// Send command
	err = client.ScrollWindow(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to scroll window: %w", err)
	}

	amountStr := "default"
	if len(s.Amount) > 0 {
		amountStr = strings.Join(s.Amount, ", ")
	}

	if s.Match != "" {
		fmt.Printf("Scrolled window(s) matching: %s (amount: %s)\n", s.Match, amountStr)
	} else {
		fmt.Printf("Scrolled window (amount: %s)\n", amountStr)
	}

	return nil
}

func NewScrollWindowCommand() (*ScrollWindowCommand, error) {
	cmdDesc := cmds.NewCommandDescription(
		"scroll-window",
		cmds.WithShort("Scroll a kitty window"),
		cmds.WithLong(`Scroll the contents of one or more kitty windows.

This command allows you to scroll window contents up or down by a specified amount.
The amount parameter is a complex 2-item list that specifies the scroll direction
and magnitude.

Amount format (2 items):
  - First item: direction ("up", "down", "page_up", "page_down", "home", "end")
  - Second item: number of lines or "all" for page/home/end operations

Match patterns examples:
  title:*vim*     - Windows with "vim" in the title
  id:1           - Window with ID 1
  pid:12345      - Window with process ID 12345
  cwd:/home/user - Windows in the /home/user directory

Examples:
  kitty-control scroll-window --amount="up,5"                    # Scroll up 5 lines
  kitty-control scroll-window --amount="down,10"                 # Scroll down 10 lines
  kitty-control scroll-window --amount="page_up,all"             # Page up
  kitty-control scroll-window --amount="home,all"                # Scroll to top
  kitty-control scroll-window --amount="end,all"                 # Scroll to bottom
  kitty-control scroll-window --match="title:*log*" --amount="down,20"  # Scroll log windows`),

		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"socket-path",
				parameters.ParameterTypeString,
				parameters.WithHelp("Path to kitty socket (defaults to KITTY_LISTEN_ON)"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"amount",
				parameters.ParameterTypeStringList,
				parameters.WithHelp("Scroll amount as 2-item list: direction,magnitude (e.g., 'up,5' or 'page_up,all')"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"match",
				parameters.ParameterTypeString,
				parameters.WithHelp("Pattern to match windows to scroll"),
				parameters.WithDefault(""),
			),
		),
	)

	return &ScrollWindowCommand{
		CommandDescription: cmdDesc,
	}, nil
}
