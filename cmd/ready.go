package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// readyCmd represents the ready command
var readyCmd = &cobra.Command{
	Use:   "ready",
	Short: "Show tasks ready to work on",
	Long: `Display all tasks that are:
- Open (not closed)
- Not blocked by any dependencies (all dependencies are closed)

These are the tasks you can start working on immediately.`,
	Example: `  sb ready`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ready, err := store.FindReady()
		if err != nil {
			return fmt.Errorf("failed to find ready tasks: %w", err)
		}

		if len(ready) == 0 {
			fmt.Println("No tasks ready to work on.")
			return nil
		}

		// Sort by priority (ascending) then by creation date
		sort.Slice(ready, func(i, j int) bool {
			if ready[i].Priority != ready[j].Priority {
				return ready[i].Priority < ready[j].Priority
			}
			return ready[i].CreatedAt.Before(ready[j].CreatedAt)
		})

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tPRIORITY")
		fmt.Fprintln(w, "--\t-----\t--------")

		for _, task := range ready {
			fmt.Fprintf(w, "%s\t%s\tP%d\n",
				task.ID,
				task.Title,
				task.Priority,
			)
		}

		w.Flush()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(readyCmd)
}
