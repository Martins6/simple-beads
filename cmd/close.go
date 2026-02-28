package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// closeCmd represents the close command
var closeCmd = &cobra.Command{
	Use:   "close [task-id...]",
	Short: "Close one or more tasks",
	Long:  `Mark one or more tasks as closed/completed.`,
	Example: `  sb close sb-abc
  sb close sb-abc sb-123 sb-456`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, taskID := range args {
			task, err := store.FindByID(taskID)
			if err != nil {
				fmt.Printf("✗ Task not found: %s\n", taskID)
				continue
			}

			if task.Status == "closed" {
				fmt.Printf("⚠ Task already closed: %s\n", taskID)
				continue
			}

			task.Close()
			if err := store.UpdateTask(task); err != nil {
				fmt.Printf("✗ Failed to close task %s: %v\n", taskID, err)
				continue
			}

			fmt.Printf("✓ Closed task: %s\n", taskID)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(closeCmd)
}
