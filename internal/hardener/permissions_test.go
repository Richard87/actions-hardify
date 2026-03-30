package hardener

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/richard87/actions-hardify/internal/workflow"
	"gopkg.in/yaml.v3"
)

func TestCheckPermissions_NoTopLevel(t *testing.T) {
	w := parseWorkflow(t, `name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
`)

	findings := checkPermissions(w)

	var foundWorkflowLevel, foundJobLevel bool
	for _, f := range findings {
		if f.Job == "" && f.Type == FindingPermissions {
			foundWorkflowLevel = true
		}
		if f.Job == "build" && f.Type == FindingPermissions {
			foundJobLevel = true
		}
	}
	if !foundWorkflowLevel {
		t.Error("expected finding for missing top-level permissions")
	}
	if !foundJobLevel {
		t.Error("expected finding for job without permissions when top-level is empty")
	}
}

func TestCheckPermissions_BroadPermissions(t *testing.T) {
	w := parseWorkflow(t, `name: CI
on: push
permissions: write-all
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
`)

	findings := checkPermissions(w)

	var found bool
	for _, f := range findings {
		if f.Type == FindingPermissions && f.Current == "write-all" {
			found = true
		}
	}
	if !found {
		t.Error("expected finding for overly broad permissions")
	}
}

func TestCheckPermissions_RestrictedPermissions(t *testing.T) {
	w := parseWorkflow(t, `name: CI
on: push
permissions:
  contents: read
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
`)

	findings := checkPermissions(w)

	for _, f := range findings {
		if f.Type == FindingPermissions && f.Job == "" {
			t.Errorf("unexpected top-level permission finding: %s", f.Message)
		}
	}
}

func TestApplyPermissionFixes_AddsMissing(t *testing.T) {
	w := parseWorkflow(t, `name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
`)

	if !w.Permissions.IsEmpty {
		t.Fatal("precondition: expected permissions to be empty")
	}

	applyPermissionFixes(w)

	root := w.Doc.Content[0]
	var foundPermissions bool
	for i := 0; i < len(root.Content)-1; i += 2 {
		if root.Content[i].Value == "permissions" {
			foundPermissions = true
			break
		}
	}
	if !foundPermissions {
		t.Error("expected permissions key to be inserted")
	}
}

func TestApplyPermissionFixes_SkipsExisting(t *testing.T) {
	w := parseWorkflow(t, `name: CI
on: push
permissions:
  contents: read
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
`)

	root := w.Doc.Content[0]
	countBefore := countKey(root, "permissions")

	applyPermissionFixes(w)

	countAfter := countKey(root, "permissions")
	if countAfter != countBefore {
		t.Errorf("permissions count changed: %d -> %d", countBefore, countAfter)
	}
}

func parseWorkflow(t *testing.T, content string) *workflow.Workflow {
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

func countKey(mapping *yaml.Node, key string) int {
	count := 0
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		if mapping.Content[i].Value == key {
			count++
		}
	}
	return count
}
