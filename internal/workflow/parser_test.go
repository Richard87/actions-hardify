package workflow

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParse_BasicWorkflow(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ci.yaml")
	content := `name: CI
on: push
permissions:
  contents: read
jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: read
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      - name: Build
        run: go build ./...
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	w, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if w.Path != path {
		t.Errorf("Path = %q, want %q", w.Path, path)
	}
	if w.Permissions == nil || w.Permissions.IsEmpty {
		t.Fatal("expected top-level permissions")
	}
	if w.Permissions.Values["contents"] != "read" {
		t.Errorf("permissions.contents = %q, want %q", w.Permissions.Values["contents"], "read")
	}
	if len(w.Jobs) != 1 {
		t.Fatalf("len(Jobs) = %d, want 1", len(w.Jobs))
	}
	job := w.Jobs[0]
	if job.ID != "build" {
		t.Errorf("Job.ID = %q, want %q", job.ID, "build")
	}
	if len(job.Steps) != 2 {
		t.Fatalf("len(Steps) = %d, want 2", len(job.Steps))
	}
	if job.Steps[0].Uses == nil {
		t.Fatal("step 0 should have Uses")
	}
	if job.Steps[0].Uses.Owner != "actions" || job.Steps[0].Uses.Repo != "checkout" {
		t.Errorf("Uses = %v, want actions/checkout", job.Steps[0].Uses)
	}
	if job.Steps[0].Uses.Ref != "v4" {
		t.Errorf("Ref = %q, want %q", job.Steps[0].Uses.Ref, "v4")
	}
	if job.Steps[1].Uses != nil {
		t.Error("step 1 (run:) should have nil Uses")
	}
}

func TestParse_NoPermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ci.yaml")
	content := `name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	w, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if !w.Permissions.IsEmpty {
		t.Error("expected IsEmpty for missing permissions")
	}
}

func TestParse_ScalarPermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ci.yaml")
	content := `name: CI
on: push
permissions: write-all
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	w, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if !w.Permissions.IsAll {
		t.Error("expected IsAll for scalar permissions")
	}
	if w.Permissions.Values["all"] != "write-all" {
		t.Errorf("Values[all] = %q, want %q", w.Permissions.Values["all"], "write-all")
	}
}

func TestParse_ActionWithSubpath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ci.yaml")
	content := `name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: aws-actions/configure-aws-credentials/assume-role@v4
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	w, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	ref := w.Jobs[0].Steps[0].Uses
	if ref == nil {
		t.Fatal("expected Uses")
	}
	if ref.Owner != "aws-actions" {
		t.Errorf("Owner = %q, want %q", ref.Owner, "aws-actions")
	}
	if ref.Repo != "configure-aws-credentials" {
		t.Errorf("Repo = %q, want %q", ref.Repo, "configure-aws-credentials")
	}
	if ref.Path != "/assume-role" {
		t.Errorf("Path = %q, want %q", ref.Path, "/assume-role")
	}
}

func TestParse_LocalAndDockerActions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ci.yaml")
	content := `name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: ./local-action
      - uses: docker://alpine:3.18
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	w, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	for i, step := range w.Jobs[0].Steps {
		if step.Uses != nil {
			t.Errorf("step %d: expected nil Uses for local/docker action", i)
		}
	}
}

func TestParse_FileNotFound(t *testing.T) {
	_, err := Parse("/nonexistent/file.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestParse_MultipleJobs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ci.yaml")
	content := `name: CI
on: push
jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	w, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if len(w.Jobs) != 2 {
		t.Fatalf("len(Jobs) = %d, want 2", len(w.Jobs))
	}
	if w.Jobs[0].ID != "lint" {
		t.Errorf("Jobs[0].ID = %q, want %q", w.Jobs[0].ID, "lint")
	}
	if w.Jobs[1].ID != "test" {
		t.Errorf("Jobs[1].ID = %q, want %q", w.Jobs[1].ID, "test")
	}
	if len(w.Jobs[1].Steps) != 2 {
		t.Errorf("Jobs[1].Steps = %d, want 2", len(w.Jobs[1].Steps))
	}
}

func TestCollectActions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ci.yaml")
	content := `name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: echo hello
      - uses: actions/setup-go@v5
      - uses: ./local
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	w, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	refs := CollectActions(w)
	if len(refs) != 2 {
		t.Fatalf("len(refs) = %d, want 2", len(refs))
	}
	if refs[0].String() != "actions/checkout@v4" {
		t.Errorf("refs[0] = %q, want %q", refs[0].String(), "actions/checkout@v4")
	}
	if refs[1].String() != "actions/setup-go@v5" {
		t.Errorf("refs[1] = %q, want %q", refs[1].String(), "actions/setup-go@v5")
	}
}
