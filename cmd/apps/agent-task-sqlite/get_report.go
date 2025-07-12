package main

import (
	"context"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/go-go-golems/glazed/pkg/types"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// GetReportCommand retrieves reports for a task
type GetReportCommand struct {
	*cmds.CommandDescription
}

// GetReportSettings holds the parameters for getting reports
type GetReportSettings struct {
	TaskID       string `glazed.parameter:"task"`
	TaskSlug     string `glazed.parameter:"task-slug"`
	ProjectID    string `glazed.parameter:"project-id"`
	WithLocations bool  `glazed.parameter:"with-locations"`
	Limit        int    `glazed.parameter:"limit"`
}

// Ensure interface implementation
var _ cmds.GlazeCommand = &GetReportCommand{}

func (c *GetReportCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	// Parse settings
	s := &GetReportSettings{}
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
		SELECT r.id, r.task_id, r.content, r.created_at,
		       t.slug as task_slug, p.name as project_name
		FROM reports r
		JOIN tasks t ON r.task_id = t.id
		JOIN projects p ON t.project_id = p.id
		WHERE r.task_id = ?
		ORDER BY r.created_at DESC
	`
	args := []interface{}{taskID}

	// Add limit if specified
	if s.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, s.Limit)
	}

	// Execute query
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to query reports")
	}
	defer rows.Close()

	// Process results
	for rows.Next() {
		var report struct {
			ID          int    `db:"id"`
			TaskID      int    `db:"task_id"`
			Content     string `db:"content"`
			CreatedAt   string `db:"created_at"`
			TaskSlug    string `db:"task_slug"`
			ProjectName string `db:"project_name"`
		}

		err := rows.Scan(
			&report.ID, &report.TaskID, &report.Content, &report.CreatedAt,
			&report.TaskSlug, &report.ProjectName,
		)
		if err != nil {
			return errors.Wrap(err, "failed to scan report")
		}

		row := types.NewRow(
			types.MRP("report_id", report.ID),
			types.MRP("task_id", report.TaskID),
			types.MRP("task_slug", report.TaskSlug),
			types.MRP("project_name", report.ProjectName),
			types.MRP("content", report.Content),
			types.MRP("content_length", len(report.Content)),
			types.MRP("created_at", report.CreatedAt),
		)

		// Add linked locations if requested
		if s.WithLocations {
			locations, err := c.getReportLocations(ctx, db, report.ID)
			if err != nil {
				return errors.Wrap(err, "failed to get report locations")
			}
			row.Set("locations", locations)
			row.Set("location_count", len(locations))
		}

		if err := gp.AddRow(ctx, row); err != nil {
			return errors.Wrap(err, "failed to add row")
		}
	}

	return rows.Err()
}

func (c *GetReportCommand) getReportLocations(ctx context.Context, db interface{}, reportID int) ([]string, error) {
	query := `
		SELECT gl.location, gl.description
		FROM gathered_locations gl
		JOIN report_locations rl ON gl.id = rl.location_id
		WHERE rl.report_id = ?
		ORDER BY gl.created_at DESC
	`

	rows, err := db.(*sqlx.DB).QueryContext(ctx, query, reportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []string
	for rows.Next() {
		var location, description string
		if err := rows.Scan(&location, &description); err != nil {
			return nil, err
		}
		locations = append(locations, location+" - "+description)
	}

	return locations, rows.Err()
}

func NewGetReportCommand() (interface{}, error) {
	// Create the Glazed layer for output formatting
	glazedLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, err
	}

	// Create command description
	cmdDesc := cmds.NewCommandDescription(
		"get-report",
		cmds.WithShort("Retrieve reports for a task"),
		cmds.WithLong(`
Retrieve completion reports written for a specific task.

Reports are returned in reverse chronological order (newest first).
Optionally include linked code locations with --with-locations.

Examples:
  # Get all reports for a task
  get-report --task=5
  
  # Get reports with linked locations
  get-report --task=auth-analysis/analysis-task --with-locations
  
  # Get latest report only
  get-report --project-id=1 --task-slug=analysis-task --limit=1
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
				"with-locations",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Include linked code locations in the output"),
				parameters.WithDefault(false),
			),
			parameters.NewParameterDefinition(
				"limit",
				parameters.ParameterTypeInteger,
				parameters.WithHelp("Maximum number of reports to return"),
				parameters.WithDefault(0),
			),
		),
		// Add parameter layers
		cmds.WithLayersList(
			glazedLayer,
		),
	)

	return &GetReportCommand{
		CommandDescription: cmdDesc,
	}, nil
}
