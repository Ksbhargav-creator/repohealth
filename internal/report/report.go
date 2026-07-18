package report

import (
	"context"
	"time"

	"github.com/Ksbhargav-creator/repohealth/internal/checks"
	"github.com/google/go-github/v66/github"
)

type RepoReport struct {
	Name          string
	HasCI         bool
	HasReadme     bool
	HasLicense    bool
	StaleBranches []string
	Score         float64
}

func Generate(ctx context.Context, client *github.Client, owner, repo string) (*RepoReport, error) {
	ci, err := checks.HasCI(ctx, client, owner, repo)
	if err != nil {
		return nil, err
	}
	readme, err := checks.HasReadme(ctx, client, owner, repo)
	if err != nil {
		return nil, err
	}
	license, err := checks.HasLicense(ctx, client, owner, repo)
	if err != nil {
		return nil, err
	}
	stale, err := checks.StaleBranches(ctx, client, owner, repo, 90*24*time.Hour)
	if err != nil {
		return nil, err
	}

	passed := 0
	total := 4
	if ci {
		passed++
	}
	if readme {
		passed++
	}
	if license {
		passed++
	}

	return &RepoReport{
		Name:          repo,
		HasCI:         ci,
		HasReadme:     readme,
		HasLicense:    license,
		StaleBranches: stale,
		Score:         float64(passed) / float64(total),
	}, nil
}
