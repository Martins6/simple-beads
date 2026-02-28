package utils

import (
	"testing"

	"github.com/user/sbeads/internal/models"
)

func TestNewDependencyChecker(t *testing.T) {
	tasks := make(map[string]*models.Task)
	dc := NewDependencyChecker(tasks)

	if dc == nil {
		t.Error("Expected non-nil DependencyChecker")
	}

	if dc.tasks == nil {
		t.Error("Expected tasks map to be initialized")
	}
}

func TestDependencyCheckerIsBlocked(t *testing.T) {
	tasks := make(map[string]*models.Task)

	// Create open dependency
	depOpen := models.NewTask("Open Dependency")
	tasks[depOpen.ID] = depOpen

	// Create closed dependency
	depClosed := models.NewTask("Closed Dependency")
	depClosed.Close()
	tasks[depClosed.ID] = depClosed

	dc := NewDependencyChecker(tasks)

	// Test task with open dependency - should be blocked
	taskBlocked := models.NewTask("Blocked Task")
	taskBlocked.AddDependency(depOpen.ID)
	tasks[taskBlocked.ID] = taskBlocked

	blocked, err := dc.IsBlocked(taskBlocked.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !blocked {
		t.Error("Expected task to be blocked")
	}

	// Test task with closed dependency - should not be blocked
	taskUnblocked := models.NewTask("Unblocked Task")
	taskUnblocked.AddDependency(depClosed.ID)
	tasks[taskUnblocked.ID] = taskUnblocked

	blocked, err = dc.IsBlocked(taskUnblocked.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if blocked {
		t.Error("Expected task to not be blocked")
	}

	// Test non-existent task
	_, err = dc.IsBlocked("sb-nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent task")
	}
}

func TestGetBlockingDependencies(t *testing.T) {
	tasks := make(map[string]*models.Task)

	// Create dependencies
	dep1 := models.NewTask("Dep 1")
	dep2 := models.NewTask("Dep 2")
	dep2.Close()
	tasks[dep1.ID] = dep1
	tasks[dep2.ID] = dep2

	task := models.NewTask("Task")
	task.AddDependency(dep1.ID)
	task.AddDependency(dep2.ID)
	tasks[task.ID] = task

	dc := NewDependencyChecker(tasks)

	blocking, err := dc.GetBlockingDependencies(task.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(blocking) != 1 {
		t.Errorf("Expected 1 blocking dependency, got %d", len(blocking))
	}

	if len(blocking) > 0 && blocking[0] != dep1.ID {
		t.Errorf("Expected blocking dep '%s', got '%s'", dep1.ID, blocking[0])
	}
}

func TestHasCircularDependency(t *testing.T) {
	tasks := make(map[string]*models.Task)

	// Create chain: A -> B -> C
	taskA := models.NewTask("Task A")
	taskB := models.NewTask("Task B")
	taskC := models.NewTask("Task C")

	taskA.AddDependency(taskB.ID)
	taskB.AddDependency(taskC.ID)

	tasks[taskA.ID] = taskA
	tasks[taskB.ID] = taskB
	tasks[taskC.ID] = taskC

	dc := NewDependencyChecker(tasks)

	// C -> A would create cycle
	if !dc.HasCircularDependency(taskC.ID, taskA.ID) {
		t.Error("Expected cycle detection for C -> A")
	}

	// A -> C would not create cycle
	if dc.HasCircularDependency(taskA.ID, taskC.ID) {
		t.Error("Did not expect cycle for A -> C")
	}

	// B -> A would create cycle
	if !dc.HasCircularDependency(taskB.ID, taskA.ID) {
		t.Error("Expected cycle detection for B -> A")
	}
}

func TestGetDependencyChain(t *testing.T) {
	tasks := make(map[string]*models.Task)

	// Create chain: A -> B -> C
	taskA := models.NewTask("Task A")
	taskB := models.NewTask("Task B")
	taskC := models.NewTask("Task C")

	taskA.AddDependency(taskB.ID)
	taskB.AddDependency(taskC.ID)

	tasks[taskA.ID] = taskA
	tasks[taskB.ID] = taskB
	tasks[taskC.ID] = taskC

	dc := NewDependencyChecker(tasks)

	chain, err := dc.GetDependencyChain(taskA.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Should have B and C in chain
	if len(chain) != 2 {
		t.Errorf("Expected 2 dependencies in chain, got %d", len(chain))
	}

	// Check non-existent task
	_, err = dc.GetDependencyChain("sb-nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent task")
	}
}

func TestGetDependents(t *testing.T) {
	tasks := make(map[string]*models.Task)

	taskA := models.NewTask("Task A")
	taskB := models.NewTask("Task B")
	taskC := models.NewTask("Task C")

	// B depends on A, C depends on A
	taskB.AddDependency(taskA.ID)
	taskC.AddDependency(taskA.ID)

	tasks[taskA.ID] = taskA
	tasks[taskB.ID] = taskB
	tasks[taskC.ID] = taskC

	dc := NewDependencyChecker(tasks)

	dependents := dc.GetDependents(taskA.ID)

	if len(dependents) != 2 {
		t.Errorf("Expected 2 dependents, got %d", len(dependents))
	}

	// Task with no dependents
	dependents = dc.GetDependents(taskB.ID)
	if len(dependents) != 0 {
		t.Errorf("Expected 0 dependents, got %d", len(dependents))
	}
}

func TestValidateDependency(t *testing.T) {
	tasks := make(map[string]*models.Task)

	taskA := models.NewTask("Task A")
	taskB := models.NewTask("Task B")
	taskC := models.NewTask("Task C")

	// A -> B already exists
	taskA.AddDependency(taskB.ID)

	tasks[taskA.ID] = taskA
	tasks[taskB.ID] = taskB
	tasks[taskC.ID] = taskC

	dc := NewDependencyChecker(tasks)

	// Valid: A -> C
	err := dc.ValidateDependency(taskA.ID, taskC.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Invalid: non-existent task
	err = dc.ValidateDependency("sb-nonexistent", taskB.ID)
	if err == nil {
		t.Error("Expected error for non-existent task")
	}

	// Invalid: non-existent dependency
	err = dc.ValidateDependency(taskA.ID, "sb-nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent dependency")
	}

	// Invalid: self-dependency
	err = dc.ValidateDependency(taskA.ID, taskA.ID)
	if err == nil {
		t.Error("Expected error for self-dependency")
	}

	// Invalid: would create cycle (C -> A, and A -> B, so C -> A -> B, not a cycle yet)
	// But A -> C would create cycle since C could depend on something that depends on A
	// Actually C -> A is fine, A -> C would create cycle
	err = dc.ValidateDependency(taskC.ID, taskA.ID)
	if err != nil {
		t.Errorf("C -> A should be valid: %v", err)
	}
}
