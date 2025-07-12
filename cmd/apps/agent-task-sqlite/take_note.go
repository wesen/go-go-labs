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

// TakeNoteCommand records a note for a task
type TakeNoteCommand struct {
	*cmds.CommandDescription
}

// TakeNoteSettings holds the parameters for taking a note
type TakeNoteSettings struct {
	TaskID     string `glazed.parameter:"task"`
	TaskSlug   string `glazed.parameter:"task-slug"`
	ProjectID  string `glazed.parameter:"project-id"`
	Type       string `glazed.parameter:"type"`
	Content    string `glazed.parameter:"content"`
}

// Ensure interface implementation
var _ cmds.GlazeCommand = &TakeNoteCommand{}

func (c *TakeNoteCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	// Parse settings
	s := &TakeNoteSettings{}
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

	// Set default note type
	noteType := s.Type
	if noteType == "" {
		noteType = "notes"
	}

	// Insert the note
	result, err := db.ExecContext(ctx, `
		INSERT INTO task_notes (task_id, type, content)
		VALUES (?, ?, ?)
	`, taskID, noteType, s.Content)
	if err != nil {
		return errors.Wrap(err, "failed to insert note")
	}

	// Get the inserted note ID
	noteID, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "failed to get note ID")
	}

	// Get task details for output
	var taskSlug, projectName string
	err = db.QueryRowContext(ctx, `
		SELECT t.slug, p.name FROM tasks t
		JOIN projects p ON t.project_id = p.id
		WHERE t.id = ?
	`, taskID).Scan(&taskSlug, &projectName)
	if err != nil {
		return errors.Wrap(err, "failed to get task details")
	}

	// Output the note result
	row := types.NewRow(
		types.MRP("note_id", noteID),
		types.MRP("task_id", taskID),
		types.MRP("task_slug", taskSlug),
		types.MRP("project_name", projectName),
		types.MRP("type", noteType),
		types.MRP("content", s.Content),
	)

	return gp.AddRow(ctx, row)
}

func NewTakeNoteCommand() (interface{}, error) {
	// Create the Glazed layer for output formatting
	glazedLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, err
	}

	// Create command description
	cmdDesc := cmds.NewCommandDescription(
		"take-note",
		cmds.WithShort("Record a note for a task"),
		cmds.WithLong(`
Record a note, observation, bug, or lesson learned while working on a task.

Notes can be of different types:
- observations: General observations during task execution
- lessons_learned: Important lessons to remember
- notes: General notes
- bugs: Bug reports and issues encountered
- ideas: Ideas for future improvements
- issues: Issues that need to be addressed

Examples:
  # Take a general note
  take-note --task=5 --content="Found the authentication logic in auth.go"
  
  # Record a bug
  take-note --task=auth-analysis/gather-patterns --type=bugs --content="Authentication middleware crashes with nil pointer on line 42"
  
  # Record an observation
  take-note --project-id=1 --task-slug=gather-patterns --type=observations --content="Most auth code is in middleware package"
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
				parameters.WithHelp("Type of note (observations, lessons_learned, notes, bugs, ideas, issues)"),
				parameters.WithDefault("notes"),
			),
			parameters.NewParameterDefinition(
				"content",
				parameters.ParameterTypeString,
				parameters.WithHelp("Content of the note"),
				parameters.WithRequired(true),
			),
		),
		// Add parameter layers
		cmds.WithLayersList(
			glazedLayer,
		),
	)

	return &TakeNoteCommand{
		CommandDescription: cmdDesc,
	}, nil
}
