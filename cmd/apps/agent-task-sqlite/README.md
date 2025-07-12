# Agent Task SQLite

A command-line tool for managing agent tasks, locations, and reports in SQLite. This tool provides a convenient interface for working with the agent task database, allowing you to organize analysis work through projects, assign tasks to agents, and track progress.

## Features

- **Project Management**: Create and organize work into projects with friendly slugs
- **Agent Management**: Register agents and assign them to tasks with slug-based identification
- **Task Management**: Create, query, and track analysis tasks with dependencies using flexible identifiers
- **Task Assignment**: Assign tasks to agents and automatically set them to in-progress
- **Location Tracking**: Store and manage code locations relevant to tasks
- **Report Generation**: Create reports linked to tasks and locations
- **Slug-based Identification**: Use human-friendly slugs instead of numeric IDs
- **Database Locking**: Safe concurrent access with automatic retry logic

## Installation

Build the tool from source:

```bash
cd go-go-labs
go build -o cmd/apptask-manager/agent-task-sqlite ./cmd/apps/agent-task-sqlite
```

## Database Configuration

By default, the tool uses `/tmp/agent-work.db` as the database file. You can override this by setting the `AGENT_SQLITE_DB` environment variable:

```bash
export AGENT_SQLITE_DB=/path/to/your/database.db
```

## Commands Overview

The tool provides several commands organized around the core entities:

- **Projects**: `create-project`, `list-projects`
- **Agents**: `create-agent`, `list-agents`
- **Tasks**: `insert-task`, `query-tasks`, `assign-task`
- **Locations**: `insert-locations`
- **Reports**: `create-report`

## Usage Examples

### 1. Setting Up a New Project

First, create a project to organize your work:

```bash
# Create a new project (slug auto-generated from name)
task-manager create-project \
  --name="Authentication Analysis" \
  --description="Analyze authentication patterns and security implementations across the codebase"

# Create a project with custom slug
task-manager create-project \
  --name="Database Security Audit" \
  --description="Comprehensive security audit of database access patterns" \
  --slug="db-security"

# List all projects
task-manager list-projects

# Show specific project by slug
task-manager list-projects --project=authentication-analysis
```

### 2. Registering Agents

Create agents that will work on tasks:

```bash
# Create a code analysis agent (slug auto-generated)
task-manager create-agent \
  --name="Code Analyzer" \
  --description="Specialized in analyzing code patterns and architecture"

# Create an oracle analysis agent with custom slug
task-manager create-agent \
  --name="Oracle Analyst" \
  --description="Performs deep analysis and generates insights from gathered data" \
  --slug="oracle"

# List all agents
task-manager list-agents

# Show specific agent by slug
task-manager list-agents --agent=code-analyzer
```

### 3. Creating Tasks

Create tasks within your project:

```bash
# Create a gather task using project slug
task-manager insert-task \
  --project=authentication-analysis \
  --type=gather_information \
  --instructions="Gather all authentication-related code patterns from the main application"

# Create an analysis task with agent assignment and dependencies
task-manager insert-task \
  --project=authentication-analysis \
  --agent=oracle \
  --type=oracle_analysis \
  --instructions="Analyze gathered authentication patterns and identify security vulnerabilities" \
  --dependencies=authentication-analysis/gather-all-authentication-related-code-patterns

# Assign a pending task to an agent (sets status to in_progress)
task-manager assign-task \
  --agent=code-analyzer \
  --task=authentication-analysis/gather-all-authentication-related-code-patterns

# Query tasks for a specific project
task-manager query-tasks --project=authentication-analysis

# Query tasks assigned to a specific agent
task-manager query-tasks --agent=code-analyzer
```

### 4. Managing Code Locations

Store relevant code locations for tasks:

```bash
# Add individual location
task-manager insert-locations \
  --task-id=1 \
  --location="src/auth/login.go" \
  --description="Main login function with password validation"

# Add multiple locations at once
task-manager insert-locations \
  --task-id=1 \
  --locations="src/auth/middleware.go:Authentication middleware" \
  --locations="src/auth/jwt.go:JWT token handling" \
  --locations="src/auth/session.go:Session management"
```

### 5. Creating Reports

Generate reports for completed analysis:

```bash
# Create a report with direct content
task-manager create-report \
  --task-id=2 \
  --content="Authentication analysis complete. Found 3 potential security issues: 1) Weak password policy, 2) Missing rate limiting, 3) Insecure session storage."

# Create a report from a file
task-manager create-report \
  --task-id=2 \
  --content-file="analysis-report.md"

# Create a report linked to specific locations
task-manager create-report \
  --task-id=2 \
  --content="Security vulnerability found in authentication flow" \
  --location-ids=1,2,3
```

### 6. Querying and Filtering

The tool provides powerful querying capabilities:

```bash
# Show all pending tasks
task-manager query-tasks --status=pending

# Show tasks assigned to a specific agent
task-manager query-tasks --agent=code-analyzer

# Show tasks for a project with dependencies
task-manager query-tasks --project=authentication-analysis --show-deps

# Show last 10 completed tasks
task-manager query-tasks --status=completed --limit=10

# Show specific task details by project/task slug
task-manager query-tasks --task=authentication-analysis/analyze-patterns --show-deps
```

## Output Formats

All commands support multiple output formats through Glazed:

```bash
# JSON output
task-manager query-tasks --output=json

# YAML output
task-manager list-projects --output=yaml

# CSV output
task-manager list-agents --output=csv

# Table output (default)
task-manager query-tasks --output=table

# Select specific fields
task-manager query-tasks --fields=id,type,status,instructions
```

## Workflow Example

Here's a complete workflow example:

```bash
# 1. Create a project
task-manager create-project \
  --name="Database Security Audit" \
  --description="Comprehensive security audit of database access patterns" \
  --slug="db-security"

# 2. Register agents
task-manager create-agent \
  --name="Security Scanner" \
  --description="Automated security pattern detection" \
  --slug="security-scanner"

# 3. Create initial gather task
task-manager insert-task \
  --project=db-security \
  --type=gather_information \
  --instructions="Collect all database query patterns and access controls" \
  --slug="gather-db-patterns"

# 4. Assign task to agent
task-manager assign-task \
  --agent=security-scanner \
  --task=db-security/gather-db-patterns

# 5. Add code locations
task-manager insert-locations \
  --task=db-security/gather-db-patterns \
  --locations="src/db/queries.go:Database query functions" \
  --locations="src/db/migrations/:Database schema migrations" \
  --locations="src/middleware/auth.go:Database access authorization"

# 6. Create analysis task
task-manager insert-task \
  --project=db-security \
  --type=oracle_analysis \
  --instructions="Analyze database access patterns for security vulnerabilities" \
  --dependencies=db-security/gather-db-patterns \
  --slug="analyze-security"

# 7. Generate report
task-manager create-report \
  --task=db-security/analyze-security \
  --content="Database security audit complete. Found SQL injection vulnerabilities in user input handling." \
  --location-ids=1,2

# 8. Review results
task-manager query-tasks --project=db-security --show-deps --output=json
```

## Database Schema

The tool uses the following database schema:

- **projects**: Store project information
- **agents**: Store agent descriptions and capabilities
- **tasks**: Store tasks with project/agent assignments and dependencies
- **task_dependencies**: Track task dependency relationships
- **gathered_locations**: Store code locations relevant to tasks
- **reports**: Store analysis reports
- **report_locations**: Link reports to specific locations
- **agent_steps**: Track agent execution steps (for future use)

## Advanced Features

### Database Locking

The tool handles database locking automatically with:
- 30-second busy timeout
- Automatic retry logic
- Single connection pool for SQLite compatibility

### Task Dependencies

Tasks can depend on other tasks, creating a dependency graph:
- Tasks with dependencies cannot be started until prerequisites are complete
- Use `--show-deps` flag to visualize dependencies
- Supports complex dependency chains
- Dependencies can be specified by ID or project_slug/task_slug format

### Task Assignment

Tasks can be assigned to agents, automatically setting them to in_progress:
- Only pending tasks can be assigned
- Agents can only work on one task at a time
- Assignment automatically updates agent's current work tracking
- Use `assign-task` command to assign tasks to agents

### Slug-based Identification

All entities support human-friendly slugs:
- **Projects**: Use project slug instead of numeric ID
- **Agents**: Use agent slug for easy identification
- **Tasks**: Use project_slug/task_slug format for clear references
- Slugs are auto-generated from names but can be customized

### Flexible Location Storage

Code locations are stored as flexible text references:
- File paths: `src/auth/login.go`
- Function names: `AuthenticateUser`
- Line ranges: `main.go:10-20`
- Any reference format that makes sense for your workflow

## Help and Documentation

Each command provides detailed help:

```bash
# General help
task-manager --help

# Command-specific help
task-manager insert-task --help
task-manager query-tasks --help
task-manager assign-task --help
task-manager create-report --help
```

## Environment Variables

- `AGENT_SQLITE_DB`: Path to SQLite database file (default: `/tmp/agent-work.db`)

## Contributing

This tool is part of the go-go-labs project. Feel free to submit issues and enhancement requests! 