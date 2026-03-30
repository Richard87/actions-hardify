package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

// FindWorkflows returns all .yml/.yaml files directly inside dir/.github/workflows.
// Only direct children are returned — subdirectories are not traversed.
func FindWorkflows(dir string) ([]string, error) {
	workflowDir := filepath.Join(dir, ".github", "workflows")

	entries, err := os.ReadDir(workflowDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext == ".yml" || ext == ".yaml" {
			files = append(files, filepath.Join(workflowDir, e.Name()))
		}
	}

	return files, nil
}
