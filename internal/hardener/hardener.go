package hardener

import (
	"context"
	"fmt"

	"github.com/richard87/actions-hardify/internal/github"
	"github.com/richard87/actions-hardify/internal/workflow"
)

// FindingType categorizes findings.
type FindingType int

const (
	FindingPermissions FindingType = iota
	FindingUnpinned
	FindingOutdated
)

func (t FindingType) String() string {
	switch t {
	case FindingPermissions:
		return "permissions"
	case FindingUnpinned:
		return "unpinned"
	case FindingOutdated:
		return "outdated"
	default:
		return "unknown"
	}
}

// Finding describes a single issue found during hardening.
type Finding struct {
	File    string
	Job     string
	Step    string
	Type    FindingType
	Current string
	Fixed   string
	Message string
}

// HardenAll runs all hardening checks on the given workflows.
// Returns findings and modifies the workflow ASTs in-place (unless dryRun).
func HardenAll(ctx context.Context, workflows []*workflow.Workflow, gh *github.Client, dryRun bool) ([]Finding, error) {
	var allFindings []Finding
	for _, w := range workflows {
		findings, err := hardenWorkflow(ctx, w, gh, dryRun)
		if err != nil {
			return allFindings, fmt.Errorf("hardening %s: %w", w.Path, err)
		}
		allFindings = append(allFindings, findings...)
	}
	return allFindings, nil
}

func hardenWorkflow(ctx context.Context, w *workflow.Workflow, gh *github.Client, dryRun bool) ([]Finding, error) {
	var findings []Finding

	// 1. Check permissions
	findings = append(findings, checkPermissions(w)...)

	// 2. Pin actions to SHA + check outdated
	actionFindings, err := checkActions(ctx, w, gh, dryRun)
	if err != nil {
		return findings, err
	}
	findings = append(findings, actionFindings...)

	// Apply permission fixes if not dry run
	if !dryRun {
		applyPermissionFixes(w)
	}

	return findings, nil
}
