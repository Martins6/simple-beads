package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show [task-id]",
	Short: "Show task details",
	Long:  `Display detailed information about a specific task including its dependencies.`,
	Example: `  sb show sb-abc
  sb show sb-abc, sb-123`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]

		task, err := store.FindByID(taskID)
		if err != nil {
			return fmt.Errorf("task not found: %s", taskID)
		}

		fmt.Printf("ID:          %s\n", task.ID)
		fmt.Printf("Title:       %s\n", task.Title)
		if task.Description != "" {
			fmt.Printf("Description: %s\n", task.Description)
		}
		fmt.Printf("Priority:    P%d\n", task.Priority)
		fmt.Printf("Status:      %s\n", task.Status)
		if task.Parent != "" {
			fmt.Printf("Parent:      %s\n", task.Parent)
		}
		if len(task.Dependencies) > 0 {
			fmt.Printf("Dependencies: %s\n", strings.Join(task.Dependencies, ", "))
		}
		fmt.Printf("Created:     %s\n", task.CreatedAt.Format("2006-01-02 15:04"))
		fmt.Printf("Updated:     %s\n", task.UpdatedAt.Format("2006-01-02 15:04"))
		if task.Status == "closed" {
			fmt.Printf("Closed:      %s\n", task.ClosedAt.Format("2006-01-02 15:04"))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
