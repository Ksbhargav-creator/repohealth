package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.1.0-dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the repohealth version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("repohealth", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
