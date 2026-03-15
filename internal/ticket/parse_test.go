package ticket

import (
	"path/filepath"
	"testing"
)

func TestParseFilename(t *testing.T) {
	tests := []struct {
		base      string
		wantULID  string
		wantSlug  string
		wantError bool
	}{
		{"01J9T8ZK1BC5A9JH56T9Y9M1DX-implement-login-api.md", "01J9T8ZK1BC5A9JH56T9Y9M1DX", "implement-login-api", false},
		{"01J9T8ZK1BC5A9JH56T9Y9M1DX-single.md", "01J9T8ZK1BC5A9JH56T9Y9M1DX", "single", false},
		{"short-md", "", "", true},
		{"01J9T8ZK1BC5A9JH56T9Y9M1DX-.md", "01J9T8ZK1BC5A9JH56T9Y9M1DX", "", false},
	}
	for _, tt := range tests {
		ulid, slug, err := parseFilename(tt.base)
		if (err != nil) != tt.wantError {
			t.Errorf("parseFilename(%q) err = %v, wantError %v", tt.base, err, tt.wantError)
			continue
		}
		if !tt.wantError && (ulid != tt.wantULID || slug != tt.wantSlug) {
			t.Errorf("parseFilename(%q) = %q, %q; want %q, %q", tt.base, ulid, slug, tt.wantULID, tt.wantSlug)
		}
	}
}

func TestParse(t *testing.T) {
	path := filepath.Join(t.TempDir(), "01J9T8ZK1BC5A9JH56T9Y9M1DX-foo.md")
	content := []byte(`---
priority: p1
status: todo
ready: true
---

body
`)
	ticket, err := Parse(path, content)
	if err != nil {
		t.Fatal(err)
	}
	if ticket.ULID != "01J9T8ZK1BC5A9JH56T9Y9M1DX" || ticket.Priority != "p1" || ticket.Status != "todo" || !ticket.Ready {
		t.Errorf("got %+v", ticket)
	}
}
