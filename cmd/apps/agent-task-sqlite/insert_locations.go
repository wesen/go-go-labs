package main

import (
	"context"
	"strings"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/go-go-golems/glazed/pkg/types"
	"github.com/pkg/errors"
)

// InsertLocationsCommand inserts code locations into the database
type InsertLocationsCommand struct {
	*cmds.CommandDescription
}

// InsertLocationsSettings holds the parameters for inserting locations
type InsertLocationsSettings struct {
	TaskID      int      `glazed.parameter:"task-id"`
	Location    string   `glazed.parameter:"location"`
	Description string   `glazed.parameter:"description"`
	Locations   []string `glazed.parameter:"locations"`
}

// Location represents a single code location
type Location struct {
	Location    string
	Description string
}

// Ensure interface implementation
var _ cmds.GlazeCommand = &InsertLocationsCommand{}

func (c *InsertLocationsCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	// Parse settings
	s := &InsertLocationsSettings{}
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

	// Parse locations
	var locations []Location
	
	// If individual parameters are provided, use them
	if s.Location != "" {
		locations = append(locations, Location{
			Location:    s.Location,
			Description: s.Description,
		})
	}
	
	// Parse bulk locations if provided
	for _, locStr := range s.Locations {
		loc, err := parseLocationString(locStr)
		if err != nil {
			return errors.Wrapf(err, "failed to parse location: %s", locStr)
		}
		locations = append(locations, loc)
	}

	if len(locations) == 0 {
		return errors.New("no locations provided")
	}

	// Start transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer tx.Rollback()

	// Insert all locations
	insertedCount := 0
	for _, loc := range locations {
		result, err := tx.ExecContext(ctx, `
			INSERT INTO gathered_locations (task_id, location, description)
			VALUES (?, ?, ?)
		`, s.TaskID, loc.Location, loc.Description)
		if err != nil {
			return errors.Wrapf(err, "failed to insert location %s", loc.Location)
		}

		// Get the inserted location ID
		locationID, err := result.LastInsertId()
		if err != nil {
			return errors.Wrap(err, "failed to get location ID")
		}

		// Output the created location
		row := types.NewRow(
			types.MRP("id", locationID),
			types.MRP("task_id", s.TaskID),
			types.MRP("location", loc.Location),
			types.MRP("description", loc.Description),
		)

		if err := gp.AddRow(ctx, row); err != nil {
			return errors.Wrap(err, "failed to add row")
		}

		insertedCount++
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// parseLocationString parses a location string in the format:
// "location:description"
func parseLocationString(locStr string) (Location, error) {
	parts := strings.SplitN(locStr, ":", 2)
	if len(parts) < 2 {
		return Location{}, errors.New("location format should be 'location:description'")
	}

	location := parts[0]
	description := parts[1]

	return Location{
		Location:    location,
		Description: description,
	}, nil
}

func NewInsertLocationsCommand() (interface{}, error) {
	// Create the Glazed layer for output formatting
	glazedLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, err
	}

	// Create command description
	cmdDesc := cmds.NewCommandDescription(
		"insert-locations",
		cmds.WithShort("Insert code locations into the database"),
		cmds.WithLong(`
Insert code locations associated with a task into the database.

You can insert locations in two ways:
1. Individual location using --location and --description
2. Multiple locations using --locations with format: "location:description"

Examples:
  # Insert a single location
  insert-locations --task-id=1 --location="main.go" --description="Main function"
  
  # Insert multiple locations
  insert-locations --task-id=1 --locations="main.go:Main function" --locations="auth.go:Authentication logic"
  
  # Mix both approaches
  insert-locations --task-id=1 --location="config.go" --description="Config struct" --locations="utils.go:Helper functions"
		`),
		// Define command flags
		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"task-id",
				parameters.ParameterTypeInteger,
				parameters.WithHelp("ID of the task to associate locations with"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"location",
				parameters.ParameterTypeString,
				parameters.WithHelp("Location reference (e.g., file path, function name, etc.)"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"description",
				parameters.ParameterTypeString,
				parameters.WithHelp("Description of the code location"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"locations",
				parameters.ParameterTypeStringList,
				parameters.WithHelp("List of locations in format 'location:description'"),
				parameters.WithDefault([]string{}),
			),
		),
		// Add parameter layers
		cmds.WithLayersList(
			glazedLayer,
		),
	)

	return &InsertLocationsCommand{
		CommandDescription: cmdDesc,
	}, nil
} 