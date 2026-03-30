package main

import (
	"context"
	"fmt"
	"os"

	"github.com/richard87/actions-hardify/internal/github"
	"github.com/richard87/actions-hardify/internal/hardener"
	"github.com/richard87/actions-hardify/internal/report"
	"github.com/richard87/actions-hardify/internal/scanner"
	"github.com/richard87/actions-hardify/internal/workflow"
	"github.com/urfave/cli/v3"
)

var version = "dev"

func main() {
	cmd := &cli.Command{
		Name:    "actions-hardify",
		Usage:   "Harden GitHub Actions workflow files",
		Version: version,
		Description: `What it does:
  • Restrict GITHUB_TOKEN permissions to least-privilege
  • Pin third-party actions to full-length commit SHAs
  • Detect outdated action versions and suggest upgrades`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dir",
				Aliases: []string{"d"},
				Value:   ".",
				Usage:   "directory to scan for GitHub Actions workflows",
			},
			&cli.StringFlag{
				Name:    "token",
				Aliases: []string{"t"},
				Usage:   "GitHub personal access token",
				Sources: cli.EnvVars("GITHUB_TOKEN"),
			},
			&cli.BoolFlag{
				Name:    "dry-run",
				Aliases: []string{"n"},
				Usage:   "print a report of findings without modifying files",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return run(ctx, cmd.String("dir"), cmd.String("token"), cmd.Bool("dry-run"))
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, dir, token string, dryRun bool) error {
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
