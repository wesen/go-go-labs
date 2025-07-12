package main

import (
	"context"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/go-go-golems/glazed/pkg/types"
	"github.com/pkg/errors"
)

// GetNotesCommand retrieves notes for a task
type GetNotesCommand struct {
	*cmds.CommandDescription
}

// GetNotesSettings holds the parameters for getting notes
type GetNotesSettings struct {
	TaskID    string `glazed.parameter:"task"`
	TaskSlug  string `glazed.parameter:"task-slug"`
	ProjectID string `glazed.parameter:"project-id"`
	Type      string `glazed.parameter:"type"`
	Limit     int    `glazed.parameter:"limit"`
}

// Ensure interface implementation
var _ cmds.GlazeCommand = &GetNotesCommand{}

func (c *GetNotesCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	// Parse settings
	s := &GetNotesSettings{}
	if err := parsedLayers.InitializeStruct(layers.DefaultSlug, s); err != nil {
		return errors.Wrap(err, "failed to parse settings")
	}

	// Initialize database
	db, err := InitDatabase()
	if err != nil {
		return errors.Wrap(err, "failed to initialize database")
	}
	defer db.Close()

	// Resolve task ID with fallback options
	var taskID int
	
	if s.TaskID != "" {
		taskID, err = ResolveTaskID(ctx, db, s.TaskID)
		if err != nil {
			return errors.Wrap(err, "failed to resolve task")
		}
	} else if s.TaskSlug != "" && s.ProjectID != "" {
		// Resolve using project-id and task-slug combination
		projectID, err := ResolveProjectID(ctx, db, s.ProjectID)
		if err != nil {
			return errors.Wrap(err, "failed to resolve project")
		}
		
		err = db.QueryRowContext(ctx, `
			SELECT id FROM tasks WHERE project_id = ? AND slug = ?
		`, projectID, s.TaskSlug).Scan(&taskID)
		if err != nil {
			return errors.Errorf("task not found: project %s, task %s", s.ProjectID, s.TaskSlug)
		}
	} else {
		return errors.New("must provide either --task or both --project-id and --task-slug")
	}

	// Build query
	query := `
		SELECT tn.id, tn.task_id, tn.type, tn.content, tn.created_at,
		       t.slug as task_slug, p.name as project_name
		FROM task_notes tn
		JOIN tasks t ON tn.task_id = t.id
		JOIN projects p ON t.project_id = p.id
		WHERE tn.task_id = ?
	`
	args := []interface{}{taskID}

	// Add type filter if specified
	if s.Type != "" {
		query += " AND tn.type = ?"
		args = append(args, s.Type)
	}

	// Add ordering
	query += " ORDER BY tn.created_at DESC"

	// Add limit if specified
	if s.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, s.Limit)
	}

	// Execute query
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to query notes")
	}
	defer rows.Close()

	// Process results
	for rows.Next() {
		var note struct {
			ID          int    `db:"id"`
			TaskID      int    `db:"task_id"`
			Type        string `db:"type"`
			Content     string `db:"content"`
			CreatedAt   string `db:"created_at"`
			TaskSlug    string `db:"task_slug"`
			ProjectName string `db:"project_name"`
		}

		err := rows.Scan(
			&note.ID, &note.TaskID, &note.Type, &note.Content, &note.CreatedAt,
			&note.TaskSlug, &note.ProjectName,
		)
		if err != nil {
			return errors.Wrap(err, "failed to scan note")
		}

		row := types.NewRow(
			types.MRP("note_id", note.ID),
			types.MRP("task_id", note.TaskID),
			types.MRP("task_slug", note.TaskSlug),
			types.MRP("project_name", note.ProjectName),
			types.MRP("type", note.Type),
			types.MRP("content", note.Content),
			types.MRP("created_at", note.CreatedAt),
		)

		if err := gp.AddRow(ctx, row); err != nil {
			return errors.Wrap(err, "failed to add row")
		}
	}

	return rows.Err()
}

func NewGetNotesCommand() (interface{}, error) {
	// Create the Glazed layer for output formatting
	glazedLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, err
	}

	// Create command description
	cmdDesc := cmds.NewCommandDescription(
		"get-notes",
		cmds.WithShort("Retrieve notes for a task"),
		cmds.WithLong(`
Retrieve notes recorded for a specific task.

You can filter by note type and limit the number of results.
Notes are returned in reverse chronological order (newest first).

Examples:
  # Get all notes for a task
  get-notes --task=5
  
  # Get only bug reports for a task
  get-notes --task=auth-analysis/gather-patterns --type=bugs
  
  # Get last 10 notes for a task
  get-notes --project-id=1 --task-slug=gather-patterns --limit=10
  
  # Get observations only
  get-notes --task=5 --type=observations
		`),
		// Define command flags
		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"task",
				parameters.ParameterTypeString,
				parameters.WithHelp("Task ID or project_slug/task_slug"),
			),
			parameters.NewParameterDefinition(
				"task-slug",
				parameters.ParameterTypeString,
				parameters.WithHelp("Task slug (use with --project-id)"),
			),
			parameters.NewParameterDefinition(
				"project-id",
				parameters.ParameterTypeString,
				parameters.WithHelp("Project ID or slug (use with --task-slug)"),
			),
			parameters.NewParameterDefinition(
				"type",
				parameters.ParameterTypeString,
				parameters.WithHelp("Filter by note type (observations, lessons_learned, notes, bugs, ideas, issues)"),
			),
			parameters.NewParameterDefinition(
				"limit",
				parameters.ParameterTypeInteger,
				parameters.WithHelp("Maximum number of notes to return"),
				parameters.WithDefault(0),
			),
		),
		// Add parameter layers
		cmds.WithLayersList(
			glazedLayer,
		),
	)

	return &GetNotesCommand{
		CommandDescription: cmdDesc,
	}, nil
}
