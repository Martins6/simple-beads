package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// blockedCmd represents the blocked command
var blockedCmd = &cobra.Command{
	Use:   "blocked",
	Short: "Show blocked tasks",
	Long: `Display all tasks that are:
- Open (not closed)
- Blocked by one or more unclosed dependencies

Shows which tasks are waiting on other tasks to be completed.`,
	Example: `  sb blocked`,
	RunE: func(cmd *cobra.Command, args []string) error {
		blocked, err := store.FindBlocked()
		if err != nil {
			return fmt.Errorf("failed to find blocked tasks: %w", err)
		}

		if len(blocked) == 0 {
			fmt.Println("No blocked tasks.")
			return nil
		}

		// Sort by priority (ascending) then by creation date
		sort.Slice(blocked, func(i, j int) bool {
			if blocked[i].Priority != blocked[j].Priority {
				return blocked[i].Priority < blocked[j].Priority
			}
			return blocked[i].CreatedAt.Before(blocked[j].CreatedAt)
		})

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tPRIORITY\tBLOCKED BY")
		fmt.Fprintln(w, "--\t-----\t--------\t----------")

		for _, task := range blocked {
			// Get list of blocking dependencies
			tasks, _ := store.LoadTasks()
			var blocking []string
			for _, depID := range task.Dependencies {
				if depTask, exists := tasks[depID]; exists {
					if depTask.Status != "closed" {
						blocking = append(blocking, depID)
					}
				}
			}

			fmt.Fprintf(w, "%s\t%s\tP%d\t%s\n",
				task.ID,
				task.Title,
				task.Priority,
				fmt.Sprintf("%v", blocking),
			)
		}

		w.Flush()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(blockedCmd)
}
