package report

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/richard87/actions-hardify/internal/hardener"
)

// Print writes a pretty table report of findings to w.
func Print(w io.Writer, findings []hardener.Finding) {
	if len(findings) == 0 {
		fmt.Fprintln(w, "No issues found — workflows are already hardened.")
		return
	}

	// Build per-location rows, aggregating permission status with pin/outdated info.
	type row struct {
		loc, perms, old, new string
	}

	// Track which locations have permission findings.
	permStatus := make(map[string]string)
	for _, f := range findings {
		if f.Type == hardener.FindingPermissions {
			loc := location(f)
			permStatus[loc] = f.Message
		}
	}

	var rows []row
	seen := make(map[string]bool)

	for _, f := range findings {
		loc := location(f)

		if f.Type == hardener.FindingPermissions {
			// Only emit a row if there's no pin/outdated finding at this location.
			// Those rows will carry the perm status themselves.
			continue
		}

		perms := "ok"
		if msg, exists := permStatus[loc]; exists {
			perms = msg
		}

		newVal := ""
		if f.Fixed != "" && f.Fixed != f.Current {
			newVal = f.Fixed
		}

		rows = append(rows, row{
			loc:   loc,
			perms: perms,
			old:   f.Current,
			new:   newVal,
		})
		seen[loc] = true
	}

	// Add standalone permission-only findings.
	for _, f := range findings {
		if f.Type != hardener.FindingPermissions {
			continue
		}
		loc := location(f)
		if seen[loc] {
			continue
		}
		rows = append(rows, row{
			loc:   loc,
			perms: f.Message,
		})
	}

	if len(rows) == 0 {
		fmt.Fprintln(w, "No issues found — workflows are already hardened.")
		return
	}

	// Compute column widths.
	locW, permW, oldW, newW := len("LOCATION"), len("PERMISSIONS"), len("OLD"), len("NEW")
	for _, r := range rows {
		locW = max(locW, len(r.loc))
		permW = max(permW, len(r.perms))
		oldW = max(oldW, len(r.old))
		newW = max(newW, len(r.new))
	}

	permW = min(permW, 60)
	oldW = min(oldW, 50)
	newW = min(newW, 50)

	sep := "+"
	for _, cw := range []int{locW, permW, oldW, newW} {
		sep += strings.Repeat("-", cw+2) + "+"
	}

	fmtStr := fmt.Sprintf("| %%-%ds | %%-%ds | %%-%ds | %%-%ds |\n",
		locW, permW, oldW, newW)

	fmt.Fprintln(w)
	fmt.Fprintln(w, sep)
	fmt.Fprintf(w, fmtStr, "LOCATION", "PERMISSIONS", "OLD", "NEW")
	fmt.Fprintln(w, sep)
	for _, r := range rows {
		fmt.Fprintf(w, fmtStr,
			truncate(r.loc, locW),
			truncate(r.perms, permW),
			truncate(r.old, oldW),
			truncate(r.new, newW),
		)
	}
	fmt.Fprintln(w, sep)
	fmt.Fprintf(w, "\nTotal: %d finding(s)\n", len(rows))
}

func location(f hardener.Finding) string {
	file := filepath.Base(f.File)
	parts := []string{file}
	if f.Job != "" {
		parts = append(parts, f.Job)
	}
	if f.Step != "" {
		parts = append(parts, f.Step)
	}
	return strings.Join(parts, " > ")
}

func truncate(s string, maxW int) string {
	if len(s) <= maxW {
		return s
	}
	if maxW <= 3 {
		return s[:maxW]
	}
	return s[:maxW-3] + "..."
}
