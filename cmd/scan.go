package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/Ksbhargav-creator/repohealth/internal/repogh"
	"github.com/Ksbhargav-creator/repohealth/internal/report"
	"github.com/spf13/cobra"
)

// Subcommand that lets you scan your repo
// Currently lets you export all of your repo names
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan repos for health signals",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := repogh.NewClient()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		repos, err := repogh.ListMyRepos(context.Background(), client)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		for _, r := range repos {
			fmt.Println(r.GetName())

			report, err := report.Generate(context.Background(), client, r.GetOwner().GetLogin(), r.GetName())
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Println(report)
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
