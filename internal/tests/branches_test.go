package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestStaleBranches_MixedAges(t *testing.T) {
	client, mux := setupTestClient(t)

	oldDate := time.Now().Add(-200 * 24 * time.Hour).UTC().Format(time.RFC3339)
	newDate := time.Now().Add(-1 * 24 * time.Hour).UTC().Format(time.RFC3339)

	mux.HandleFunc("/repos/testowner/testrepo/branches", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[
			{"name": "old-feature", "commit": {"sha": "oldsha"}},
			{"name": "main", "commit": {"sha": "newsha"}}
		]`))
	})

	mux.HandleFunc("/repos/testowner/testrepo/commits/oldsha", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"sha": "oldsha", "commit": {"committer": {"date": "%s"}}}`, oldDate)
	})

	mux.HandleFunc("/repos/testowner/testrepo/commits/newsha", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"sha": "newsha", "commit": {"committer": {"date": "%s"}}}`, newDate)
	})

	stale, err := StaleBranches(context.Background(), client, "testowner", "testrepo", 90*24*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(stale) != 1 || stale[0] != "old-feature" {
		t.Errorf("expected only [old-feature] to be stale, got %v", stale)
	}
}

func TestStaleBranches_NoBranches(t *testing.T) {
	client, mux := setupTestClient(t)

	mux.HandleFunc("/repos/testowner/testrepo/branches", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[]`))
	})

	stale, err := StaleBranches(context.Background(), client, "testowner", "testrepo", 90*24*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stale) != 0 {
		t.Errorf("expected no stale branches, got %v", stale)
	}
}
