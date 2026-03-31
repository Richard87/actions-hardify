package workflow

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Parse reads a workflow YAML file and returns a Workflow.
func Parse(path string) (*Workflow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	w := &Workflow{
		Path:      path,
		Raw:       data,
		Doc:       doc,
		Snapshots: snapshotNodes(&doc),
	}
	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 {
		return w, nil
	}
	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return w, nil
	}
	w.Permissions = findPermissions(root)
	w.Jobs = findJobs(root)
	return w, nil
}

func findPermissions(mapping *yaml.Node) *Permissions {
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		key := mapping.Content[i]
		val := mapping.Content[i+1]
		if key.Value == "permissions" {
			return parsePermissions(val)
		}
	}
	return &Permissions{IsEmpty: true}
}

func parsePermissions(node *yaml.Node) *Permissions {
	p := &Permissions{Node: node}
	switch node.Kind {
	case yaml.ScalarNode:
		p.IsAll = true
		p.Values = map[string]string{"all": node.Value}
	case yaml.MappingNode:
		p.Values = make(map[string]string, len(node.Content)/2)
		for i := 0; i < len(node.Content)-1; i += 2 {
			p.Values[node.Content[i].Value] = node.Content[i+1].Value
		}
	}
	return p
}

func findJobs(root *yaml.Node) []Job {
	for i := 0; i < len(root.Content)-1; i += 2 {
		key := root.Content[i]
		val := root.Content[i+1]
		if key.Value != "jobs" || val.Kind != yaml.MappingNode {
			continue
		}
		return parseJobs(val)
	}
	return nil
}

func parseJobs(jobsNode *yaml.Node) []Job {
	var jobs []Job
	for i := 0; i < len(jobsNode.Content)-1; i += 2 {
		key := jobsNode.Content[i]
		val := jobsNode.Content[i+1]
		if val.Kind != yaml.MappingNode {
			continue
		}
		job := Job{
			ID:          key.Value,
			Node:        val,
			Permissions: findPermissions(val),
			Steps:       findSteps(val),
		}
		jobs = append(jobs, job)
	}
	return jobs
}

func findSteps(jobNode *yaml.Node) []Step {
	for i := 0; i < len(jobNode.Content)-1; i += 2 {
		key := jobNode.Content[i]
		val := jobNode.Content[i+1]
		if key.Value != "steps" || val.Kind != yaml.SequenceNode {
			continue
		}
		return parseSteps(val)
	}
	return nil
}

func parseSteps(stepsNode *yaml.Node) []Step {
	var steps []Step
	for _, item := range stepsNode.Content {
		if item.Kind != yaml.MappingNode {
			continue
		}
		s := Step{Node: item}
		for i := 0; i < len(item.Content)-1; i += 2 {
			k := item.Content[i]
			v := item.Content[i+1]
			switch k.Value {
			case "name":
				s.Name = v.Value
			case "uses":
				s.Uses = parseActionRef(v)
			}
		}
		steps = append(steps, s)
	}
	return steps
}

// parseActionRef parses "owner/repo@ref" or "owner/repo/path@ref".
func parseActionRef(node *yaml.Node) *ActionRef {
	raw := node.Value
	// Skip local actions (./path)
	if strings.HasPrefix(raw, ".") || strings.HasPrefix(raw, "docker://") {
		return nil
	}
	atIdx := strings.LastIndex(raw, "@")
	if atIdx < 0 {
		return nil
	}
	slug := raw[:atIdx]
	ref := raw[atIdx+1:]
	parts := strings.SplitN(slug, "/", 3)
	if len(parts) < 2 {
		return nil
	}
	a := &ActionRef{
		Owner: parts[0],
		Repo:  parts[1],
		Ref:   ref,
		Node:  node,
	}
	if len(parts) == 3 {
		a.Path = "/" + parts[2]
	}
	return a
}

// CollectActions returns all unique action references from a workflow.
func CollectActions(w *Workflow) []*ActionRef {
	var refs []*ActionRef
	seen := make(map[string]bool)
	for _, job := range w.Jobs {
		for _, step := range job.Steps {
			if step.Uses == nil {
				continue
			}
			key := fmt.Sprintf("%s@%s:%d", step.Uses.Full(), step.Uses.Ref, step.Uses.Node.Line)
			if seen[key] {
				continue
			}
			seen[key] = true
			refs = append(refs, step.Uses)
		}
	}
	return refs
}

// snapshotNodes records every node's Value and LineComment at parse time
// so the writer can detect which nodes were modified.
func snapshotNodes(n *yaml.Node) map[*yaml.Node]NodeSnapshot {
	m := make(map[*yaml.Node]NodeSnapshot)
	var walk func(*yaml.Node)
	walk = func(n *yaml.Node) {
		if n == nil {
			return
		}
		m[n] = NodeSnapshot{Value: n.Value, LineComment: n.LineComment}
		for _, c := range n.Content {
			walk(c)
		}
	}
	walk(n)
	return m
}
