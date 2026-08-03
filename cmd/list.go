package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/Martins6/simple-beads/internal/models"
	"github.com/spf13/cobra"
)

// dateFlag is a pflag.Value that parses a YYYY-MM-DD date string in local
// time. The zero value indicates the flag was not set.
type dateFlag struct {
	raw   string
	day   time.Time
	isSet bool
}

func (d *dateFlag) String() string {
	return d.raw
}

func (d *dateFlag) Set(v string) error {
	t, err := time.ParseInLocation("2006-01-02", v, time.Local)
	if err != nil {
		return fmt.Errorf("invalid date %q: expected format YYYY-MM-DD", v)
	}
	d.raw = v
	d.day = t
	d.isSet = true
	return nil
}

func (d *dateFlag) Type() string {
	return "date"
}

func (d *dateFlag) IsSet() bool {
	return d.isSet
}

// startOfDay returns t truncated to the local-time day boundary (00:00:00).
func startOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// endOfDayExclusive returns the start of the next local-time day after t.
// Combined with startOfDay it forms a half-open interval [start, end) that
// is immune to DST transitions.
func endOfDayExclusive(t time.Time) time.Time {
	return startOfDay(t).AddDate(0, 0, 1)
}

var (
	listAll      bool
	listPriority int
	listStatus   string
	listParent   string
	listOn       dateFlag
	listAfter    dateFlag
	listBefore   dateFlag
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Long: `List tasks with optional filtering.

By default, only open tasks are shown. Use --all to show all tasks including closed ones.

Date filters (--on, --after, --before) accept ISO 8601 dates in YYYY-MM-DD format
and are interpreted in local time. They are mutually exclusive with each other:
--on cannot be combined with --after or --before. The relevant date column is
selected automatically: closed_at for closed tasks and created_at for open tasks.
Bounds are inclusive.`,
	Example: `  sb list
  sb list --all
  sb list -p 1
  sb list --status closed
  sb list --parent sb-abc
  sb list --status closed --on 2026-08-03
  sb list --after 2026-08-01 --before 2026-08-04`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if listOn.IsSet() && (listAfter.IsSet() || listBefore.IsSet()) {
			return fmt.Errorf("--on cannot be combined with --after or --before")
		}

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

			// Filter by date (smart by status: closed_at for closed, created_at for open)
			if listOn.IsSet() || listAfter.IsSet() || listBefore.IsSet() {
				relevant := task.CreatedAt
				if task.Status == models.StatusClosed {
					relevant = task.ClosedAt
				}

				if listOn.IsSet() {
					onStart := startOfDay(listOn.day)
					onEnd := endOfDayExclusive(listOn.day)
					if relevant.Before(onStart) || !relevant.Before(onEnd) {
						continue
					}
				}
				if listAfter.IsSet() {
					if relevant.Before(startOfDay(listAfter.day)) {
						continue
					}
				}
				if listBefore.IsSet() {
					if !relevant.Before(endOfDayExclusive(listBefore.day)) {
						continue
					}
				}
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
	listCmd.Flags().Var(&listOn, "on", "Filter by exact date (YYYY-MM-DD, local time)")
	listCmd.Flags().Var(&listAfter, "after", "Filter by date >= value (YYYY-MM-DD, local time, inclusive)")
	listCmd.Flags().Var(&listBefore, "before", "Filter by date <= value (YYYY-MM-DD, local time, inclusive)")
}
