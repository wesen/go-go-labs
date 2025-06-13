package cmd

import (
	"context"
	"fmt"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/rs/zerolog/log"

	"github.com/go-go-golems/go-go-labs/cmd/apps/kitty-control/pkg"
)

type SendTextCommand struct {
	*cmds.CommandDescription
}

type SendTextSettings struct {
	SocketPath     string `glazed.parameter:"socket-path"`
	Text           string `glazed.parameter:"text"`
	Match          string `glazed.parameter:"match"`
	MatchTab       string `glazed.parameter:"match-tab"`
	All            bool   `glazed.parameter:"all"`
	ExcludeActive  bool   `glazed.parameter:"exclude-active"`
	SessionID      string `glazed.parameter:"session-id"`
	BracketedPaste string `glazed.parameter:"bracketed-paste"`
}

var _ cmds.BareCommand = &SendTextCommand{}

func (c *SendTextCommand) Run(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
) error {
	log.Debug().Msg("executing send-text command")
	s := &SendTextSettings{}
	if err := parsedLayers.InitializeStruct(layers.DefaultSlug, s); err != nil {
		log.Error().Err(err).Msg("failed to initialize settings")
		return fmt.Errorf("failed to initialize settings: %w", err)
	}
	log.Debug().Str("text", s.Text).Str("match", s.Match).Bool("all", s.All).Msg("send-text settings initialized")

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
	payload := &pkg.SendTextPayload{
		Data: s.Text,
	}

	if s.Match != "" {
		payload.Match = &s.Match
	}
	if s.MatchTab != "" {
		payload.MatchTab = &s.MatchTab
	}
	if s.All {
		payload.All = &s.All
	}
	if s.ExcludeActive {
		payload.ExcludeActive = &s.ExcludeActive
	}
	if s.SessionID != "" {
		payload.SessionID = &s.SessionID
	}
	if s.BracketedPaste != "" {
		payload.BracketedPaste = &s.BracketedPaste
	}

	// Send command
	log.Debug().Str("text", s.Text).Interface("payload", payload).Msg("sending text to kitty")
	err = client.SendText(ctx, s.Text, payload)
	if err != nil {
		log.Error().Err(err).Str("text", s.Text).Msg("failed to send text")
		return fmt.Errorf("failed to send text: %w", err)
	}

	log.Info().Str("text", s.Text).Msg("text sent successfully to kitty window(s)")
	fmt.Printf("Text sent successfully to kitty window(s)\n")
	return nil
}

func NewSendTextCommand() (*SendTextCommand, error) {
	cmdDesc := cmds.NewCommandDescription(
		"send-text",
		cmds.WithShort("Send text to kitty windows"),
		cmds.WithLong(`Send text to one or more kitty windows.

This command allows you to send text input to kitty windows, as if it was typed
by the user. You can target specific windows using match patterns, send to all
windows, or exclude the active window.

Examples:
  kitty-control send-text --text="Hello World"              # Send to active window
  kitty-control send-text --text="ls -la" --all            # Send to all windows
  kitty-control send-text --text="exit" --match="title:*sh*" # Send to shell windows
  kitty-control send-text --text="command" --bracketed-paste=enable # Use bracketed paste`),

		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"socket-path",
				parameters.ParameterTypeString,
				parameters.WithHelp("Path to kitty socket (defaults to KITTY_LISTEN_ON)"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"text",
				parameters.ParameterTypeString,
				parameters.WithHelp("Text to send to the window(s)"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"match",
				parameters.ParameterTypeString,
				parameters.WithHelp("Pattern to match windows"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"match-tab",
				parameters.ParameterTypeString,
				parameters.WithHelp("Pattern to match tabs"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"all",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Send to all windows"),
				parameters.WithDefault(false),
			),
			parameters.NewParameterDefinition(
				"exclude-active",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Exclude the active window"),
				parameters.WithDefault(false),
			),
			parameters.NewParameterDefinition(
				"session-id",
				parameters.ParameterTypeString,
				parameters.WithHelp("Session ID for broadcast sessions"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"bracketed-paste",
				parameters.ParameterTypeChoice,
				parameters.WithHelp("Whether to use bracketed paste mode"),
				parameters.WithChoices("disable", "enable"),
				parameters.WithDefault("disable"),
			),
		),
	)

	return &SendTextCommand{
		CommandDescription: cmdDesc,
	}, nil
}
