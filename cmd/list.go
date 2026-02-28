package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/Martins6/simple-beads/internal/models"
)

var (
	listAll      bool
	listPriority int
	listStatus   string
	listParent   string
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Long: `List tasks with optional filtering.

By default, only open tasks are shown. Use --all to show all tasks including closed ones.`,
	Example: `  sb list
  sb list --all
  sb list -p 1
  sb list --status closed
  sb list --parent sb-abc`,
	RunE: func(cmd *cobra.Command, args []string) error {
		tasks, err := store.GetAllTasks()
		if err != nil {
			return fmt.Errorf("failed to load tasks: %w", err)
		}

		// Filter tasks
		var filtered []*models.Task
		for _, task := range tasks {
			// Filter by status
			if !listAll && task.Status != models.StatusOpen {
				continue
			}
			if listStatus != "" && string(task.Status) != listStatus {
				continue
			}

			// Filter by priority
			if cmd.Flags().Changed("priority") && int(task.Priority) != listPriority {
				continue
			}

			// Filter by parent
			if listParent != "" && task.Parent != listParent {
				continue
			}

			filtered = append(filtered, task)
		}

		// Sort by priority (ascending) then by creation date
		sort.Slice(filtered, func(i, j int) bool {
			if filtered[i].Priority != filtered[j].Priority {
				return filtered[i].Priority < filtered[j].Priority
			}
			return filtered[i].CreatedAt.Before(filtered[j].CreatedAt)
		})

		// Display tasks
		if len(filtered) == 0 {
			fmt.Println("No tasks found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tPRIORITY\tSTATUS\tPARENT")
		fmt.Fprintln(w, "--\t-----\t--------\t------\t------")

		for _, task := range filtered {
			parent := task.Parent
			if parent == "" {
				parent = "-"
			}
			fmt.Fprintf(w, "%s\t%s\tP%d\t%s\t%s\n",
				task.ID,
				task.Title,
				task.Priority,
				task.Status,
				parent,
			)
		}

		w.Flush()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVar(&listAll, "all", false, "Show all tasks including closed ones")
	listCmd.Flags().IntVarP(&listPriority, "priority", "p", -1, "Filter by priority (0-4)")
	listCmd.Flags().StringVar(&listStatus, "status", "", "Filter by status (open, closed)")
	listCmd.Flags().StringVar(&listParent, "parent", "", "Filter by parent task ID")
}
