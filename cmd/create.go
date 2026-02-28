package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/sbeads/internal/models"
)

var (
	createDescription string
	createPriority    int
	createParent      string
	createDeps        []string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new task",
	Long: `Create a new task with the given title.

You can specify additional details like description, priority, parent task,
and dependencies using flags.`,
	Example: `  sb create "Fix login bug"
  sb create "Update docs" -d "Add API examples" -p 1
  sb create "Subtask" --parent sb-abc
  sb create "Blocked task" --deps sb-123,sb-456`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := args[0]

		// Create new task
		task := models.NewTask(title)
		task.Description = createDescription
		task.Priority = models.Priority(createPriority)
		task.Parent = createParent

		// Parse dependencies
		if len(createDeps) > 0 {
			// Handle comma-separated dependencies
			for _, dep := range createDeps {
				parts := strings.Split(dep, ",")
				for _, p := range parts {
					p = strings.TrimSpace(p)
					if p != "" {
						task.Dependencies = append(task.Dependencies, p)
					}
				}
			}
		}

		// Validate task
		if err := task.Validate(); err != nil {
			return fmt.Errorf("invalid task: %w", err)
		}

		// Save task
		if err := store.SaveTask(task); err != nil {
			return fmt.Errorf("failed to save task: %w", err)
		}

		fmt.Printf("✓ Created task: %s\n", task.String())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&createDescription, "description", "d", "", "Task description")
	createCmd.Flags().IntVarP(&createPriority, "priority", "p", 2, "Task priority (0-4, where 0 is highest)")
	createCmd.Flags().StringVar(&createParent, "parent", "", "Parent task ID")
	createCmd.Flags().StringSliceVar(&createDeps, "deps", []string{}, "Comma-separated list of dependency task IDs")
}
