package checks

import (
	"context"
	"net/http"

	"github.com/google/go-github/v66/github"
)

// HasCI reports whether a repo has at least one CI workflow file
// under .github/workflows
func HasCI(ctx context.Context, client *github.Client, owner, repo string) (bool, error) {
	//GetContents has something like (fileContents *[]RepositoryContent, resp *Response,
	//err error) as it's return value.
	_, dirContents, resp, err := client.Repositories.GetContents(ctx, owner, repo, ".github/workflows", nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return false, nil // no workflows dir = fail
		}
		return false, err
	}
	return len(dirContents) > 0, nil
}
