package ticket

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"implement login api", "implement-login-api"},
		{"  spaces  ", "spaces"},
		{"UPPERCASE", "uppercase"},
		{"Already-Slug", "already-slug"},
		{"a", "a"},
	}
	for _, tt := range tests {
		got := Slugify(tt.in)
		if got != tt.want {
			t.Errorf("Slugify(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
