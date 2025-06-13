package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"

	"github.com/go-go-golems/go-go-labs/cmd/apps/kitty-control/pkg"
)

type EnvCommand struct {
	*cmds.CommandDescription
}

type EnvSettings struct {
	SocketPath string            `glazed.parameter:"socket-path"`
	Env        map[string]string `glazed.parameter:"env"`
}

var _ cmds.BareCommand = &EnvCommand{}

func (c *EnvCommand) Run(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
) error {
	s := &EnvSettings{}
	if err := parsedLayers.InitializeStruct(layers.DefaultSlug, s); err != nil {
		return fmt.Errorf("failed to initialize settings: %w", err)
	}

	if len(s.Env) == 0 {
		return fmt.Errorf("at least one environment variable must be specified")
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
	payload := &pkg.EnvPayload{
		Env: s.Env,
	}

	// Send command
	err = client.SetEnv(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to set environment variables: %w", err)
	}

	// Print confirmation
	var envVars []string
	for key, value := range s.Env {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}
	fmt.Printf("Set environment variables: %s\n", strings.Join(envVars, ", "))

	return nil
}

func NewEnvCommand() (*EnvCommand, error) {
	cmdDesc := cmds.NewCommandDescription(
		"env",
		cmds.WithShort("Set environment variables in kitty"),
		cmds.WithLong(`Set environment variables that will be available to new processes launched in kitty.

This command allows you to set environment variables that will be inherited by
any new processes (windows, tabs, etc.) launched in kitty after this command
is executed. This is useful for setting up development environments, API keys,
or configuration variables.

The environment variables are specified as key=value pairs using the --env flag.
You can use multiple --env flags to set multiple variables at once.

Examples:
  kitty-control env --env=PATH=/usr/local/bin:$PATH --env=DEBUG=1
  kitty-control env --env=API_KEY=secret123 --env=NODE_ENV=development
  kitty-control env --env=EDITOR=vim --env=BROWSER=firefox`),

		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"socket-path",
				parameters.ParameterTypeString,
				parameters.WithHelp("Path to kitty socket (defaults to KITTY_LISTEN_ON)"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"env",
				parameters.ParameterTypeKeyValue,
				parameters.WithHelp("Environment variables to set (key=value format, can be used multiple times)"),
				parameters.WithRequired(true),
			),
		),
	)

	return &EnvCommand{
		CommandDescription: cmdDesc,
	}, nil
}
