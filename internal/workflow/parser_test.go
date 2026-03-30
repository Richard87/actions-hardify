package workflow
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
























































































































































































































































}	}		t.Errorf("refs[1] = %q, want %q", refs[1].String(), "actions/setup-go@v5")	if refs[1].String() != "actions/setup-go@v5" {	}		t.Errorf("refs[0] = %q, want %q", refs[0].String(), "actions/checkout@v4")	if refs[0].String() != "actions/checkout@v4" {	}		t.Fatalf("len(refs) = %d, want 2", len(refs))	if len(refs) != 2 {	refs := CollectActions(w)	}		t.Fatalf("Parse() error: %v", err)	if err != nil {	w, err := Parse(path)	}		t.Fatal(err)	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {`      - uses: ./local      - uses: actions/setup-go@v5      - run: echo hello      - uses: actions/checkout@v4    steps:    runs-on: ubuntu-latest  build:jobs:on: push	content := `name: CI	path := filepath.Join(dir, "ci.yaml")	dir := t.TempDir()func TestCollectActions(t *testing.T) {}	}		t.Errorf("Jobs[1].Steps = %d, want 2", len(w.Jobs[1].Steps))	if len(w.Jobs[1].Steps) != 2 {	}		t.Errorf("Jobs[1].ID = %q, want %q", w.Jobs[1].ID, "test")	if w.Jobs[1].ID != "test" {	}		t.Errorf("Jobs[0].ID = %q, want %q", w.Jobs[0].ID, "lint")	if w.Jobs[0].ID != "lint" {	}		t.Fatalf("len(Jobs) = %d, want 2", len(w.Jobs))	if len(w.Jobs) != 2 {	}		t.Fatalf("Parse() error: %v", err)	if err != nil {	w, err := Parse(path)	}		t.Fatal(err)	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {`      - uses: actions/setup-go@v5      - uses: actions/checkout@v4    steps:    runs-on: ubuntu-latest  test:      - uses: actions/checkout@v4    steps:    runs-on: ubuntu-latest  lint:jobs:on: push	content := `name: CI	path := filepath.Join(dir, "ci.yaml")	dir := t.TempDir()func TestParse_MultipleJobs(t *testing.T) {}	}		t.Fatal("expected error for invalid YAML")	if err == nil {	_, err := Parse(path)	}		t.Fatal(err)	if err := os.WriteFile(path, []byte(":\n  :\n   - }{"), 0o644); err != nil {	path := filepath.Join(dir, "bad.yaml")	dir := t.TempDir()func TestParse_InvalidYAML(t *testing.T) {}	}		t.Fatal("expected error for nonexistent file")	if err == nil {	_, err := Parse("/nonexistent/file.yaml")func TestParse_FileNotFound(t *testing.T) {}	}		}			t.Errorf("step %d: expected nil Uses for local/docker action", i)		if step.Uses != nil {	for i, step := range w.Jobs[0].Steps {	}		t.Fatalf("Parse() error: %v", err)	if err != nil {	w, err := Parse(path)	}		t.Fatal(err)	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {`      - uses: docker://alpine:3.18      - uses: ./local-action    steps:    runs-on: ubuntu-latest  build:jobs:on: push	content := `name: CI	path := filepath.Join(dir, "ci.yaml")	dir := t.TempDir()func TestParse_LocalAndDockerActions(t *testing.T) {}	}		t.Errorf("Path = %q, want %q", ref.Path, "/assume-role")	if ref.Path != "/assume-role" {	}		t.Errorf("Repo = %q, want %q", ref.Repo, "configure-aws-credentials")	if ref.Repo != "configure-aws-credentials" {	}		t.Errorf("Owner = %q, want %q", ref.Owner, "aws-actions")	if ref.Owner != "aws-actions" {	}		t.Fatal("expected Uses")	if ref == nil {	ref := w.Jobs[0].Steps[0].Uses	}		t.Fatalf("Parse() error: %v", err)	if err != nil {	w, err := Parse(path)	}		t.Fatal(err)	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {`      - uses: aws-actions/configure-aws-credentials/assume-role@v4    steps:    runs-on: ubuntu-latest  build:jobs:on: push	content := `name: CI	path := filepath.Join(dir, "ci.yaml")	dir := t.TempDir()func TestParse_ActionWithSubpath(t *testing.T) {}	}		t.Errorf("Values[all] = %q, want %q", w.Permissions.Values["all"], "write-all")	if w.Permissions.Values["all"] != "write-all" {	}		t.Error("expected IsAll for scalar permissions")	if !w.Permissions.IsAll {	}		t.Fatalf("Parse() error: %v", err)	if err != nil {	w, err := Parse(path)	}		t.Fatal(err)	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {`      - uses: actions/checkout@v4    steps:    runs-on: ubuntu-latest  build:jobs:permissions: write-allon: push	content := `name: CI	path := filepath.Join(dir, "ci.yaml")	dir := t.TempDir()func TestParse_ScalarPermissions(t *testing.T) {}	}		t.Error("expected IsEmpty for missing permissions")	if !w.Permissions.IsEmpty {	}		t.Fatalf("Parse() error: %v", err)	if err != nil {	w, err := Parse(path)	}		t.Fatal(err)	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {`      - uses: actions/checkout@v4    steps:    runs-on: ubuntu-latest  build:jobs:on: push	content := `name: CI	path := filepath.Join(dir, "ci.yaml")	dir := t.TempDir()func TestParse_NoPermissions(t *testing.T) {}	}		t.Error("step 1 (run:) should have nil Uses")	if job.Steps[1].Uses != nil {	}		t.Errorf("Ref = %q, want %q", job.Steps[0].Uses.Ref, "v4")	if job.Steps[0].Uses.Ref != "v4" {	}		t.Errorf("Uses = %v, want actions/checkout", job.Steps[0].Uses)	if job.Steps[0].Uses.Owner != "actions" || job.Steps[0].Uses.Repo != "checkout" {	}		t.Fatal("step 0 should have Uses")	if job.Steps[0].Uses == nil {	}		t.Fatalf("len(Steps) = %d, want 2", len(job.Steps))	if len(job.Steps) != 2 {	}		t.Errorf("Job.ID = %q, want %q", job.ID, "build")	if job.ID != "build" {	job := w.Jobs[0]	}		t.Fatalf("len(Jobs) = %d, want 1", len(w.Jobs))	if len(w.Jobs) != 1 {	}		t.Errorf("permissions.contents = %q, want %q", w.Permissions.Values["contents"], "read")	if w.Permissions.Values["contents"] != "read" {	}		t.Fatal("expected top-level permissions")	if w.Permissions == nil || w.Permissions.IsEmpty {	}		t.Errorf("Path = %q, want %q", w.Path, path)	if w.Path != path {	}		t.Fatalf("Parse() error: %v", err)	if err != nil {	w, err := Parse(path)