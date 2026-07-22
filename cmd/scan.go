package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Ksbhargav-creator/repohealth/internal/config"
	"github.com/Ksbhargav-creator/repohealth/internal/repogh"
	"github.com/Ksbhargav-creator/repohealth/internal/report"
	"github.com/google/go-github/v66/github"
	"github.com/spf13/cobra"
)

var format string
var configPath string

// Subcommand that scans your repos and prints a health report,
// either as a table (default) or as JSON via --format json.
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan repos for health signals",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load(configPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		client, err := repogh.NewClient()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		var targets []config.RepoRef

		if len(cfg.Repos) > 0 || len(cfg.Orgs) > 0 {
			targets = append(targets, cfg.Repos...)

			for _, org := range cfg.Orgs {
				orgRepos, err := repogh.ListOrgRepos(context.Background(), client, org)
				if err != nil {
					var rateErr *github.RateLimitError
					var errResp *github.ErrorResponse
					if errors.As(err, &rateErr) {
						//Github API rate limiting error handling
						fmt.Fprintf(os.Stderr, "rate limit hit, resets at %s\n", rateErr.Rate.Reset.Time)
						os.Exit(1)
					} else if errors.As(err, &errResp) && errResp.Response.StatusCode == http.StatusNotFound {
						//404 not found error handling
						fmt.Fprintln(os.Stderr, "Skipping org %s: inaccessible", org)
						continue
					} else {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}
				}
				for _, r := range orgRepos {
					targets = append(targets, config.RepoRef{Owner: r.GetOwner().GetLogin(), Name: r.GetName()})
				}
			}
		} else {
			myRepos, err := repogh.ListMyRepos(context.Background(), client)
			if err != nil {
				var rateErr *github.RateLimitError
				var errResp *github.ErrorResponse
				if errors.As(err, &rateErr) {
					//Github API rate limiting error handling
					fmt.Fprintf(os.Stderr, "rate limit hit, resets at %s\n", rateErr.Rate.Reset.Time)
					os.Exit(1)
				} else {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}
			for _, r := range myRepos {
				targets = append(targets, config.RepoRef{Owner: r.GetOwner().GetLogin(), Name: r.GetName()})
			}
		}

		var reports []*report.RepoReport
		for _, t := range targets {
			rep, err := report.Generate(context.Background(), client, t.Owner, t.Name)
			if err != nil {
				var rateErr *github.RateLimitError
				var errResp *github.ErrorResponse
				if errors.As(err, &rateErr) {
					//Github API rate limiting error handling
					fmt.Fprintf(os.Stderr, "rate limit hit, resets at %s\n", rateErr.Rate.Reset.Time)
					os.Exit(1)
				} else if errors.As(err, &errResp) && errResp.Response.StatusCode == http.StatusNotFound {
					//404 not found error handling
					fmt.Fprintln(os.Stderr, "Skipping this repo as it is inaccessible")
					continue
				} else {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
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
	scanCmd.Flags().StringVar(&configPath, "config", "repohealth.yaml", "path to config file")
}
