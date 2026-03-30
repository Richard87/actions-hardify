package workflow

import "testing"

func TestActionRef_IsSHA(t *testing.T) {
	tests := []struct {
		ref  string
		want bool
	}{
		{"abc123", false},
		{"v4", false},
		{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", true},
		{"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", false},
		{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", false},
		{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", false},
		{"ghijklmnopqrstuvwxyzghijklmnopqrstuvwxyz", false},
	}
	for _, tt := range tests {
		a := &ActionRef{Ref: tt.ref}
		if got := a.IsSHA(); got != tt.want {
			t.Errorf("IsSHA(%q) = %v, want %v", tt.ref, got, tt.want)
		}
	}
}

func TestActionRef_Full(t *testing.T) {
	tests := []struct {
		owner, repo, path, want string
	}{
		{"actions", "checkout", "", "actions/checkout"},
		{"aws-actions", "configure-aws-credentials", "/assume-role", "aws-actions/configure-aws-credentials/assume-role"},
	}
	for _, tt := range tests {
		a := &ActionRef{Owner: tt.owner, Repo: tt.repo, Path: tt.path}
		if got := a.Full(); got != tt.want {
			t.Errorf("Full() = %q, want %q", got, tt.want)
		}
	}
}

func TestActionRef_String(t *testing.T) {
	a := &ActionRef{Owner: "actions", Repo: "checkout", Ref: "v4"}
	want := "actions/checkout@v4"
	if got := a.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
