package workflow
package workflow

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWrite(t *testing.T) {
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

	if err := Write(w); err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	// Read back and verify it's valid YAML
	written, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	if len(written) == 0 {












}	}		t.Errorf("re-parsed Jobs = %d, want 1", len(w2.Jobs))	if len(w2.Jobs) != 1 {	}		t.Fatalf("re-Parse() error: %v", err)	if err != nil {	w2, err := Parse(path)	// Re-parse to verify	}		t.Fatal("written file is empty")