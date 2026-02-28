package utils

import (
	"fmt"

	"github.com/user/sbeads/internal/models"
)

// DependencyChecker provides utility functions for checking task dependencies
type DependencyChecker struct {
	tasks map[string]*models.Task
}

// NewDependencyChecker creates a new DependencyChecker
func NewDependencyChecker(tasks map[string]*models.Task) *DependencyChecker {
	return &DependencyChecker{tasks: tasks}
}

// IsBlocked checks if a task is blocked by unclosed dependencies
func (dc *DependencyChecker) IsBlocked(taskID string) (bool, error) {
	task, exists := dc.tasks[taskID]
	if !exists {
		return false, fmt.Errorf("task %s not found", taskID)
	}

	return task.IsBlocked(dc.tasks), nil
}

// GetBlockingDependencies returns the list of dependencies that are blocking a task
func (dc *DependencyChecker) GetBlockingDependencies(taskID string) ([]string, error) {
	task, exists := dc.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task %s not found", taskID)
	}

	var blocking []string
	for _, depID := range task.Dependencies {
		if depTask, exists := dc.tasks[depID]; exists {
			if depTask.Status != models.StatusClosed {
				blocking = append(blocking, depID)
			}
		}
	}

	return blocking, nil
}

// HasCircularDependency checks if adding a dependency would create a cycle
func (dc *DependencyChecker) HasCircularDependency(taskID, depID string) bool {
	visited := make(map[string]bool)
	var check func(string) bool
	check = func(currentID string) bool {
		if currentID == taskID {
			return true // Found a cycle
		}
		if visited[currentID] {
			return false
		}
		visited[currentID] = true

		task, exists := dc.tasks[currentID]
		if !exists {
			return false
		}

		for _, dep := range task.Dependencies {
			if check(dep) {
				return true
			}
		}
		return false
	}

	return check(depID)
}

// GetDependencyChain returns the full chain of dependencies for a task
func (dc *DependencyChecker) GetDependencyChain(taskID string) ([]string, error) {
	task, exists := dc.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task %s not found", taskID)
	}

	var chain []string
	visited := make(map[string]bool)

	var collect func(string)
	collect = func(currentID string) {
		if visited[currentID] {
			return
		}
		visited[currentID] = true

		currentTask, exists := dc.tasks[currentID]
		if !exists {
			return
		}

		for _, depID := range currentTask.Dependencies {
			chain = append(chain, depID)
			collect(depID)
		}
	}

	// Collect all dependencies (direct and indirect)
	for _, depID := range task.Dependencies {
		chain = append(chain, depID) // Add direct dependency
		collect(depID)               // Recursively collect indirect dependencies
	}

	return chain, nil
}

// GetDependents returns all tasks that depend on the given task
func (dc *DependencyChecker) GetDependents(taskID string) []string {
	var dependents []string
	for id, task := range dc.tasks {
		if task.HasDependency(taskID) {
			dependents = append(dependents, id)
		}
	}
	return dependents
}

// ValidateDependency checks if a dependency can be added
func (dc *DependencyChecker) ValidateDependency(taskID, depID string) error {
	// Check if both tasks exist
	if _, exists := dc.tasks[taskID]; !exists {
		return fmt.Errorf("task %s not found", taskID)
	}
	if _, exists := dc.tasks[depID]; !exists {
		return fmt.Errorf("dependency task %s not found", depID)
	}

	// Check for self-dependency
	if taskID == depID {
		return fmt.Errorf("task cannot depend on itself")
	}

	// Check for circular dependency
	if dc.HasCircularDependency(taskID, depID) {
		return fmt.Errorf("adding this dependency would create a circular dependency")
	}

	return nil
}
