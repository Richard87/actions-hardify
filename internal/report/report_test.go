package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/richard87/actions-hardify/internal/hardener"
)

func TestPrint_NoFindings(t *testing.T) {
	var buf bytes.Buffer
	Print(&buf, nil)
	if !strings.Contains(buf.String(), "No issues found") {
		t.Errorf("expected 'No issues found', got %q", buf.String())
	}
}

func TestPrint_WithFindings(t *testing.T) {
	findings := []hardener.Finding{
		{
			File:    "ci.yaml",
			Job:     "build",
			Step:    "Checkout",
			Type:    hardener.FindingUnpinned,
			Current: "actions/checkout@v4",
			Fixed:   "actions/checkout@abc123 # v4",
			Message: "pin actions/checkout to commit SHA",
		},
		{
			File:    "ci.yaml",
			Type:    hardener.FindingPermissions,
			Message: "workflow has no top-level permissions",
		},
	}

	var buf bytes.Buffer
	Print(&buf, findings)
	output := buf.String()

	if !strings.Contains(output, "LOCATION") {
		t.Error("expected table header 'LOCATION'")
	}
	if !strings.Contains(output, "PERMISSIONS") {
		t.Error("expected table header 'PERMISSIONS'")
	}
	if !strings.Contains(output, "actions/checkout@v4") {
		t.Error("expected action reference in output")
	}
	if !strings.Contains(output, "Total:") {
		t.Error("expected 'Total:' in output")
	}
}

func TestPrint_PermissionsOnly(t *testing.T) {
	findings := []hardener.Finding{
		{
			File:    "ci.yaml",
			Type:    hardener.FindingPermissions,
			Message: "workflow has no top-level permissions",
		},
	}

	var buf bytes.Buffer
	Print(&buf, findings)
	output := buf.String()

	if !strings.Contains(output, "ci.yaml") {
		t.Error("expected file name in output")
	}
	if !strings.Contains(output, "Total: 1") {
		t.Errorf("expected 'Total: 1', got %q", output)
	}
}

func TestPrint_OutdatedFinding(t *testing.T) {
	findings := []hardener.Finding{
		{
			File:    "ci.yaml",
			Job:     "build",
			Step:    "Setup Go",
			Type:    hardener.FindingOutdated,
			Current: "v4",
			Fixed:   "v5",
			Message: "actions/setup-go can be upgraded from v4 to v5",
		},
	}

	var buf bytes.Buffer
	Print(&buf, findings)
	output := buf.String()

	if !strings.Contains(output, "v4") {
		t.Error("expected old version in output")
	}
	if !strings.Contains(output, "v5") {
		t.Error("expected new version in output")
	}
}
