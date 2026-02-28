package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Status represents the task status
type Status string

const (
	StatusOpen   Status = "open"
	StatusClosed Status = "closed"
)

// Priority represents task priority (0-4, where 0 is highest)
type Priority int

const (
	PriorityP0 Priority = 0 // Highest
	PriorityP1 Priority = 1
	PriorityP2 Priority = 2
	PriorityP3 Priority = 3
	PriorityP4 Priority = 4 // Lowest
)

// Task represents a single task in the system
type Task struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	Priority     Priority  `json:"priority"`
	Status       Status    `json:"status"`
	Parent       string    `json:"parent,omitempty"`
	Dependencies []string  `json:"dependencies,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ClosedAt     time.Time `json:"closed_at,omitempty"`
}

// NewTask creates a new task with the given title
func NewTask(title string) *Task {
	now := time.Now()
	return &Task{
		ID:           GenerateID(),
		Title:        title,
		Priority:     PriorityP2, // Default priority
		Status:       StatusOpen,
		Dependencies: []string{},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// GenerateID generates a unique task ID in format sb-XXXX
func GenerateID() string {
	// Generate short unique ID (first 4 chars of UUID)
	id := uuid.New().String()[:4]
	return fmt.Sprintf("sb-%s", id)
}

// IsBlocked checks if this task is blocked by unclosed dependencies
func (t *Task) IsBlocked(tasks map[string]*Task) bool {
	for _, depID := range t.Dependencies {
		if depTask, exists := tasks[depID]; exists {
			if depTask.Status != StatusClosed {
				return true
			}
		}
	}
	return false
}

// Close marks the task as closed
func (t *Task) Close() {
	t.Status = StatusClosed
	t.ClosedAt = time.Now()
	t.UpdatedAt = time.Now()
}

// Reopen marks the task as open
func (t *Task) Reopen() {
	t.Status = StatusOpen
	t.ClosedAt = time.Time{}
	t.UpdatedAt = time.Now()
}

// AddDependency adds a dependency to the task
func (t *Task) AddDependency(depID string) error {
	// Check if dependency already exists
	for _, existing := range t.Dependencies {
		if existing == depID {
			return fmt.Errorf("dependency %s already exists", depID)
		}
	}
	t.Dependencies = append(t.Dependencies, depID)
	t.UpdatedAt = time.Now()
	return nil
}

// RemoveDependency removes a dependency from the task
func (t *Task) RemoveDependency(depID string) error {
	for i, dep := range t.Dependencies {
		if dep == depID {
			// Remove the dependency
			t.Dependencies = append(t.Dependencies[:i], t.Dependencies[i+1:]...)
			t.UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("dependency %s not found", depID)
}

// HasDependency checks if task has a specific dependency
func (t *Task) HasDependency(depID string) bool {
	for _, dep := range t.Dependencies {
		if dep == depID {
			return true
		}
	}
	return false
}

// Validate checks if the task is valid
func (t *Task) Validate() error {
	if t.ID == "" {
		return fmt.Errorf("task ID is required")
	}
	if t.Title == "" {
		return fmt.Errorf("task title is required")
	}
	if t.Priority < PriorityP0 || t.Priority > PriorityP4 {
		return fmt.Errorf("priority must be between 0 and 4")
	}
	if t.Status != StatusOpen && t.Status != StatusClosed {
		return fmt.Errorf("invalid status: %s", t.Status)
	}
	return nil
}

// String returns a string representation of the task
func (t *Task) String() string {
	return fmt.Sprintf("[%s] %s (P%d) - %s", t.ID, t.Title, t.Priority, t.Status)
}
