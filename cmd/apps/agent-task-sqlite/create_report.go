package main

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/go-go-golems/glazed/pkg/types"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// CreateReportCommand creates a report for a task
type CreateReportCommand struct {
	*cmds.CommandDescription
}

// CreateReportSettings holds the parameters for creating a report
type CreateReportSettings struct {
	TaskID      int      `glazed.parameter:"task-id"`
	Content     string   `glazed.parameter:"content"`
	ContentFile string   `glazed.parameter:"content-file"`
	LocationIDs []int    `glazed.parameter:"location-ids"`
}

// Ensure interface implementation
var _ cmds.GlazeCommand = &CreateReportCommand{}

func (c *CreateReportCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	// Parse settings
	s := &CreateReportSettings{}
	if err := parsedLayers.InitializeStruct(layers.DefaultSlug, s); err != nil {
		return errors.Wrap(err, "failed to parse settings")
	}

	// Initialize database
	db, err := InitDatabase()
	if err != nil {
		return errors.Wrap(err, "failed to initialize database")
	}
	defer db.Close()

	// Verify task exists
	var taskExists bool
	err = db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM tasks WHERE id = ?)", s.TaskID).Scan(&taskExists)
	if err != nil {
		return errors.Wrap(err, "failed to check task existence")
	}
	if !taskExists {
		return errors.Errorf("task with ID %d does not exist", s.TaskID)
	}

	// Get report content
	content := s.Content
	if s.ContentFile != "" {
		fileContent, err := readContentFromFile(s.ContentFile)
		if err != nil {
			return errors.Wrapf(err, "failed to read content from file: %s", s.ContentFile)
		}
		if content != "" {
			content += "\n\n" + fileContent
		} else {
			content = fileContent
		}
	}

	if content == "" {
		return errors.New("no content provided for report")
	}

	// Validate location IDs if provided
	if len(s.LocationIDs) > 0 {
		for _, locationID := range s.LocationIDs {
			var locationExists bool
			err = db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM gathered_locations WHERE id = ?)", locationID).Scan(&locationExists)
			if err != nil {
				return errors.Wrapf(err, "failed to check location existence for ID %d", locationID)
			}
			if !locationExists {
				return errors.Errorf("location with ID %d does not exist", locationID)
			}
		}
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
	`, s.TaskID, content)
	if err != nil {
		return errors.Wrap(err, "failed to insert report")
	}

	// Get the inserted report ID
	reportID, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "failed to get report ID")
	}

	// Link report to locations if provided
	for _, locationID := range s.LocationIDs {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO report_locations (report_id, location_id)
			VALUES (?, ?)
		`, reportID, locationID)
		if err != nil {
			return errors.Wrapf(err, "failed to link report to location %d", locationID)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	// Get linked locations for output
	linkedLocations, err := getLinkedLocations(ctx, db, reportID)
	if err != nil {
		return errors.Wrap(err, "failed to get linked locations")
	}

	// Output the created report
	row := types.NewRow(
		types.MRP("id", reportID),
		types.MRP("task_id", s.TaskID),
		types.MRP("content_length", len(content)),
		types.MRP("content_preview", getContentPreview(content, 100)),
		types.MRP("linked_locations", len(linkedLocations)),
		types.MRP("location_ids", linkedLocations),
	)

	return gp.AddRow(ctx, row)
}

// readContentFromFile reads content from a file, handling stdin if filename is "-"
func readContentFromFile(filename string) (string, error) {
	var reader io.Reader
	
	if filename == "-" {
		reader = os.Stdin
	} else {
		file, err := os.Open(filename)
		if err != nil {
			return "", err
		}
		defer file.Close()
		reader = file
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// getContentPreview returns a preview of the content, truncated to maxLength
func getContentPreview(content string, maxLength int) string {
	// Replace newlines with spaces for preview
	preview := strings.ReplaceAll(content, "\n", " ")
	preview = strings.ReplaceAll(preview, "\r", " ")
	
	// Collapse multiple spaces
	for strings.Contains(preview, "  ") {
		preview = strings.ReplaceAll(preview, "  ", " ")
	}
	
	preview = strings.TrimSpace(preview)
	
	if len(preview) <= maxLength {
		return preview
	}
	
	return preview[:maxLength-3] + "..."
}

// getLinkedLocations retrieves the location IDs linked to a report
func getLinkedLocations(ctx context.Context, db *sqlx.DB, reportID int64) ([]int, error) {
	query := `
		SELECT location_id
		FROM report_locations
		WHERE report_id = ?
		ORDER BY location_id
	`

	rows, err := db.QueryContext(ctx, query, reportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locationIDs []int
	for rows.Next() {
		var locationID int
		if err := rows.Scan(&locationID); err != nil {
			return nil, err
		}
		locationIDs = append(locationIDs, locationID)
	}

	return locationIDs, rows.Err()
}

func NewCreateReportCommand() (interface{}, error) {
	// Create the Glazed layer for output formatting
	glazedLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, err
	}

	// Create command description
	cmdDesc := cmds.NewCommandDescription(
		"create-report",
		cmds.WithShort("Create a report for a task"),
		cmds.WithLong(`
Create a report associated with a task in the database.

Report content can be provided directly via --content or read from a file using --content-file.
Use "-" as the filename to read from stdin.

You can optionally link the report to specific code locations using --location-ids.

Examples:
  # Create a report with direct content
  create-report --task-id=1 --content="Analysis complete. Found 3 authentication patterns."
  
  # Create a report from a file
  create-report --task-id=1 --content-file="analysis.md"
  
  # Create a report from stdin
  echo "Report content" | create-report --task-id=1 --content-file="-"
  
  # Create a report linked to specific locations
  create-report --task-id=1 --content="Analysis of auth patterns" --location-ids=1,2,3
  
  # Combine content and file
  create-report --task-id=1 --content="Summary:" --content-file="details.txt"
		`),
		// Define command flags
		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"task-id",
				parameters.ParameterTypeInteger,
				parameters.WithHelp("ID of the task to create report for"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"content",
				parameters.ParameterTypeString,
				parameters.WithHelp("Report content"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"content-file",
				parameters.ParameterTypeString,
				parameters.WithHelp("File to read content from (use '-' for stdin)"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"location-ids",
				parameters.ParameterTypeIntegerList,
				parameters.WithHelp("List of location IDs to link to this report"),
				parameters.WithDefault([]int{}),
			),
		),
		// Add parameter layers
		cmds.WithLayersList(
			glazedLayer,
		),
	)

	return &CreateReportCommand{
		CommandDescription: cmdDesc,
	}, nil
} 