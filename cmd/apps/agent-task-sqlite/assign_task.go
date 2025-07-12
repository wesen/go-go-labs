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

// AssignTaskCommand assigns a task to an agent
type AssignTaskCommand struct {
	*cmds.CommandDescription
}

// AssignTaskSettings holds the parameters for assigning a task
type AssignTaskSettings struct {
	AgentSlug string `glazed.parameter:"agent"`
	TaskID    string `glazed.parameter:"task"`
}

// Ensure interface implementation
var _ cmds.GlazeCommand = &AssignTaskCommand{}

func (c *AssignTaskCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	// Parse settings
	s := &AssignTaskSettings{}
	if err := parsedLayers.InitializeStruct(layers.DefaultSlug, s); err != nil {
		return errors.Wrap(err, "failed to parse settings")
	}

	// Initialize database
	db, err := InitDatabase()
	if err != nil {
		return errors.Wrap(err, "failed to initialize database")
	}
	defer db.Close()

	// Resolve agent ID
	agentID, err := ResolveAgentID(db, s.AgentSlug)
	if err != nil {
		return errors.Wrap(err, "failed to resolve agent")
	}

	// Resolve task ID
	taskID, err := ResolveTaskID(db, s.TaskID)
	if err != nil {
		return errors.Wrap(err, "failed to resolve task")
	}

	// Check if task is available (pending status)
	var currentStatus string
	var currentAgentID sql.NullInt64
	var projectID int
	err = db.QueryRowContext(ctx, `
		SELECT status, agent_id, project_id FROM tasks WHERE id = ?
	`, taskID).Scan(&currentStatus, &currentAgentID, &projectID)
	if err != nil {
		return errors.Wrap(err, "failed to check task status")
	}

	if currentStatus != "pending" {
		return errors.Errorf("task is not available for assignment (current status: %s)", currentStatus)
	}

	// Check if agent is already working on another task
	var currentTaskID sql.NullInt64
	err = db.QueryRowContext(ctx, `
		SELECT current_task_id FROM agents WHERE id = ?
	`, agentID).Scan(&currentTaskID)
	if err != nil {
		return errors.Wrap(err, "failed to check agent current task")
	}

	if currentTaskID.Valid {
		return errors.Errorf("agent is already assigned to task %d", currentTaskID.Int64)
	}

	// Start transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer tx.Rollback()

	// Assign task to agent and set status to in_progress
	_, err = tx.ExecContext(ctx, `
		UPDATE tasks 
		SET agent_id = ?, status = 'in_progress', started_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, agentID, taskID)
	if err != nil {
		return errors.Wrap(err, "failed to assign task")
	}

	// Update agent's current work
	_, err = tx.ExecContext(ctx, `
		UPDATE agents 
		SET current_project_id = ?, current_task_id = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, projectID, taskID, agentID)
	if err != nil {
		return errors.Wrap(err, "failed to update agent")
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	// Get task details for output
	var taskSlug, taskType, instructions string
	err = db.QueryRowContext(ctx, `
		SELECT slug, type, instructions FROM tasks WHERE id = ?
	`, taskID).Scan(&taskSlug, &taskType, &instructions)
	if err != nil {
		return errors.Wrap(err, "failed to get task details")
	}

	// Get agent details for output
	var agentSlug, agentName string
	err = db.QueryRowContext(ctx, `
		SELECT slug, name FROM agents WHERE id = ?
	`, agentID).Scan(&agentSlug, &agentName)
	if err != nil {
		return errors.Wrap(err, "failed to get agent details")
	}

	// Output the assignment result
	row := types.NewRow(
		types.MRP("task_id", taskID),
		types.MRP("task_slug", taskSlug),
		types.MRP("task_type", taskType),
		types.MRP("agent_id", agentID),
		types.MRP("agent_slug", agentSlug),
		types.MRP("agent_name", agentName),
		types.MRP("status", "in_progress"),
		types.MRP("instructions", instructions),
	)

	return gp.AddRow(ctx, row)
}

func NewAssignTaskCommand() (interface{}, error) {
	// Create the Glazed layer for output formatting
	glazedLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, err
	}

	// Create command description
	cmdDesc := cmds.NewCommandDescription(
		"assign-task",
		cmds.WithShort("Assign a task to an agent"),
		cmds.WithLong(`
Assign a task to an agent and set it to in_progress status.

The task must be in pending status to be assigned. The agent must not be currently 
working on another task.

Examples:
  # Assign task by ID to agent by slug
  assign-task --agent=code-analyzer --task=1
  
  # Assign task by project/task slug to agent by slug
  assign-task --agent=code-analyzer --task=auth-analysis/gather-patterns
		`),
		// Define command flags
		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"agent",
				parameters.ParameterTypeString,
				parameters.WithHelp("Agent slug or ID to assign the task to"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"task",
				parameters.ParameterTypeString,
				parameters.WithHelp("Task ID or project_slug/task_slug to assign"),
				parameters.WithRequired(true),
			),
		),
		// Add parameter layers
		cmds.WithLayersList(
			glazedLayer,
		),
	)

	return &AssignTaskCommand{
		CommandDescription: cmdDesc,
	}, nil
} 