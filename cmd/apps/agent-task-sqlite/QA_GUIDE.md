# QA Guide: Agent Task SQLite Testing

This guide provides step-by-step testing scenarios to validate the complete agent task management system. Follow these scenarios to ensure all features work correctly.

## Prerequisites

1. Build the task-manager:
   ```bash
   cd go-go-labs
   go build -o task-manager ./cmd/apps/agent-task-sqlite
   ```

2. Set up test database (optional):
   ```bash
   export AGENT_SQLITE_DB=/tmp/test-agent-work.db
   # Or use default: /tmp/agent-work.db
   ```

3. Clean slate (if needed):
   ```bash
   rm -f /tmp/agent-work.db /tmp/test-agent-work.db
   ```

## Test Scenario 1: Basic Project and Agent Setup

### Step 1.1: Create Projects
```bash
# Create first project with auto-generated slug
task-manager create-project \
  --name="Authentication Analysis" \
  --description="Analyze authentication patterns and security implementations"

# Expected output: Shows project with id=1, slug="authentication-analysis"
```

```bash
# Create second project with custom slug
task-manager create-project \
  --name="Database Security Audit" \
  --description="Comprehensive security audit of database access patterns" \
  --slug="db-security"

# Expected output: Shows project with id=2, slug="db-security"
```

### Step 1.2: List Projects
```bash
# List all projects
task-manager list-projects

# Expected output: Shows both projects with their slugs
```

```bash
# Show specific project by slug
task-manager list-projects --project=db-security

# Expected output: Shows only the db-security project
```

### Step 1.3: Create Agents (Now Project-Specific)
```bash
# Create code analyzer agent for authentication project
task-manager create-agent \
  --name="Code Analyzer" \
  --description="Specialized in analyzing code patterns and architecture" \
  --project=authentication-analysis

# Expected output: Shows agent with id=1, slug="code-analyzer", project_id=1
```

```bash
# Create oracle agent for authentication project with custom slug
task-manager create-agent \
  --name="Oracle Analyst" \
  --description="Performs deep analysis and generates insights from gathered data" \
  --project=authentication-analysis \
  --slug="oracle"

# Expected output: Shows agent with id=2, slug="oracle", project_id=1
```

```bash
# Create database expert for db-security project
task-manager create-agent \
  --name="Database Expert" \
  --description="Database analysis specialist" \
  --project=db-security \
  --slug="db-expert"

# Expected output: Shows agent with id=3, slug="db-expert", project_id=2
```

### Step 1.4: List Agents
```bash
# List all agents
task-manager list-agents

# Expected output: Shows all agents with their current_project_id assignments
```

## Test Scenario 2: Task Creation and Dependencies

### Step 2.1: Create Initial Gather Tasks
```bash
# Create gather task for auth project
task-manager insert-task \
  --project=authentication-analysis \
  --type=gather_information \
  --instructions="Gather all authentication-related code patterns from the main application"

# Expected output: Shows task with project_id=1, slug auto-generated, status="pending"
```

```bash
# Create gather task for db project with custom slug
task-manager insert-task \
  --project=db-security \
  --type=gather_information \
  --instructions="Collect all database query patterns and access controls" \
  --slug="gather-db-patterns"

# Expected output: Shows task with project_id=2, slug="gather-db-patterns"
```

### Step 2.2: Create Analysis Tasks with Dependencies
```bash
# Create analysis task that depends on the auth gather task (by ID)
task-manager insert-task \
  --project=authentication-analysis \
  --type=oracle_analysis \
  --instructions="Analyze gathered authentication patterns and identify security vulnerabilities" \
  --dependencies=1

# Expected output: Shows task with dependencies=[1] listed
```

```bash
# Create analysis task for db project using slug-based dependency
task-manager insert-task \
  --project=db-security \
  --type=oracle_analysis \
  --instructions="Analyze database access patterns for security vulnerabilities" \
  --dependencies=db-security/gather-db-patterns \
  --slug="analyze-security"

# Expected output: Shows task with dependency on db gather task
```

### Step 2.3: Query Tasks with Dependencies
```bash
# Show all tasks
task-manager query-tasks

# Expected output: Shows all 4 tasks with their details
```

```bash
# Show tasks for specific project
task-manager query-tasks --project-id=1

# Expected output: Shows only auth project tasks
```

```bash
# Show tasks with dependencies (Note: --show-deps may have timeout issues)
task-manager query-tasks --task-id=3

# Expected output: Shows analysis task details
```

## Test Scenario 3: Task Assignment and Status Management

### Step 3.1: Assign Tasks to Agents
```bash
# Assign auth gather task to code analyzer (same project)
task-manager assign-task \
  --agent=code-analyzer \
  --task=1

# Expected output: Shows assignment details, task status should be "in_progress"
```

```bash
# Assign db gather task to db expert (same project)
task-manager assign-task \
  --agent=db-expert \
  --task=2

# Expected output: Shows assignment details
```

### Step 3.2: Test Project Validation
```bash
# Try to assign db task to auth agent (should fail due to project mismatch)
task-manager assign-task \
  --agent=code-analyzer \
  --task=db-security/analyze-security

# Expected output: Error about agent belonging to different project
```

### Step 3.3: Test Dependency Blocking
```bash
# Mark task 1 as completed (simulate completion)
sqlite3 /tmp/agent-work.db "UPDATE tasks SET status='completed', completed_at=CURRENT_TIMESTAMP WHERE id=1;"

# Try to assign analysis task that depends on completed task (should work)
task-manager assign-task \
  --agent=oracle \
  --task=3

# Expected output: Successful assignment
```

```bash
# Try to assign analysis task with incomplete dependencies (should fail)
task-manager assign-task \
  --agent=db-expert \
  --task=4

# Expected output: Error about incomplete dependencies
```

### Step 3.4: Test Assignment Validation
```bash
# Try to assign another task to busy agent (should fail)
task-manager assign-task \
  --agent=oracle \
  --task=1

# Expected output: Error message about agent already being assigned
```

### Step 3.5: Test Force Reassignment
```bash
# Create a new agent for testing force reassignment
task-manager create-agent \
  --name="Backup Analyst" \
  --description="Backup analyst for testing" \
  --project=authentication-analysis \
  --slug="backup"

# Force reassign task from one agent to another
task-manager assign-task \
  --agent=backup \
  --task=3 \
  --force

# Expected output: Successful reassignment with debug logs
```

## Test Scenario 4: Location Management

### Step 4.1: Add Code Locations
```bash
# Add individual location using task ID
task-manager insert-locations \
  --task-id=1 \
  --location="src/auth/login.go" \
  --description="Main login function with password validation"

# Expected output: Shows inserted location with task association
```

```bash
# Add multiple locations at once using task slug
task-manager insert-locations \
  --task=db-security/gather-db-patterns \
  --locations="src/db/queries.go:Database query functions" \
  --locations="src/db/migrations/:Database schema migrations" \
  --locations="src/middleware/auth.go:Database access authorization"

# Expected output: Shows multiple inserted locations
```

### Step 4.2: Test Location Validation
```bash
# Try to add location to non-existent task (should fail)
task-manager insert-locations \
  --task-id=999 \
  --location="test.go" \
  --description="Test file"

# Expected output: Error about task not found
```

## Test Scenario 5: Report Generation

### Step 5.1: Create Reports
```bash
# Create report with direct content using task ID
task-manager create-report \
  --task-id=1 \
  --content="Gathered 15 authentication patterns. Found potential vulnerabilities in password reset flow."

# Expected output: Shows created report with content preview
```

```bash
# Create report from file (create test file first)
echo "# Database Pattern Analysis

## Summary
Detailed analysis of database patterns shows:
- 25 query patterns analyzed
- 3 potential SQL injection points found
- Authentication bypass vulnerability in admin queries

## Recommendations
1. Use parameterized queries
2. Implement proper input validation  
3. Add query result sanitization" > /tmp/test-report.md

task-manager create-report \
  --task=db-security/gather-db-patterns \
  --content-file=/tmp/test-report.md \
  --location-ids=2,3

# Expected output: Shows report with linked locations
```

## Test Scenario 6: Output Formats and Filtering

### Step 6.1: Test Different Output Formats
```bash
# JSON output
task-manager query-tasks --output=json

# Expected output: Valid JSON with all task data
```

```bash
# CSV output
task-manager list-projects --output=csv

# Expected output: CSV format with headers
```

```bash
# YAML output
task-manager list-agents --output=yaml

# Expected output: Valid YAML format
```

### Step 6.2: Test Field Selection
```bash
# Select specific fields
task-manager query-tasks --fields=id,slug,status,type

# Expected output: Only specified fields shown
```

### Step 6.3: Test Filtering Options
```bash
# Filter by status
task-manager query-tasks --status=in_progress

# Expected output: Only in_progress tasks
```

```bash
# Filter by type
task-manager query-tasks --type=gather_information

# Expected output: Only gather_information tasks
```

```bash
# Filter by agent
task-manager query-tasks --agent-id=1

# Expected output: Only tasks assigned to agent 1
```

```bash
# Limit results
task-manager query-tasks --limit=2

# Expected output: Maximum 2 tasks
```

## Test Scenario 7: Error Handling and Edge Cases

### Step 7.1: Test Invalid References
```bash
# Invalid project reference
task-manager insert-task \
  --project=nonexistent \
  --type=gather_information \
  --instructions="Test task"

# Expected output: Error about project not found
```

```bash
# Invalid agent reference
task-manager assign-task \
  --agent=nonexistent \
  --task=1

# Expected output: Error about agent not found
```

```bash
# Invalid task reference
task-manager assign-task \
  --agent=code-analyzer \
  --task=nonexistent/task

# Expected output: Error about task not found
```

### Step 7.2: Test Agent Creation Without Project
```bash
# Try to create agent without specifying project (should fail)
task-manager create-agent \
  --name="Orphan Agent" \
  --description="Agent without project"

# Expected output: Error about missing required project parameter
```

### Step 7.3: Test Slug Uniqueness
```bash
# Try to create project with duplicate name (should get unique slug)
task-manager create-project \
  --name="Authentication Analysis" \
  --description="Another auth analysis project"

# Expected output: Should create with slug like "authentication-analysis-2"
```

### Step 7.4: Test Database Locking
```bash
# Run multiple commands simultaneously (in different terminals)
# Terminal 1:
task-manager query-tasks &

# Terminal 2:
task-manager list-projects &

# Expected output: Both should complete successfully without locking issues
```

## Test Scenario 8: Complete Workflow Validation

### Step 8.1: End-to-End Workflow
```bash
# 1. Create new project
task-manager create-project \
  --name="API Security Review" \
  --description="Review API endpoints for security vulnerabilities" \
  --slug="api-security"

# 2. Create specialized agents for the project
task-manager create-agent \
  --name="API Security Expert" \
  --description="Specialized in API security analysis" \
  --project=api-security \
  --slug="api-expert"

task-manager create-agent \
  --name="Code Gatherer" \
  --description="Collects API endpoint information" \
  --project=api-security \
  --slug="api-gatherer"

# 3. Create gather task
task-manager insert-task \
  --project=api-security \
  --type=gather_information \
  --instructions="Collect all API endpoint definitions and security configurations" \
  --slug="gather-apis"

# 4. Assign gather task
task-manager assign-task \
  --agent=api-gatherer \
  --task=api-security/gather-apis

# 5. Add locations
task-manager insert-locations \
  --task=api-security/gather-apis \
  --locations="src/api/routes.go:API route definitions" \
  --locations="src/middleware/security.go:Security middleware"

# 6. Create analysis task with dependencies
task-manager insert-task \
  --project=api-security \
  --type=oracle_analysis \
  --instructions="Analyze API security patterns and identify vulnerabilities" \
  --dependencies=api-security/gather-apis \
  --slug="analyze-api-security"

# 7. Complete gather task (simulate)
sqlite3 /tmp/agent-work.db "UPDATE tasks SET status='completed', completed_at=CURRENT_TIMESTAMP WHERE slug='gather-apis';"

# 8. Assign analysis task
task-manager assign-task \
  --agent=api-expert \
  --task=api-security/analyze-api-security

# 9. Generate report
task-manager create-report \
  --task=api-security/gather-apis \
  --content="API gathering complete. Found 25 endpoints, 3 without authentication."

# 10. Review complete workflow
task-manager query-tasks --project-id=3 --output=json
```

## Validation Checklist

After running all scenarios, verify:

- [ ] **Projects**: Created with auto-generated and custom slugs
- [ ] **Agents**: Created with proper slug generation and **required project assignment**
- [ ] **Project Binding**: Agents cannot be assigned to tasks from different projects
- [ ] **Tasks**: Created with dependencies and proper slug handling
- [ ] **Assignment**: Tasks properly assigned and status updated to in_progress
- [ ] **Dependency Blocking**: Tasks with incomplete dependencies cannot be assigned
- [ ] **Force Reassignment**: --force flag allows reassignment with proper validation
- [ ] **Agent Tracking**: Agents show current work assignments
- [ ] **Locations**: Code locations properly linked to tasks
- [ ] **Reports**: Reports created and linked to locations
- [ ] **Slug Resolution**: All commands work with both IDs and slugs
- [ ] **Dependencies**: Task dependencies properly resolved and tracked
- [ ] **Output Formats**: JSON, CSV, YAML outputs work correctly
- [ ] **Filtering**: All query filters work as expected
- [ ] **Error Handling**: Appropriate errors for invalid references
- [ ] **Database Locking**: Concurrent access works without issues

## New Features Validated

### Dependency Blocking System
- [x] Tasks with incomplete dependencies cannot be assigned
- [x] Clear error messages showing which dependencies are not completed
- [x] Debug logging shows dependency checking process
- [x] Assignments succeed when all dependencies are completed

### Force Reassignment System  
- [x] --force flag allows reassignment of already-assigned tasks
- [x] Clear error messages when trying to reassign without --force
- [x] Proper cleanup of previous agent assignments
- [x] Transaction safety during reassignment

### Project-Agent Binding
- [x] Agents must be created with a specific project
- [x] Agents cannot be assigned to tasks from different projects
- [x] Clear error messages for project mismatches
- [x] Project validation works even with --force

### Enhanced Debugging
- [x] --log-level=debug provides detailed operation logs
- [x] Dependencies are listed during validation
- [x] Assignment process is fully traced
- [x] Timeout protection prevents hanging commands

## Cleanup

After testing, clean up test data:
```bash
rm -f /tmp/agent-work.db /tmp/test-agent-work.db /tmp/test-report.md
```

## Troubleshooting

### Common Issues

1. **Database locked errors**: Wait a moment and retry, or check for hung processes
2. **Task not found errors**: Verify slug format is correct (project_slug/task_slug)
3. **Agent already assigned errors**: Check agent status with `list-agents`
4. **Dependency errors**: Ensure parent tasks exist and are completed before assignment
5. **Project mismatch errors**: Ensure agent and task belong to the same project
6. **Missing project parameter**: All agents must be created with --project specified

### Debug Commands

```bash
# Check database file exists
ls -la /tmp/agent-work.db

# Verify database schema (requires sqlite3)
sqlite3 /tmp/agent-work.db ".schema"

# Check all data (requires sqlite3)
sqlite3 /tmp/agent-work.db "SELECT * FROM projects;"
sqlite3 /tmp/agent-work.db "SELECT * FROM agents;"
sqlite3 /tmp/agent-work.db "SELECT * FROM tasks;"
sqlite3 /tmp/agent-work.db "SELECT * FROM task_dependencies;"

# Debug dependency checking
task-manager --log-level=debug assign-task --agent=agent-slug --task=task-id

# Check specific task dependencies manually
sqlite3 /tmp/agent-work.db "SELECT td.task_id, td.parent_task_id, t.status FROM task_dependencies td JOIN tasks t ON td.parent_task_id = t.id WHERE td.task_id = YOUR_TASK_ID;"
```

This comprehensive QA guide ensures all features of the agent task management system work correctly, including the new dependency blocking, force reassignment, and project-agent binding features.
