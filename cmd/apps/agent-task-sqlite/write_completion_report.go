package main

import (
	"context"
	"io/ioutil"
	"database/sql"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/go-go-golems/glazed/pkg/types"
	"github.com/pkg/errors"
)

// WriteCompletionReportCommand creates a completion report for a task
type WriteCompletionReportCommand struct {
	*cmds.CommandDescription
}

// WriteCompletionReportSettings holds the parameters for writing a completion report
type WriteCompletionReportSettings struct {
	TaskID      string   `glazed.parameter:"task"`
	TaskSlug    string   `glazed.parameter:"task-slug"`
	ProjectID   string   `glazed.parameter:"project-id"`
	Content     string   `glazed.parameter:"content"`
	ReportFile  string   `glazed.parameter:"report-file"`
	LocationIDs []string `glazed.parameter:"location-ids"`
}

// Ensure interface implementation
var _ cmds.GlazeCommand = &WriteCompletionReportCommand{}

func (c *WriteCompletionReportCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	// Parse settings
	s := &WriteCompletionReportSettings{}
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

	// Get content from file or direct input
	var content string
	if s.ReportFile != "" {
		fileContent, err := ioutil.ReadFile(s.ReportFile)
		if err != nil {
			return errors.Wrap(err, "failed to read report file")
		}
		content = string(fileContent)
	} else if s.Content != "" {
		content = s.Content
	} else {
		return errors.New("must provide either --content or --report-file")
	}

	// Check if task exists and is completed
	var taskStatus string
	err = db.QueryRowContext(ctx, `
		SELECT status FROM tasks WHERE id = ?
	`, taskID).Scan(&taskStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.Errorf("task not found: %d", taskID)
		}
		return errors.Wrap(err, "failed to check task status")
	}

	if taskStatus != "completed" {
		return errors.Errorf("task must be completed before writing a report (current status: %s)", taskStatus)
	}

	// Start transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer tx.Rollback()

	// Insert the report
	result, err := tx.ExecContext(ctx, `
		INSERT INTO reports (task_id, content)
		VALUES (?, ?)
	`, taskID, content)
	if err != nil {
		return errors.Wrap(err, "failed to insert report")
	}

	// Get the inserted report ID
	reportID, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "failed to get report ID")
	}

	// Link locations to the report if provided
	if len(s.LocationIDs) > 0 {
		for _, locationIDStr := range s.LocationIDs {
			var locationID int
			var locationTaskID int
			
			// Verify location exists and belongs to the same task
			err = tx.QueryRowContext(ctx, `
				SELECT id, task_id FROM gathered_locations WHERE id = ?
			`, locationIDStr).Scan(&locationID, &locationTaskID)
			if err != nil {
				if err == sql.ErrNoRows {
					return errors.Errorf("location not found: %s", locationIDStr)
				}
				return errors.Wrap(err, "failed to verify location")
			}
			
			if locationTaskID != taskID {
				return errors.Errorf("location %d belongs to different task (%d vs %d)", locationID, locationTaskID, taskID)
			}

			// Link location to report
			_, err = tx.ExecContext(ctx, `
				INSERT INTO report_locations (report_id, location_id)
				VALUES (?, ?)
			`, reportID, locationID)
			if err != nil {
				return errors.Wrap(err, "failed to link location to report")
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
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

	// Output the report result
	row := types.NewRow(
		types.MRP("report_id", reportID),
		types.MRP("task_id", taskID),
		types.MRP("task_slug", taskSlug),
		types.MRP("project_name", projectName),
		types.MRP("content_length", len(content)),
		types.MRP("linked_locations", len(s.LocationIDs)),
	)

	if s.ReportFile != "" {
		row.Set("source_file", s.ReportFile)
	}

	return gp.AddRow(ctx, row)
}

func NewWriteCompletionReportCommand() (interface{}, error) {
	// Create the Glazed layer for output formatting
	glazedLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, err
	}

	// Create command description
	cmdDesc := cmds.NewCommandDescription(
		"write-completion-report",
		cmds.WithShort("Write a completion report for a finished task"),
		cmds.WithLong(`
Write a comprehensive completion report for a finished task.

The task must be in 'completed' status before a report can be written.
You can provide the report content directly or read it from a file.
Optionally link specific code locations to the report.

Examples:
  # Write a report with direct content
  write-completion-report --task=5 --content="Analysis complete. Found 3 security issues in authentication middleware."
  
  # Write a report from a file
  write-completion-report --task=auth-analysis/analysis-task --report-file="./analysis-report.md"
  
  # Write a report with linked locations
  write-completion-report --task=5 --content="Security analysis complete." --location-ids=1,2,3
  
  # Write a report using project/task slug
  write-completion-report --project-id=1 --task-slug=analysis-task --report-file="report.md"
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
				"content",
				parameters.ParameterTypeString,
				parameters.WithHelp("Report content (use instead of --report-file)"),
			),
			parameters.NewParameterDefinition(
				"report-file",
				parameters.ParameterTypeString,
				parameters.WithHelp("Path to file containing the report content"),
			),
			parameters.NewParameterDefinition(
				"location-ids",
				parameters.ParameterTypeStringList,
				parameters.WithHelp("Comma-separated list of location IDs to link to this report"),
			),
		),
		// Add parameter layers
		cmds.WithLayersList(
			glazedLayer,
		),
	)

	return &WriteCompletionReportCommand{
		CommandDescription: cmdDesc,
	}, nil
}
