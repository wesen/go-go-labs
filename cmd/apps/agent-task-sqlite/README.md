# Agent Task SQLite

A command-line tool for managing agent tasks, locations, and reports in SQLite. This tool provides a convenient interface for working with the agent task database, allowing you to organize analysis work through projects, assign tasks to agents, and track progress with dependency management.

## Features

- **Project Management**: Create and organize work into projects with friendly slugs and guidelines
- **Project Guidelines**: Store concise and detailed guidelines displayed during task operations
- **Agent Management**: Register agents tied to specific projects with slug-based identification
- **Task Management**: Create, query, and track analysis tasks with dependencies using flexible identifiers
- **Task Notes**: Record observations, bugs, lessons learned, and ideas during task execution
- **Dependency Blocking**: Tasks with incomplete dependencies cannot be assigned until prerequisites are completed
- **Task Assignment**: Assign tasks to agents with validation and force reassignment options
- **Project-Agent Binding**: Agents are tied to specific projects and cannot work on tasks from other projects
- **Location Tracking**: Store and manage code locations relevant to tasks
- **Report Generation**: Create comprehensive completion reports linked to tasks and locations
- **Web Interface**: Monitor and manage tasks through a web dashboard with quick actions
- **Slug-based Identification**: Use human-friendly slugs instead of numeric IDs
- **Database Locking**: Safe concurrent access with automatic retry logic
- **Debug Logging**: Comprehensive debug logging for troubleshooting

## Installation

Build the tool from source:

```bash
cd go-go-labs
go build -o task-manager ./cmd/apps/agent-task-sqlite
```

## Database Configuration

By default, the tool uses `/tmp/agent-work.db` as the database file. You can override this by setting the `AGENT_SQLITE_DB` environment variable:

```bash
export AGENT_SQLITE_DB=/path/to/your/database.db
```

## Commands Overview

The tool provides several commands organized around the core entities:

- **Projects**: `create-project`, `list-projects`, `get-guidelines`
- **Agents**: `create-agent`, `list-agents`
- **Tasks**: `insert-task`, `query-tasks`, `assign-task`, `complete-task`
- **Notes**: `take-note`, `get-notes`
- **Locations**: `insert-locations`
- **Reports**: `create-report`, `write-completion-report`, `get-report`
- **Web Interface**: `serve`

## Usage Examples

### 1. Setting Up a New Project

First, create a project to organize your work:

```bash
# Create a new project (slug auto-generated from name)
task-manager create-project \
  --name="Authentication Analysis" \
  --description="Analyze authentication patterns and security implementations across the codebase"

# Create a project with custom slug and guidelines
task-manager create-project \
  --name="Database Security Audit" \
  --description="Comprehensive security audit of database access patterns" \
  --slug="db-security" \
  --concise-guidelines="Focus on SQL injection and access control vulnerabilities" \
  --full-guidelines="Use static analysis tools, review all database queries, check user permissions, document findings in detail"

# List all projects
task-manager list-projects

# Show specific project by slug
task-manager list-projects --project=authentication-analysis

# Display project guidelines
task-manager get-guidelines --project=db-security

# Display only concise guidelines
task-manager get-guidelines --project=db-security --concise
```

### 2. Registering Agents (Project-Specific)

Create agents that will work on tasks within a specific project:

```bash
# Create a code analysis agent for the authentication project
task-manager create-agent \
  --name="Code Analyzer" \
  --description="Specialized in analyzing code patterns and architecture" \
  --project=authentication-analysis

# Create an oracle analysis agent for the authentication project with custom slug
task-manager create-agent \
  --name="Oracle Analyst" \
  --description="Performs deep analysis and generates insights from gathered data" \
  --project=authentication-analysis \
  --slug="oracle"

# Create a database expert for the database security project  
task-manager create-agent \
  --name="Database Expert" \
  --description="Database analysis specialist" \
  --project=db-security \
  --slug="db-expert"

# List all agents (shows project assignments)
task-manager list-agents
```

### 3. Creating Tasks with Dependencies

Create tasks within your project and set up dependency relationships:

```bash
# Create a gather task using project slug
task-manager insert-task \
  --project=authentication-analysis \
  --type=gather_information \
  --instructions="Gather all authentication-related code patterns from the main application"

# Create an analysis task with dependencies
task-manager insert-task \
  --project=authentication-analysis \
  --type=oracle_analysis \
  --instructions="Analyze gathered authentication patterns and identify security vulnerabilities" \
  --dependencies=1

# Query tasks for a specific project
task-manager query-tasks --project-id=1
```

### 4. Task Assignment with Validation

Assign tasks to agents with automatic validation:

```bash
# Assign a pending task to an agent (sets status to in_progress)
task-manager assign-task \
  --agent=code-analyzer \
  --task=1

# Try to assign task with incomplete dependencies (will fail)
task-manager assign-task \
  --agent=oracle \
  --task=2

# Expected error: Cannot assign task with incomplete dependencies

# Complete the dependency task using the complete-task command
task-manager complete-task \
  --task=1 \
  --notes="Gathered authentication patterns from main application"

# Now assign the analysis task (will succeed)
task-manager assign-task \
  --agent=oracle \
  --task=2

# Force reassign a task from one agent to another
task-manager assign-task \
  --agent=other-agent \
  --task=2 \
  --force
```

### 5. Task Completion

Complete tasks when work is finished:

```bash
# Complete a task with notes
task-manager complete-task \
  --task=1 \
  --notes="Authentication patterns gathered. Found 15 code files with auth logic."

# Complete a task using project/task slug
task-manager complete-task \
  --task=authentication-analysis/gather-patterns \
  --notes="Comprehensive analysis complete. 3 security vulnerabilities identified."

# Complete a task without notes
task-manager complete-task \
  --task=2
```

### 6. Recording Task Notes

Record observations, bugs, and lessons learned during task execution:

```bash
# Take a general note
task-manager take-note \
  --task=5 \
  --content="Found the authentication logic in auth.go"

# Record a bug
task-manager take-note \
  --task=authentication-analysis/gather-patterns \
  --type=bugs \
  --content="Authentication middleware crashes with nil pointer on line 42"

# Record observations
task-manager take-note \
  --project-id=1 \
  --task-slug=gather-patterns \
  --type=observations \
  --content="Most auth code is in middleware package"

# Record lessons learned
task-manager take-note \
  --task=5 \
  --type=lessons_learned \
  --content="Always check for null pointers in middleware chains"

# Get all notes for a task
task-manager get-notes --task=5

# Get only bug reports
task-manager get-notes --task=authentication-analysis/gather-patterns --type=bugs
```

### 7. Project-Agent Validation

The system enforces that agents can only work on tasks from their assigned project:

```bash
# This will fail - agent from project 1 cannot work on project 2 task
task-manager assign-task \
  --agent=code-analyzer \
  --task=db-security/gather-patterns

# Expected error: Agent belongs to different project
```

### 8. Managing Code Locations

Store relevant code locations for tasks:

```bash
# Add individual location using task ID
task-manager insert-locations \
  --task-id=1 \
  --location="src/auth/login.go" \
  --description="Main login function with password validation"

# Add multiple locations at once using task slug
task-manager insert-locations \
  --task=authentication-analysis/gather-patterns \
  --locations="src/auth/middleware.go:Authentication middleware" \
  --locations="src/auth/jwt.go:JWT token handling" \
  --locations="src/auth/session.go:Session management"
```

### 9. Creating Reports

Generate reports for completed analysis:

```bash
# Create a report with direct content
task-manager create-report \
  --task-id=2 \
  --content="Authentication analysis complete. Found 3 potential security issues: 1) Weak password policy, 2) Missing rate limiting, 3) Insecure session storage."

# Create a report from a file
task-manager create-report \
  --task=authentication-analysis/analysis-task \
  --content-file="analysis-report.md"

# Create a report linked to specific locations
task-manager create-report \
  --task-id=2 \
  --content="Security vulnerability found in authentication flow" \
  --location-ids=1,2,3

# Write a completion report for a finished task
task-manager write-completion-report \
  --task=5 \
  --content="Analysis complete. Found 3 security issues in authentication middleware."

# Write a completion report from file
task-manager write-completion-report \
  --task=authentication-analysis/analysis-task \
  --report-file="./analysis-report.md"

# Get all reports for a task
task-manager get-report --task=5

# Get reports with linked locations
task-manager get-report --task=authentication-analysis/analysis-task --with-locations
```

### 10. Querying and Filtering

The tool provides powerful querying capabilities:

```bash
# Show all pending tasks
task-manager query-tasks --status=pending

# Show tasks assigned to a specific agent
task-manager query-tasks --agent-id=1

# Show tasks for a project
task-manager query-tasks --project-id=1

# Show tasks with dependencies (may timeout due to known issue)
task-manager query-tasks --show-deps

# Show last 10 completed tasks
task-manager query-tasks --status=completed --limit=10

# Show specific task details
task-manager query-tasks --task-id=5

# Show tasks with notes included
task-manager query-tasks --with-notes

# Show tasks with reports included
task-manager query-tasks --with-reports

# Show tasks with both notes and reports
task-manager query-tasks --with-notes --with-reports
```

### 11. Web Interface

Start the web dashboard to monitor and manage tasks:

```bash
# Start web server on default port 8080
task-manager serve

# Start on custom port
task-manager serve --port=8081
```

The web interface provides:
- Task overview with status indicators
- Project reports and statistics
- Task detail pages with notes and reports
- Quick action buttons for taking notes and writing reports
- Real-time auto-refresh functionality

Access the dashboard at: http://localhost:8080

### 12. Debug Logging

Use debug logging to troubleshoot issues:

```bash
# Enable debug logging for assignment operations
task-manager --log-level=debug assign-task \
  --agent=code-analyzer \
  --task=1

# Debug output shows:
# - Dependency checking process
# - Project validation
# - Agent availability checks
# - Transaction operations
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

## Complete Workflow Example

Here's a complete workflow example demonstrating all features:

```bash
# 1. Create a project with guidelines
task-manager create-project \
  --name="API Security Review" \
  --description="Review API endpoints for security vulnerabilities" \
  --slug="api-security" \
  --concise-guidelines="Focus on authentication, authorization, and input validation" \
  --full-guidelines="Use OWASP Top 10 as reference, test all endpoints for common vulnerabilities, document findings with PoC examples"

# 2. Register project-specific agents
task-manager create-agent \
  --name="Code Gatherer" \
  --description="Collects API endpoint information" \
  --project=api-security \
  --slug="gatherer"

task-manager create-agent \
  --name="Security Analyst" \
  --description="Analyzes API security patterns" \
  --project=api-security \
  --slug="analyst"

# 3. Create initial gather task
task-manager insert-task \
  --project=api-security \
  --type=gather_information \
  --instructions="Collect all API endpoint definitions and security configurations" \
  --slug="gather-apis"

# 4. Assign gather task
task-manager assign-task \
  --agent=gatherer \
  --task=api-security/gather-apis

# 5. Add code locations and take notes during gathering
task-manager insert-locations \
  --task=api-security/gather-apis \
  --locations="src/api/routes.go:API route definitions" \
  --locations="src/middleware/security.go:Security middleware"

# Take notes during gathering
task-manager take-note \
  --task=api-security/gather-apis \
  --type=observations \
  --content="Found 25 API endpoints, most use JWT authentication"

task-manager take-note \
  --task=api-security/gather-apis \
  --type=bugs \
  --content="POST /admin/users endpoint missing authentication check"

# 6. Create analysis task with dependency
task-manager insert-task \
  --project=api-security \
  --type=oracle_analysis \
  --instructions="Analyze API security patterns and identify vulnerabilities" \
  --dependencies=api-security/gather-apis \
  --slug="analyze-api-security"

# 7. Try to assign analysis task (will fail due to incomplete dependency)
task-manager assign-task \
  --agent=analyst \
  --task=api-security/analyze-api-security

# Expected error: Dependency not completed

# 8. Complete the gather task
task-manager complete-task \
  --task=api-security/gather-apis \
  --notes="API gathering complete. Found 25 endpoints, documented security configurations"

# 9. Now assign analysis task (will succeed)
task-manager assign-task \
  --agent=analyst \
  --task=api-security/analyze-api-security

# 10. Take notes during analysis
task-manager take-note \
  --task=api-security/analyze-api-security \
  --type=lessons_learned \
  --content="Input validation should be implemented at both API gateway and application level"

# 11. Complete the analysis task
task-manager complete-task \
  --task=api-security/analyze-api-security \
  --notes="API security analysis complete. Found 3 endpoints without authentication, 2 potential SQL injection points"

# 12. Write comprehensive completion report
task-manager write-completion-report \
  --task=api-security/analyze-api-security \
  --content="Comprehensive API security analysis complete. Found 3 endpoints without authentication, 2 potential SQL injection points. Recommendations: implement OAuth2, use parameterized queries." \
  --location-ids=1,2

# 13. Review complete workflow with notes and reports
task-manager query-tasks --project-id=1 --with-notes --with-reports --output=json

# 14. Start web interface to view results
task-manager serve --port=8080
```

## Database Schema

The tool uses the following database schema:

- **projects**: Store project information with unique slugs and guidelines
- **agents**: Store agent descriptions tied to specific projects
- **tasks**: Store tasks with project/agent assignments and dependencies  
- **task_dependencies**: Track task dependency relationships
- **task_notes**: Store notes, observations, bugs, and lessons learned during task execution
- **gathered_locations**: Store code locations relevant to tasks
- **reports**: Store analysis reports
- **report_locations**: Link reports to specific locations
- **agent_steps**: Track agent execution steps (for future use)

## Advanced Features

### Project Guidelines System

Projects can include both concise and detailed guidelines:
- **Concise Guidelines**: Short instructions displayed when creating/completing tasks
- **Full Guidelines**: Detailed working instructions accessible via `get-guidelines` command
- Guidelines help maintain consistency across project work
- Automatically displayed during task lifecycle operations

### Task Notes System

Record various types of notes during task execution:
- **observations**: General observations during task work
- **lessons_learned**: Important lessons to remember for future tasks
- **notes**: General notes and comments
- **bugs**: Bug reports and issues encountered
- **ideas**: Ideas for future improvements
- **issues**: Issues that need to be addressed

Notes are color-coded in the web interface and can be filtered by type.

### Dependency Management

Tasks can depend on other tasks, creating a dependency graph:
- Tasks with incomplete dependencies cannot be assigned until prerequisites are completed
- Dependencies are validated during assignment with clear error messages
- Supports both numeric IDs and project_slug/task_slug format for dependencies
- Debug logging shows dependency checking process

### Task Assignment and Completion

Tasks can be assigned to agents with comprehensive validation:
- Only agents from the same project can be assigned to tasks
- Agents can only work on one task at a time
- Tasks with incomplete dependencies cannot be assigned
- Assignment automatically updates task status to in_progress
- Use `--force` flag to reassign tasks with proper cleanup

Tasks can be completed when work is finished:
- Only tasks in 'in_progress' status can be completed
- Completion automatically updates task status to 'completed'
- Agent assignment is automatically cleared upon completion
- Optional completion notes can be added to record results

### Project-Agent Binding

All agents are tied to specific projects:
- Agents must be created with a --project parameter
- Agents cannot be assigned to tasks from different projects
- Project validation is enforced even with --force reassignment
- Clear error messages for project mismatches

### Database Locking

The tool handles database locking automatically with:
- 30-second busy timeout
- Automatic retry logic  
- Single connection pool for SQLite compatibility
- Transaction safety with rollback protection

### Slug-based Identification

All entities support human-friendly slugs:
- **Projects**: Use project slug instead of numeric ID
- **Agents**: Use agent slug for easy identification  
- **Tasks**: Use project_slug/task_slug format for clear references
- Slugs are auto-generated from names but can be customized

### Debug Logging

Comprehensive debug logging available:
- Use --log-level=debug for detailed operation logs
- Shows dependency validation process
- Traces assignment operations step-by-step
- Logs project validation and agent availability checks
- Helpful for troubleshooting assignment failures

## Help and Documentation

Each command provides detailed help:

```bash
# General help
task-manager --help

# Command-specific help
task-manager insert-task --help
task-manager query-tasks --help
task-manager assign-task --help
task-manager complete-task --help
task-manager create-agent --help
task-manager create-report --help
task-manager take-note --help
task-manager write-completion-report --help
task-manager get-notes --help
task-manager get-report --help
task-manager get-guidelines --help
task-manager serve --help
```

## Environment Variables

- `AGENT_SQLITE_DB`: Path to SQLite database file (default: `/tmp/agent-work.db`)

## Known Issues

1. **--show-deps timeout**: The `--show-deps` flag may cause timeouts due to transaction handling issues
2. **Manual status updates**: Some test scenarios require manual status updates via SQLite commands

## Troubleshooting

### Common Error Messages

1. **"cannot assign task: dependency task X is not completed"**
   - Solution: Complete the dependency task first using `complete-task` command or check task status

2. **"agent belongs to project X but task belongs to project Y"**  
   - Solution: Use an agent from the correct project or create a new agent

3. **"task is already assigned to agent X (use --force to reassign)"**
   - Solution: Use --force flag to reassign or choose a different agent

4. **"agent is already assigned to task X"**
   - Solution: Complete the current task using `complete-task` or use a different agent

5. **"task cannot be completed (current status: pending). Task must be in_progress to be completed"**
   - Solution: Assign the task to an agent first using `assign-task` command

### Debug Commands

```bash
# Check task dependencies
task-manager --log-level=debug assign-task --agent=agent-slug --task=task-id

# Check database state
sqlite3 /tmp/agent-work.db "SELECT * FROM tasks WHERE id = X;"
sqlite3 /tmp/agent-work.db "SELECT * FROM task_dependencies WHERE task_id = X;"

# Verify agent-project assignments  
task-manager list-agents
```

## Contributing

This tool is part of the go-go-labs project. Feel free to submit issues and enhancement requests!

Key areas for contribution:
- Fix --show-deps timeout issue
- Implement agent step tracking
- Add more sophisticated dependency visualization
- Add task pause/resume functionality
