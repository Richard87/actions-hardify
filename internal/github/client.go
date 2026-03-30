package github

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

// Client communicates with the GitHub REST API.
type Client struct {
	HTTPClient *http.Client
	Token      string
	BaseURL    string // default: https://api.github.com
}

// NewClient creates a Client, optionally using $GITHUB_TOKEN for authentication.
func NewClient(token string) *Client {
	if token == "" {
		log.Println("warning: no GitHub token provided; API rate limits may apply")
	}
	return &Client{
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		Token:      token,
		BaseURL:    "https://api.github.com",
	}
}

// Tag represents a Git tag from the GitHub API.
type Tag struct {
	Name   string `json:"name"`
	Commit struct {
		SHA string `json:"sha"`
	} `json:"commit"`
}

// Release represents a GitHub release.
type Release struct {
	TagName     string `json:"tag_name"`
	Prerelease  bool   `json:"prerelease"`
	Draft       bool   `json:"draft"`
	PublishedAt string `json:"published_at"`
}

// ResolveTagSHA looks up the full commit SHA for a tag or branch ref.
func (c *Client) ResolveTagSHA(ctx context.Context, owner, repo, ref string) (string, error) {
	// Try as a Git ref first (works for both tags and branches).
	url := fmt.Sprintf("%s/repos/%s/%s/git/ref/tags/%s", c.BaseURL, owner, repo, ref)
	var refResp struct {
		Object struct {
			SHA  string `json:"sha"`
			Type string `json:"type"`
		} `json:"object"`
	}
	if err := c.get(ctx, url, &refResp); err == nil {
		sha := refResp.Object.SHA
		// If it's an annotated tag, dereference to the commit.
		if refResp.Object.Type == "tag" {
			return c.dereferenceTag(ctx, owner, repo, sha)
		}
		return sha, nil
	}

	// Fallback: search in tags list (handles lightweight tags).
	tags, err := c.ListTags(ctx, owner, repo)
	if err != nil {
		return "", fmt.Errorf("resolving ref %s for %s/%s: %w", ref, owner, repo, err)
	}
	for _, t := range tags {
		if t.Name == ref {
			return t.Commit.SHA, nil
		}
	}
	return "", fmt.Errorf("ref %q not found in %s/%s", ref, owner, repo)
}

func (c *Client) dereferenceTag(ctx context.Context, owner, repo, sha string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/git/tags/%s", c.BaseURL, owner, repo, sha)
	var tagObj struct {
		Object struct {
			SHA string `json:"sha"`
		} `json:"object"`
	}
	if err := c.get(ctx, url, &tagObj); err != nil {
		return sha, nil // fallback to the tag SHA
	}
	return tagObj.Object.SHA, nil
}

// ListTags returns all tags for a repo (paginated, up to 100).
func (c *Client) ListTags(ctx context.Context, owner, repo string) ([]Tag, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/tags?per_page=100", c.BaseURL, owner, repo)
	var tags []Tag
	if err := c.get(ctx, url, &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

// LatestRelease returns the latest non-prerelease, non-draft release tag.
func (c *Client) LatestRelease(ctx context.Context, owner, repo string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", c.BaseURL, owner, repo)
	var rel Release
	if err := c.get(ctx, url, &rel); err != nil {
		return "", err
	}
	return rel.TagName, nil
}

// LatestTag returns the most recent tag by version sorting.
func (c *Client) LatestTag(ctx context.Context, owner, repo string) (string, error) {
	tags, err := c.ListTags(ctx, owner, repo)
	if err != nil {
		return "", err
	}
	if len(tags) == 0 {
		return "", fmt.Errorf("no tags found for %s/%s", owner, repo)
	}
	// Sort tags with semver-like ordering (latest first by name).
	sort.Slice(tags, func(i, j int) bool {
		return compareSemverish(tags[i].Name, tags[j].Name) > 0
	})
	return tags[0].Name, nil
}

func (c *Client) get(ctx context.Context, url string, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("GitHub API %s returned %d", url, resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// compareSemverish does a rough semver comparison; returns >0 if a > b.
func compareSemverish(a, b string) int {
	a = strings.TrimPrefix(a, "v")
	b = strings.TrimPrefix(b, "v")
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	for i := 0; i < len(aParts) && i < len(bParts); i++ {
		if aParts[i] != bParts[i] {
			if aParts[i] > bParts[i] {
				return 1
			}
			return -1
		}
	}
	return len(aParts) - len(bParts)
}
