# Overview

Simple Beads (sb) is a lightweight, local-first task management CLI. It stores tasks in SQLite with support for parent-child relationships, dependencies, and priorities.

The application provides a simple command-line interface for managing tasks with features like task creation, updates, filtering, dependency management, and status tracking. It uses SQLite with WAL mode for multi-process safety and ACID transactions.

# Files

- feat-storage.md - SQLite storage with WAL mode, transactions, and database schema
- feat-task-management.md - Task CRUD operations and lifecycle management
- feat-cli-commands.md - CLI interface using Cobra framework
- feat-task-relationships.md - Parent-child hierarchies and dependency management
- feat-priority-filtering.md - Priority levels and status filtering
