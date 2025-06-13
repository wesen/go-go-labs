package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/go-go-golems/glazed/pkg/types"
	"github.com/rs/zerolog/log"

	"github.com/go-go-golems/go-go-labs/cmd/apps/kitty-control/pkg"
)

type ListCommand struct {
	*cmds.CommandDescription
}

type ListSettings struct {
	SocketPath string `glazed.parameter:"socket-path"`
	AllEnvVars bool   `glazed.parameter:"all-env-vars"`
	Match      string `glazed.parameter:"match"`
	MatchTab   string `glazed.parameter:"match-tab"`
	Self       bool   `glazed.parameter:"self"`
}

var _ cmds.GlazeCommand = &ListCommand{}

func (c *ListCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	log.Debug().Msg("executing list command")
	s := &ListSettings{}
	if err := parsedLayers.InitializeStruct(layers.DefaultSlug, s); err != nil {
		log.Error().Err(err).Msg("failed to initialize settings")
		return fmt.Errorf("failed to initialize settings: %w", err)
	}
	log.Debug().Interface("settings", s).Msg("initialized list settings")

	// Create client
	var client *pkg.Client
	var err error

	if s.SocketPath != "" {
		log.Debug().Str("socket_path", s.SocketPath).Msg("using provided socket path")
		client = pkg.NewClient(s.SocketPath)
	} else {
		log.Debug().Msg("using socket path from environment")
		client, err = pkg.NewClientFromEnv()
		if err != nil {
			log.Error().Err(err).Msg("failed to create client from environment")
			return fmt.Errorf("failed to create client from environment: %w", err)
		}
	}

	// Create payload
	payload := &pkg.ListPayload{}
	if s.AllEnvVars {
		payload.AllEnvVars = &s.AllEnvVars
	}
	if s.Match != "" {
		payload.Match = &s.Match
	}
	if s.MatchTab != "" {
		payload.MatchTab = &s.MatchTab
	}
	if s.Self {
		payload.Self = &s.Self
	}

	// Send command
	response, err := client.List(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to list windows: %w", err)
	}

	// Parse response
	var kittyResponse struct {
		OK   bool            `json:"ok"`
		Data json.RawMessage `json:"data"`
	}

	// Extract JSON from kitty protocol response
	responseStr := string(response)
	if len(responseStr) > 12 && responseStr[:12] == "\x1bP@kitty-cmd" {
		// Find the JSON part
		jsonStart := 12
		jsonEnd := len(responseStr) - 2 // Remove \x1b\\
		if jsonEnd > jsonStart {
			responseStr = responseStr[jsonStart:jsonEnd]
		}
	}

	if err := json.Unmarshal([]byte(responseStr), &kittyResponse); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !kittyResponse.OK {
		return fmt.Errorf("kitty command failed")
	}

	// Parse the data as a list of windows/tabs
	var windowsData interface{}
	if err := json.Unmarshal(kittyResponse.Data, &windowsData); err != nil {
		return fmt.Errorf("failed to parse windows data: %w", err)
	}

	// Convert to rows for structured output
	if windowsList, ok := windowsData.([]interface{}); ok {
		for _, window := range windowsList {
			if windowMap, ok := window.(map[string]interface{}); ok {
				row := types.NewRowFromMap(windowMap)
				if err := gp.AddRow(ctx, row); err != nil {
					return err
				}
			}
		}
	} else {
		// Single window case
		if windowMap, ok := windowsData.(map[string]interface{}); ok {
			row := types.NewRowFromMap(windowMap)
			if err := gp.AddRow(ctx, row); err != nil {
				return err
			}
		}
	}

	return nil
}

func NewListCommand() (*ListCommand, error) {
	glazedLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, fmt.Errorf("failed to create glazed layer: %w", err)
	}

	cmdDesc := cmds.NewCommandDescription(
		"list",
		cmds.WithShort("List kitty windows and tabs"),
		cmds.WithLong(`List all windows and tabs in the kitty terminal emulator.

This command connects to kitty using its remote control protocol and retrieves
information about all open windows and tabs. The output can be formatted as
JSON, YAML, CSV, or table format using standard Glazed flags.

Examples:
  kitty-control list                                    # List all windows
  kitty-control list --output=json                     # Output as JSON
  kitty-control list --match="title:*vim*"            # Filter by title
  kitty-control list --self                           # List only current window
  kitty-control list --all-env-vars                   # Include all environment variables`),

		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"socket-path",
				parameters.ParameterTypeString,
				parameters.WithHelp("Path to kitty socket (defaults to KITTY_LISTEN_ON)"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"all-env-vars",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Include all environment variables for each window"),
				parameters.WithDefault(false),
			),
			parameters.NewParameterDefinition(
				"match",
				parameters.ParameterTypeString,
				parameters.WithHelp("Only list windows matching this pattern"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"match-tab",
				parameters.ParameterTypeString,
				parameters.WithHelp("Only list tabs matching this pattern"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"self",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Only list the window this command is run in"),
				parameters.WithDefault(false),
			),
		),

		cmds.WithLayersList(glazedLayer),
	)

	return &ListCommand{
		CommandDescription: cmdDesc,
	}, nil
}
