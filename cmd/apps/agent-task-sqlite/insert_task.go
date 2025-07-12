package main

import (
	"context"
	"database/sql"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/go-go-golems/glazed/pkg/types"
	"github.com/pkg/errors"
)

// InsertTaskCommand inserts a new task into the database
type InsertTaskCommand struct {
	*cmds.CommandDescription
}

// InsertTaskSettings holds the parameters for inserting a task
type InsertTaskSettings struct {
	ProjectID    int      `glazed.parameter:"project-id"`
	AgentID      int      `glazed.parameter:"agent-id"`
	Type         string   `glazed.parameter:"type"`
	Instructions string   `glazed.parameter:"instructions"`
	Dependencies []int    `glazed.parameter:"dependencies"`
}

// Ensure interface implementation
var _ cmds.GlazeCommand = &InsertTaskCommand{}

func (c *InsertTaskCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	// Parse settings
	s := &InsertTaskSettings{}
	if err := parsedLayers.InitializeStruct(layers.DefaultSlug, s); err != nil {
		return errors.Wrap(err, "failed to parse settings")
	}

	// Initialize database
	db, err := InitDatabase()
	if err != nil {
		return errors.Wrap(err, "failed to initialize database")
	}
	defer db.Close()

	// Verify project exists
	var projectExists bool
	err = db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM projects WHERE id = ?)", s.ProjectID).Scan(&projectExists)
	if err != nil {
		return errors.Wrap(err, "failed to check project existence")
	}
	if !projectExists {
		return errors.Errorf("project with ID %d does not exist", s.ProjectID)
	}

	// Verify agent exists if specified
	if s.AgentID > 0 {
		var agentExists bool
		err = db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM agents WHERE id = ?)", s.AgentID).Scan(&agentExists)
		if err != nil {
			return errors.Wrap(err, "failed to check agent existence")
		}
		if !agentExists {
			return errors.Errorf("agent with ID %d does not exist", s.AgentID)
		}
	}

	// Start transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer tx.Rollback()

	// Insert the task
	var result sql.Result
	if s.AgentID > 0 {
		result, err = tx.ExecContext(ctx, `
			INSERT INTO tasks (project_id, agent_id, type, instructions, status)
			VALUES (?, ?, ?, ?, 'pending')
		`, s.ProjectID, s.AgentID, s.Type, s.Instructions)
	} else {
		result, err = tx.ExecContext(ctx, `
			INSERT INTO tasks (project_id, type, instructions, status)
			VALUES (?, ?, ?, 'pending')
		`, s.ProjectID, s.Type, s.Instructions)
	}
	if err != nil {
		return errors.Wrap(err, "failed to insert task")
	}

	// Get the inserted task ID
	taskID, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "failed to get task ID")
	}

	// Insert dependencies if any
	for _, parentID := range s.Dependencies {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO task_dependencies (task_id, parent_task_id)
			VALUES (?, ?)
		`, taskID, parentID)
		if err != nil {
			return errors.Wrapf(err, "failed to insert dependency on task %d", parentID)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	// Output the created task
	row := types.NewRow(
		types.MRP("id", taskID),
		types.MRP("project_id", s.ProjectID),
		types.MRP("agent_id", s.AgentID),
		types.MRP("type", s.Type),
		types.MRP("instructions", s.Instructions),
		types.MRP("status", "pending"),
		types.MRP("dependencies", s.Dependencies),
	)

	return gp.AddRow(ctx, row)
}

func NewInsertTaskCommand() (interface{}, error) {
	// Create the Glazed layer for output formatting
	glazedLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, err
	}

	// Create command description
	cmdDesc := cmds.NewCommandDescription(
		"insert-task",
		cmds.WithShort("Insert a new task into the database"),
		cmds.WithLong(`
Insert a new task into the agent task database.

Tasks can be of two types:
- gather_information: Tasks for collecting source code information
- oracle_analysis: Tasks for analyzing collected information

Dependencies can be specified as a list of parent task IDs that must be completed
before this task can be started.

Examples:
  # Insert a simple gather task
  insert-task --project-id=1 --type=gather_information --instructions="Gather authentication patterns"
  
  # Insert an analysis task that depends on task 1
  insert-task --project-id=1 --type=oracle_analysis --instructions="Analyze auth patterns" --dependencies=1
  
  # Insert a task with agent assignment and dependencies
  insert-task --project-id=1 --agent-id=2 --type=oracle_analysis --instructions="Compare patterns" --dependencies=1,2,3
		`),
		// Define command flags
		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"project-id",
				parameters.ParameterTypeInteger,
				parameters.WithHelp("ID of the project this task belongs to"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"agent-id",
				parameters.ParameterTypeInteger,
				parameters.WithHelp("ID of the agent assigned to this task (optional)"),
				parameters.WithDefault(0),
			),
			parameters.NewParameterDefinition(
				"type",
				parameters.ParameterTypeChoice,
				parameters.WithHelp("Type of task to create"),
				parameters.WithChoices("gather_information", "oracle_analysis"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"instructions",
				parameters.ParameterTypeString,
				parameters.WithHelp("Instructions for the task"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"dependencies",
				parameters.ParameterTypeIntegerList,
				parameters.WithHelp("List of parent task IDs this task depends on"),
				parameters.WithDefault([]int{}),
			),
		),
		// Add parameter layers
		cmds.WithLayersList(
			glazedLayer,
		),
	)

	return &InsertTaskCommand{
		CommandDescription: cmdDesc,
	}, nil
} 