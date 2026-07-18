package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/Ksbhargav-creator/repohealth/internal/repogh"
	"github.com/Ksbhargav-creator/repohealth/internal/report"
	"github.com/spf13/cobra"
)

var format string

// Subcommand that scans your repos and prints a health report,
// either as a table (default) or as JSON via --format json.
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

		var reports []*report.RepoReport
		for _, r := range repos {
			rep, err := report.Generate(context.Background(), client, r.GetOwner().GetLogin(), r.GetName())
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			reports = append(reports, rep)
		}

		switch format {
		case "json":
			if err := report.WriteJSON(os.Stdout, reports); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		case "table":
			if err := report.WriteTable(os.Stdout, reports); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		default:
			fmt.Fprintf(os.Stderr, "unknown format %q (want table or json)\n", format)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVar(&format, "format", "table", "output format: table or json")
}
