package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/go-go-golems/glazed/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// InsertTaskCommand inserts a new task into the database
type InsertTaskCommand struct {
	*cmds.CommandDescription
}

// InsertTaskSettings holds the parameters for inserting a task
type InsertTaskSettings struct {
	Project      string   `glazed.parameter:"project"`
	Agent        string   `glazed.parameter:"agent"`
	Type         string   `glazed.parameter:"type"`
	Instructions string   `glazed.parameter:"instructions"`
	Slug         string   `glazed.parameter:"slug"`
	Dependencies []string `glazed.parameter:"dependencies"`
}

// Ensure interface implementation
var _ cmds.GlazeCommand = &InsertTaskCommand{}

func (c *InsertTaskCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	log.Debug().Msg("Starting insert-task command")
	
	// Create a timeout context for the entire operation
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	// Parse settings
	s := &InsertTaskSettings{}
	if err := parsedLayers.InitializeStruct(layers.DefaultSlug, s); err != nil {
		return errors.Wrap(err, "failed to parse settings")
	}
	log.Debug().Interface("settings", s).Msg("Parsed settings")

	// Initialize database
	log.Debug().Msg("Initializing database")
	db, err := InitDatabase()
	if err != nil {
		return errors.Wrap(err, "failed to initialize database")
	}
	defer db.Close()
	log.Debug().Msg("Database initialized successfully")

	// Resolve project ID
	log.Debug().Str("project", s.Project).Msg("Resolving project ID")
	projectID, err := ResolveProjectID(timeoutCtx, db, s.Project)
	if err != nil {
		return errors.Wrap(err, "failed to resolve project")
	}
	log.Debug().Int("project_id", projectID).Msg("Project ID resolved")

	// Resolve agent ID if specified
	var agentID int
	if s.Agent != "" {
		log.Debug().Str("agent", s.Agent).Msg("Resolving agent ID")
		agentID, err = ResolveAgentID(timeoutCtx, db, s.Agent)
		if err != nil {
			return errors.Wrap(err, "failed to resolve agent")
		}
		log.Debug().Int("agent_id", agentID).Msg("Agent ID resolved")
	}

	// Generate slug if not provided
	slug := s.Slug
	if slug == "" {
		slug = GenerateSlug(s.Instructions)
		log.Debug().Str("generated_slug", slug).Msg("Generated slug from instructions")
	}

	// Ensure slug is unique within the project
	log.Debug().Str("slug", slug).Int("project_id", projectID).Msg("Checking slug uniqueness")
	originalSlug := slug
	counter := 1
	for {
		var exists bool
		err := db.QueryRowContext(timeoutCtx, "SELECT EXISTS(SELECT 1 FROM tasks WHERE project_id = ? AND slug = ?)", projectID, slug).Scan(&exists)
		if err != nil {
			return errors.Wrap(err, "failed to check slug uniqueness")
		}
		
		if !exists {
			log.Debug().Str("final_slug", slug).Msg("Slug is unique")
			break
		}
		
		counter++
		slug = originalSlug + "-" + string(rune('0'+counter-1))
		log.Debug().Str("new_slug", slug).Int("counter", counter).Msg("Slug exists, trying new variant")
	}

	// Start transaction
	log.Debug().Msg("Starting database transaction")
	tx, err := db.BeginTx(timeoutCtx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer tx.Rollback()
	log.Debug().Msg("Transaction started successfully")

	// Insert the task
	log.Debug().Str("type", s.Type).Str("slug", slug).Msg("Inserting task")
	var result sql.Result
	if agentID > 0 {
		result, err = tx.ExecContext(timeoutCtx, `
			INSERT INTO tasks (slug, project_id, agent_id, type, instructions, status)
			VALUES (?, ?, ?, ?, ?, 'pending')
		`, slug, projectID, agentID, s.Type, s.Instructions)
	} else {
		result, err = tx.ExecContext(timeoutCtx, `
			INSERT INTO tasks (slug, project_id, type, instructions, status)
			VALUES (?, ?, ?, ?, 'pending')
		`, slug, projectID, s.Type, s.Instructions)
	}
	if err != nil {
		return errors.Wrap(err, "failed to insert task")
	}
	log.Debug().Msg("Task inserted successfully")

	// Get the inserted task ID
	taskID, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "failed to get task ID")
	}
	log.Debug().Int64("task_id", taskID).Msg("Got task ID")

	// Insert dependencies if any
	if len(s.Dependencies) > 0 {
		log.Debug().Interface("dependencies", s.Dependencies).Msg("Processing dependencies")
		for _, parentIdentifier := range s.Dependencies {
			log.Debug().Str("parent_identifier", parentIdentifier).Msg("Resolving dependency task ID")
			parentID, err := ResolveTaskID(timeoutCtx, tx, parentIdentifier)
			if err != nil {
				return errors.Wrapf(err, "failed to resolve dependency task: %s", parentIdentifier)
			}
			log.Debug().Int("parent_id", parentID).Msg("Dependency task ID resolved")
			
			_, err = tx.ExecContext(timeoutCtx, `
				INSERT INTO task_dependencies (task_id, parent_task_id)
				VALUES (?, ?)
			`, taskID, parentID)
			if err != nil {
				return errors.Wrapf(err, "failed to insert dependency on task %d", parentID)
			}
			log.Debug().Int64("task_id", taskID).Int("parent_id", parentID).Msg("Dependency inserted")
		}
	}

	// Commit transaction
	log.Debug().Msg("Committing transaction")
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}
	log.Debug().Msg("Transaction committed successfully")

	// Output the created task
	row := types.NewRow(
		types.MRP("id", taskID),
		types.MRP("slug", slug),
		types.MRP("project_id", projectID),
		types.MRP("agent_id", agentID),
		types.MRP("type", s.Type),
		types.MRP("instructions", s.Instructions),
		types.MRP("status", "pending"),
		types.MRP("dependencies", s.Dependencies),
	)

	return gp.AddRow(ctx, row)
}

func NewInsertTaskCommand() (interface{}, error) {
	// Create the Glazed layer for output formatting
	glazedLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, err
	}

	// Create command description
	cmdDesc := cmds.NewCommandDescription(
		"insert-task",
		cmds.WithShort("Insert a new task into the database"),
		cmds.WithLong(`
Insert a new task into the agent task database.

Tasks can be of two types:
- gather_information: Tasks for collecting source code information
- oracle_analysis: Tasks for analyzing collected information

Dependencies can be specified as a list of parent task IDs that must be completed
before this task can be started.

Examples:
  # Insert a simple gather task
  insert-task --project=auth-analysis --type=gather_information --instructions="Gather authentication patterns"
  
  # Insert an analysis task that depends on task 1
  insert-task --project=auth-analysis --type=oracle_analysis --instructions="Analyze auth patterns" --dependencies=1
  
  # Insert a task with agent assignment and dependencies using slugs
  insert-task --project=auth-analysis --agent=code-analyzer --type=oracle_analysis --instructions="Compare patterns" --dependencies=auth-analysis/gather-patterns
		`),
		// Define command flags
		cmds.WithFlags(
			parameters.NewParameterDefinition(
				"project",
				parameters.ParameterTypeString,
				parameters.WithHelp("Project ID or slug this task belongs to"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"agent",
				parameters.ParameterTypeString,
				parameters.WithHelp("Agent ID or slug assigned to this task (optional)"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"type",
				parameters.ParameterTypeChoice,
				parameters.WithHelp("Type of task to create"),
				parameters.WithChoices("gather_information", "oracle_analysis"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"instructions",
				parameters.ParameterTypeString,
				parameters.WithHelp("Instructions for the task"),
				parameters.WithRequired(true),
			),
			parameters.NewParameterDefinition(
				"slug",
				parameters.ParameterTypeString,
				parameters.WithHelp("URL-friendly slug for the task (auto-generated if not provided)"),
				parameters.WithDefault(""),
			),
			parameters.NewParameterDefinition(
				"dependencies",
				parameters.ParameterTypeStringList,
				parameters.WithHelp("List of parent task IDs or project_slug/task_slug this task depends on"),
				parameters.WithDefault([]string{}),
			),
		),
		// Add parameter layers
		cmds.WithLayersList(
			glazedLayer,
		),
	)

	return &InsertTaskCommand{
		CommandDescription: cmdDesc,
	}, nil
} 