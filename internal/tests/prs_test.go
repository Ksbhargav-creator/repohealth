package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestStalePRs_MixedAges(t *testing.T) {
	client, mux := setupTestClient(t)

	oldDate := time.Now().Add(-120 * 24 * time.Hour).UTC().Format(time.RFC3339)
	newDate := time.Now().Add(-2 * 24 * time.Hour).UTC().Format(time.RFC3339)

	mux.HandleFunc("/repos/testowner/testrepo/pulls", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `[
			{"title": "old PR", "created_at": "%s"},
			{"title": "new PR", "created_at": "%s"}
		]`, oldDate, newDate)
	})

	stale, err := StalePRs(context.Background(), client, "testowner", "testrepo", 90*24*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(stale) != 1 || stale[0] != "old PR" {
		t.Errorf("expected only [\"old PR\"] to be stale, got %v", stale)
	}
}

func TestStalePRs_NoOpenPRs(t *testing.T) {
	client, mux := setupTestClient(t)

	mux.HandleFunc("/repos/testowner/testrepo/pulls", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[]`))
	})

	stale, err := StalePRs(context.Background(), client, "testowner", "testrepo", 90*24*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stale) != 0 {
		t.Errorf("expected no stale PRs, got %v", stale)
	}
}
