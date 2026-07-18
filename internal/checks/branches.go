package checks

import (
	"context"
	"time"

	"github.com/google/go-github/v66/github"
)

func StaleBranches(ctx context.Context, client *github.Client, owner, repo string, maxAge time.Duration) ([]string, error) {
	branches, _, err := client.Repositories.ListBranches(ctx, owner, repo, nil)
	if err != nil {
		return nil, err
	}

	var stale []string
	cutoff := time.Now().Add(-maxAge)

	for _, b := range branches {
		commit, _, err := client.Repositories.GetCommit(ctx, owner, repo, b.GetCommit().GetSHA(), nil)
		if err != nil {
			return nil, err
		}
		if commit.GetCommit().GetCommitter().GetDate().Before(cutoff) {
			stale = append(stale, b.GetName())
		}
	}
	return stale, nil
}
