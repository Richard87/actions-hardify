package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FindWorkflowsInRoots returns workflow files from the given project root
// folders. Each root is expected to contain a .github/workflows directory.
// Only .yml/.yaml files directly inside each root's .github/workflows
// directory are returned.
func FindWorkflowsInRoots(roots []string) ([]string, error) {
	var files []string
	for _, root := range roots {
		wfDir := filepath.Join(root, ".github", "workflows")
		entries, err := os.ReadDir(wfDir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("reading %s: %w", wfDir, err)
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			ext := strings.ToLower(filepath.Ext(e.Name()))
			if ext == ".yml" || ext == ".yaml" {
				files = append(files, filepath.Join(wfDir, e.Name()))
			}
		}
	}
	return files, nil
}

// FindWorkflows returns all .yml/.yaml workflow files found in any
// .github/workflows directory under dir. It walks the directory tree
// recursively, matching the pattern **/.github/workflows/*.{yml,yaml}.
// Only direct children of each .github/workflows directory are returned —
// subdirectories within workflows are not traversed.
func FindWorkflows(dir string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable entries
		}

		// Skip .git directories for performance.
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}

		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(d.Name()))
		if ext != ".yml" && ext != ".yaml" {
			return nil
		}

		// Check that the file sits directly inside a .github/workflows directory.
		parentDir := filepath.Dir(path)
		if filepath.Base(parentDir) == "workflows" && filepath.Base(filepath.Dir(parentDir)) == ".github" {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
