# Overview

SQLite storage layer with WAL mode for multi-process safety and ACID transactions.

# Details

- Uses SQLite database stored in `.sbeads/tasks.db`
- WAL mode enabled for better concurrency and multi-process access
- Busy timeout set to 5000ms to prevent lock conflicts
- Supports transaction-based operations for data integrity
- Database schema includes tasks and dependencies tables
- Creates indexes on status, priority, parent, and dependencies for performance
- Provides CRUD operations: LoadTasks, SaveTask, UpdateTask, DeleteTask
- Query methods: FindByID, FindByParent, FindBlocked, FindReady, GetAllTasks
- HasCircularDependency uses recursive CTE to detect dependency cycles

# File Paths

- internal/storage/storage.go
