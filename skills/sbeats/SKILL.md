# sbeads - Simple Beads Task Management

## Overview

sbeads (Simple Beads) is a lightweight, local-only task management CLI that stores tasks in SQLite database. It provides hierarchical task relationships, dependency tracking, and priority-based workflow management with full multi-agent concurrency support.

## Installation

**Note:** sbeads is designed for Unix systems (macOS and Linux). Windows is not supported.

### Quick Install (Unix - No sudo required)

```bash
# Using curl (recommended)
curl -fsSL https://raw.githubusercontent.com/user/sbeads/main/install.sh | bash

# Or clone and install manually
git clone https://github.com/user/sbeads.git
cd sbeads
./install.sh              # Install to ~/.local/bin (no sudo)
```

The install script will:
- Install to `~/.local/bin` by default (no sudo needed)
- Tell you how to add it to your PATH
- Never ask for administrator privileges

### Using Go

```bash
go install github.com/user/sbeads@latest
```

Make sure `$GOPATH/bin` is in your PATH:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Quick Start

```bash
# Initialize sbeads in your project directory
sb init

# Create your first task
sb create "Fix navigation bug" -p 1

# List all open tasks
sb list

# See what you can work on now
sb ready
```

## Core Commands

### Task Management

| Command | Description | Example |
|---------|-------------|---------|
| `sb init` | Initialize sbeads storage | `sb init` |
| `sb create` | Create a new task | `sb create "Title" -d "Description" -p 0` |
| `sb list` | List tasks with filters | `sb list --all -p 0 --status open` |
| `sb show` | Show task details | `sb show sb-abc123` |
| `sb update` | Update task fields | `sb update sb-abc123 --title "New title"` |
| `sb close` | Mark task(s) as closed | `sb close sb-abc123 sb-def456` |
| `sb delete` | Delete task(s) | `sb delete sb-abc123` |

### Workflow Commands

| Command | Description | Example |
|---------|-------------|---------|
| `sb ready` | Show unblocked tasks | `sb ready` |
| `sb blocked` | Show blocked tasks | `sb blocked` |

### Dependency Commands

| Command | Description | Example |
|---------|-------------|---------|
| `sb dep add` | Add dependency | `sb dep add sb-task sb-dependency` |
| `sb dep remove` | Remove dependency | `sb dep remove sb-task sb-dependency` |

## Command Flags

### Create Task Flags

- `-d, --description`: Task description
- `-p, --priority`: Priority (0-4, where 0 is highest). Default: 2
- `--parent`: Parent task ID (for hierarchical tasks)
- `--deps`: Comma-separated list of dependency IDs

### List Task Flags

- `--all`: Show all tasks including closed
- `-p, --priority`: Filter by priority (0-4)
- `--status`: Filter by status (open, closed)
- `--parent`: Filter by parent task ID

## Examples

### Basic Task Creation

```bash
# Simple task
sb create "Fix login bug"

# With description and high priority
sb create "Security audit" -d "Review authentication code" -p 0

# Low priority maintenance task
sb create "Update dependencies" -p 4
```

### Parent-Child Tasks (Epics)

```bash
# Create epic (parent task)
sb create "Launch v2.0" -d "Major release with new features" -p 0
# Result: sb-abc123

# Create child tasks
sb create "Design UI" --parent sb-abc123 -p 1
sb create "Implement backend" --parent sb-abc123 -p 1
sb create "Write tests" --parent sb-abc123 -p 2
sb create "Deploy" --parent sb-abc123 -p 0

# View all children
sb list --parent sb-abc123
```

### Managing Dependencies

```bash
# Create tasks
sb create "Design database schema"
# Result: sb-design

sb create "Write API endpoints"
# Result: sb-api

sb create "Build frontend"
# Result: sb-frontend

# Set dependencies
sb dep add sb-api sb-design        # API depends on design
sb dep add sb-frontend sb-api      # Frontend depends on API

# Check what's blocked
sb blocked
# Shows: sb-api (blocked by sb-design)
#        sb-frontend (blocked by sb-api)

# Close design task
sb close sb-design

# Check what's ready now
sb ready
# Shows: sb-api (now unblocked)
```

### Daily Workflow

```bash
# Morning: Check what's ready to work on
sb ready

# Work on highest priority task...
# ...complete work...

# Mark as done
sb close sb-task-id

# Check what became unblocked
sb ready

# Afternoon: Create new tasks discovered during work
sb create "Refactor auth middleware" -d "Extract from main.go" -p 1
sb create "Add rate limiting" --parent sb-epic-id -p 2

# End of day: Review all tasks
sb list
sb blocked  # See what's waiting on others
```

### Filtering and Organization

```bash
# Show only high priority (P0-P1) tasks
sb list -p 0
sb list -p 1

# Show all closed tasks
sb list --status closed

# Show everything including closed
sb list --all

# Show tasks by parent
sb list --parent sb-epic-123

# Combine filters
sb list --parent sb-epic-123 --status open -p 0
```

### Batch Operations

```bash
# Close multiple tasks at once
sb close sb-abc123 sb-def456 sb-ghi789

# Delete multiple tasks
sb delete sb-old1 sb-old2 sb-old3
```

### Complex Project Example

```bash
# Initialize
sb init

# Create project epic
sb create "Build e-commerce platform" -d "Full stack e-commerce solution" -p 0
# Epic ID: sb-epic1

# Phase 1: Design
sb create "Create wireframes" --parent sb-epic1 -p 0
# Result: sb-design1
sb create "Design database schema" --parent sb-epic1 -p 0
# Result: sb-db1

# Phase 2: Backend (depends on Phase 1)
sb create "Setup PostgreSQL" --parent sb-epic1 -p 1 --deps sb-db1
# Result: sb-setup-db
sb create "Implement user auth" --parent sb-epic1 -p 1 --deps sb-setup-db
# Result: sb-auth
sb create "Build product API" --parent sb-epic1 -p 1 --deps sb-setup-db
# Result: sb-products

# Phase 3: Frontend (depends on Phase 2)
sb create "Setup React app" --parent sb-epic1 -p 2
# Result: sb-react
sb create "Build login page" --parent sb-epic1 -p 2 --deps sb-auth,sb-react
# Result: sb-login
sb create "Build product catalog" --parent sb-epic1 -p 2 --deps sb-products,sb-react
# Result: sb-catalog

# Phase 4: Deployment (depends on everything)
sb create "Deploy to production" --parent sb-epic1 -p 0 --deps sb-login,sb-catalog
# Result: sb-deploy

# Check workflow
sb ready      # Shows: wireframes, DB schema, React setup
sb blocked    # Shows: everything else waiting on dependencies
```

## Data Storage

Tasks are stored in `.sbeads/tasks.db` using SQLite with WAL (Write-Ahead Logging) mode:

### Database Schema

```sql
-- Tasks table
CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    priority INTEGER NOT NULL DEFAULT 2,
    status TEXT NOT NULL DEFAULT 'open',
    parent TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    closed_at DATETIME
);

-- Dependencies table
CREATE TABLE dependencies (
    task_id TEXT NOT NULL,
    depends_on TEXT NOT NULL,
    PRIMARY KEY (task_id, depends_on),
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (depends_on) REFERENCES tasks(id) ON DELETE CASCADE
);
```

### Benefits:
- **ACID transactions**: Atomic operations guarantee data consistency
- **Multi-agent safe**: SQLite uses file-level locking - multiple processes can safely access the database
- **No data corruption**: If an agent crashes mid-write, the database remains consistent
- **Clear errors**: "database is locked" tells you exactly what's happening
- **Fast queries**: Indexed columns for efficient filtering
- **Portable**: Single `.db` file that's easy to backup

## Best Practices

### 1. Use Epics for Large Projects

Create a parent task for large initiatives, then break down into child tasks:

```bash
sb create "Q1 Goals" -d "Quarterly objectives" -p 0
sb create "Objective 1" --parent sb-epic -p 0
sb create "Objective 2" --parent sb-epic -p 0
```

### 2. Set Dependencies for Natural Flow

Don't just create a list - create a workflow:

```bash
# Good: Tasks flow naturally
Design → Backend → Frontend → Deploy

# Not as good: Everything at same priority with no dependencies
Design, Backend, Frontend, Deploy (all P0, no deps)
```

### 3. Use Priority Levels

- **P0**: Critical, blocking other work (releases, hotfixes)
- **P1**: Important, high value features
- **P2**: Normal work, standard tasks
- **P3**: Nice to have, improvements
- **P4**: Low priority, maintenance, chores

### 4. Morning Workflow

Start your day by checking `sb ready` - it shows only unblocked tasks sorted by priority.

### 5. Keep Tasks Small

A task should be completable in a reasonable time (few hours to a day). Break down large tasks into smaller ones with dependencies.

### 6. Clean Up Regularly

```bash
# Review closed tasks
sb list --status closed

# Archive old completed work (manually backup and delete)
sb delete sb-old-task1 sb-old-task2
```

## Integration with Git

Since data is stored in `.sbeads/tasks.db` (SQLite), you can track tasks with your code:

```bash
# Add to git
git add .sbeads/tasks.db
git commit -m "Update tasks: completed auth system"

# Share with team
git push

# Team members pull and see updated tasks
git pull
sb list
```

**Note:** SQLite databases are binary files, so git diffs won't be human-readable like they were with JSONL. However, you can export to SQL or use git's binary diff handling.

## Troubleshooting

### Task not found
```bash
# Check the exact ID
sb list
sb show sb-abc123
```

### Circular dependency error
```bash
# List dependencies to understand the chain
sb show sb-task-id
# Then remove the problematic dependency
sb dep remove sb-task-id sb-other-id
```

### Database locked errors
If you see "database is locked" errors, it means another process is currently writing to the database. This is expected with multiple agents. SQLite will retry automatically, or you can wait a moment and retry your command.

### Data inspection
```bash
# View raw database (requires sqlite3 CLI)
sqlite3 .sbeads/tasks.db "SELECT * FROM tasks;"

# Export to SQL
sqlite3 .sbeads/tasks.db ".dump" > tasks_backup.sql

# Check database integrity
sqlite3 .sbeads/tasks.db "PRAGMA integrity_check;"
```

## Comparison with beads

| Feature | sb (Simple Beats) | beads |
|---------|-------------------|-------|
| Storage | SQLite (WAL mode) | SQLite + JSONL |
| Sync | Manual (git) | Automatic |
| Complexity | Simple | Full-featured |
| Setup | Single binary | Requires setup |
| Best for | Personal/small team | Team workflows |

sb is designed for simplicity - it's beads without the complexity of automatic sync, web UI, and advanced features.

## Tips for AI Assistants

When working with a codebase that uses sbeads:

1. Always check `sb ready` first to see what tasks are actionable
2. Use `sb blocked` to understand what's waiting
3. Create child tasks for complex features
4. Set dependencies to model the natural order of work
5. Close tasks as you complete work
6. Use `--parent` to organize related tasks under epics

Example workflow when implementing a feature:
```bash
# Check what's ready
sb ready

# Create implementation tasks
sb create "Implement auth middleware" --parent sb-feature-epic -p 1
sb create "Add JWT validation" --parent sb-feature-epic -p 2 --deps sb-middleware
sb create "Write tests" --parent sb-feature-epic -p 2 --deps sb-middleware

# Work on tasks and close as completed
sb close sb-middleware
sb ready  # See what became unblocked
```
