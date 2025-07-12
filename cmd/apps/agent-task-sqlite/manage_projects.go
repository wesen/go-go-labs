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

// CreateProjectCommand creates a new project
type CreateProjectCommand struct {
	*cmds.CommandDescription
}

// CreateProjectSettings holds the parameters for creating a project
type CreateProjectSettings struct {
	Name              string `glazed.parameter:"name"`
	Description       string `glazed.parameter:"description"`
	Slug              string `glazed.parameter:"slug"`
	ConciseGuidelines string `glazed.parameter:"concise-guidelines"`
	FullGuidelines    string `glazed.parameter:"full-guidelines"`
}

// Ensure interface implementation
var _ cmds.GlazeCommand = &CreateProjectCommand{}

func (c *CreateProjectCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	// Parse settings
	s := &CreateProjectSettings{}
	if err := parsedLayers.InitializeStruct(layers.DefaultSlug, s); err != nil {
		return errors.Wrap(err, "failed to parse settings")
	}

	// Initialize database
	db, err := InitDatabase()
	if err != nil {
		return errors.Wrap(err, "failed to initialize database")
	}
	defer db.Close()

	// Generate slug if not provided
	slug := s.Slug
	if slug == "" {
		slug = GenerateSlug(s.Name)
	}

	// Ensure slug is unique
	slug, err = EnsureUniqueSlug(db, "projects", slug)
	if err != nil {
		return errors.Wrap(err, "failed to ensure unique slug")
	}

	// Insert the project
	result, err := db.ExecContext(ctx, `
		INSERT INTO projects (slug, name, description, concise_guidelines, full_guidelines)
		VALUES (?, ?, ?, ?, ?)
	`, slug, s.Name, s.Description, s.ConciseGuidelines, s.FullGuidelines)
	if err != nil {
		return errors.Wrap(err, "failed to insert project")
	}

	// Get the inserted project ID
	projectID, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "failed to get project ID")
	}

	// Output the created project
	row := types.NewRow(
		types.MRP("id", projectID),
		types.MRP("slug", slug),
		types.MRP("name", s.Name),
		types.MRP("description", s.Description),
		types.MRP("concise_guidelines", s.ConciseGuidelines),
		types.MRP("full_guidelines", s.FullGuidelines),
	)

	return gp.AddRow(ctx, row)
}

func NewCreateProjectCommand() (interface{}, error) {
	// Create the Glazed layer for output formatting
	glazedLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, err
	}

	// Create command description
	cmdDesc := cmds.NewCommandDescription(
		"create-project",
		cmds.WithShort("Create a new project"),
		cmds.WithLong(`
Create a new project in the database.

Projects are used to organize tasks and provide context for analysis work.

Examples:
  # Create a new project
  create-project --name="Authentication Analysis" --description="Analyze authentication patterns across the codebase"
		`),
		// Define command flags
		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"name",
				parameters.ParameterTypeString,
				parameters.WithHelp("Name of the project"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"description",
				parameters.ParameterTypeString,
				parameters.WithHelp("Description of the project"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"slug",
				parameters.ParameterTypeString,
				parameters.WithHelp("URL-friendly slug for the project (auto-generated if not provided)"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"concise-guidelines",
				parameters.ParameterTypeString,
				parameters.WithHelp("Short guidelines displayed when creating/completing tasks"),
			),
			parameters.NewParameterDefinition(
				"full-guidelines",
				parameters.ParameterTypeString,
				parameters.WithHelp("Detailed project guidelines and working instructions"),
			),
		),
		// Add parameter layers
		cmds.WithLayersList(
			glazedLayer,
		),
	)

	return &CreateProjectCommand{
		CommandDescription: cmdDesc,
	}, nil
}

// ListProjectsCommand lists all projects
type ListProjectsCommand struct {
	*cmds.CommandDescription
}

// ListProjectsSettings holds the parameters for listing projects
type ListProjectsSettings struct {
	ProjectID string `glazed.parameter:"project"`
}

// Ensure interface implementation
var _ cmds.GlazeCommand = &ListProjectsCommand{}

func (c *ListProjectsCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	// Parse settings
	s := &ListProjectsSettings{}
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
		SELECT id, slug, name, description, created_at, updated_at
		FROM projects
		WHERE 1=1
	`
	args := []interface{}{}

	// Add filters
	if s.ProjectID != "" {
		query += " AND (id = ? OR slug = ?)"
		args = append(args, s.ProjectID, s.ProjectID)
	}

	query += " ORDER BY created_at DESC"

	// Execute query
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to query projects")
	}
	defer rows.Close()

	// Process results
	for rows.Next() {
		var id int
		var slug, name, description, createdAt, updatedAt string

		err := rows.Scan(&id, &slug, &name, &description, &createdAt, &updatedAt)
		if err != nil {
			return errors.Wrap(err, "failed to scan project row")
		}

		// Create row data
		rowData := types.NewRow(
			types.MRP("id", id),
			types.MRP("slug", slug),
			types.MRP("name", name),
			types.MRP("description", description),
			types.MRP("created_at", createdAt),
			types.MRP("updated_at", updatedAt),
		)

		if err := gp.AddRow(ctx, rowData); err != nil {
			return errors.Wrap(err, "failed to add row")
		}
	}

	return rows.Err()
}

func NewListProjectsCommand() (interface{}, error) {
	// Create the Glazed layer for output formatting
	glazedLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, err
	}

	// Create command description
	cmdDesc := cmds.NewCommandDescription(
		"list-projects",
		cmds.WithShort("List projects"),
		cmds.WithLong(`
List all projects in the database.

Examples:
  # List all projects
  list-projects
  
  # Show specific project by ID
  list-projects --project=1
  
  # Show specific project by slug
  list-projects --project=auth-analysis
		`),
		// Define command flags
		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"project",
				parameters.ParameterTypeString,
				parameters.WithHelp("Show specific project by ID or slug"),
				parameters.WithDefault(""),
			),
		),
		// Add parameter layers
		cmds.WithLayersList(
			glazedLayer,
		),
	)

	return &ListProjectsCommand{
		CommandDescription: cmdDesc,
	}, nil
}
