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

	// Create workflow files
	for _, name := range []string{"ci.yml", "release.yaml"} {
		if err := os.WriteFile(filepath.Join(wfDir, name), []byte("name: test"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	// Create non-YAML file (should be ignored)
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



































}	}		t.Fatalf("len(paths) = %d, want 4", len(paths))	if len(paths) != 4 {	}		t.Fatalf("FindWorkflows() error: %v", err)	if err != nil {	paths, err := FindWorkflows(dir)	}		}			t.Fatal(err)		if err := os.WriteFile(filepath.Join(wfDir, name), []byte("name: test"), 0o644); err != nil {	for _, name := range files {	files := []string{"a.yml", "b.yaml", "c.YML", "d.YAML"}	}		t.Fatal(err)	if err := os.MkdirAll(wfDir, 0o755); err != nil {	wfDir := filepath.Join(dir, ".github", "workflows")	dir := t.TempDir()func TestFindWorkflows_YMLAndYAML(t *testing.T) {}	}		t.Fatalf("len(paths) = %d, want 0", len(paths))	if len(paths) != 0 {	}		t.Fatalf("FindWorkflows() error: %v", err)	if err != nil {	paths, err := FindWorkflows(dir)	dir := t.TempDir()func TestFindWorkflows_Empty(t *testing.T) {