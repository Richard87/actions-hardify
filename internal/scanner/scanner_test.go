package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindWorkflows_GitHubDir(t *testing.T) {
	dir := t.TempDir()
	wfDir := filepath.Join(dir, ".github", "workflows")
	if err := os.MkdirAll(wfDir, 0o755); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"ci.yml", "release.yaml"} {
		if err := os.WriteFile(filepath.Join(wfDir, name), []byte("name: test"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if err := os.WriteFile(filepath.Join(wfDir, "README.md"), []byte("# test"), 0o644); err != nil {
		t.Fatal(err)
	}

	paths, err := FindWorkflows(dir)
	if err != nil {
		t.Fatalf("FindWorkflows() error: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("len(paths) = %d, want 2", len(paths))
	}
}

func TestFindWorkflows_FallbackToDir(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "workflow.yml"), []byte("name: test"), 0o644); err != nil {
		t.Fatal(err)
	}

	paths, err := FindWorkflows(dir)
	if err != nil {
		t.Fatalf("FindWorkflows() error: %v", err)
	}
	if len(paths) != 1 {
		t.Fatalf("len(paths) = %d, want 1", len(paths))
	}
}

func TestFindWorkflows_Empty(t *testing.T) {
	dir := t.TempDir()

	paths, err := FindWorkflows(dir)
	if err != nil {
		t.Fatalf("FindWorkflows() error: %v", err)
	}
	if len(paths) != 0 {
		t.Fatalf("len(paths) = %d, want 0", len(paths))
	}
}

func TestFindWorkflows_YMLAndYAML(t *testing.T) {
	dir := t.TempDir()
	wfDir := filepath.Join(dir, ".github", "workflows")
	if err := os.MkdirAll(wfDir, 0o755); err != nil {
		t.Fatal(err)
	}

	files := []string{"a.yml", "b.yaml", "c.YML", "d.YAML"}
	for _, name := range files {
		if err := os.WriteFile(filepath.Join(wfDir, name), []byte("name: test"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	paths, err := FindWorkflows(dir)
	if err != nil {
		t.Fatalf("FindWorkflows() error: %v", err)
	}
	if len(paths) != 4 {
		t.Fatalf("len(paths) = %d, want 4", len(paths))
	}
}
