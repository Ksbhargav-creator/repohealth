package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan repos for health signals",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("scan: not implemented yet")
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
