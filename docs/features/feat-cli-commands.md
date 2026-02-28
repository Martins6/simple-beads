# Overview

Command-line interface built with Cobra framework providing task management commands.

# Details

- Root command: `sb` with help and global storage initialization
- `sb init` - Initialize `.sbeads` directory and database
- `sb create [title]` - Create new task with flags for description, priority, parent, deps
- `sb list` - List tasks with filters for status, priority, parent, all
- `sb show [task-id]` - Display detailed task information
- `sb update [task-id]` - Update task fields selectively
- `sb close [task-id...]` - Mark tasks as closed/completed
- `sb delete [task-id...]` - Permanently delete tasks
- `sb ready` - Show tasks ready to work on (open, not blocked)
- `sb blocked` - Show blocked tasks with blocking dependencies
- `sb dep add/remove` - Manage task dependencies
- Global storage instance initialized in root command

# File Paths

- cmd/root.go
- cmd/init.go
- cmd/create.go
- cmd/list.go
- cmd/show.go
- cmd/update.go
- cmd/close.go
- cmd/delete.go
- cmd/ready.go
- cmd/blocked.go
- cmd/dep.go
- main.go
