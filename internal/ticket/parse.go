package ticket

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/oklog/ulid/v2"
	"gopkg.in/yaml.v3"
)

// Frontmatter holds YAML frontmatter (known fields). Unknown fields are ignored by yaml.
type Frontmatter struct {
	Schema    string   `yaml:"$schema,omitempty"`
	Title     string   `yaml:"title"`
	Priority  string   `yaml:"priority"`
	Status    string   `yaml:"status"`
	Ready     bool     `yaml:"ready"`
	Assignee  string   `yaml:"assignee,omitempty"`
	BlockedBy []string `yaml:"blocked_by"`
	Subtasks  []string `yaml:"subtasks"`
	Tags      []string `yaml:"tags"`
}

// ParseFile reads a ticket file and returns a Ticket. The path is used for Ticket.Path and errors.
func ParseFile(path string) (*Ticket, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return Parse(path, data)
}

// Parse parses ticket content (path is the ticket file path for Ticket.Path and errors).
func Parse(path string, data []byte) (*Ticket, error) {
	base := filepath.Base(path)
	ulidStr, slug, err := parseFilename(base)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}

	var fm Frontmatter
	rest, err := frontmatter.Parse(strings.NewReader(string(data)), &fm)
	if err != nil {
		return nil, fmt.Errorf("%s: invalid frontmatter: %w", path, err)
	}
	_ = rest // body not used for Ticket struct

	t := &Ticket{
		ULID:      ulidStr,
		Title:     fm.Title,
		Priority:  fm.Priority,
		Status:    fm.Status,
		Ready:     fm.Ready,
		Assignee:  fm.Assignee,
		BlockedBy: fm.BlockedBy,
		Subtasks:  fm.Subtasks,
		Tags:      fm.Tags,
		Path:      path,
	}
	if t.Title == "" {
		t.Title = slug
	}
	if err := Validate(t); err != nil {
		return nil, err
	}
	return t, nil
}

// parseFilename extracts ULID and slug from a basename like "01J9T8ZK1BC5A9JH56T9Y9M1DX-implement-login-api.md".
func parseFilename(base string) (ulidStr, slug string, err error) {
	if len(base) < 28 || !strings.HasSuffix(base, ".md") {
		return "", "", fmt.Errorf("filename must be [ulid]-[slug].md")
	}
	base = base[:len(base)-3] // drop .md
	idx := strings.Index(base, "-")
	if idx != 26 {
		return "", "", fmt.Errorf("filename must be [ulid]-[slug].md (ulid is 26 chars)")
	}
	ulidStr = base[:26]
	slug = base[27:]
	if _, err := ulid.Parse(ulidStr); err != nil {
		return "", "", fmt.Errorf("invalid ULID in filename: %w", err)
	}
	return ulidStr, slug, nil
}

// ParseFrontmatterOnly parses only the YAML frontmatter from data (e.g. for schema or defaults).
// Used when we need to read/write frontmatter without full ticket validation.
func ParseFrontmatterOnly(data []byte) (*Frontmatter, error) {
	var fm Frontmatter
	_, err := frontmatter.Parse(strings.NewReader(string(data)), &fm)
	if err != nil {
		return nil, err
	}
	return &fm, nil
}

// MarshalFrontmatter writes frontmatter to a YAML block (without --- delimiters; caller adds them).
func MarshalFrontmatter(fm *Frontmatter) ([]byte, error) {
	return yaml.Marshal(fm)
}

// ReadFrontmatterAndBody reads a ticket file and returns frontmatter and body (content after the closing ---).
func ReadFrontmatterAndBody(path string) (*Frontmatter, []byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	var fm Frontmatter
	rest, err := frontmatter.Parse(strings.NewReader(string(data)), &fm)
	if err != nil {
		return nil, nil, err
	}
	return &fm, []byte(rest), nil
}

// WriteFrontmatterAndBody writes a ticket file with the given frontmatter and body.
func WriteFrontmatterAndBody(path string, fm *Frontmatter, body []byte) error {
	raw, err := MarshalFrontmatter(fm)
	if err != nil {
		return err
	}
	out := []byte("---\n")
	out = append(out, raw...)
	out = append(out, "---\n\n"...)
	out = append(out, body...)
	return os.WriteFile(path, out, 0644)
}
