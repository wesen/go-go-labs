# Agent Task SQLite

A command-line tool for managing agent tasks, locations, and reports in SQLite. This tool provides a convenient interface for working with the agent task database, allowing you to organize analysis work through projects, assign tasks to agents, and track progress.

## Features

- **Project Management**: Create and organize work into projects
- **Agent Management**: Register agents and assign them to tasks
- **Task Management**: Create, query, and track analysis tasks with dependencies
- **Location Tracking**: Store and manage code locations relevant to tasks
- **Report Generation**: Create reports linked to tasks and locations
- **Database Locking**: Safe concurrent access with automatic retry logic

## Installation

Build the tool from source:

```bash
cd go-go-labs
go build -o cmd/apps/agent-task-sqlite/agent-task-sqlite ./cmd/apps/agent-task-sqlite
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
- **Tasks**: `insert-task`, `query-tasks`
- **Locations**: `insert-locations`
- **Reports**: `create-report`

## Usage Examples

### 1. Setting Up a New Project

First, create a project to organize your work:

```bash
# Create a new project
./agent-task-sqlite create-project \
  --name="Authentication Analysis" \
  --description="Analyze authentication patterns and security implementations across the codebase"

# List all projects
./agent-task-sqlite list-projects
```

### 2. Registering Agents

Create agents that will work on tasks:

```bash
# Create a code analysis agent
./agent-task-sqlite create-agent \
  --name="Code Analyzer" \
  --description="Specialized in analyzing code patterns and architecture"

# Create an oracle analysis agent
./agent-task-sqlite create-agent \
  --name="Oracle Analyst" \
  --description="Performs deep analysis and generates insights from gathered data"

# List all agents
./agent-task-sqlite list-agents
```

### 3. Creating Tasks

Create tasks within your project:

```bash
# Create a gather task (assuming project ID 1)
./agent-task-sqlite insert-task \
  --project-id=1 \
  --type=gather_information \
  --instructions="Gather all authentication-related code patterns from the main application"

# Create an analysis task that depends on the gather task (assuming task ID 1)
./agent-task-sqlite insert-task \
  --project-id=1 \
  --agent-id=2 \
  --type=oracle_analysis \
  --instructions="Analyze gathered authentication patterns and identify security vulnerabilities" \
  --dependencies=1

# Query tasks for a specific project
./agent-task-sqlite query-tasks --project-id=1
```

### 4. Managing Code Locations

Store relevant code locations for tasks:

```bash
# Add individual location
./agent-task-sqlite insert-locations \
  --task-id=1 \
  --location="src/auth/login.go" \
  --description="Main login function with password validation"

# Add multiple locations at once
./agent-task-sqlite insert-locations \
  --task-id=1 \
  --locations="src/auth/middleware.go:Authentication middleware" \
  --locations="src/auth/jwt.go:JWT token handling" \
  --locations="src/auth/session.go:Session management"
```

### 5. Creating Reports

Generate reports for completed analysis:

```bash
# Create a report with direct content
./agent-task-sqlite create-report \
  --task-id=2 \
  --content="Authentication analysis complete. Found 3 potential security issues: 1) Weak password policy, 2) Missing rate limiting, 3) Insecure session storage."

# Create a report from a file
./agent-task-sqlite create-report \
  --task-id=2 \
  --content-file="analysis-report.md"

# Create a report linked to specific locations
./agent-task-sqlite create-report \
  --task-id=2 \
  --content="Security vulnerability found in authentication flow" \
  --location-ids=1,2,3
```

### 6. Querying and Filtering

The tool provides powerful querying capabilities:

```bash
# Show all pending tasks
./agent-task-sqlite query-tasks --status=pending

# Show tasks assigned to a specific agent
./agent-task-sqlite query-tasks --agent-id=1

# Show tasks for a project with dependencies
./agent-task-sqlite query-tasks --project-id=1 --show-deps

# Show last 10 completed tasks
./agent-task-sqlite query-tasks --status=completed --limit=10

# Show specific task details
./agent-task-sqlite query-tasks --task-id=5 --show-deps
```

## Output Formats

All commands support multiple output formats through Glazed:

```bash
# JSON output
./agent-task-sqlite query-tasks --output=json

# YAML output
./agent-task-sqlite list-projects --output=yaml

# CSV output
./agent-task-sqlite list-agents --output=csv

# Table output (default)
./agent-task-sqlite query-tasks --output=table

# Select specific fields
./agent-task-sqlite query-tasks --fields=id,type,status,instructions
```

## Workflow Example

Here's a complete workflow example:

```bash
# 1. Create a project
./agent-task-sqlite create-project \
  --name="Database Security Audit" \
  --description="Comprehensive security audit of database access patterns"

# 2. Register agents
./agent-task-sqlite create-agent \
  --name="Security Scanner" \
  --description="Automated security pattern detection"

# 3. Create initial gather task
./agent-task-sqlite insert-task \
  --project-id=1 \
  --type=gather_information \
  --instructions="Collect all database query patterns and access controls"

# 4. Add code locations
./agent-task-sqlite insert-locations \
  --task-id=1 \
  --locations="src/db/queries.go:Database query functions" \
  --locations="src/db/migrations/:Database schema migrations" \
  --locations="src/middleware/auth.go:Database access authorization"

# 5. Create analysis task
./agent-task-sqlite insert-task \
  --project-id=1 \
  --agent-id=1 \
  --type=oracle_analysis \
  --instructions="Analyze database access patterns for security vulnerabilities" \
  --dependencies=1

# 6. Generate report
./agent-task-sqlite create-report \
  --task-id=2 \
  --content="Database security audit complete. Found SQL injection vulnerabilities in user input handling." \
  --location-ids=1,2

# 7. Review results
./agent-task-sqlite query-tasks --project-id=1 --show-deps --output=json
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
./agent-task-sqlite --help

# Command-specific help
./agent-task-sqlite insert-task --help
./agent-task-sqlite query-tasks --help
./agent-task-sqlite create-report --help
```

## Environment Variables

- `AGENT_SQLITE_DB`: Path to SQLite database file (default: `/tmp/agent-work.db`)

## Contributing

This tool is part of the go-go-labs project. Feel free to submit issues and enhancement requests! 