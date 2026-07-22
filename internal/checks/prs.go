package checks

import (
	"context"
	"time"

	"github.com/google/go-github/v66/github"
)

func StalePRs(ctx context.Context, client *github.Client, owner, repo string, maxAge time.Duration) ([]string, error) {
	opts := &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var all_prs []*github.PullRequest
	for {
		prs, resp, err := client.PullRequests.List(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}
		all_prs = append(all_prs, prs...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	var stale []string
	cutoff := time.Now().Add(-maxAge)

	for _, pr := range all_prs {
		if pr.GetCreatedAt().Before(cutoff) {
			stale = append(stale, pr.GetTitle())
		}
	}

	return stale, nil
}
