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
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// QueryTasksCommand queries tasks from the database
type QueryTasksCommand struct {
	*cmds.CommandDescription
}

// QueryTasksSettings holds the parameters for querying tasks
type QueryTasksSettings struct {
	Status       string `glazed.parameter:"status"`
	Type         string `glazed.parameter:"type"`
	ProjectID    int    `glazed.parameter:"project-id"`
	AgentID      int    `glazed.parameter:"agent-id"`
	Limit        int    `glazed.parameter:"limit"`
	ShowDeps     bool   `glazed.parameter:"show-deps"`
	TaskID       int    `glazed.parameter:"task-id"`
}

// Ensure interface implementation
var _ cmds.GlazeCommand = &QueryTasksCommand{}

func (c *QueryTasksCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	// Parse settings
	s := &QueryTasksSettings{}
	if err := parsedLayers.InitializeStruct(layers.DefaultSlug, s); err != nil {
		return errors.Wrap(err, "failed to parse settings")
	}

	// Initialize database
	db, err := InitDatabase()
	if err != nil {
		return errors.Wrap(err, "failed to initialize database")
	}
	defer db.Close()

	// Build query
	query := `
		SELECT id, project_id, agent_id, type, status, instructions, created_at, started_at, completed_at
		FROM tasks
		WHERE 1=1
	`
	args := []interface{}{}

	// Add filters
	if s.Status != "" {
		query += " AND status = ?"
		args = append(args, s.Status)
	}
	if s.Type != "" {
		query += " AND type = ?"
		args = append(args, s.Type)
	}
	if s.ProjectID > 0 {
		query += " AND project_id = ?"
		args = append(args, s.ProjectID)
	}
	if s.AgentID > 0 {
		query += " AND agent_id = ?"
		args = append(args, s.AgentID)
	}
	if s.TaskID > 0 {
		query += " AND id = ?"
		args = append(args, s.TaskID)
	}

	// Add ordering and limit
	query += " ORDER BY created_at DESC"
	if s.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, s.Limit)
	}

	// Execute query
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to query tasks")
	}
	defer rows.Close()

	// Process results
	for rows.Next() {
		var id, projectID int
		var agentID sql.NullInt64
		var taskType, status, instructions string
		var createdAt string
		var startedAt, completedAt sql.NullString

		err := rows.Scan(&id, &projectID, &agentID, &taskType, &status, &instructions, &createdAt, &startedAt, &completedAt)
		if err != nil {
			return errors.Wrap(err, "failed to scan task row")
		}

		// Create row data
		rowData := types.NewRow(
			types.MRP("id", id),
			types.MRP("project_id", projectID),
			types.MRP("type", taskType),
			types.MRP("status", status),
			types.MRP("instructions", instructions),
			types.MRP("created_at", createdAt),
		)

		// Add agent_id if present
		if agentID.Valid {
			rowData.Set("agent_id", agentID.Int64)
		}

		// Add optional fields
		if startedAt.Valid {
			rowData.Set("started_at", startedAt.String)
		}
		if completedAt.Valid {
			rowData.Set("completed_at", completedAt.String)
		}

		// Add dependencies if requested
		if s.ShowDeps {
			deps, err := getTaskDependencies(ctx, db, id)
			if err != nil {
				return errors.Wrapf(err, "failed to get dependencies for task %d", id)
			}
			rowData.Set("dependencies", deps)
		}

		if err := gp.AddRow(ctx, rowData); err != nil {
			return errors.Wrap(err, "failed to add row")
		}
	}

	return rows.Err()
}

// getTaskDependencies retrieves the parent task IDs for a given task
func getTaskDependencies(ctx context.Context, db *sqlx.DB, taskID int) ([]int, error) {
	query := `
		SELECT parent_task_id
		FROM task_dependencies
		WHERE task_id = ?
		ORDER BY parent_task_id
	`

	rows, err := db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deps []int
	for rows.Next() {
		var parentID int
		if err := rows.Scan(&parentID); err != nil {
			return nil, err
		}
		deps = append(deps, parentID)
	}

	return deps, rows.Err()
}

func NewQueryTasksCommand() (interface{}, error) {
	// Create the Glazed layer for output formatting
	glazedLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, err
	}

	// Create command description
	cmdDesc := cmds.NewCommandDescription(
		"query-tasks",
		cmds.WithShort("Query tasks from the database"),
		cmds.WithLong(`
Query tasks from the agent task database with optional filtering.

You can filter by status, type, or specific task ID. The results are ordered
by creation time (newest first) and can be limited.

Examples:
  # Show all tasks
  query-tasks
  
  # Show only pending tasks
  query-tasks --status=pending
  
  # Show tasks for a specific project
  query-tasks --project-id=1
  
  # Show tasks assigned to a specific agent
  query-tasks --agent-id=2
  
  # Show gather tasks with dependencies
  query-tasks --type=gather_information --show-deps
  
  # Show specific task with full details
  query-tasks --task-id=5 --show-deps
  
  # Show last 10 completed tasks for project 1
  query-tasks --project-id=1 --status=completed --limit=10
		`),
		// Define command flags
		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"status",
				parameters.ParameterTypeChoice,
				parameters.WithHelp("Filter by task status"),
				parameters.WithChoices("", "pending", "in_progress", "completed", "failed"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"type",
				parameters.ParameterTypeChoice,
				parameters.WithHelp("Filter by task type"),
				parameters.WithChoices("", "gather_information", "oracle_analysis"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"project-id",
				parameters.ParameterTypeInteger,
				parameters.WithHelp("Filter by project ID"),
				parameters.WithDefault(0),
			),
			parameters.NewParameterDefinition(
				"agent-id",
				parameters.ParameterTypeInteger,
				parameters.WithHelp("Filter by agent ID"),
				parameters.WithDefault(0),
			),
			parameters.NewParameterDefinition(
				"limit",
				parameters.ParameterTypeInteger,
				parameters.WithHelp("Maximum number of tasks to return"),
				parameters.WithDefault(0),
			),
			parameters.NewParameterDefinition(
				"show-deps",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Show task dependencies"),
				parameters.WithDefault(false),
			),
			parameters.NewParameterDefinition(
				"task-id",
				parameters.ParameterTypeInteger,
				parameters.WithHelp("Show specific task by ID"),
				parameters.WithDefault(0),
			),
		),
		// Add parameter layers
		cmds.WithLayersList(
			glazedLayer,
		),
	)

	return &QueryTasksCommand{
		CommandDescription: cmdDesc,
	}, nil
} 