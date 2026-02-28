package models

import (
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	task := NewTask("Test Task")

	if task.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got '%s'", task.Title)
	}

	if task.Status != StatusOpen {
		t.Errorf("Expected status 'open', got '%s'", task.Status)
	}

	if task.Priority != PriorityP2 {
		t.Errorf("Expected default priority P2, got P%d", task.Priority)
	}

	if task.ID == "" {
		t.Error("Expected non-empty ID")
	}

	if len(task.Dependencies) != 0 {
		t.Errorf("Expected empty dependencies, got %v", task.Dependencies)
	}
}

func TestGenerateID(t *testing.T) {
	id1 := GenerateID()
	id2 := GenerateID()

	if id1 == id2 {
		t.Error("Generated IDs should be unique")
	}

	if len(id1) != 7 { // "sb-" + 4 chars
		t.Errorf("Expected ID length 7, got %d", len(id1))
	}

	if id1[:3] != "sb-" {
		t.Errorf("Expected ID to start with 'sb-', got '%s'", id1[:3])
	}
}

func TestTaskValidate(t *testing.T) {
	tests := []struct {
		name    string
		task    *Task
		wantErr bool
	}{
		{
			name:    "valid task",
			task:    NewTask("Valid Task"),
			wantErr: false,
		},
		{
			name: "missing ID",
			task: &Task{
				ID:    "",
				Title: "Task",
			},
			wantErr: true,
		},
		{
			name: "missing title",
			task: &Task{
				ID:    "sb-test",
				Title: "",
			},
			wantErr: true,
		},
		{
			name: "invalid priority low",
			task: &Task{
				ID:       "sb-test",
				Title:    "Task",
				Priority: -1,
			},
			wantErr: true,
		},
		{
			name: "invalid priority high",
			task: &Task{
				ID:       "sb-test",
				Title:    "Task",
				Priority: 5,
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			task: &Task{
				ID:     "sb-test",
				Title:  "Task",
				Status: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaskClose(t *testing.T) {
	task := NewTask("Test Task")
	originalUpdatedAt := task.UpdatedAt

	time.Sleep(10 * time.Millisecond) // Ensure time difference
	task.Close()

	if task.Status != StatusClosed {
		t.Errorf("Expected status 'closed', got '%s'", task.Status)
	}

	if task.ClosedAt.IsZero() {
		t.Error("Expected ClosedAt to be set")
	}

	if !task.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}

func TestTaskReopen(t *testing.T) {
	task := NewTask("Test Task")
	task.Close()
	task.Reopen()

	if task.Status != StatusOpen {
		t.Errorf("Expected status 'open', got '%s'", task.Status)
	}

	if !task.ClosedAt.IsZero() {
		t.Error("Expected ClosedAt to be cleared")
	}
}

func TestTaskAddDependency(t *testing.T) {
	task := NewTask("Test Task")

	err := task.AddDependency("sb-123")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(task.Dependencies) != 1 {
		t.Errorf("Expected 1 dependency, got %d", len(task.Dependencies))
	}

	if task.Dependencies[0] != "sb-123" {
		t.Errorf("Expected dependency 'sb-123', got '%s'", task.Dependencies[0])
	}

	// Test adding duplicate
	err = task.AddDependency("sb-123")
	if err == nil {
		t.Error("Expected error for duplicate dependency")
	}
}

func TestTaskRemoveDependency(t *testing.T) {
	task := NewTask("Test Task")
	task.AddDependency("sb-123")
	task.AddDependency("sb-456")

	err := task.RemoveDependency("sb-123")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(task.Dependencies) != 1 {
		t.Errorf("Expected 1 dependency, got %d", len(task.Dependencies))
	}

	// Test removing non-existent
	err = task.RemoveDependency("sb-999")
	if err == nil {
		t.Error("Expected error for non-existent dependency")
	}
}

func TestTaskHasDependency(t *testing.T) {
	task := NewTask("Test Task")
	task.AddDependency("sb-123")

	if !task.HasDependency("sb-123") {
		t.Error("Expected HasDependency to return true for existing dependency")
	}

	if task.HasDependency("sb-999") {
		t.Error("Expected HasDependency to return false for non-existent dependency")
	}
}

func TestTaskIsBlocked(t *testing.T) {
	tasks := make(map[string]*Task)

	// Create a closed dependency
	dep1 := NewTask("Dependency 1")
	dep1.Close()
	tasks[dep1.ID] = dep1

	// Create an open dependency
	dep2 := NewTask("Dependency 2")
	tasks[dep2.ID] = dep2

	// Create task with closed dependency only
	task1 := NewTask("Task 1")
	task1.AddDependency(dep1.ID)
	tasks[task1.ID] = task1

	if task1.IsBlocked(tasks) {
		t.Error("Expected task1 to not be blocked (all deps closed)")
	}

	// Create task with open dependency
	task2 := NewTask("Task 2")
	task2.AddDependency(dep2.ID)
	tasks[task2.ID] = task2

	if !task2.IsBlocked(tasks) {
		t.Error("Expected task2 to be blocked (has open dep)")
	}

	// Create task with both dependencies
	task3 := NewTask("Task 3")
	task3.AddDependency(dep1.ID)
	task3.AddDependency(dep2.ID)
	tasks[task3.ID] = task3

	if !task3.IsBlocked(tasks) {
		t.Error("Expected task3 to be blocked (one dep open)")
	}
}

func TestTaskString(t *testing.T) {
	task := NewTask("Test Task")
	task.Priority = PriorityP1

	str := task.String()
	expected := "[" + task.ID + "] Test Task (P1) - open"
	if str != expected {
		t.Errorf("Expected '%s', got '%s'", expected, str)
	}
}
