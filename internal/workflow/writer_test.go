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

	written, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	if len(written) == 0 {
		t.Fatal("written file is empty")
	}

	w2, err := Parse(path)
	if err != nil {
		t.Fatalf("re-Parse() error: %v", err)
	}
	if len(w2.Jobs) != 1 {
		t.Errorf("re-parsed Jobs = %d, want 1", len(w2.Jobs))
	}
}

func TestWrite_PreservesWhitespace(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ci.yaml")
	content := `name: CI
on: push

# Run checks
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
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

	// Write without modifications — output must be byte-for-byte identical.
	if err := Write(w); err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	written, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(written) != content {
		t.Errorf("Write changed the file.\nwant:\n%s\ngot:\n%s", content, written)
	}
}

func TestWrite_ModifiedValuePreservesOtherLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ci.yaml")
	content := `name: CI
on: push

# Run checks
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
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

	// Modify the action ref (simulating what the pinner does).
	ref := w.Jobs[0].Steps[0].Uses
	if ref == nil {
		t.Fatal("expected Uses ref")
	}
	ref.Node.Value = "actions/checkout@aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	ref.Node.LineComment = "v4"

	if err := Write(w); err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	written, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	want := `name: CI
on: push

# Run checks
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa # v4
      - name: Build
        run: go build ./...
`
	if string(written) != want {
		t.Errorf("Write output mismatch.\nwant:\n%s\ngot:\n%s", want, written)
	}
}
