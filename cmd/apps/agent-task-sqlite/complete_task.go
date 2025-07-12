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
	"github.com/rs/zerolog/log"
)

// CompleteTaskCommand marks a task as completed
type CompleteTaskCommand struct {
	*cmds.CommandDescription
}

// CompleteTaskSettings holds the parameters for completing a task
type CompleteTaskSettings struct {
	TaskID string `glazed.parameter:"task"`
	Notes  string `glazed.parameter:"notes"`
}

// Ensure interface implementation
var _ cmds.GlazeCommand = &CompleteTaskCommand{}

func (c *CompleteTaskCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	// Parse settings
	s := &CompleteTaskSettings{}
	if err := parsedLayers.InitializeStruct(layers.DefaultSlug, s); err != nil {
		return errors.Wrap(err, "failed to parse settings")
	}

	// Initialize database
	db, err := InitDatabase()
	if err != nil {
		return errors.Wrap(err, "failed to initialize database")
	}
	defer db.Close()

	// Resolve task ID
	taskID, err := ResolveTaskID(ctx, db, s.TaskID)
	if err != nil {
		return errors.Wrap(err, "failed to resolve task")
	}

	// Check current task status and get agent information
	var currentStatus string
	var currentAgentID sql.NullInt64
	var projectID int
	var taskSlug, taskType, instructions string
	err = db.QueryRowContext(ctx, `
		SELECT status, agent_id, project_id, slug, type, instructions 
		FROM tasks WHERE id = ?
	`, taskID).Scan(&currentStatus, &currentAgentID, &projectID, &taskSlug, &taskType, &instructions)
	if err != nil {
		return errors.Wrap(err, "failed to check task status")
	}

	// Only allow completion if task is in_progress
	if currentStatus != "in_progress" {
		return errors.Errorf("task cannot be completed (current status: %s). Task must be in_progress to be completed", currentStatus)
	}

	// Start transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer tx.Rollback()

	// Mark task as completed
	log.Debug().Int("task_id", taskID).Str("notes", s.Notes).Msg("Marking task as completed")
	_, err = tx.ExecContext(ctx, `
		UPDATE tasks 
		SET status = 'completed', completed_at = CURRENT_TIMESTAMP, completion_notes = ?
		WHERE id = ?
	`, s.Notes, taskID)
	if err != nil {
		return errors.Wrap(err, "failed to complete task")
	}

	// Clear agent's current work assignment if assigned
	if currentAgentID.Valid && currentAgentID.Int64 != 0 {
		log.Debug().Int64("agent_id", currentAgentID.Int64).Msg("Clearing agent assignment after task completion")
		_, err = tx.ExecContext(ctx, `
			UPDATE agents 
			SET current_project_id = NULL, current_task_id = NULL, updated_at = CURRENT_TIMESTAMP
			WHERE id = ?
		`, currentAgentID.Int64)
		if err != nil {
			return errors.Wrap(err, "failed to clear agent assignment")
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	// Get agent details for output (if assigned)
	var agentSlug, agentName sql.NullString
	if currentAgentID.Valid && currentAgentID.Int64 != 0 {
		err = db.QueryRowContext(ctx, `
			SELECT slug, name FROM agents WHERE id = ?
		`, currentAgentID.Int64).Scan(&agentSlug, &agentName)
		if err != nil {
			log.Warn().Err(err).Int64("agent_id", currentAgentID.Int64).Msg("Failed to get agent details for output")
		}
	}

	// Output the completion result
	row := types.NewRow(
		types.MRP("task_id", taskID),
		types.MRP("task_slug", taskSlug),
		types.MRP("task_type", taskType),
		types.MRP("status", "completed"),
		types.MRP("instructions", instructions),
	)

	// Add agent info if available
	if currentAgentID.Valid && currentAgentID.Int64 != 0 {
		row.Set("agent_id", currentAgentID.Int64)
		if agentSlug.Valid {
			row.Set("agent_slug", agentSlug.String)
		}
		if agentName.Valid {
			row.Set("agent_name", agentName.String)
		}
	}

	// Add completion notes if provided
	if s.Notes != "" {
		row.Set("completion_notes", s.Notes)
	}

	// Get project guidelines to display
	var conciseGuidelines sql.NullString
	err = db.QueryRowContext(ctx, `
		SELECT concise_guidelines FROM projects 
		WHERE id = (SELECT project_id FROM tasks WHERE id = ?)
	`, taskID).Scan(&conciseGuidelines)
	if err != nil {
		log.Warn().Err(err).Int("task_id", taskID).Msg("Failed to get project guidelines")
	}

	// Add project guidelines if available
	if conciseGuidelines.Valid && conciseGuidelines.String != "" {
		row.Set("project_guidelines", conciseGuidelines.String)
	}

	// Add reminder
	row.Set("reminder", "Don't forget to provide a full report using write-completion-report.")

	return gp.AddRow(ctx, row)
}

func NewCompleteTaskCommand() (interface{}, error) {
	// Create the Glazed layer for output formatting
	glazedLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, err
	}

	// Create command description
	cmdDesc := cmds.NewCommandDescription(
		"complete-task",
		cmds.WithShort("Mark a task as completed"),
		cmds.WithLong(`
Mark a task as completed and clear agent assignment.

The task must be in 'in_progress' status to be completed. This will:
- Set the task status to 'completed'
- Set the completion timestamp
- Clear the agent's current task assignment
- Optionally record completion notes

Examples:
  # Complete task by ID
  complete-task --task=1
  
  # Complete task by project/task slug
  complete-task --task=auth-analysis/gather-patterns
  
  # Complete task with notes
  complete-task --task=1 --notes="Found 3 security issues in authentication flow"
		`),
		// Define command flags
		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"task",
				parameters.ParameterTypeString,
				parameters.WithHelp("Task ID or project_slug/task_slug to complete"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"notes",
				parameters.ParameterTypeString,
				parameters.WithHelp("Optional completion notes or summary"),
				parameters.WithDefault(""),
			),
		),
		// Add parameter layers
		cmds.WithLayersList(
			glazedLayer,
		),
	)

	return &CompleteTaskCommand{
		CommandDescription: cmdDesc,
	}, nil
}
