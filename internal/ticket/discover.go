package ticket

import (
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/notwillk/swimlane/internal/config"
)

// Discover finds all ticket files matching the config glob and parses them.
// The glob is resolved relative to the config file's directory (ConfigDir).
// Parse errors (e.g. invalid frontmatter) cause the function to return an error.
func Discover(cfg *config.Config) ([]*Ticket, error) {
	pattern := filepath.ToSlash(cfg.Tickets)
	matches, err := doublestar.Glob(os.DirFS(cfg.ConfigDir), pattern, doublestar.WithFilesOnly())
	if err != nil {
		return nil, err
	}
	var tickets []*Ticket
	for _, m := range matches {
		path := filepath.Join(cfg.ConfigDir, filepath.FromSlash(m))
		t, err := ParseFile(path)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}
