package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Setting up cobra-cli skeleton
var rootCmd = &cobra.Command{
	Use:   "repohealth",
	Short: "Audit GitHub repos for health signals",
	Long: `repohealth scans GitHub repositories for health signals
			like CI presence, README/LICENSE, stale branches, and aging PRs,
			then reports a health score per repo.`,
}

// Executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
