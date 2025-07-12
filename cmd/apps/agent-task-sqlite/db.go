package main

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

// GetDatabasePath returns the database path from environment or default
func GetDatabasePath() string {
	if path := os.Getenv("AGENT_SQLITE_DB"); path != "" {
		return path
	}
	return "/tmp/agent-work.db"
}

// InitDatabase creates and initializes the database with schema
func InitDatabase() (*sql.DB, error) {
	dbPath := GetDatabasePath()
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, errors.Wrap(err, "failed to create database directory")
	}

	// Open database connection with busy timeout for locking
	db, err := sql.Open("sqlite3", dbPath+"?_busy_timeout=30000")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}

	// Set connection pool settings
	db.SetMaxOpenConns(1) // SQLite works best with single connection
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	// Test connection with retry logic
	if err := pingWithRetry(db, 10, time.Second); err != nil {
		db.Close()
		return nil, errors.Wrap(err, "failed to ping database")
	}

	// Create schema
	if err := createSchema(db); err != nil {
		db.Close()
		return nil, errors.Wrap(err, "failed to create schema")
	}

	return db, nil
}

// pingWithRetry attempts to ping the database with retry logic
func pingWithRetry(db *sql.DB, maxRetries int, delay time.Duration) error {
	for i := 0; i < maxRetries; i++ {
		if err := db.Ping(); err == nil {
			return nil
		}
		if i < maxRetries-1 {
			time.Sleep(delay)
		}
	}
	return errors.New("database ping failed after retries")
}

// createSchema creates all the required tables
func createSchema(db *sql.DB) error {
	schema := `
-- Table for big picture projects
CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    slug TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table for agent descriptions
CREATE TABLE IF NOT EXISTS agents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    slug TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    current_project_id INTEGER,
    current_task_id INTEGER,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (current_project_id) REFERENCES projects(id),
    FOREIGN KEY (current_task_id) REFERENCES tasks(id)
);

-- Tables for managing the dependency graph of analysis tasks
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    slug TEXT NOT NULL,
    project_id INTEGER NOT NULL,
    agent_id INTEGER,
    type TEXT NOT NULL CHECK(type IN ('gather_information', 'oracle_analysis')),
    status TEXT NOT NULL CHECK(status IN ('pending', 'in_progress', 'completed', 'failed')) DEFAULT 'pending',
    instructions TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at DATETIME,
    completed_at DATETIME,
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (agent_id) REFERENCES agents(id),
    UNIQUE(project_id, slug)
);

CREATE TABLE IF NOT EXISTS task_dependencies (
    task_id INTEGER NOT NULL,
    parent_task_id INTEGER NOT NULL,
    PRIMARY KEY (task_id, parent_task_id),
    FOREIGN KEY (task_id) REFERENCES tasks(id),
    FOREIGN KEY (parent_task_id) REFERENCES tasks(id)
);

-- Table for tracking agent execution steps
CREATE TABLE IF NOT EXISTS agent_steps (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id INTEGER NOT NULL,
    step_type TEXT NOT NULL CHECK(step_type IN ('gather', 'analyze', 'report')),
    details TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks(id)
);

-- Table for storing code locations and their context (simplified)
CREATE TABLE IF NOT EXISTS gathered_locations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id INTEGER NOT NULL,
    location TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks(id)
);

-- Table for analysis reports
CREATE TABLE IF NOT EXISTS reports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks(id)
);

-- Junction table for report-location relationships
CREATE TABLE IF NOT EXISTS report_locations (
    report_id INTEGER NOT NULL,
    location_id INTEGER NOT NULL,
    PRIMARY KEY (report_id, location_id),
    FOREIGN KEY (report_id) REFERENCES reports(id),
    FOREIGN KEY (location_id) REFERENCES gathered_locations(id)
);

-- Indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_project ON tasks(project_id);
CREATE INDEX IF NOT EXISTS idx_tasks_agent ON tasks(agent_id);
CREATE INDEX IF NOT EXISTS idx_tasks_slug ON tasks(slug);
CREATE INDEX IF NOT EXISTS idx_projects_slug ON projects(slug);
CREATE INDEX IF NOT EXISTS idx_agents_slug ON agents(slug);
CREATE INDEX IF NOT EXISTS idx_agents_current_task ON agents(current_task_id);
CREATE INDEX IF NOT EXISTS idx_steps_task ON agent_steps(task_id);
CREATE INDEX IF NOT EXISTS idx_locations_task ON gathered_locations(task_id);
CREATE INDEX IF NOT EXISTS idx_reports_task ON reports(task_id);
`

	_, err := db.Exec(schema)
	return err
}

// GenerateSlug creates a URL-friendly slug from a string
func GenerateSlug(text string) string {
	// Convert to lowercase
	slug := strings.ToLower(text)
	
	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")
	
	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")
	
	// Limit length
	if len(slug) > 50 {
		slug = slug[:50]
		slug = strings.Trim(slug, "-")
	}
	
	return slug
}

// EnsureUniqueSlug ensures a slug is unique by appending a counter if needed
func EnsureUniqueSlug(db *sql.DB, table, slug string) (string, error) {
	originalSlug := slug
	counter := 1
	
	for {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM "+table+" WHERE slug = ?)", slug).Scan(&exists)
		if err != nil {
			return "", err
		}
		
		if !exists {
			return slug, nil
		}
		
		counter++
		slug = originalSlug + "-" + string(rune('0'+counter-1))
	}
}

// ResolveProjectID resolves a project identifier (ID or slug) to an ID
func ResolveProjectID(db *sql.DB, identifier string) (int, error) {
	// Try to parse as integer first
	var id int
	err := db.QueryRow("SELECT id FROM projects WHERE id = ? OR slug = ?", identifier, identifier).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.Errorf("project not found: %s", identifier)
		}
		return 0, err
	}
	return id, nil
}

// ResolveAgentID resolves an agent identifier (ID or slug) to an ID
func ResolveAgentID(db *sql.DB, identifier string) (int, error) {
	// Try to parse as integer first
	var id int
	err := db.QueryRow("SELECT id FROM agents WHERE id = ? OR slug = ?", identifier, identifier).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.Errorf("agent not found: %s", identifier)
		}
		return 0, err
	}
	return id, nil
}

// ResolveTaskID resolves a task identifier (ID or project_slug/task_slug) to an ID
func ResolveTaskID(ctx context.Context, db *sql.DB, identifier string) (int, error) {
	// Check if it contains a slash (project/task format)
	if strings.Contains(identifier, "/") {
		parts := strings.SplitN(identifier, "/", 2)
		if len(parts) != 2 {
			return 0, errors.New("invalid task identifier format, use project_slug/task_slug")
		}
		
		projectSlug := parts[0]
		taskSlug := parts[1]
		
		var id int
		err := db.QueryRowContext(ctx, `
			SELECT t.id FROM tasks t
			JOIN projects p ON t.project_id = p.id
			WHERE p.slug = ? AND t.slug = ?
		`, projectSlug, taskSlug).Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				return 0, errors.Errorf("task not found: %s", identifier)
			}
			return 0, err
		}
		return id, nil
	}
	
	// Try to parse as integer
	var id int
	err := db.QueryRowContext(ctx, "SELECT id FROM tasks WHERE id = ?", identifier).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.Errorf("task not found: %s", identifier)
		}
		return 0, err
	}
	return id, nil
} 