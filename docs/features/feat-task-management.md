# Overview

Task model and lifecycle management including creation, validation, and state transitions.

# Details

- Task struct with ID, Title, Description, Priority, Status, Parent, Dependencies
- Task IDs generated in format `sb-XXXX` (first 4 chars of UUID)
- Priority levels: P0 (highest) to P4 (lowest), default is P2
- Status states: open and closed
- Timestamps: CreatedAt, UpdatedAt, ClosedAt
- Methods: NewTask, Validate, Close, Reopen
- Dependency methods: AddDependency, RemoveDependency, HasDependency, IsBlocked
- Validation checks for required fields and valid enum values

# File Paths

- internal/models/task.go
