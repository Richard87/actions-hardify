package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

// FindWorkflows returns all .yml/.yaml files under dir/.github/workflows.
// Falls back to scanning dir directly if no .github/workflows exists.
func FindWorkflows(dir string) ([]string, error) {
	var files []string

	workflowDir := filepath.Join(dir, ".github", "workflows")
	if err := collectYAML(workflowDir, &files); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if len(files) == 0 {
		if err := collectYAML(dir, &files); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}

	return files, nil
}

func collectYAML(dir string, out *[]string) error {
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".yml" || ext == ".yaml" {
			*out = append(*out, path)
		}
		return nil
	})
}
