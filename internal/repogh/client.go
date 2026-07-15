package repogh

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v66/github"
)

// NewClient builds an authenticated GitHub client using GITHUB_TOKEN
// from the environment
func NewClient() (*github.Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable not set")
	}
	return github.NewClient(nil).WithAuthToken(token), nil
}

// ListMyRepos returns the authenticated user's repositories
func ListMyRepos(ctx context.Context, client *github.Client) ([]*github.Repository, error) {
	repos, _, err := client.Repositories.ListByAuthenticatedUser(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("listing repos: %w", err)
	}
	return repos, nil
}
