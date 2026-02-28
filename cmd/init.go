package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize sbeads in the current directory",
	Long: `Initialize sbeads by creating the .sbeads directory and empty tasks.jsonl file.

This command must be run before using other sbeads commands in a new directory.
It creates the local storage structure for your tasks.`,
	Example: `  sb init`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := store.Init(); err != nil {
			return fmt.Errorf("failed to initialize sbeads: %w", err)
		}
		fmt.Printf("✓ Initialized sbeads in %s\n", store.GetDataPath())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
