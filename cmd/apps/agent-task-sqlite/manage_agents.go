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

// CreateAgentCommand creates a new agent
type CreateAgentCommand struct {
	*cmds.CommandDescription
}

// CreateAgentSettings holds the parameters for creating an agent
type CreateAgentSettings struct {
	Name        string `glazed.parameter:"name"`
	Description string `glazed.parameter:"description"`
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

	// Insert the agent
	result, err := db.ExecContext(ctx, `
		INSERT INTO agents (name, description)
		VALUES (?, ?)
	`, s.Name, s.Description)
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

Agents represent different AI assistants or workers that can be assigned to tasks.

Examples:
  # Create a new agent
  create-agent --name="Code Analyzer" --description="Specialized in analyzing code patterns and architecture"
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
	AgentID int `glazed.parameter:"agent-id"`
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
		SELECT id, name, description, created_at, updated_at
		FROM agents
		WHERE 1=1
	`
	args := []interface{}{}

	// Add filters
	if s.AgentID > 0 {
		query += " AND id = ?"
		args = append(args, s.AgentID)
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
		var name, description, createdAt, updatedAt string

		err := rows.Scan(&id, &name, &description, &createdAt, &updatedAt)
		if err != nil {
			return errors.Wrap(err, "failed to scan agent row")
		}

		// Create row data
		rowData := types.NewRow(
			types.MRP("id", id),
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
  
  # Show specific agent
  list-agents --agent-id=1
		`),
		// Define command flags
		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"agent-id",
				parameters.ParameterTypeInteger,
				parameters.WithHelp("Show specific agent by ID"),
				parameters.WithDefault(0),
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