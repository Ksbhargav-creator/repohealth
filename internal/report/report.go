package report

import (
	"context"
	"time"

	"github.com/Ksbhargav-creator/repohealth/internal/checks"
	"github.com/google/go-github/v66/github"
)

type RepoReport struct {
	Name          string   `json:"name"`
	HasCI         bool     `json:"has_ci"`
	HasReadme     bool     `json:"has_readme"`
	HasLicense    bool     `json:"has_license"`
	StaleBranches []string `json:"stale_branches"`
	Score         float64  `json:"score"`
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
