package checks

import (
	"context"
	"time"

	"github.com/google/go-github/v66/github"
)

func StalePRs(ctx context.Context, client *github.Client, owner, repo string, maxAge time.Duration) ([]string, error) {
	opts := &github.PullRequestListOptions{State: "open"}
	prs, _, err := client.PullRequests.List(ctx, owner, repo, opts)
	if err != nil {
		return nil, err
	}

	var stale []string
	cutoff := time.Now().Add(-maxAge)

	for _, pr := range prs {
		if pr.GetCreatedAt().Before(cutoff) {
			stale = append(stale, pr.GetTitle())
		}
	}

	return stale, nil
}
