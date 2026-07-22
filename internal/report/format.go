package report

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
)

// WriteTable renders a slice of reports as an aligned, human-readable
// table to w. Uses the standard library's tabwriter rather than a
// third-party table library to avoid adding a new dependency.
func WriteTable(w io.Writer, reports []*RepoReport) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)

	fmt.Fprintln(tw, "REPO\tCI\tREADME\tLICENSE\tSTALE BRANCHES\tSCORE")
	for _, r := range reports {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%d\t%.0f%%\n",
			r.Name,
			yesNo(r.HasCI),
			yesNo(r.HasReadme),
			yesNo(r.HasLicense),
			len(r.StaleBranches),
			len(r.StalePRs),
			r.Score*100,
		)
	}

	return tw.Flush()
}

func yesNo(ok bool) string {
	if ok {
		return "yes"
	}
	return "no"
}

// WriteJSON renders a slice of reports as indented, machine-readable
// JSON to w.
func WriteJSON(w io.Writer, reports []*RepoReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(reports)
}
