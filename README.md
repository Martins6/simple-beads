# Simple Beads (sb)

A lightweight, local-only task management CLI inspired by [beads](https://github.com/steveyegge/beads). Store tasks in SQLite with support for parent-child relationships, dependencies, and priorities.

## Features

- **Local-first**: All data stored in `.sbeads/tasks.db` (SQLite) - no cloud required
- **Multi-agent safe**: SQLite with WAL mode supports concurrent access from multiple processes
- **ACID transactions**: Atomic operations guarantee data consistency
- **Hierarchical tasks**: Create parent-child task relationships
- **Dependencies**: Block tasks until their dependencies are completed
- **Priorities**: 5 priority levels (P0-P4, where P0 is highest)
- **Fast**: Built with Go and SQLite for speed and efficiency
- **Cross-platform**: Works on Linux and macOS (both Intel and Apple Silicon)

## Installation

**For Unix systems only (macOS and Linux). Windows is not supported.**

### Quick Install

```bash
# Using curl (recommended)
curl -fsSL https://raw.githubusercontent.com/Martins6/simple-beads/main/install.sh | bash

# Or clone and install
git clone https://github.com/Martins6/simple-beads.git
cd sbeads
./install.sh              # Install to /usr/local/bin
./install.sh --local      # Or install to ~/.local/bin (no sudo)
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/Martins6/simple-beads.git
cd sbeads

# Build
go build -o sb .

# Install (choose one)
go install .                    # Recommended - installs to $GOPATH/bin
# Or copy manually:
mkdir -p ~/.local/bin && cp sb ~/.local/bin/  # User local (no sudo)
# Or system-wide (only if /usr/local/bin is writable):
cp sb /usr/local/bin/
```

### Using Go

```bash
go install github.com/Martins6/simple-beads@latest
```

**Note:** Make sure `$GOPATH/bin` or `$HOME/go/bin` is in your PATH. Add to your `.bashrc` or `.zshrc`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Quick Start

```bash
# Initialize sbeads in your project directory
sb init

# Create a task
sb create "Fix login bug"

# Create a high-priority task with description
sb create "Update documentation" -d "Add API examples" -p 0

# List all open tasks
sb list

# Create a child task
sb create "Write tests" --parent sb-abc

# Create a task that depends on another
sb create "Deploy to prod" --deps sb-abc

# See what you can work on now
sb ready

# Mark a task as done
sb close sb-abc
```

## Commands

### Core Commands

| Command     | Description                            | Example                            |
| ----------- | -------------------------------------- | ---------------------------------- |
| `sb init`   | Initialize sbeads in current directory | `sb init`                          |
| `sb create` | Create a new task                      | `sb create "Title" -d "Desc" -p 1` |
| `sb list`   | List tasks with filters                | `sb list --all -p 0`               |
| `sb show`   | Show task details                      | `sb show sb-abc`                   |
| `sb update` | Update task fields                     | `sb update sb-abc --title "New"`   |
| `sb close`  | Mark task(s) as closed                 | `sb close sb-abc sb-123`           |
| `sb delete` | Delete task(s) permanently             | `sb delete sb-abc`                 |

### Workflow Commands

| Command      | Description                 | Example      |
| ------------ | --------------------------- | ------------ |
| `sb ready`   | Show tasks ready to work on | `sb ready`   |
| `sb blocked` | Show blocked tasks          | `sb blocked` |

### Dependency Commands

| Command         | Description            | Example                       |
| --------------- | ---------------------- | ----------------------------- |
| `sb dep add`    | Add dependency to task | `sb dep add sb-abc sb-123`    |
| `sb dep remove` | Remove dependency      | `sb dep remove sb-abc sb-123` |

## Command Flags

### `sb create`

- `-d, --description`: Task description
- `-p, --priority`: Priority (0-4, where 0 is highest). Default: 2
- `--parent`: Parent task ID
- `--deps`: Comma-separated list of dependency IDs

### `sb list`

- `--all`: Show all tasks including closed
- `-p, --priority`: Filter by priority (0-4)
- `--status`: Filter by status (open, closed)
- `--parent`: Filter by parent task ID

### `sb update`

- `--title`: New task title
- `-d, --description`: New description
- `-p, --priority`: New priority (0-4)
- `--parent`: New parent task ID

## Examples

### Creating Tasks

```bash
# Simple task
sb create "Fix navigation bug"

# With description and high priority
sb create "Security audit" -d "Review authentication code" -p 0

# Child task
sb create "Research" --parent sb-abc
sb create "Implementation" --parent sb-abc
sb create "Testing" --parent sb-abc

# Task with dependencies
sb create "Deploy" --deps sb-123,sb-456 -p 1
```

### Managing Dependencies

```bash
# Task B depends on Task A
sb create "Task A"
# Output: ✓ Created task: [sb-abc] Task A (P2) - open

sb create "Task B"
# Output: ✓ Created task: [sb-def] Task B (P2) - open

# Make B depend on A
sb dep add sb-def sb-abc

# Now B is blocked
sb blocked
# Output shows sb-def is blocked by sb-abc

# Complete A
sb close sb-abc

# Now B is ready
sb ready
# Output shows sb-def
```

### Workflow Example

```bash
# Initialize
sb init

# Create epic (parent task)
sb create "Launch v2.0" -d "Major release with new features" -p 0
# Assume it gets ID: sb-epic

# Create child tasks
sb create "Design UI" --parent sb-epic
sb create "Implement backend" --parent sb-epic
sb create "Write tests" --parent sb-epic

# Add dependencies
sb dep add sb-backend sb-design

# Check what's ready to work on
sb ready
# Shows: Design UI (no dependencies)

# Work on design, then close it
sb close sb-design

# Now backend is unblocked
sb ready
# Shows: Implement backend

# Continue...
```

## Data Format

Tasks are stored in `.sbeads/tasks.db` (SQLite database):

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

-- Dependencies table (many-to-many relationship)
CREATE TABLE dependencies (
    task_id TEXT NOT NULL,
    depends_on TEXT NOT NULL,
    PRIMARY KEY (task_id, depends_on),
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (depends_on) REFERENCES tasks(id) ON DELETE CASCADE
);
```

### Benefits of SQLite:

- **ACID transactions**: Atomic operations guarantee data consistency
- **Multi-process safe**: Uses WAL (Write-Ahead Logging) mode for concurrent access
- **No data corruption**: File-level locking prevents race conditions
- **Fast queries**: Indexed columns for efficient filtering
- **Portable**: Single file that's easy to backup or version control

## Task Structure

```go
type Task struct {
    ID           string    // Unique ID (format: sb-XXXX)
    Title        string    // Task title
    Description  string    // Optional description
    Priority     int       // 0-4 (0=highest)
    Status       string    // "open" or "closed"
    Parent       string    // Optional parent task ID
    Dependencies []string  // Optional dependency IDs
    CreatedAt    time.Time // Creation timestamp
    UpdatedAt    time.Time // Last update timestamp
    ClosedAt     time.Time // Optional close timestamp
}
```

## Development

### Prerequisites

- Go 1.21 or later
- Task (optional, for builds): `go install github.com/go-task/task/v3/cmd/task@latest`

### Build

```bash
# Build for current platform
go build -o sb .

# Run tests
go test -v ./...

# Install (choose one)
go install .                    # Recommended - installs to $GOPATH/bin
mkdir -p ~/.local/bin && cp sb ~/.local/bin/  # Or user-local install
```

**Optional:** The project includes a Taskfile.yml for those who prefer using [Task](https://taskfile.dev):

```bash
# If you have Task installed
task build         # Build binary
task test          # Run tests
task install       # Install to $GOPATH/bin via go install
task install-local # Install to ~/.local/bin
```

### Project Structure

```
sbeads/
├── cmd/                    # Cobra commands
│   ├── root.go            # Root command
│   ├── init.go            # Init command
│   ├── create.go          # Create command
│   ├── list.go            # List command
│   ├── show.go            # Show command
│   ├── update.go          # Update command
│   ├── close.go           # Close command
│   ├── delete.go          # Delete command
│   ├── dep.go             # Dependency commands
│   ├── ready.go           # Ready command
│   └── blocked.go         # Blocked command
├── internal/
│   ├── models/
│   │   └── task.go        # Task struct
│   ├── storage/
│   │   └── storage.go     # SQLite storage
│   └── utils/
│       └── dependencies.go # Dependency utilities
├── Taskfile.yml           # Build automation
├── install.sh             # Installation script
├── go.mod
├── main.go
└── README.md
```

## Comparison with beads

| Feature    | sb (Simple Beads) | beads              |
| ---------- | ----------------- | ------------------ |
| Storage    | SQLite (WAL mode) | SQLite + JSONL     |
| Sync       | Manual (git)      | Automatic git sync |
| Complexity | Simple            | Full-featured      |
| Setup      | Single binary     | Requires setup     |
| Best for   | Personal projects | Team workflows     |

sb is designed for simplicity - it's beads without the complexity of sync, web UI, and advanced features.

## License

MIT License - see LICENSE file for details

## Contributing

Contributions welcome! Please feel free to submit a Pull Request.

## Acknowledgments

Inspired by [beads](https://github.com/anomalyco/beads) - an excellent issue tracker with first-class dependency support.
