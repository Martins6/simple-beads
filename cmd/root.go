/*
Copyright © 2026 Adriel Martins
*/
package cmd

import (
	"os"

	"github.com/Martins6/simple-beads/internal/storage"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sb",
	Short: "Simple Beads - A lightweight task management CLI",
	Long: `Simple Beads (sb) is a lightweight, local-only task management CLI.

Store tasks in SQLite database (.sbeads/tasks.db), with support for:
- Multi-agent concurrent access (ACID transactions)
- Parent-child task relationships
- Dependencies between tasks
- Priority levels (0-4, where 0 is highest)
- Task status tracking (open/closed)

Examples:
  sb init                              # Initialize .sbeads directory
  sb create "Fix bug" -p 1             # Create high priority task
  sb create "Subtask" --parent sb-abc  # Create child task
  sb list                              # List all open tasks
  sb ready                             # Show tasks ready to work on
  sb close sb-abc                      # Mark task as done`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// store is the global storage instance
var store *storage.Storage

func init() {
	// Initialize storage
	store = storage.NewStorage("")
}
