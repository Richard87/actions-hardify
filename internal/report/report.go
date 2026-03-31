package report

import (
	"fmt"
	"io"
	"os"
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
		file, loc, perms, old, new string
	}

	// Track which locations have permission findings.
	type locKey struct{ file, loc string }
	permStatus := make(map[locKey]string)
	for _, f := range findings {
		if f.Type == hardener.FindingPermissions {
			file, loc := splitLocation(f)
			permStatus[locKey{file, loc}] = f.Message
		}
	}

	var rows []row
	seen := make(map[locKey]bool)

	for _, f := range findings {
		file, loc := splitLocation(f)
		key := locKey{file, loc}

		if f.Type == hardener.FindingPermissions {
			// Only emit a row if there's no pin/outdated finding at this location.
			// Those rows will carry the perm status themselves.
			continue
		}

		perms := "ok"
		if msg, exists := permStatus[key]; exists {
			perms = msg
		}

		newVal := ""
		if f.Fixed != "" && f.Fixed != f.Current {
			newVal = f.Fixed
		}

		rows = append(rows, row{
			file:  file,
			loc:   loc,
			perms: perms,
			old:   f.Current,
			new:   newVal,
		})
		seen[key] = true
	}

	// Add standalone permission-only findings.
	for _, f := range findings {
		if f.Type != hardener.FindingPermissions {
			continue
		}
		file, loc := splitLocation(f)
		key := locKey{file, loc}
		if seen[key] {
			continue
		}
		rows = append(rows, row{
			file:  file,
			loc:   loc,
			perms: f.Message,
		})
	}

	if len(rows) == 0 {
		fmt.Fprintln(w, "No issues found — workflows are already hardened.")
		return
	}

	// Compute column widths.
	fileW, locW := len("FILE"), len("LOCATION")
	permW, oldW, newW := len("PERMISSIONS"), len("OLD"), len("NEW")
	for _, r := range rows {
		fileW = max(fileW, len(r.file))
		locW = max(locW, len(r.loc))
		permW = max(permW, len(r.perms))
		oldW = max(oldW, len(r.old))
		newW = max(newW, len(r.new))
	}

	permW = min(permW, 60)
	oldW = min(oldW, 50)
	newW = min(newW, 50)

	sep := "+"
	for _, cw := range []int{fileW, locW, permW, oldW, newW} {
		sep += strings.Repeat("-", cw+2) + "+"
	}

	fmtStr := fmt.Sprintf("| %%-%ds | %%-%ds | %%-%ds | %%-%ds | %%-%ds |\n",
		fileW, locW, permW, oldW, newW)

	fmt.Fprintln(w)
	fmt.Fprintln(w, sep)
	fmt.Fprintf(w, fmtStr, "FILE", "LOCATION", "PERMISSIONS", "OLD", "NEW")
	fmt.Fprintln(w, sep)
	for _, r := range rows {
		fmt.Fprintf(w, fmtStr,
			truncate(r.file, fileW),
			truncate(r.loc, locW),
			truncate(r.perms, permW),
			truncate(r.old, oldW),
			truncate(r.new, newW),
		)
	}
	fmt.Fprintln(w, sep)
	fmt.Fprintf(w, "\nTotal: %d finding(s)\n", len(rows))
}

func splitLocation(f hardener.Finding) (file, loc string) {
	file = relPath(f.File)
	var parts []string
	if f.Job != "" {
		parts = append(parts, f.Job)
	}
	if f.Step != "" {
		parts = append(parts, f.Step)
	}
	loc = strings.Join(parts, " > ")
	return file, loc
}

func relPath(path string) string {
	if cwd, err := os.Getwd(); err == nil {
		if rel, err := filepath.Rel(cwd, path); err == nil {
			return rel
		}
	}
	return path
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
