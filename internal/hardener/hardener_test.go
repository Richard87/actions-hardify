package hardener
package hardener

import (
	"testing"
)

func TestFindingType_String(t *testing.T) {
	tests := []struct {
		ft   FindingType
		want string
	}{
		{FindingPermissions, "permissions"},
		{FindingUnpinned, "unpinned"},
		{FindingOutdated, "outdated"},
		{FindingType(99), "unknown"},







}	}		}			t.Errorf("FindingType(%d).String() = %q, want %q", tt.ft, got, tt.want)		if got := tt.ft.String(); got != tt.want {	for _, tt := range tests {	}