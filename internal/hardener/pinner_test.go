package hardener

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/richard87/actions-hardify/internal/github"
	"github.com/richard87/actions-hardify/internal/workflow"
)

func TestCheckActions_PinsUnpinned(t *testing.T) {
	fakeSHA := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/repos/actions/checkout/git/ref/tags/v4":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"object": map[string]string{"sha": fakeSHA, "type": "commit"},
			})
		case r.URL.Path == "/repos/actions/checkout/releases/latest":
			json.NewEncoder(w).Encode(github.Release{TagName: "v4"})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	gh := &github.Client{
		HTTPClient: srv.Client(),
		BaseURL:    srv.URL,
	}

	wf := parseTestWorkflow(t, `name: CI
on: push
permissions:
  contents: read
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
`)

	findings, err := checkActions(context.Background(), wf, gh, true)
	if err != nil {
		t.Fatalf("checkActions() error: %v", err)
	}

	var foundUnpinned bool
	for _, f := range findings {
		if f.Type == FindingUnpinned {
			foundUnpinned = true
			if f.Current != "v4" {
				t.Errorf("Current = %q, want %q", f.Current, "v4")
			}
		}
	}
	if !foundUnpinned {
		t.Error("expected unpinned finding")
	}
}

func TestCheckActions_AlreadyPinned(t *testing.T) {
	fakeSHA := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/repos/actions/checkout/releases/latest":
			json.NewEncoder(w).Encode(github.Release{TagName: "v4"})
		case r.URL.Path == "/repos/actions/checkout/tags":
			json.NewEncoder(w).Encode([]github.Tag{{Name: "v4", Commit: struct {
				SHA string `json:"sha"`
			}{SHA: fakeSHA}}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	gh := &github.Client{
		HTTPClient: srv.Client(),
		BaseURL:    srv.URL,
	}

	wf := parseTestWorkflow(t, `name: CI
on: push
permissions:
  contents: read
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@`+fakeSHA+`
`)

	findings, err := checkActions(context.Background(), wf, gh, true)
	if err != nil {
		t.Fatalf("checkActions() error: %v", err)
	}

	for _, f := range findings {
		if f.Type == FindingUnpinned {
			t.Error("should not flag already-pinned action as unpinned")
		}
	}
}

func parseTestWorkflow(t *testing.T, content string) *workflow.Workflow {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	w, err := workflow.Parse(path)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	return w
}

func TestCheckActions_ReusableWorkflow(t *testing.T) {
	fakeSHA := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/repos/equinor/radix-reusable-workflows/git/ref/tags/v1.0.1":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"object": map[string]string{"sha": fakeSHA, "type": "commit"},
			})
		case r.URL.Path == "/repos/equinor/radix-reusable-workflows/releases/latest":
			json.NewEncoder(w).Encode(github.Release{TagName: "v1.0.1"})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	gh := &github.Client{
		HTTPClient: srv.Client(),
		BaseURL:    srv.URL,
	}

	wf := parseTestWorkflow(t, `name: CI
on: push
permissions: {}
jobs:
  reusable:
    permissions:
      contents: read
    uses: equinor/radix-reusable-workflows/.github/workflows/template.yml@v1.0.1
`)

	findings, err := checkActions(context.Background(), wf, gh, true)
	if err != nil {
		t.Fatalf("checkActions() error: %v", err)
	}

	var foundUnpinned bool
	for _, f := range findings {
		if f.Type == FindingUnpinned {
			foundUnpinned = true
			if f.Current != "v1.0.1" {
				t.Errorf("Current = %q, want %q", f.Current, "v1.0.1")
			}
		}
	}
	if !foundUnpinned {
		t.Error("expected unpinned finding for reusable workflow")
	}
}
