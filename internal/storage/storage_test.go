package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Martins6/simple-beads/internal/models"
)

func setupTestStorage(t *testing.T) (*Storage, func()) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	store := NewStorageWithPath(dbPath)

	if err := store.Init(); err != nil {
		t.Fatalf("Failed to init storage: %v", err)
	}

	cleanup := func() {
		store.Close()
	}

	return store, cleanup
}

func TestNewStorage(t *testing.T) {
	store := NewStorage("")
	if store.GetDataPath() != ".sbeads/tasks.db" {
		t.Errorf("Expected default path '.sbeads/tasks.db', got '%s'", store.GetDataPath())
	}

	store = NewStorage("/custom/path")
	expected := filepath.Join("/custom/path", "tasks.db")
	if store.GetDataPath() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, store.GetDataPath())
	}
}

func TestStorageInit(t *testing.T) {
	store, cleanup := setupTestStorage(t)
	defer cleanup()

	// Check database was created
	if _, err := os.Stat(store.GetDataPath()); os.IsNotExist(err) {
		t.Error("Expected database file to be created")
	}
}

func TestSaveAndLoadTasks(t *testing.T) {
	store, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create and save tasks
	task1 := models.NewTask("Task 1")
	task1.Priority = models.PriorityP0
	if err := store.SaveTask(task1); err != nil {
		t.Errorf("Failed to save task: %v", err)
	}

	task2 := models.NewTask("Task 2")
	task2.Priority = models.PriorityP1
	if err := store.SaveTask(task2); err != nil {
		t.Errorf("Failed to save task: %v", err)
	}

	// Load tasks
	tasks, err := store.LoadTasks()
	if err != nil {
		t.Errorf("Failed to load tasks: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	// Verify tasks
	if tasks[task1.ID] == nil {
		t.Error("Task 1 not found")
	}
	if tasks[task2.ID] == nil {
		t.Error("Task 2 not found")
	}
}

func TestFindByID(t *testing.T) {
	store, cleanup := setupTestStorage(t)
	defer cleanup()

	task := models.NewTask("Test Task")
	if err := store.SaveTask(task); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}

	// Find existing task
	found, err := store.FindByID(task.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if found.ID != task.ID {
		t.Errorf("Expected task ID '%s', got '%s'", task.ID, found.ID)
	}

	// Find non-existent task
	_, err = store.FindByID("sb-nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent task")
	}
}

func TestUpdateTask(t *testing.T) {
	store, cleanup := setupTestStorage(t)
	defer cleanup()

	task := models.NewTask("Original Title")
	if err := store.SaveTask(task); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}

	// Update task
	task.Title = "Updated Title"
	task.Priority = models.PriorityP0
	if err := store.UpdateTask(task); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify update
	updated, err := store.FindByID(task.ID)
	if err != nil {
		t.Fatalf("Failed to find task: %v", err)
	}
	if updated.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got '%s'", updated.Title)
	}
	if updated.Priority != models.PriorityP0 {
		t.Errorf("Expected priority P0, got P%d", updated.Priority)
	}
}

func TestDeleteTask(t *testing.T) {
	store, cleanup := setupTestStorage(t)
	defer cleanup()

	task := models.NewTask("Task to Delete")
	if err := store.SaveTask(task); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}

	// Delete task
	if err := store.DeleteTask(task.ID); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify deletion
	_, err := store.FindByID(task.ID)
	if err == nil {
		t.Error("Expected error for deleted task")
	}

	// Delete non-existent task
	err = store.DeleteTask("sb-nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent task")
	}
}

func TestFindByParent(t *testing.T) {
	store, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create parent task
	parent := models.NewTask("Parent Task")
	if err := store.SaveTask(parent); err != nil {
		t.Fatalf("Failed to save parent: %v", err)
	}

	// Create child tasks
	child1 := models.NewTask("Child 1")
	child1.Parent = parent.ID
	if err := store.SaveTask(child1); err != nil {
		t.Fatalf("Failed to save child: %v", err)
	}

	child2 := models.NewTask("Child 2")
	child2.Parent = parent.ID
	if err := store.SaveTask(child2); err != nil {
		t.Fatalf("Failed to save child: %v", err)
	}

	// Create unrelated task
	other := models.NewTask("Other Task")
	if err := store.SaveTask(other); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}

	// Find children
	children, err := store.FindByParent(parent.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(children))
	}
}

func TestFindBlocked(t *testing.T) {
	store, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create open dependency
	depOpen := models.NewTask("Open Dependency")
	if err := store.SaveTask(depOpen); err != nil {
		t.Fatalf("Failed to save dep: %v", err)
	}

	// Create closed dependency
	depClosed := models.NewTask("Closed Dependency")
	depClosed.Close()
	if err := store.SaveTask(depClosed); err != nil {
		t.Fatalf("Failed to save dep: %v", err)
	}

	// Create blocked task
	blocked := models.NewTask("Blocked Task")
	blocked.Dependencies = []string{depOpen.ID}
	if err := store.SaveTask(blocked); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}

	// Create unblocked task
	unblocked := models.NewTask("Unblocked Task")
	unblocked.Dependencies = []string{depClosed.ID}
	if err := store.SaveTask(unblocked); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}

	// Create task with no dependencies
	nodeps := models.NewTask("No Dependencies")
	if err := store.SaveTask(nodeps); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}

	// Find blocked
	blockedTasks, err := store.FindBlocked()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(blockedTasks) != 1 {
		t.Errorf("Expected 1 blocked task, got %d", len(blockedTasks))
	}

	if len(blockedTasks) > 0 && blockedTasks[0].ID != blocked.ID {
		t.Errorf("Expected blocked task '%s', got '%s'", blocked.ID, blockedTasks[0].ID)
	}
}

func TestFindReady(t *testing.T) {
	store, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create closed task (should not be in ready)
	closed := models.NewTask("Closed Task")
	closed.Close()
	if err := store.SaveTask(closed); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}

	// Create open task with no dependencies
	ready1 := models.NewTask("Ready Task 1")
	if err := store.SaveTask(ready1); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}

	// Create blocked task (should not be in ready)
	dep := models.NewTask("Dependency")
	if err := store.SaveTask(dep); err != nil {
		t.Fatalf("Failed to save dep: %v", err)
	}
	blocked := models.NewTask("Blocked Task")
	blocked.Dependencies = []string{dep.ID}
	if err := store.SaveTask(blocked); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}

	// Create another ready task
	ready2 := models.NewTask("Ready Task 2")
	if err := store.SaveTask(ready2); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}

	// Close the dependency to unblock the blocked task
	dep.Close()
	if err := store.UpdateTask(dep); err != nil {
		t.Fatalf("Failed to close dep: %v", err)
	}

	// Find ready
	readyTasks, err := store.FindReady()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(readyTasks) != 3 {
		t.Errorf("Expected 3 ready tasks, got %d", len(readyTasks))
	}
}

func TestHasCircularDependency(t *testing.T) {
	store, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create tasks
	taskA := models.NewTask("Task A")
	taskB := models.NewTask("Task B")
	taskC := models.NewTask("Task C")

	if err := store.SaveTask(taskA); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}
	if err := store.SaveTask(taskB); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}
	if err := store.SaveTask(taskC); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}

	// A -> B (A depends on B)
	taskA.Dependencies = []string{taskB.ID}
	if err := store.UpdateTask(taskA); err != nil {
		t.Fatalf("Failed to update task: %v", err)
	}

	// B -> C (B depends on C)
	taskB.Dependencies = []string{taskC.ID}
	if err := store.UpdateTask(taskB); err != nil {
		t.Fatalf("Failed to update task: %v", err)
	}

	// Check: C -> A would create cycle (C -> A -> B -> C)
	hasCycle, err := store.HasCircularDependency(taskC.ID, taskA.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !hasCycle {
		t.Error("Expected cycle detection for C -> A")
	}

	// Check: A -> C would not create cycle
	hasCycle, err = store.HasCircularDependency(taskA.ID, taskC.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if hasCycle {
		t.Error("Did not expect cycle for A -> C")
	}
}

func TestGetAllTasks(t *testing.T) {
	store, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create tasks
	for i := 0; i < 3; i++ {
		task := models.NewTask("Task")
		if err := store.SaveTask(task); err != nil {
			t.Fatalf("Failed to save task: %v", err)
		}
	}

	// Get all tasks
	tasks, err := store.GetAllTasks()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}
}

func TestDependencies(t *testing.T) {
	store, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create tasks
	taskA := models.NewTask("Task A")
	taskB := models.NewTask("Task B")
	taskC := models.NewTask("Task C")

	if err := store.SaveTask(taskA); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}
	if err := store.SaveTask(taskB); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}
	if err := store.SaveTask(taskC); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}

	// Add multiple dependencies
	taskA.Dependencies = []string{taskB.ID, taskC.ID}
	if err := store.UpdateTask(taskA); err != nil {
		t.Fatalf("Failed to update task: %v", err)
	}

	// Verify dependencies were saved
	found, err := store.FindByID(taskA.ID)
	if err != nil {
		t.Fatalf("Failed to find task: %v", err)
	}

	if len(found.Dependencies) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(found.Dependencies))
	}

	// Remove a dependency
	taskA.Dependencies = []string{taskB.ID}
	if err := store.UpdateTask(taskA); err != nil {
		t.Fatalf("Failed to update task: %v", err)
	}

	// Verify dependency was removed
	found, err = store.FindByID(taskA.ID)
	if err != nil {
		t.Fatalf("Failed to find task: %v", err)
	}

	if len(found.Dependencies) != 1 {
		t.Errorf("Expected 1 dependency after removal, got %d", len(found.Dependencies))
	}
}

func TestConcurrentAccess(t *testing.T) {
	store, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create initial task
	task := models.NewTask("Concurrent Test")
	if err := store.SaveTask(task); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}

	// Simulate concurrent reads
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := store.FindByID(task.ID)
			if err != nil {
				t.Errorf("Concurrent read failed: %v", err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}
