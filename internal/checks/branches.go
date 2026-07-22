package checks

import (
	"context"
	"time"

	"github.com/google/go-github/v66/github"
)

func StaleBranches(ctx context.Context, client *github.Client, owner, repo string, maxAge time.Duration) ([]string, error) {
	opts := &github.BranchListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var all_branches []*github.Branch
	for {
		branches, resp, err := client.Repositories.ListBranches(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}
		all_branches = append(all_branches, branches...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	var stale []string
	cutoff := time.Now().Add(-maxAge)

	for _, b := range all_branches {
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
