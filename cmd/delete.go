package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [task-id...]",
	Short: "Delete one or more tasks",
	Long: `Delete one or more tasks permanently.

Warning: This action cannot be undone.`,
	Example: `  sb delete sb-abc
  sb delete sb-abc sb-123`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, taskID := range args {
			if err := store.DeleteTask(taskID); err != nil {
				fmt.Printf("✗ Failed to delete task %s: %v\n", taskID, err)
				continue
			}
			fmt.Printf("✓ Deleted task: %s\n", taskID)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
