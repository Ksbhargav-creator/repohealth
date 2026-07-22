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
	//Pagination error handling
	opts := &github.RepositoryListByAuthenticatedUserOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var all_repos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByAuthenticatedUser(ctx, opts)
		if err != nil {
			return nil, err
		}
		all_repos = append(all_repos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return all_repos, nil
}

// ListOrgRepos returns the org repositories
func ListOrgRepos(ctx context.Context, client *github.Client, org string) ([]*github.Repository, error) {
	opts := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var all_repos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, org, opts)
		if err != nil {
			return nil, err
		}
		all_repos = append(all_repos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return all_repos, nil
}
