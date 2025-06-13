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

type LaunchCommand struct {
	*cmds.CommandDescription
}

type LaunchSettings struct {
	SocketPath  string   `glazed.parameter:"socket-path"`
	Args        []string `glazed.parameter:"args"`
	Match       string   `glazed.parameter:"match"`
	WindowTitle string   `glazed.parameter:"window-title"`
	Cwd         string   `glazed.parameter:"cwd"`
	Env         []string `glazed.parameter:"env"`
	Type        string   `glazed.parameter:"type"`
	KeepFocus   bool     `glazed.parameter:"keep-focus"`
	Location    string   `glazed.parameter:"location"`
	Hold        bool     `glazed.parameter:"hold"`
}

var _ cmds.BareCommand = &LaunchCommand{}

func (c *LaunchCommand) Run(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
) error {
	s := &LaunchSettings{}
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
	payload := &pkg.LaunchPayload{
		Args: s.Args,
	}

	if s.Match != "" {
		payload.Match = &s.Match
	}
	if s.WindowTitle != "" {
		payload.WindowTitle = &s.WindowTitle
	}
	if s.Cwd != "" {
		payload.Cwd = &s.Cwd
	}
	if len(s.Env) > 0 {
		payload.Env = s.Env
	}
	if s.Type != "" {
		payload.Type = &s.Type
	}
	if s.KeepFocus {
		payload.KeepFocus = &s.KeepFocus
	}
	if s.Location != "" {
		payload.Location = &s.Location
	}
	if s.Hold {
		payload.Hold = &s.Hold
	}

	// Send command
	err = client.Launch(ctx, s.Args, payload)
	if err != nil {
		return fmt.Errorf("failed to launch: %w", err)
	}

	fmt.Printf("Launched command: %s\n", strings.Join(s.Args, " "))
	return nil
}

func NewLaunchCommand() (*LaunchCommand, error) {
	cmdDesc := cmds.NewCommandDescription(
		"launch",
		cmds.WithShort("Launch a new window/tab with a command"),
		cmds.WithLong(`Launch a new window, tab, or OS window running the specified command.

This command allows you to create new kitty windows running any command you specify.
You can control where the window is created, its title, working directory, environment
variables, and more.

Window Types:
  window    - Create a new window in the current tab (default)
  tab       - Create a new tab with a window
  os-window - Create a new OS window

Location Options (for type=window):
  default - Default location
  hsplit  - Horizontal split (side by side)
  vsplit  - Vertical split (stacked)
  after   - After current window
  before  - Before current window
  neighbor - Near current window

Examples:
  kitty-control launch "htop"                              # Launch htop in new window
  kitty-control launch "vim" "file.txt"                   # Launch vim with file.txt
  kitty-control launch --type=tab "bash"                  # Launch bash in new tab
  kitty-control launch --type=os-window "firefox"         # Launch firefox in new OS window
  kitty-control launch --location=hsplit "tail" "-f" "/var/log/syslog"  # Split horizontally
  kitty-control launch --cwd="/tmp" "ls" "-la"            # Launch in /tmp directory
  kitty-control launch --window-title="System Monitor" "htop"  # Set window title
  kitty-control launch --env="DEBUG=1" --env="LOG_LEVEL=trace" "myapp"  # Set environment
  kitty-control launch --hold "date"                      # Keep window open after command exits
  kitty-control launch --keep-focus "background-task"     # Don't switch focus to new window`),

		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"socket-path",
				parameters.ParameterTypeString,
				parameters.WithHelp("Path to kitty socket (defaults to KITTY_LISTEN_ON)"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"args",
				parameters.ParameterTypeStringList,
				parameters.WithHelp("Command and arguments to launch"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"match",
				parameters.ParameterTypeString,
				parameters.WithHelp("Which tab to open the new window in"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"window-title",
				parameters.ParameterTypeString,
				parameters.WithHelp("Title for the new window"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"cwd",
				parameters.ParameterTypeString,
				parameters.WithHelp("Working directory for the command"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"env",
				parameters.ParameterTypeStringList,
				parameters.WithHelp("Environment variables (format: KEY=value)"),
			),
			parameters.NewParameterDefinition(
				"type",
				parameters.ParameterTypeChoice,
				parameters.WithHelp("Type of window to create"),
				parameters.WithChoices("window", "tab", "os-window"),
				parameters.WithDefault("window"),
			),
			parameters.NewParameterDefinition(
				"keep-focus",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Keep focus on current window"),
				parameters.WithDefault(false),
			),
			parameters.NewParameterDefinition(
				"location",
				parameters.ParameterTypeChoice,
				parameters.WithHelp("Where in tab to open new window"),
				parameters.WithChoices("default", "hsplit", "vsplit", "after", "before", "neighbor"),
				parameters.WithDefault("default"),
			),
			parameters.NewParameterDefinition(
				"hold",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Keep window open after command exits"),
				parameters.WithDefault(false),
			),
		),
	)

	return &LaunchCommand{
		CommandDescription: cmdDesc,
	}, nil
}
