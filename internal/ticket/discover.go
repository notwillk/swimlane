package ticket

import (
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/notwillk/swimlane/internal/config"
)

// GlobPaths returns paths of all files matching the config tickets glob.
// The glob is resolved relative to the config file's directory (ConfigDir).
func GlobPaths(cfg *config.Config) ([]string, error) {
	pattern := filepath.ToSlash(cfg.Tickets)
	matches, err := doublestar.Glob(os.DirFS(cfg.ConfigDir), pattern, doublestar.WithFilesOnly())
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(matches))
	for _, m := range matches {
		paths = append(paths, filepath.Join(cfg.ConfigDir, filepath.FromSlash(m)))
	}
	return paths, nil
}

// Discover finds all ticket files matching the config glob and parses them.
// The glob is resolved relative to the config file's directory (ConfigDir).
// Parse errors (e.g. invalid frontmatter) cause the function to return an error.
func Discover(cfg *config.Config) ([]*Ticket, error) {
	paths, err := GlobPaths(cfg)
	if err != nil {
		return nil, err
	}
	var tickets []*Ticket
	for _, path := range paths {
		t, err := ParseFile(path)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}
