package checks

import (
	"context"
	"net/http"

	"github.com/google/go-github/v66/github"
)

// HasReadme reports whether the repo has a readme or not
func HasReadme(ctx context.Context, client *github.Client, owner, repo string) (bool, error) {
	_, resp, err := client.Repositories.GetReadme(ctx, owner, repo, nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// HasLicense reports whether the repo has a License or not
func HasLicense(ctx context.Context, client *github.Client, owner, repo string) (bool, error) {
	_, resp, err := client.Repositories.License(ctx, owner, repo)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
