package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/richard87/actions-hardify/internal/github"
	"github.com/richard87/actions-hardify/internal/hardener"
	"github.com/richard87/actions-hardify/internal/report"
	"github.com/richard87/actions-hardify/internal/scanner"
	"github.com/richard87/actions-hardify/internal/workflow"
)

func main() {
	dir := flag.String("dir", ".", "directory to scan for GitHub Actions workflows")
	token := flag.String("token", os.Getenv("GITHUB_TOKEN"), "GitHub personal access token (defaults to $GITHUB_TOKEN)")
	dryRun := flag.Bool("dry-run", false, "print a report of findings without modifying files")
	flag.Parse()

	if err := run(*dir, *token, *dryRun); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(dir, token string, dryRun bool) error {
	ctx := context.Background()

	// 1. Find workflow files
	paths, err := scanner.FindWorkflows(dir)
	if err != nil {
		return fmt.Errorf("scanning %s: %w", dir, err)
	}
	if len(paths) == 0 {
		fmt.Println("No workflow files found.")
		return nil
	}
	fmt.Printf("Found %d workflow file(s)\n", len(paths))

	// 2. Parse workflows
	var workflows []*workflow.Workflow
	for _, p := range paths {
		w, err := workflow.Parse(p)
		if err != nil {
			return err
		}
		workflows = append(workflows, w)
	}

	// 3. Run hardening checks
	gh := github.NewClient(token)
	findings, err := hardener.HardenAll(ctx, workflows, gh, dryRun)
	if err != nil {
		return err
	}

	// 4. Print report
	report.Print(os.Stdout, findings)

	// 5. Write modified files (unless dry-run)
	if !dryRun {
		for _, w := range workflows {
			if err := workflow.Write(w); err != nil {
				return err
			}
		}
		fmt.Println("\n✅ Workflows hardened successfully.")
	} else {
		fmt.Println("\n(dry-run mode — no files were modified)")
	}

	return nil
}
