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

// CreateAgentCommand creates a new agent
type CreateAgentCommand struct {
	*cmds.CommandDescription
}

// CreateAgentSettings holds the parameters for creating an agent
type CreateAgentSettings struct {
	Name        string `glazed.parameter:"name"`
	Description string `glazed.parameter:"description"`
	Slug        string `glazed.parameter:"slug"`
	Project     string `glazed.parameter:"project"`
}

// Ensure interface implementation
var _ cmds.GlazeCommand = &CreateAgentCommand{}

func (c *CreateAgentCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	// Parse settings
	s := &CreateAgentSettings{}
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
	projectID, err := ResolveProjectID(ctx, db, s.Project)
	if err != nil {
		return errors.Wrap(err, "failed to resolve project")
	}

	// Generate slug if not provided
	slug := s.Slug
	if slug == "" {
		slug = GenerateSlug(s.Name)
	}

	// Ensure slug is unique
	slug, err = EnsureUniqueSlug(db, "agents", slug)
	if err != nil {
		return errors.Wrap(err, "failed to ensure unique slug")
	}

	// Insert the agent
	result, err := db.ExecContext(ctx, `
		INSERT INTO agents (slug, name, description, current_project_id)
		VALUES (?, ?, ?, ?)
	`, slug, s.Name, s.Description, projectID)
	if err != nil {
		return errors.Wrap(err, "failed to insert agent")
	}

	// Get the inserted agent ID
	agentID, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "failed to get agent ID")
	}

	// Output the created agent
	row := types.NewRow(
		types.MRP("id", agentID),
		types.MRP("slug", slug),
		types.MRP("name", s.Name),
		types.MRP("description", s.Description),
	)

	return gp.AddRow(ctx, row)
}

func NewCreateAgentCommand() (interface{}, error) {
	// Create the Glazed layer for output formatting
	glazedLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, err
	}

	// Create command description
	cmdDesc := cmds.NewCommandDescription(
		"create-agent",
		cmds.WithShort("Create a new agent"),
		cmds.WithLong(`
Create a new agent in the database.

Agents represent different AI assistants or workers that can be assigned to tasks
within a specific project.

Examples:
  # Create a new agent for a project
  create-agent --name="Code Analyzer" --description="Specialized in analyzing code patterns and architecture" --project=authentication-analysis
		`),
		// Define command flags
		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"name",
				parameters.ParameterTypeString,
				parameters.WithHelp("Name of the agent"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"description",
				parameters.ParameterTypeString,
				parameters.WithHelp("Description of the agent's capabilities"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"project",
				parameters.ParameterTypeString,
				parameters.WithHelp("Project slug or ID that this agent will work on"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"slug",
				parameters.ParameterTypeString,
				parameters.WithHelp("URL-friendly slug for the agent (auto-generated if not provided)"),
				parameters.WithDefault(""),
			),
		),
		// Add parameter layers
		cmds.WithLayersList(
			glazedLayer,
		),
	)

	return &CreateAgentCommand{
		CommandDescription: cmdDesc,
	}, nil
}

// ListAgentsCommand lists all agents
type ListAgentsCommand struct {
	*cmds.CommandDescription
}

// ListAgentsSettings holds the parameters for listing agents
type ListAgentsSettings struct {
	AgentID string `glazed.parameter:"agent"`
}

// Ensure interface implementation
var _ cmds.GlazeCommand = &ListAgentsCommand{}

func (c *ListAgentsCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	// Parse settings
	s := &ListAgentsSettings{}
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
		SELECT id, slug, name, description, current_project_id, current_task_id, created_at, updated_at
		FROM agents
		WHERE 1=1
	`
	args := []interface{}{}

	// Add filters
	if s.AgentID != "" {
		query += " AND (id = ? OR slug = ?)"
		args = append(args, s.AgentID, s.AgentID)
	}

	query += " ORDER BY created_at DESC"

	// Execute query
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to query agents")
	}
	defer rows.Close()

	// Process results
	for rows.Next() {
		var id int
		var slug, name, description, createdAt, updatedAt string
		var currentProjectID, currentTaskID sql.NullInt64

		err := rows.Scan(&id, &slug, &name, &description, &currentProjectID, &currentTaskID, &createdAt, &updatedAt)
		if err != nil {
			return errors.Wrap(err, "failed to scan agent row")
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

		// Add current work if assigned
		if currentProjectID.Valid {
			rowData.Set("current_project_id", currentProjectID.Int64)
		}
		if currentTaskID.Valid {
			rowData.Set("current_task_id", currentTaskID.Int64)
		}

		if err := gp.AddRow(ctx, rowData); err != nil {
			return errors.Wrap(err, "failed to add row")
		}
	}

	return rows.Err()
}

func NewListAgentsCommand() (interface{}, error) {
	// Create the Glazed layer for output formatting
	glazedLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, err
	}

	// Create command description
	cmdDesc := cmds.NewCommandDescription(
		"list-agents",
		cmds.WithShort("List agents"),
		cmds.WithLong(`
List all agents in the database.

Examples:
  # List all agents
  list-agents
  
  # Show specific agent by ID
  list-agents --agent=1
  
  # Show specific agent by slug
  list-agents --agent=code-analyzer
		`),
		// Define command flags
		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"agent",
				parameters.ParameterTypeString,
				parameters.WithHelp("Show specific agent by ID or slug"),
				parameters.WithDefault(""),
			),
		),
		// Add parameter layers
		cmds.WithLayersList(
			glazedLayer,
		),
	)

	return &ListAgentsCommand{
		CommandDescription: cmdDesc,
	}, nil
}
