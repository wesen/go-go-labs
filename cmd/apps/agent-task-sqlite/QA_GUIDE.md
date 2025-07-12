# QA Guide: Agent Task SQLite Testing

This guide provides step-by-step testing scenarios to validate the complete agent task management system. Follow these scenarios to ensure all features work correctly.

## Prerequisites

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

### Step 1.3: Create Agents
```bash
# Create code analyzer agent
task-manager create-agent \
  --name="Code Analyzer" \
  --description="Specialized in analyzing code patterns and architecture"

# Expected output: Shows agent with id=1, slug="code-analyzer"
```

```bash
# Create oracle agent with custom slug
task-manager create-agent \
  --name="Oracle Analyst" \
  --description="Performs deep analysis and generates insights from gathered data" \
  --slug="oracle"

# Expected output: Shows agent with id=2, slug="oracle"
```

### Step 1.4: List Agents
```bash
# List all agents
task-manager list-agents

# Expected output: Shows both agents, current_project_id and current_task_id should be null
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
# Create analysis task that depends on the auth gather task
task-manager insert-task \
  --project=authentication-analysis \
  --type=oracle_analysis \
  --instructions="Analyze gathered authentication patterns and identify security vulnerabilities" \
  --dependencies=authentication-analysis/gather-all-authentication-related-code-patterns

# Expected output: Shows task with dependencies listed
```

```bash
# Create analysis task for db project using numeric ID dependency
task-manager insert-task \
  --project=db-security \
  --type=oracle_analysis \
  --instructions="Analyze database access patterns for security vulnerabilities" \
  --dependencies=2 \
  --slug="analyze-security"

# Expected output: Shows task with dependency on task ID 2
```

### Step 2.3: Query Tasks
```bash
# Show all tasks
task-manager query-tasks

# Expected output: Shows all 4 tasks with their details
```

```bash
# Show tasks for specific project
task-manager query-tasks --project=authentication-analysis

# Expected output: Shows only auth project tasks
```

```bash
# Show tasks with dependencies
task-manager query-tasks --show-deps

# Expected output: Shows all tasks with dependency information
```

## Test Scenario 3: Task Assignment and Status Management

### Step 3.1: Assign Tasks to Agents
```bash
# Assign auth gather task to code analyzer
task-manager assign-task \
  --agent=code-analyzer \
  --task=authentication-analysis/gather-all-authentication-related-code-patterns

# Expected output: Shows assignment details, task status should be "in_progress"
```

```bash
# Assign db gather task to oracle agent
task-manager assign-task \
  --agent=oracle \
  --task=db-security/gather-db-patterns

# Expected output: Shows assignment details
```

### Step 3.2: Verify Agent Status Updates
```bash
# Check agent status after assignment
task-manager list-agents

# Expected output: Both agents should show current_project_id and current_task_id
```

```bash
# Query tasks by agent
task-manager query-tasks --agent=code-analyzer

# Expected output: Shows only tasks assigned to code-analyzer
```

### Step 3.3: Test Assignment Validation
```bash
# Try to assign another task to busy agent (should fail)
task-manager assign-task \
  --agent=code-analyzer \
  --task=db-security/analyze-security

# Expected output: Error message about agent already being assigned
```

```bash
# Try to assign analysis task that has pending dependencies (should work but task stays pending)
task-manager assign-task \
  --agent=oracle \
  --task=authentication-analysis/analyze-gathered-authentication-patterns

# Expected output: Error about agent already being assigned to another task
```

## Test Scenario 4: Location Management

### Step 4.1: Add Code Locations
```bash
# Add individual location
task-manager insert-locations \
  --task=authentication-analysis/gather-all-authentication-related-code-patterns \
  --location="src/auth/login.go" \
  --description="Main login function with password validation"

# Expected output: Shows inserted location with task association
```

```bash
# Add multiple locations at once
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
  --task=999 \
  --location="test.go" \
  --description="Test file"

# Expected output: Error about task not found
```

## Test Scenario 5: Report Generation

### Step 5.1: Create Reports
```bash
# Create report with direct content
task-manager create-report \
  --task=authentication-analysis/gather-all-authentication-related-code-patterns \
  --content="Gathered 15 authentication patterns. Found potential vulnerabilities in password reset flow."

# Expected output: Shows created report with content preview
```

```bash
# Create report from file (create test file first)
echo "Detailed analysis of database patterns..." > /tmp/test-report.md
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

### Step 7.2: Test Slug Uniqueness
```bash
# Try to create project with duplicate name (should get unique slug)
task-manager create-project \
  --name="Authentication Analysis" \
  --description="Another auth analysis project"

# Expected output: Should create with slug like "authentication-analysis-2"
```

### Step 7.3: Test Database Locking
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

# 2. Create specialized agent
task-manager create-agent \
  --name="API Security Expert" \
  --description="Specialized in API security analysis" \
  --slug="api-expert"

# 3. Create gather task
task-manager insert-task \
  --project=api-security \
  --type=gather_information \
  --instructions="Collect all API endpoint definitions and security configurations" \
  --slug="gather-apis"

# 4. Assign task
task-manager assign-task \
  --agent=api-expert \
  --task=api-security/gather-apis

# 5. Add locations
task-manager insert-locations \
  --task=api-security/gather-apis \
  --locations="src/api/routes.go:API route definitions" \
  --locations="src/middleware/security.go:Security middleware"

# 6. Create analysis task
task-manager insert-task \
  --project=api-security \
  --type=oracle_analysis \
  --instructions="Analyze API security patterns and identify vulnerabilities" \
  --dependencies=api-security/gather-apis \
  --slug="analyze-api-security"

# 7. Generate report
task-manager create-report \
  --task=api-security/gather-apis \
  --content="API gathering complete. Found 25 endpoints, 3 without authentication."

# 8. Review complete workflow
task-manager query-tasks --project=api-security --show-deps --output=json
```

## Validation Checklist

After running all scenarios, verify:

- [ ] **Projects**: Created with auto-generated and custom slugs
- [ ] **Agents**: Created with proper slug generation and status tracking
- [ ] **Tasks**: Created with dependencies and proper slug handling
- [ ] **Assignment**: Tasks properly assigned and status updated to in_progress
- [ ] **Agent Tracking**: Agents show current work assignments
- [ ] **Locations**: Code locations properly linked to tasks
- [ ] **Reports**: Reports created and linked to locations
- [ ] **Slug Resolution**: All commands work with both IDs and slugs
- [ ] **Dependencies**: Task dependencies properly resolved and tracked
- [ ] **Output Formats**: JSON, CSV, YAML outputs work correctly
- [ ] **Filtering**: All query filters work as expected
- [ ] **Error Handling**: Appropriate errors for invalid references
- [ ] **Database Locking**: Concurrent access works without issues

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
4. **Dependency errors**: Ensure parent tasks exist before creating dependent tasks

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
```

This comprehensive QA guide ensures all features of the agent task management system work correctly and provides a reliable testing framework for future development. 