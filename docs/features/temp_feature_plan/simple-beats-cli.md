# Implementation Plan: Simple Beats CLI (sbeads)

## Problem Statement
Build a lightweight, local-only task management CLI similar to beads but simplified. Store tasks in JSON Lines format (jazzy lines), support parent-child relationships, dependencies, priorities, and be easily installable on Unix systems.

## Solution Overview
A Go-based CLI using Cobra framework that stores tasks in `.sbeads/tasks.jsonl` file. Tasks include ID, title, description, priority (0-4), parent ID, and dependency list. The tool provides essential CRUD operations with dependency tracking.

## Project Structure
```
sbeads/
├── cmd/
│   ├── root.go         # Root command & config
│   ├── init.go         # Initialize .sbeads directory
│   ├── create.go       # Create new tasks
│   ├── list.go         # List tasks with filters
│   ├── show.go         # Show task details
│   ├── update.go       # Update task fields
│   ├── close.go        # Close tasks
│   ├── delete.go       # Delete tasks
│   └── dep.go          # Manage dependencies
├── internal/
│   ├── storage/
│   │   └── storage.go  # JSONL file operations
│   ├── models/
│   │   └── task.go     # Task struct definitions
│   └── utils/
│       └── utils.go    # Helper functions
├── Taskfile.yml        # Build automation
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## Steps

### Step 1: Initialize Go Project & Dependencies
Set up Go module and install Cobra CLI tools

Commands:
```bash
go mod init github.com/Martins6/simple-beads
go get github.com/spf13/cobra@latest
go install github.com/spf13/cobra-cli@latest
cobra-cli init
```

### Step 2: Define Task Model
Create Task struct with all required fields

Fields:
- ID (auto-generated, format: sb-XXXX)
- Title (string)
- Description (string)
- Priority (int, 0-4, 0=highest)
- Status (string: open, closed)
- Parent (string, optional parent ID)
- Dependencies ([]string, optional dependency IDs)
- CreatedAt (time.Time)
- UpdatedAt (time.Time)
- ClosedAt (time.Time, optional)

### Step 3: Implement JSONL Storage Layer
Create file-based storage using JSON Lines format

Operations:
- `LoadTasks()`: Read all tasks from `.sbeads/tasks.jsonl`
- `SaveTask(task)`: Append single task to file
- `UpdateTask(task)`: Rewrite file with updated task
- `DeleteTask(id)`: Rewrite file excluding task
- `FindByID(id)`: Return single task
- `FindByParent(parentID)`: Return child tasks
- `FindBlocked()`: Return tasks with unmet dependencies

### Step 4: Create CLI Commands
Implement all user-facing commands using Cobra

Commands to implement:
1. `sb init` - Create `.sbeads/` directory and empty `tasks.jsonl`
2. `sb create "title" [flags]` - Create new task
   - Flags: `-d, --description`, `-p, --priority`, `--parent`, `--deps`
3. `sb list [flags]` - List tasks
   - Flags: `--all`, `-p, --priority`, `--status`, `--parent`
4. `sb show <id>` - Show task details including dependencies
5. `sb update <id> [flags]` - Update task fields
6. `sb close <id>` - Mark task as closed
7. `sb delete <id>` - Remove task
8. `sb dep add <task> <dependency>` - Add dependency
9. `sb dep remove <task> <dependency>` - Remove dependency
10. `sb ready` - Show tasks with no blockers

### Step 5: Implement Dependency Resolution
Add logic to check if dependencies are met

Logic:
- Check if all dependency tasks are closed
- Detect circular dependencies when adding
- Show blocked tasks in `sb blocked` command

### Step 6: Create Taskfile for Building
Setup Taskfile for cross-platform builds and installation

Targets:
- `task build` - Build for current platform
- `task build-all` - Build for Linux, macOS (amd64, arm64)
- `task install` - Install to `/usr/local/bin`
- `task clean` - Clean build artifacts
- `task test` - Run tests

### Step 7: Create Installation Script
Unix installation script for easy setup

Features:
- Detect OS and architecture
- Download appropriate binary from releases
- Or build from source if Go is available
- Install to `/usr/local/bin` or `~/.local/bin`
- Verify installation

### Step 8: Add Tests
Unit tests for storage and models

Coverage:
- Task CRUD operations
- Dependency resolution
- JSONL file format
- Priority filtering

## Challenges
1. JSONL File Concurrency: Use file locking or warn users
2. Dependency Cycles: Check cycle before adding dependency
3. Large File Performance: Document as limitation
4. Cross-Platform Paths: Use filepath.Join()

## Risk Assessment
- Level: Low
- Simple file-based storage, well-understood patterns, no external dependencies
