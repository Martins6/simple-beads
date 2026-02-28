# Overview

Parent-child task hierarchies and dependency management with circular dependency detection.

# Details

- Parent field allows creating task hierarchies and subtasks
- Dependencies define task prerequisites that must be closed first
- DependencyChecker utility for in-memory dependency validation
- Circular dependency detection using recursive traversal in memory
- Database-level circular detection using recursive CTE
- GetBlockingDependencies returns list of unclosed dependencies
- GetDependencyChain returns full transitive dependency list
- GetDependents finds all tasks that depend on a given task
- ValidateDependency checks existence, self-dependency, and cycles
- Self-dependency prevention enforced at application layer

# File Paths

- internal/utils/dependencies.go
- internal/models/task.go
- internal/storage/storage.go
