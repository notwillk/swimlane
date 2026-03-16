package commands

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"

	"github.com/notwillk/swimlane/internal/ticket"
	"github.com/oklog/ulid/v2"
	"github.com/spf13/cobra"
)

func NewNew() *cobra.Command {
	return &cobra.Command{
		Use:   "new [title]",
		Short: "Create a new ticket",
		Args:  cobra.ExactArgs(1),
		RunE:  runNew,
	}
}

func runNew(cmd *cobra.Command, args []string) error {
	cfg, err := getConfig(cmd)
	if err != nil {
		return err
	}
	title := args[0]
	id := ulid.MustNew(ulid.Now(), rand.Reader)
	slug := ticket.Slugify(title)
	if slug == "" {
		slug = id.String()
	}
	filename := id.String() + "-" + slug + ".md"
	path := filepath.Join(cfg.ConfigDir, cfg.DefaultPath, filename)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	priority := cfg.Default.Priority
	if priority == "" {
		priority = "p2"
	}
	fm := &ticket.Frontmatter{
		Title:    title,
		Priority: priority,
		Status:   "todo",
		Ready:    cfg.Default.Ready,
		Tags:     cfg.Default.Tags,
	}
	if cfg.Default.Schema != "" {
		fm.Schema = cfg.Default.Schema
	}
	if fm.Tags == nil {
		fm.Tags = []string{}
	}
	body := "---\n"
	raw, err := ticket.MarshalFrontmatter(fm)
	if err != nil {
		return err
	}
	body += string(raw)
	body += "---\n\n"

	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		return err
	}
	fmt.Println(path)
	return nil
}
