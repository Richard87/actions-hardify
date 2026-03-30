package github

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompareSemverish(t *testing.T) {
	tests := []struct {
		a, b string
		want int // >0 means a > b, <0 means a < b, 0 means equal
	}{
		{"v1.0.0", "v0.9.0", 1},
		{"v0.9.0", "v1.0.0", -1},
		{"v1.0.0", "v1.0.0", 0},
		{"v2.0.0", "v1.9.9", 1},
		{"v1.0", "v1.0.0", -1},
		{"1.0.0", "0.9.0", 1},
	}
	for _, tt := range tests {
		got := compareSemverish(tt.a, tt.b)
		switch {
		case tt.want > 0 && got <= 0:
			t.Errorf("compareSemverish(%q, %q) = %d, want >0", tt.a, tt.b, got)
		case tt.want < 0 && got >= 0:
			t.Errorf("compareSemverish(%q, %q) = %d, want <0", tt.a, tt.b, got)
		case tt.want == 0 && got != 0:
			t.Errorf("compareSemverish(%q, %q) = %d, want 0", tt.a, tt.b, got)
		}
	}
}

func TestResolveTagSHA(t *testing.T) {
	expectedSHA := "abc123def456abc123def456abc123def456abc1"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/actions/checkout/git/ref/tags/v4" {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"object": map[string]string{"sha": expectedSHA, "type": "commit"},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := &Client{
		HTTPClient: srv.Client(),
		BaseURL:    srv.URL,
	}

	sha, err := c.ResolveTagSHA(context.Background(), "actions", "checkout", "v4")
	if err != nil {
		t.Fatalf("ResolveTagSHA() error: %v", err)
	}
	if sha != expectedSHA {
		t.Errorf("SHA = %q, want %q", sha, expectedSHA)
	}
}

func TestResolveTagSHA_AnnotatedTag(t *testing.T) {
	tagSHA := "tag0000000000000000000000000000000000000"
	commitSHA := "commit00000000000000000000000000000000000"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/actions/checkout/git/ref/tags/v4":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"object": map[string]string{"sha": tagSHA, "type": "tag"},
			})
		case "/repos/actions/checkout/git/tags/" + tagSHA:
			json.NewEncoder(w).Encode(map[string]interface{}{
				"object": map[string]string{"sha": commitSHA},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	c := &Client{
		HTTPClient: srv.Client(),
		BaseURL:    srv.URL,
	}

	sha, err := c.ResolveTagSHA(context.Background(), "actions", "checkout", "v4")
	if err != nil {
		t.Fatalf("ResolveTagSHA() error: %v", err)
	}
	if sha != commitSHA {
		t.Errorf("SHA = %q, want %q", sha, commitSHA)
	}
}

func TestResolveTagSHA_FallbackToTags(t *testing.T) {
	expectedSHA := "fallback0000000000000000000000000000000"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/actions/checkout/tags":
			json.NewEncoder(w).Encode([]Tag{
				{Name: "v4", Commit: struct {
					SHA string `json:"sha"`
				}{SHA: expectedSHA}},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	c := &Client{
		HTTPClient: srv.Client(),
		BaseURL:    srv.URL,
	}

	sha, err := c.ResolveTagSHA(context.Background(), "actions", "checkout", "v4")
	if err != nil {
		t.Fatalf("ResolveTagSHA() error: %v", err)
	}
	if sha != expectedSHA {
		t.Errorf("SHA = %q, want %q", sha, expectedSHA)
	}
}

func TestLatestRelease(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/actions/checkout/releases/latest" {
			json.NewEncoder(w).Encode(Release{TagName: "v5"})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := &Client{
		HTTPClient: srv.Client(),
		BaseURL:    srv.URL,
	}

	tag, err := c.LatestRelease(context.Background(), "actions", "checkout")
	if err != nil {
		t.Fatalf("LatestRelease() error: %v", err)
	}
	if tag != "v5" {
		t.Errorf("tag = %q, want %q", tag, "v5")
	}
}

func TestLatestTag(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/actions/checkout/tags" {
			json.NewEncoder(w).Encode([]Tag{
				{Name: "v3"},
				{Name: "v5"},
				{Name: "v4"},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := &Client{
		HTTPClient: srv.Client(),
		BaseURL:    srv.URL,
	}

	tag, err := c.LatestTag(context.Background(), "actions", "checkout")
	if err != nil {
		t.Fatalf("LatestTag() error: %v", err)
	}
	if tag != "v5" {
		t.Errorf("tag = %q, want %q", tag, "v5")
	}
}

func TestLatestTag_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/actions/checkout/tags" {
			json.NewEncoder(w).Encode([]Tag{})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := &Client{
		HTTPClient: srv.Client(),
		BaseURL:    srv.URL,
	}

	_, err := c.LatestTag(context.Background(), "actions", "checkout")
	if err == nil {
		t.Fatal("expected error for empty tags")
	}
}

func TestClient_AuthHeader(t *testing.T) {
	var authHeader string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		json.NewEncoder(w).Encode(Release{TagName: "v1"})
	}))
	defer srv.Close()

	c := &Client{
		HTTPClient: srv.Client(),
		Token:      "test-token-123",
		BaseURL:    srv.URL,
	}

	_, _ = c.LatestRelease(context.Background(), "owner", "repo")
	if authHeader != "Bearer test-token-123" {
		t.Errorf("Authorization = %q, want %q", authHeader, "Bearer test-token-123")
	}
}

func TestClient_NoAuth(t *testing.T) {
	var authHeader string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		json.NewEncoder(w).Encode(Release{TagName: "v1"})
	}))
	defer srv.Close()

	c := &Client{
		HTTPClient: srv.Client(),
		BaseURL:    srv.URL,
	}

	_, _ = c.LatestRelease(context.Background(), "owner", "repo")
	if authHeader != "" {
		t.Errorf("Authorization = %q, want empty", authHeader)
	}
}
