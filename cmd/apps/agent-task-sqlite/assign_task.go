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

// AssignTaskCommand assigns a task to an agent
type AssignTaskCommand struct {
	*cmds.CommandDescription
}

// AssignTaskSettings holds the parameters for assigning a task
type AssignTaskSettings struct {
	AgentSlug string `glazed.parameter:"agent"`
	TaskID    string `glazed.parameter:"task"`
	Force     bool   `glazed.parameter:"force"`
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
	agentID, err := ResolveAgentID(ctx, db, s.AgentSlug)
	if err != nil {
		return errors.Wrap(err, "failed to resolve agent")
	}

	// Resolve task ID
	taskID, err := ResolveTaskID(ctx, db, s.TaskID)
	if err != nil {
		return errors.Wrap(err, "failed to resolve task")
	}

	// Check if task is available (pending status) or if force is being used
	var currentStatus string
	var currentAgentID sql.NullInt64
	var projectID int
	err = db.QueryRowContext(ctx, `
		SELECT status, agent_id, project_id FROM tasks WHERE id = ?
	`, taskID).Scan(&currentStatus, &currentAgentID, &projectID)
	if err != nil {
		return errors.Wrap(err, "failed to check task status")
	}

	// Check if task is already assigned
	if currentAgentID.Valid && currentAgentID.Int64 != 0 {
		if !s.Force {
			return errors.Errorf("task is already assigned to agent %d (use --force to reassign)", currentAgentID.Int64)
		}
		log.Debug().Int64("current_agent_id", currentAgentID.Int64).Int("new_agent_id", agentID).Msg("Force reassigning task from one agent to another")
	}

	if currentStatus != "pending" && currentStatus != "in_progress" {
		return errors.Errorf("task cannot be assigned (current status: %s)", currentStatus)
	}

	// Check if all dependencies are completed
	log.Debug().Int("task_id", taskID).Msg("Checking task dependencies")
	dependencies, err := getTaskDependencies(ctx, db, taskID)
	if err != nil {
		return errors.Wrap(err, "failed to get task dependencies")
	}
	
	if len(dependencies) > 0 {
		log.Debug().Interface("dependencies", dependencies).Msg("Task has dependencies, checking completion status")
		for _, depID := range dependencies {
			var depStatus string
			err := db.QueryRowContext(ctx, "SELECT status FROM tasks WHERE id = ?", depID).Scan(&depStatus)
			if err != nil {
				return errors.Wrapf(err, "failed to check status of dependency task %d", depID)
			}
			log.Debug().Int("dependency_id", depID).Str("status", depStatus).Msg("Dependency status")
			
			if depStatus != "completed" {
				return errors.Errorf("cannot assign task: dependency task %d is not completed (current status: %s)", depID, depStatus)
			}
		}
		log.Debug().Msg("All dependencies are completed")
	} else {
		log.Debug().Msg("Task has no dependencies")
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

	// Clear any existing agent assignment for this task if force is used or task is being reassigned
	if s.Force || (currentAgentID.Valid && currentAgentID.Int64 != 0) {
		log.Debug().Int("task_id", taskID).Bool("force", s.Force).Msg("Clearing existing agent assignment for task")
		_, err = tx.ExecContext(ctx, `
			UPDATE agents 
			SET current_project_id = NULL, current_task_id = NULL, updated_at = CURRENT_TIMESTAMP
			WHERE current_task_id = ?
		`, taskID)
		if err != nil {
			return errors.Wrap(err, "failed to clear existing agent assignment")
		}
	}

	// Assign task to agent and set status to in_progress
	log.Debug().Int("agent_id", agentID).Int("task_id", taskID).Msg("Assigning task to agent")
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

Use --force to reassign a task that is already assigned to another agent.

Examples:
  # Assign task by ID to agent by slug
  assign-task --agent=code-analyzer --task=1
  
  # Assign task by project/task slug to agent by slug
  assign-task --agent=code-analyzer --task=auth-analysis/gather-patterns
  
  # Force reassign a task from one agent to another
  assign-task --agent=new-agent --task=1 --force
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
			parameters.NewParameterDefinition(
				"force",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Force reassignment even if task is already assigned to another agent"),
				parameters.WithDefault(false),
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