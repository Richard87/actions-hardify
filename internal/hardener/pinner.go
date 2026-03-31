package hardener

import (
	"context"
	"fmt"
	"strings"

	"github.com/richard87/actions-hardify/internal/github"
	"github.com/richard87/actions-hardify/internal/workflow"
)

// checkActions pins unpinned actions to SHA and flags outdated versions.
func checkActions(ctx context.Context, w *workflow.Workflow, gh *github.Client, dryRun bool) ([]Finding, error) {
	var findings []Finding

	// Deduplicate: resolve each unique owner/repo@ref once.
	type resolved struct {
		sha       string
		latest    string
		latestSHA string
	}
	cache := make(map[string]*resolved)

	for _, job := range w.Jobs {
		// Collect all action refs: job-level uses (reusable workflows) + step-level uses.
		var refs []*workflow.ActionRef
		if job.Uses != nil {
			refs = append(refs, job.Uses)
		}
		for _, step := range job.Steps {
			if step.Uses != nil {
				refs = append(refs, step.Uses)
			}
		}
		for _, ref := range refs {
			cacheKey := ref.String()
			res, ok := cache[cacheKey]
			if !ok {
				cache[cacheKey] = &resolved{}
				res = cache[cacheKey]

				// Resolve the ref to a full SHA
				if !ref.IsSHA() {
					sha, err := gh.ResolveTagSHA(ctx, ref.Owner, ref.Repo, ref.Ref)
					if err != nil {
						return findings, fmt.Errorf("resolving %s: %w", ref, err)
					}
					res.sha = sha
				}

				// Check for newer versions
				latest, err := latestVersion(ctx, gh, ref)
				if err == nil && latest != "" && latest != ref.Ref {
					res.latest = latest
					// Resolve the latest version's SHA so we can pin to it.
					latestSHA, err := gh.ResolveTagSHA(ctx, ref.Owner, ref.Repo, latest)
					if err == nil {
						res.latestSHA = latestSHA
					}
				}
			}

			// Determine the SHA and tag to pin to: prefer latest if available.
			pinSHA := res.sha
			pinTag := ref.Ref
			if res.latest != "" && res.latestSHA != "" {
				pinSHA = res.latestSHA
				pinTag = res.latest
			}

			// Report unpinned action
			if !ref.IsSHA() && pinSHA != "" {
				findings = append(findings, Finding{
					File:    w.Path,
					Job:     job.ID,
					Type:    FindingUnpinned,
					Current: ref.Ref,
					Fixed:   pinTag,
					Message: fmt.Sprintf("pin %s to commit SHA", ref),
				})
				if !dryRun {
					ref.Node.Value = ref.Full() + "@" + pinSHA
					ref.Node.LineComment = pinTag
				}
			}

			// Report outdated action
			if res.latest != "" {
				currentTag := ref.Ref
				if ref.IsSHA() && ref.Node.LineComment != "" {
					currentTag = strings.TrimSpace(ref.Node.LineComment)
				}
				findings = append(findings, Finding{
					File:    w.Path,
					Job:     job.ID,
					Type:    FindingOutdated,
					Current: currentTag,
					Fixed:   res.latest,
					Message: fmt.Sprintf("%s can be upgraded from %s to %s", ref.Full(), currentTag, res.latest),
				})
				// Update already-pinned actions to the latest version.
				if !dryRun && ref.IsSHA() && res.latestSHA != "" {
					ref.Node.Value = ref.Full() + "@" + res.latestSHA
					ref.Node.LineComment = res.latest
				}
			}
		}
	}

	return findings, nil
}

func latestVersion(ctx context.Context, gh *github.Client, ref *workflow.ActionRef) (string, error) {
	// Try the releases endpoint first.
	latest, err := gh.LatestRelease(ctx, ref.Owner, ref.Repo)
	if err == nil && latest != "" {
		return latest, nil
	}
	// Fallback to tags.
	return gh.LatestTag(ctx, ref.Owner, ref.Repo)
}
