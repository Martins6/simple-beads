package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// depCmd represents the dep command
var depCmd = &cobra.Command{
	Use:   "dep",
	Short: "Manage task dependencies",
	Long:  `Add or remove dependencies between tasks.`,
}

// depAddCmd represents the dep add command
var depAddCmd = &cobra.Command{
	Use:     "add [task-id] [dependency-id]",
	Short:   "Add a dependency to a task",
	Long:    `Add a dependency to a task. The task will be blocked until the dependency is closed.`,
	Example: `  sb dep add sb-abc sb-123`,
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]
		depID := args[1]

		// Check if task exists
		task, err := store.FindByID(taskID)
		if err != nil {
			return fmt.Errorf("task not found: %s", taskID)
		}

		// Check if dependency exists
		if _, err := store.FindByID(depID); err != nil {
			return fmt.Errorf("dependency task not found: %s", depID)
		}

		// Check for self-dependency
		if taskID == depID {
			return fmt.Errorf("task cannot depend on itself")
		}

		// Check for circular dependency
		hasCycle, err := store.HasCircularDependency(taskID, depID)
		if err != nil {
			return fmt.Errorf("failed to check for circular dependency: %w", err)
		}
		if hasCycle {
			return fmt.Errorf("adding this dependency would create a circular dependency")
		}

		// Add dependency
		if err := task.AddDependency(depID); err != nil {
			return fmt.Errorf("failed to add dependency: %w", err)
		}

		// Save task
		if err := store.UpdateTask(task); err != nil {
			return fmt.Errorf("failed to update task: %w", err)
		}

		fmt.Printf("✓ Added dependency: %s now depends on %s\n", taskID, depID)
		return nil
	},
}

// depRemoveCmd represents the dep remove command
var depRemoveCmd = &cobra.Command{
	Use:     "remove [task-id] [dependency-id]",
	Short:   "Remove a dependency from a task",
	Long:    `Remove a dependency from a task.`,
	Example: `  sb dep remove sb-abc sb-123`,
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]
		depID := args[1]

		// Check if task exists
		task, err := store.FindByID(taskID)
		if err != nil {
			return fmt.Errorf("task not found: %s", taskID)
		}

		// Remove dependency
		if err := task.RemoveDependency(depID); err != nil {
			return fmt.Errorf("failed to remove dependency: %w", err)
		}

		// Save task
		if err := store.UpdateTask(task); err != nil {
			return fmt.Errorf("failed to update task: %w", err)
		}

		fmt.Printf("✓ Removed dependency: %s no longer depends on %s\n", taskID, depID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(depCmd)
	depCmd.AddCommand(depAddCmd)
	depCmd.AddCommand(depRemoveCmd)
}
