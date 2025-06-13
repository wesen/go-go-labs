package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"

	"github.com/go-go-golems/go-go-labs/cmd/apps/kitty-control/pkg"
)

type GetTextCommand struct {
	*cmds.CommandDescription
}

type GetTextSettings struct {
	SocketPath     string `glazed.parameter:"socket-path"`
	Match          string `glazed.parameter:"match"`
	Extent         string `glazed.parameter:"extent"`
	Ansi           bool   `glazed.parameter:"ansi"`
	Cursor         bool   `glazed.parameter:"cursor"`
	WrapMarkers    bool   `glazed.parameter:"wrap-markers"`
	ClearSelection bool   `glazed.parameter:"clear-selection"`
	Self           bool   `glazed.parameter:"self"`
}

var _ cmds.WriterCommand = &GetTextCommand{}

func (c *GetTextCommand) RunIntoWriter(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
	w io.Writer,
) error {
	s := &GetTextSettings{}
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
	payload := &pkg.GetTextPayload{}

	if s.Match != "" {
		payload.Match = &s.Match
	}
	if s.Extent != "" {
		payload.Extent = &s.Extent
	}
	if s.Ansi {
		payload.Ansi = &s.Ansi
	}
	if s.Cursor {
		payload.Cursor = &s.Cursor
	}
	if s.WrapMarkers {
		payload.WrapMarkers = &s.WrapMarkers
	}
	if s.ClearSelection {
		payload.ClearSelection = &s.ClearSelection
	}
	if s.Self {
		payload.Self = &s.Self
	}

	// Send command
	response, err := client.GetText(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to get text: %w", err)
	}

	// Parse response
	var kittyResponse struct {
		OK   bool   `json:"ok"`
		Data string `json:"data"`
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

	// Write the text content to the writer
	_, err = w.Write([]byte(kittyResponse.Data))
	if err != nil {
		return fmt.Errorf("failed to write text content: %w", err)
	}

	return nil
}

func NewGetTextCommand() (*GetTextCommand, error) {
	cmdDesc := cmds.NewCommandDescription(
		"get-text",
		cmds.WithShort("Get text from kitty windows"),
		cmds.WithLong(`Get text content from one or more kitty windows.

This command retrieves text content from kitty windows based on the specified
extent and options. You can get screen content, command output, or selections
with optional ANSI formatting codes and cursor information.

The extent parameter determines what text to retrieve:
- screen: Currently visible screen content
- first_cmd_output_on_screen: First command output visible on screen
- last_cmd_output: Last command output regardless of visibility
- last_visited_cmd_output: Last visited command output
- all: All scrollback buffer content
- selection: Currently selected text

Examples:
  kitty-control get-text                                    # Get screen content from active window
  kitty-control get-text --extent=all                      # Get all scrollback content
  kitty-control get-text --match="title:*vim*" --ansi      # Get text with ANSI codes from vim windows
  kitty-control get-text --extent=selection                # Get currently selected text
  kitty-control get-text --extent=last_cmd_output --cursor # Get last command output with cursor info
  kitty-control get-text --self --wrap-markers             # Get text from current window with wrap markers`),

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
				parameters.WithHelp("Pattern to match windows"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"extent",
				parameters.ParameterTypeChoice,
				parameters.WithHelp("What text to retrieve"),
				parameters.WithChoices("screen", "first_cmd_output_on_screen", "last_cmd_output", "last_visited_cmd_output", "all", "selection"),
				parameters.WithDefault("screen"),
			),
			parameters.NewParameterDefinition(
				"ansi",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Include ANSI formatting codes in output"),
				parameters.WithDefault(false),
			),
			parameters.NewParameterDefinition(
				"cursor",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Include cursor position/style as ANSI codes"),
				parameters.WithDefault(false),
			),
			parameters.NewParameterDefinition(
				"wrap-markers",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Add wrap markers to indicate line wrapping"),
				parameters.WithDefault(false),
			),
			parameters.NewParameterDefinition(
				"clear-selection",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Clear the selection in the matched window after getting text"),
				parameters.WithDefault(false),
			),
			parameters.NewParameterDefinition(
				"self",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Get text from the window this command is run in"),
				parameters.WithDefault(false),
			),
		),
	)

	return &GetTextCommand{
		CommandDescription: cmdDesc,
	}, nil
}
