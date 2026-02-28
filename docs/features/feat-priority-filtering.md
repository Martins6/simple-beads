# Overview

Priority levels and status filtering for organizing and finding tasks.

# Details

- Priority levels 0-4 where 0 is highest priority
- Default priority is P2 (medium)
- Status states: open, closed
- List command filters: --all, --status, --priority, --parent
- Ready command shows open tasks not blocked by dependencies
- Blocked command shows open tasks with unclosed dependencies
- Sorting by priority (ascending) then creation date
- Tabwriter formatting for aligned output
- Priority displayed as P0, P1, P2, P3, P4 in listings

# File Paths

- cmd/list.go
- cmd/ready.go
- cmd/blocked.go
- internal/models/task.go
