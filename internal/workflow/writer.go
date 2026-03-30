package workflow

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Write serializes the workflow's modified AST back to disk.
func Write(w *Workflow) error {
	out, err := marshalPreserving(&w.Doc)
	if err != nil {
		return fmt.Errorf("marshaling %s: %w", w.Path, err)
	}
	return os.WriteFile(w.Path, out, 0o644)
}

func marshalPreserving(doc *yaml.Node) ([]byte, error) {
	return yaml.Marshal(doc)
}
