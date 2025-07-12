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

// GetGuidelinesCommand displays project guidelines
type GetGuidelinesCommand struct {
	*cmds.CommandDescription
}

// GetGuidelinesSettings holds the parameters for getting guidelines
type GetGuidelinesSettings struct {
	ProjectID string `glazed.parameter:"project"`
	Concise   bool   `glazed.parameter:"concise"`
}

// Ensure interface implementation
var _ cmds.GlazeCommand = &GetGuidelinesCommand{}

func (c *GetGuidelinesCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	// Parse settings
	s := &GetGuidelinesSettings{}
	if err := parsedLayers.InitializeStruct(layers.DefaultSlug, s); err != nil {
		return errors.Wrap(err, "failed to parse settings")
	}

	// Initialize database
	db, err := InitDatabase()
	if err != nil {
		return errors.Wrap(err, "failed to initialize database")
	}
	defer db.Close()

	// Resolve project ID
	projectID, err := ResolveProjectID(ctx, db, s.ProjectID)
	if err != nil {
		return errors.Wrap(err, "failed to resolve project")
	}

	// Get project details
	var projectName, projectSlug, description, conciseGuidelines, fullGuidelines sql.NullString
	err = db.QueryRowContext(ctx, `
		SELECT name, slug, description, concise_guidelines, full_guidelines
		FROM projects
		WHERE id = ?
	`, projectID).Scan(&projectName, &projectSlug, &description, &conciseGuidelines, &fullGuidelines)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.Errorf("project not found: %s", s.ProjectID)
		}
		return errors.Wrap(err, "failed to get project details")
	}

	// Create row data
	row := types.NewRow(
		types.MRP("project_id", projectID),
		types.MRP("project_slug", projectSlug.String),
		types.MRP("project_name", projectName.String),
		types.MRP("description", description.String),
	)

	// Add guidelines based on what's requested
	if s.Concise {
		if conciseGuidelines.Valid && conciseGuidelines.String != "" {
			row.Set("concise_guidelines", conciseGuidelines.String)
		} else {
			row.Set("concise_guidelines", "No concise guidelines set")
		}
	} else {
		// Show both by default
		if conciseGuidelines.Valid && conciseGuidelines.String != "" {
			row.Set("concise_guidelines", conciseGuidelines.String)
		} else {
			row.Set("concise_guidelines", "No concise guidelines set")
		}
		
		if fullGuidelines.Valid && fullGuidelines.String != "" {
			row.Set("full_guidelines", fullGuidelines.String)
		} else {
			row.Set("full_guidelines", "No full guidelines set")
		}
	}

	return gp.AddRow(ctx, row)
}

func NewGetGuidelinesCommand() (interface{}, error) {
	// Create the Glazed layer for output formatting
	glazedLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, err
	}

	// Create command description
	cmdDesc := cmds.NewCommandDescription(
		"get-guidelines",
		cmds.WithShort("Display project guidelines"),
		cmds.WithLong(`
Display project guidelines for working on the project.

Shows both concise guidelines (displayed during task operations) 
and full guidelines (detailed working instructions) by default.
Use --concise to show only the concise guidelines.

Examples:
  # Show all guidelines for a project
  get-guidelines --project=auth-analysis
  
  # Show only concise guidelines
  get-guidelines --project=1 --concise
		`),
		// Define command flags
		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"project",
				parameters.ParameterTypeString,
				parameters.WithHelp("Project ID or slug"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"concise",
				parameters.ParameterTypeBool,
				parameters.WithHelp("Show only concise guidelines"),
				parameters.WithDefault(false),
			),
		),
		// Add parameter layers
		cmds.WithLayersList(
			glazedLayer,
		),
	)

	return &GetGuidelinesCommand{
		CommandDescription: cmdDesc,
	}, nil
}
