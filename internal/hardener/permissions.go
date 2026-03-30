package hardener

import (
	"github.com/richard87/actions-hardify/internal/workflow"
	"gopkg.in/yaml.v3"
)

// checkPermissions verifies that the workflow and each job have restricted permissions.
func checkPermissions(w *workflow.Workflow) []Finding {
	var findings []Finding

	// Check top-level permissions
	if w.Permissions.IsEmpty {
		findings = append(findings, Finding{
			File:    w.Path,
			Type:    FindingPermissions,
			Message: "workflow has no top-level permissions; add 'permissions: {}' to restrict GITHUB_TOKEN",
		})
	} else if w.Permissions.IsAll {
		findings = append(findings, Finding{
			File:    w.Path,
			Type:    FindingPermissions,
			Current: w.Permissions.Values["all"],
			Fixed:   "{}",
			Message: "workflow uses overly broad permissions",
		})
	}

	// Check each job
	for _, job := range w.Jobs {
		if job.Permissions.IsEmpty && w.Permissions.IsEmpty {
			findings = append(findings, Finding{
				File:    w.Path,
				Job:     job.ID,
				Type:    FindingPermissions,
				Message: "job has no permissions and no top-level permissions are set",
			})
		}
	}

	return findings
}

// applyPermissionFixes adds a restrictive top-level permissions block if missing.
func applyPermissionFixes(w *workflow.Workflow) {
	if !w.Permissions.IsEmpty {
		return
	}

	root := w.Doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return
	}

	// Insert "permissions: {}" right after the "name" or "on" key, or at the top.
	insertIdx := 0
	for i := 0; i < len(root.Content)-1; i += 2 {
		key := root.Content[i].Value
		if key == "name" || key == "on" {
			insertIdx = i + 2
		}
	}

	keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: "permissions", Tag: "!!str"}
	valNode := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}

	newContent := make([]*yaml.Node, 0, len(root.Content)+2)
	newContent = append(newContent, root.Content[:insertIdx]...)
	newContent = append(newContent, keyNode, valNode)
	newContent = append(newContent, root.Content[insertIdx:]...)
	root.Content = newContent
}
