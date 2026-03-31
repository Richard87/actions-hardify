package workflow

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Write serializes the workflow's modified AST back to disk,
// preserving original whitespace by patching only changed values.
func Write(w *Workflow) error {
	out, err := patchRaw(w)
	if err != nil {
		return fmt.Errorf("patching %s: %w", w.Path, err)
	}
	return os.WriteFile(w.Path, out, 0o644)
}

type byteEdit struct {
	start int
	end   int
	text  []byte
}

func patchRaw(w *Workflow) ([]byte, error) {
	if w.Raw == nil || w.Snapshots == nil {
		return yaml.Marshal(&w.Doc)
	}

	raw := w.Raw
	lineOffsets := buildLineOffsets(raw)
	var edits []byteEdit

	// Collect edits for modified scalar nodes.
	walkNodes(&w.Doc, func(n *yaml.Node) {
		snap, ok := w.Snapshots[n]
		if !ok || n.Line == 0 || n.Kind != yaml.ScalarNode {
			return
		}
		if n.Value == snap.Value && n.LineComment == snap.LineComment {
			return
		}
		start := nodeByteOffset(lineOffsets, n)
		end := lineContentEnd(raw, lineOffsets, n.Line)

		newText := n.Value
		if n.LineComment != "" {
			newText += " # " + n.LineComment
		}
		edits = append(edits, byteEdit{start: start, end: end, text: []byte(newText)})
	})

	// Collect insertions for new nodes (e.g. added permissions).
	edits = append(edits, collectInsertions(w, lineOffsets)...)

	if len(edits) == 0 {
		return raw, nil
	}

	// Sort by start offset descending so splices don't shift earlier offsets.
	sort.Slice(edits, func(i, j int) bool {
		return edits[i].start > edits[j].start
	})

	result := make([]byte, len(raw))
	copy(result, raw)
	for _, e := range edits {
		result = splice(result, e.start, e.end, e.text)
	}
	return result, nil
}

// buildLineOffsets returns a slice where lineOffsets[i] is the byte offset
// of the start of line i+1 (1-indexed in yaml.Node terms).
func buildLineOffsets(raw []byte) []int {
	offsets := []int{0}
	for i, b := range raw {
		if b == '\n' {
			offsets = append(offsets, i+1)
		}
	}
	return offsets
}

// nodeByteOffset converts a yaml.Node's 1-based Line/Column to a byte offset.
func nodeByteOffset(lineOffsets []int, n *yaml.Node) int {
	if n.Line <= 0 || n.Line > len(lineOffsets) {
		return 0
	}
	return lineOffsets[n.Line-1] + n.Column - 1
}

// lineContentEnd returns the byte offset just past the last content byte on
// the given 1-based line (i.e. the position of \n, or EOF).
func lineContentEnd(raw []byte, lineOffsets []int, line int) int {
	if line <= 0 || line > len(lineOffsets) {
		return len(raw)
	}
	start := lineOffsets[line-1]
	idx := bytes.IndexByte(raw[start:], '\n')
	if idx < 0 {
		return len(raw)
	}
	return start + idx
}

// splice replaces data[start:end] with text.
func splice(data []byte, start, end int, text []byte) []byte {
	var buf bytes.Buffer
	buf.Grow(len(data) - (end - start) + len(text))
	buf.Write(data[:start])
	buf.Write(text)
	buf.Write(data[end:])
	return buf.Bytes()
}

// collectInsertions finds key-value pairs added to the root mapping after
// parsing and returns byte edits to insert their YAML text.
func collectInsertions(w *Workflow, lineOffsets []int) []byteEdit {
	if w.Doc.Kind != yaml.DocumentNode || len(w.Doc.Content) == 0 {
		return nil
	}
	root := w.Doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return nil
	}

	var edits []byteEdit
	for i := 0; i < len(root.Content)-1; i += 2 {
		key := root.Content[i]
		val := root.Content[i+1]
		if _, ok := w.Snapshots[key]; ok {
			continue // existing node, skip
		}

		// Determine indentation from a neighbouring existing key.
		indent := ""
		if i >= 2 {
			prevKey := root.Content[i-2]
			if prevKey.Column > 1 {
				indent = strings.Repeat(" ", prevKey.Column-1)
			}
		}

		// Generate the YAML text for the new key-value pair.
		var text string
		if val.Kind == yaml.MappingNode && len(val.Content) == 0 {
			text = indent + key.Value + ": {}\n"
		} else if val.Kind == yaml.ScalarNode {
			text = indent + key.Value + ": " + val.Value + "\n"
		} else {
			out, _ := yaml.Marshal(&yaml.Node{
				Kind:    yaml.MappingNode,
				Content: []*yaml.Node{key, val},
			})
			text = string(out)
		}

		// Insert just before the next existing key's line.
		var insertAt int
		if i+2 < len(root.Content) {
			nextKey := root.Content[i+2]
			if nextKey.Line > 0 && nextKey.Line <= len(lineOffsets) {
				insertAt = lineOffsets[nextKey.Line-1]
			}
		} else {
			insertAt = len(w.Raw)
		}

		edits = append(edits, byteEdit{start: insertAt, end: insertAt, text: []byte(text)})
	}

	return edits
}

func walkNodes(n *yaml.Node, fn func(*yaml.Node)) {
	if n == nil {
		return
	}
	fn(n)
	for _, c := range n.Content {
		walkNodes(c, fn)
	}
}
