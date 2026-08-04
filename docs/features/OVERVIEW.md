# Overview

Simple Beads (sb) is a lightweight, local-first task management CLI. It stores tasks in SQLite with support for parent-child relationships, dependencies, and priorities.

The application provides a simple command-line interface for managing tasks with features like task creation, updates, filtering, dependency management, and status tracking. It uses SQLite with WAL mode for multi-process safety and ACID transactions.

## Architecture

### Original Design
The project was designed as a Go-based CLI using the Cobra framework. The original plan used JSON Lines format for storage but was later migrated to SQLite for better data integrity and multi-process safety.

### Project Structure
```
sbeads/
├── cmd/              # CLI commands (Cobra)
│   ├── root.go
│   ├── init.go
│   ├── create.go
│   ├── list.go
│   ├── show.go
│   ├── update.go
│   ├── close.go
│   ├── delete.go
│   └── dep.go
├── internal/
│   ├── storage/      # Database operations
│   ├── models/       # Task struct definitions
│   └── utils/        # Helper functions
├── Taskfile.yml      # Build automation
└── main.go
```

### Key Design Decisions
- **Storage**: Migrated from JSONL to SQLite with WAL mode for ACID compliance
- **Task IDs**: Format `sb-XXXX` (first 4 characters of UUID)
- **Priority Levels**: P0 (highest) to P4 (lowest)
- **Dependencies**: Tasks can have multiple dependencies; circular dependency detection implemented
- **Parent-Child**: Tasks can have one parent, forming hierarchical structures

# Files

- feat-storage.md - SQLite storage with WAL mode, transactions, and database schema
- feat-task-management.md - Task CRUD operations and lifecycle management
- feat-cli-commands.md - CLI interface using Cobra framework
- feat-task-relationships.md - Parent-child hierarchies and dependency management
- feat-priority-filtering.md - Priority levels and status filtering
- feat-date-filtering.md - Date filtering for `sb list` with --on/--after/--before flags
