package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/sbeads/internal/models"
)

var (
	updateTitle       string
	updateDescription string
	updatePriority    int
	updateParent      string
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update [task-id]",
	Short: "Update a task",
	Long:  `Update one or more fields of an existing task.`,
	Example: `  sb update sb-abc --title "New title"
  sb update sb-abc -p 1
  sb update sb-abc --parent sb-def`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]

		task, err := store.FindByID(taskID)
		if err != nil {
			return fmt.Errorf("task not found: %s", taskID)
		}

		// Update fields only if flags are provided
		if cmd.Flags().Changed("title") {
			task.Title = updateTitle
		}
		if cmd.Flags().Changed("description") {
			task.Description = updateDescription
		}
		if cmd.Flags().Changed("priority") {
			task.Priority = models.Priority(updatePriority)
		}
		if cmd.Flags().Changed("parent") {
			task.Parent = updateParent
		}

		// Validate updated task
		if err := task.Validate(); err != nil {
			return fmt.Errorf("invalid task: %w", err)
		}

		// Save updated task
		if err := store.UpdateTask(task); err != nil {
			return fmt.Errorf("failed to update task: %w", err)
		}

		fmt.Printf("✓ Updated task: %s\n", task.String())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVar(&updateTitle, "title", "", "New task title")
	updateCmd.Flags().StringVarP(&updateDescription, "description", "d", "", "New task description")
	updateCmd.Flags().IntVarP(&updatePriority, "priority", "p", 2, "New task priority (0-4)")
	updateCmd.Flags().StringVar(&updateParent, "parent", "", "New parent task ID")
}
