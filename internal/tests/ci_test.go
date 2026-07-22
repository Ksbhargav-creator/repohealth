package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/v66/github"
)

// setupTestClient spins up a fake HTTP server and returns a github.Client
// pointed at it instead of the real GitHub API, plus a mux to register
// fake endpoint responses on. The server is closed automatically when
// the test finishes.
func setupTestClient(t *testing.T) (*github.Client, *http.ServeMux) {
	t.Helper()

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := github.NewClient(nil)
	baseURL, err := url.Parse(server.URL + "/")
	if err != nil {
		t.Fatalf("parsing test server URL: %v", err)
	}
	client.BaseURL = baseURL

	return client, mux
}

func TestHasCI_WithWorkflows(t *testing.T) {
	client, mux := setupTestClient(t)

	mux.HandleFunc("/repos/testowner/testrepo/contents/.github/workflows", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"name": "ci.yml", "type": "file"}]`))
	})

	ok, err := HasCI(context.Background(), client, "testowner", "testrepo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Errorf("expected HasCI to return true, got false")
	}
}

func TestHasCI_NoWorkflows(t *testing.T) {
	client, mux := setupTestClient(t)

	mux.HandleFunc("/repos/testowner/testrepo/contents/.github/workflows", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	ok, err := HasCI(context.Background(), client, "testowner", "testrepo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Errorf("expected HasCI to return false, got true")
	}
}

func TestHasCI_EmptyWorkflowsDir(t *testing.T) {
	client, mux := setupTestClient(t)

	mux.HandleFunc("/repos/testowner/testrepo/contents/.github/workflows", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[]`))
	})

	ok, err := HasCI(context.Background(), client, "testowner", "testrepo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Errorf("expected HasCI to return false for an empty workflows dir, got true")
	}
}
