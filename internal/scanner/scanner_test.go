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

func TestFindWorkflows_NoFallbackToDir(t *testing.T) {
	dir := t.TempDir()

	// A YAML file in the root should NOT be found (no fallback).
	if err := os.WriteFile(filepath.Join(dir, "workflow.yml"), []byte("name: test"), 0o644); err != nil {
		t.Fatal(err)
	}

	paths, err := FindWorkflows(dir)
	if err != nil {
		t.Fatalf("FindWorkflows() error: %v", err)
	}
	if len(paths) != 0 {
		t.Fatalf("len(paths) = %d, want 0", len(paths))
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

func TestFindWorkflows_IgnoresNonWorkflowYAML(t *testing.T) {
	dir := t.TempDir()

	// Create a valid workflow file.
	wfDir := filepath.Join(dir, ".github", "workflows")
	if err := os.MkdirAll(wfDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(wfDir, "ci.yml"), []byte("name: ci"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create YAML files outside .github/workflows that should be ignored.
	chartDir := filepath.Join(dir, "charts", "templates")
	if err := os.MkdirAll(chartDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(chartDir, "deployment.yaml"), []byte("apiVersion: apps/v1"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte("key: value"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create a YAML in a subdirectory of .github/workflows that should also be ignored.
	subDir := filepath.Join(wfDir, "subfolder")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "nested.yml"), []byte("name: nested"), 0o644); err != nil {
		t.Fatal(err)
	}

	paths, err := FindWorkflows(dir)
	if err != nil {
		t.Fatalf("FindWorkflows() error: %v", err)
	}
	if len(paths) != 1 {
		t.Fatalf("len(paths) = %d, want 1; got %v", len(paths), paths)
	}
	if filepath.Base(paths[0]) != "ci.yml" {
		t.Fatalf("paths[0] = %s, want ci.yml", paths[0])
	}
}

func TestFindWorkflows_NoGitHubDir(t *testing.T) {
	// Without .github/workflows, no files should be returned — no fallback.
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "workflow.yml"), []byte("name: test"), 0o644); err != nil {
		t.Fatal(err)
	}

	paths, err := FindWorkflows(dir)
	if err != nil {
		t.Fatalf("FindWorkflows() error: %v", err)
	}
	if len(paths) != 0 {
		t.Fatalf("len(paths) = %d, want 0; got %v", len(paths), paths)
	}
}

func TestFindWorkflows_Subdirectories(t *testing.T) {
	dir := t.TempDir()

	// Root-level workflows.
	rootWf := filepath.Join(dir, ".github", "workflows")
	if err := os.MkdirAll(rootWf, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(rootWf, "ci.yml"), []byte("name: ci"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Subdirectory with its own .github/workflows.
	subWf := filepath.Join(dir, "services", "api", ".github", "workflows")
	if err := os.MkdirAll(subWf, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subWf, "build.yaml"), []byte("name: build"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Deeply nested subdirectory.
	deepWf := filepath.Join(dir, "packages", "lib", "core", ".github", "workflows")
	if err := os.MkdirAll(deepWf, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(deepWf, "test.yml"), []byte("name: test"), 0o644); err != nil {
		t.Fatal(err)
	}

	paths, err := FindWorkflows(dir)
	if err != nil {
		t.Fatalf("FindWorkflows() error: %v", err)
	}
	if len(paths) != 3 {
		t.Fatalf("len(paths) = %d, want 3; got %v", len(paths), paths)
	}

	// Verify all expected files are found.
	bases := make(map[string]bool)
	for _, p := range paths {
		bases[filepath.Base(p)] = true
	}
	for _, want := range []string{"ci.yml", "build.yaml", "test.yml"} {
		if !bases[want] {
			t.Errorf("expected %s in results, got %v", want, paths)
		}
	}
}

func TestFindWorkflowsInRoots_MultipleRoots(t *testing.T) {
	dir := t.TempDir()

	// Create two project roots with workflows.
	for _, sub := range []string{"projectA", "projectB"} {
		wfDir := filepath.Join(dir, sub, ".github", "workflows")
		if err := os.MkdirAll(wfDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(wfDir, "ci.yml"), []byte("name: "+sub), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	roots := []string{
		filepath.Join(dir, "projectA"),
		filepath.Join(dir, "projectB"),
	}
	paths, err := FindWorkflowsInRoots(roots)
	if err != nil {
		t.Fatalf("FindWorkflowsInRoots() error: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("len(paths) = %d, want 2; got %v", len(paths), paths)
	}
}

func TestFindWorkflowsInRoots_SkipsMissingWorkflowsDir(t *testing.T) {
	dir := t.TempDir()

	// One root with workflows, one without.
	wfDir := filepath.Join(dir, "hasWf", ".github", "workflows")
	if err := os.MkdirAll(wfDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(wfDir, "ci.yml"), []byte("name: ci"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "noWf"), 0o755); err != nil {
		t.Fatal(err)
	}

	paths, err := FindWorkflowsInRoots([]string{
		filepath.Join(dir, "hasWf"),
		filepath.Join(dir, "noWf"),
	})
	if err != nil {
		t.Fatalf("FindWorkflowsInRoots() error: %v", err)
	}
	if len(paths) != 1 {
		t.Fatalf("len(paths) = %d, want 1; got %v", len(paths), paths)
	}
}

func TestFindWorkflowsInRoots_Empty(t *testing.T) {
	paths, err := FindWorkflowsInRoots(nil)
	if err != nil {
		t.Fatalf("FindWorkflowsInRoots() error: %v", err)
	}
	if len(paths) != 0 {
		t.Fatalf("len(paths) = %d, want 0", len(paths))
	}
}
