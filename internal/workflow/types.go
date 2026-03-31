package workflow

import "gopkg.in/yaml.v3"

// NodeSnapshot stores the original value of a yaml.Node at parse time.
type NodeSnapshot struct {
	Value       string
	LineComment string
}

// Workflow represents a parsed GitHub Actions workflow file.
type Workflow struct {
	Path      string // filesystem path to the YAML file
	Raw       []byte // original bytes from disk
	Doc       yaml.Node
	Snapshots map[*yaml.Node]NodeSnapshot // original node values

	Permissions *Permissions
	Jobs        []Job
}

// Permissions models the top-level or job-level "permissions" block.
type Permissions struct {
	Node    *yaml.Node
	IsEmpty bool              // true when the key is absent
	Values  map[string]string // e.g. {"contents": "read", "packages": "write"}
	IsAll   bool              // true when set to a bare string like "write-all"
}

// Job is a single job inside the workflow.
type Job struct {
	ID          string
	Node        *yaml.Node
	Permissions *Permissions
	Uses        *ActionRef // non-nil for reusable workflow jobs
	Steps       []Step
}

// Step represents one step inside a job.
type Step struct {
	Name string
	Uses *ActionRef // nil when the step is a `run:` step
	Node *yaml.Node
}

// ActionRef is a parsed `uses:` value like "actions/checkout@v4".
type ActionRef struct {
	Owner string // e.g. "actions"
	Repo  string // e.g. "checkout"
	Path  string // e.g. "" or "/subdir"
	Ref   string // e.g. "v4" or a full SHA
	Node  *yaml.Node
}

// IsSHA returns true when the ref is already a 40-character hex SHA.
func (a *ActionRef) IsSHA() bool {
	if len(a.Ref) != 40 {
		return false
	}
	for _, c := range a.Ref {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}

// Full returns "owner/repo" (with optional path).
func (a *ActionRef) Full() string {
	return a.Owner + "/" + a.Repo + a.Path
}

// String returns the full uses value like "actions/checkout@v4".
func (a *ActionRef) String() string {
	return a.Full() + "@" + a.Ref
}
