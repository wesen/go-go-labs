package main

import (
	"fmt"
	"os"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/logging"
	"github.com/go-go-golems/glazed/pkg/help"
	help_cmd "github.com/go-go-golems/glazed/pkg/help/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	// Create root command
	rootCmd := &cobra.Command{
		Use:   "agent-task-sqlite",
		Short: "Agent task database management tool",
		Long: `A command-line tool for managing agent tasks, locations, and reports in SQLite.
		
This tool provides a convenient interface for working with the agent task database,
allowing you to insert tasks, query their status, manage code locations, and create reports.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := logging.InitLoggerFromViper()
			cobra.CheckErr(err)
		},
	}

	// Initialize help system
	helpSystem := help.NewHelpSystem()
	help_cmd.SetupCobraRootCommand(helpSystem, rootCmd)

	// Add logging flags
	err := logging.AddLoggingLayerToRootCommand(rootCmd, "agent-task-sqlite")
	cobra.CheckErr(err)

	err = viper.BindPFlags(rootCmd.PersistentFlags())
	cobra.CheckErr(err)

	err = logging.InitLoggerFromViper()
	cobra.CheckErr(err)

	// Create and add all commands
	commands := []func() (interface{}, error){
		NewInsertTaskCommand,
		NewQueryTasksCommand,
		NewInsertLocationsCommand,
		NewCreateReportCommand,
		NewCreateProjectCommand,
		NewListProjectsCommand,
		NewCreateAgentCommand,
		NewListAgentsCommand,
		NewAssignTaskCommand,
		NewCompleteTaskCommand,
		NewServeCommand,
		NewTakeNoteCommand,
		NewWriteCompletionReportCommand,
		NewGetNotesCommand,
		NewGetReportCommand,
		NewGetGuidelinesCommand,
	}

	for _, cmdFactory := range commands {
		cmdInterface, err := cmdFactory()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating command: %v\n", err)
			os.Exit(1)
		}

		// Cast to the appropriate command interface
		var cmd cmds.Command

		switch c := cmdInterface.(type) {
		case *InsertTaskCommand:
			cmd = c
		case *QueryTasksCommand:
			cmd = c
		case *InsertLocationsCommand:
			cmd = c
		case *CreateReportCommand:
			cmd = c
		case *CreateProjectCommand:
			cmd = c
		case *ListProjectsCommand:
			cmd = c
		case *CreateAgentCommand:
			cmd = c
		case *ListAgentsCommand:
			cmd = c
		case *AssignTaskCommand:
			cmd = c
		case *CompleteTaskCommand:
			cmd = c
		case *ServeCommand:
			cmd = c
		case *TakeNoteCommand:
			cmd = c
		case *WriteCompletionReportCommand:
			cmd = c
		case *GetNotesCommand:
			cmd = c
		case *GetReportCommand:
			cmd = c
		case *GetGuidelinesCommand:
			cmd = c
		default:
			fmt.Fprintf(os.Stderr, "Unknown command type: %T\n", cmdInterface)
			os.Exit(1)
		}

		// Convert to Cobra command
		cobraCmd, err := cli.BuildCobraCommandFromCommand(cmd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error building command: %v\n", err)
			os.Exit(1)
		}

		rootCmd.AddCommand(cobraCmd)
	}

	// Execute
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
