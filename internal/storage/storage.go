package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/Martins6/simple-beads/internal/models"
)

const (
	// DefaultDataDir is the default directory for storing tasks
	DefaultDataDir = ".sbeads"
	// DefaultDataFile is the default filename for tasks database
	DefaultDataFile = "tasks.db"
)

// Storage handles all task persistence operations using SQLite
type Storage struct {
	dbPath string
	db     *sql.DB
}

// NewStorage creates a new Storage instance
func NewStorage(dataDir string) *Storage {
	if dataDir == "" {
		dataDir = DefaultDataDir
	}
	return &Storage{
		dbPath: filepath.Join(dataDir, DefaultDataFile),
	}
}

// NewStorageWithPath creates a new Storage instance with a specific database path
func NewStorageWithPath(dbPath string) *Storage {
	return &Storage{
		dbPath: dbPath,
	}
}

// Init creates the data directory and initializes the SQLite database
func (s *Storage) Init() error {
	// Create data directory
	dir := filepath.Dir(s.dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Open database with WAL mode for better concurrency
	db, err := sql.Open("sqlite3", s.dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	s.db = db

	// Create tables
	if err := s.createTables(); err != nil {
		db.Close()
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

// createTables creates the database schema
func (s *Storage) createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS tasks (
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

	CREATE TABLE IF NOT EXISTS dependencies (
		task_id TEXT NOT NULL,
		depends_on TEXT NOT NULL,
		PRIMARY KEY (task_id, depends_on),
		FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
		FOREIGN KEY (depends_on) REFERENCES tasks(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
	CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
	CREATE INDEX IF NOT EXISTS idx_tasks_parent ON tasks(parent);
	CREATE INDEX IF NOT EXISTS idx_deps_task ON dependencies(task_id);
	CREATE INDEX IF NOT EXISTS idx_deps_depends ON dependencies(depends_on);
	`

	_, err := s.db.Exec(schema)
	return err
}

// Close closes the database connection
func (s *Storage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// LoadTasks loads all tasks from the database
func (s *Storage) LoadTasks() (map[string]*models.Task, error) {
	if s.db == nil {
		if err := s.Init(); err != nil {
			return nil, err
		}
	}

	rows, err := s.db.Query(`
		SELECT id, title, description, priority, status, parent, created_at, updated_at, closed_at
		FROM tasks
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	tasks := make(map[string]*models.Task)
	for rows.Next() {
		var task models.Task
		var closedAt sql.NullTime

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.Parent,
			&task.CreatedAt,
			&task.UpdatedAt,
			&closedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		if closedAt.Valid {
			task.ClosedAt = closedAt.Time
		}

		tasks[task.ID] = &task
	}

	// Load dependencies
	depRows, err := s.db.Query(`SELECT task_id, depends_on FROM dependencies`)
	if err != nil {
		return nil, fmt.Errorf("failed to query dependencies: %w", err)
	}
	defer depRows.Close()

	for depRows.Next() {
		var taskID, dependsOn string
		if err := depRows.Scan(&taskID, &dependsOn); err != nil {
			return nil, fmt.Errorf("failed to scan dependency: %w", err)
		}

		if task, exists := tasks[taskID]; exists {
			task.Dependencies = append(task.Dependencies, dependsOn)
		}
	}

	return tasks, nil
}

// SaveTask saves a new task to the database
func (s *Storage) SaveTask(task *models.Task) error {
	if s.db == nil {
		if err := s.Init(); err != nil {
			return err
		}
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert task
	_, err = tx.Exec(`
		INSERT INTO tasks (id, title, description, priority, status, parent, created_at, updated_at, closed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, task.ID, task.Title, task.Description, task.Priority, task.Status, task.Parent,
		task.CreatedAt, task.UpdatedAt, sql.NullTime{Time: task.ClosedAt, Valid: !task.ClosedAt.IsZero()})

	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	// Insert dependencies
	for _, depID := range task.Dependencies {
		_, err = tx.Exec(`
			INSERT INTO dependencies (task_id, depends_on)
			VALUES (?, ?)
		`, task.ID, depID)
		if err != nil {
			return fmt.Errorf("failed to insert dependency: %w", err)
		}
	}

	return tx.Commit()
}

// UpdateTask updates an existing task in the database
func (s *Storage) UpdateTask(task *models.Task) error {
	if s.db == nil {
		if err := s.Init(); err != nil {
			return err
		}
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update task
	result, err := tx.Exec(`
		UPDATE tasks
		SET title = ?, description = ?, priority = ?, status = ?, parent = ?, updated_at = ?, closed_at = ?
		WHERE id = ?
	`, task.Title, task.Description, task.Priority, task.Status, task.Parent,
		task.UpdatedAt, sql.NullTime{Time: task.ClosedAt, Valid: !task.ClosedAt.IsZero()},
		task.ID)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task %s not found", task.ID)
	}

	// Delete old dependencies
	_, err = tx.Exec(`DELETE FROM dependencies WHERE task_id = ?`, task.ID)
	if err != nil {
		return fmt.Errorf("failed to delete old dependencies: %w", err)
	}

	// Insert new dependencies
	for _, depID := range task.Dependencies {
		_, err = tx.Exec(`
			INSERT INTO dependencies (task_id, depends_on)
			VALUES (?, ?)
		`, task.ID, depID)
		if err != nil {
			return fmt.Errorf("failed to insert dependency: %w", err)
		}
	}

	return tx.Commit()
}

// DeleteTask deletes a task from the database
func (s *Storage) DeleteTask(taskID string) error {
	if s.db == nil {
		if err := s.Init(); err != nil {
			return err
		}
	}

	result, err := s.db.Exec(`DELETE FROM tasks WHERE id = ?`, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task %s not found", taskID)
	}

	return nil
}

// FindByID finds a task by its ID
func (s *Storage) FindByID(taskID string) (*models.Task, error) {
	if s.db == nil {
		if err := s.Init(); err != nil {
			return nil, err
		}
	}

	var task models.Task
	var closedAt sql.NullTime

	err := s.db.QueryRow(`
		SELECT id, title, description, priority, status, parent, created_at, updated_at, closed_at
		FROM tasks
		WHERE id = ?
	`, taskID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Priority,
		&task.Status,
		&task.Parent,
		&task.CreatedAt,
		&task.UpdatedAt,
		&closedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task %s not found", taskID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query task: %w", err)
	}

	if closedAt.Valid {
		task.ClosedAt = closedAt.Time
	}

	// Load dependencies
	depRows, err := s.db.Query(`SELECT depends_on FROM dependencies WHERE task_id = ?`, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to query dependencies: %w", err)
	}
	defer depRows.Close()

	for depRows.Next() {
		var depID string
		if err := depRows.Scan(&depID); err != nil {
			return nil, fmt.Errorf("failed to scan dependency: %w", err)
		}
		task.Dependencies = append(task.Dependencies, depID)
	}

	return &task, nil
}

// FindByParent finds all tasks with the given parent ID
func (s *Storage) FindByParent(parentID string) ([]*models.Task, error) {
	if s.db == nil {
		if err := s.Init(); err != nil {
			return nil, err
		}
	}

	rows, err := s.db.Query(`
		SELECT id, title, description, priority, status, parent, created_at, updated_at, closed_at
		FROM tasks
		WHERE parent = ?
	`, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var task models.Task
		var closedAt sql.NullTime

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.Parent,
			&task.CreatedAt,
			&task.UpdatedAt,
			&closedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		if closedAt.Valid {
			task.ClosedAt = closedAt.Time
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

// FindBlocked finds all tasks that are blocked by unclosed dependencies
func (s *Storage) FindBlocked() ([]*models.Task, error) {
	if s.db == nil {
		if err := s.Init(); err != nil {
			return nil, err
		}
	}

	rows, err := s.db.Query(`
		SELECT t.id, t.title, t.description, t.priority, t.status, t.parent, t.created_at, t.updated_at, t.closed_at
		FROM tasks t
		WHERE t.status = 'open'
		AND EXISTS (
			SELECT 1 FROM dependencies d
			JOIN tasks dep ON d.depends_on = dep.id
			WHERE d.task_id = t.id AND dep.status != 'closed'
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query blocked tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var task models.Task
		var closedAt sql.NullTime

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.Parent,
			&task.CreatedAt,
			&task.UpdatedAt,
			&closedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		if closedAt.Valid {
			task.ClosedAt = closedAt.Time
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

// FindReady finds all tasks that are open and not blocked
func (s *Storage) FindReady() ([]*models.Task, error) {
	if s.db == nil {
		if err := s.Init(); err != nil {
			return nil, err
		}
	}

	rows, err := s.db.Query(`
		SELECT t.id, t.title, t.description, t.priority, t.status, t.parent, t.created_at, t.updated_at, t.closed_at
		FROM tasks t
		WHERE t.status = 'open'
		AND NOT EXISTS (
			SELECT 1 FROM dependencies d
			JOIN tasks dep ON d.depends_on = dep.id
			WHERE d.task_id = t.id AND dep.status != 'closed'
		)
		ORDER BY t.priority ASC, t.created_at ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query ready tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var task models.Task
		var closedAt sql.NullTime

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.Parent,
			&task.CreatedAt,
			&task.UpdatedAt,
			&closedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		if closedAt.Valid {
			task.ClosedAt = closedAt.Time
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

// GetAllTasks returns all tasks as a slice
func (s *Storage) GetAllTasks() ([]*models.Task, error) {
	if s.db == nil {
		if err := s.Init(); err != nil {
			return nil, err
		}
	}

	rows, err := s.db.Query(`
		SELECT id, title, description, priority, status, parent, created_at, updated_at, closed_at
		FROM tasks
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var task models.Task
		var closedAt sql.NullTime

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.Parent,
			&task.CreatedAt,
			&task.UpdatedAt,
			&closedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		if closedAt.Valid {
			task.ClosedAt = closedAt.Time
		}

		tasks = append(tasks, &task)
	}

	// Load dependencies for all tasks
	for _, task := range tasks {
		depRows, err := s.db.Query(`SELECT depends_on FROM dependencies WHERE task_id = ?`, task.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to query dependencies: %w", err)
		}

		for depRows.Next() {
			var depID string
			if err := depRows.Scan(&depID); err != nil {
				depRows.Close()
				return nil, fmt.Errorf("failed to scan dependency: %w", err)
			}
			task.Dependencies = append(task.Dependencies, depID)
		}
		depRows.Close()
	}

	return tasks, nil
}

// HasCircularDependency checks if adding a dependency would create a cycle
func (s *Storage) HasCircularDependency(taskID, depID string) (bool, error) {
	if s.db == nil {
		if err := s.Init(); err != nil {
			return false, err
		}
	}

	// Use CTE (Common Table Expression) to traverse dependency graph
	query := `
		WITH RECURSIVE dependency_chain AS (
			-- Base case: direct dependencies of depID
			SELECT depends_on, task_id
			FROM dependencies
			WHERE task_id = ?
			
			UNION
			
			-- Recursive case: dependencies of dependencies
			SELECT d.depends_on, d.task_id
			FROM dependencies d
			JOIN dependency_chain dc ON d.task_id = dc.depends_on
		)
		SELECT COUNT(*) FROM dependency_chain WHERE depends_on = ?
	`

	var count int
	err := s.db.QueryRow(query, depID, taskID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check circular dependency: %w", err)
	}

	return count > 0, nil
}

// GetDataPath returns the current database file path
func (s *Storage) GetDataPath() string {
	return s.dbPath
}
