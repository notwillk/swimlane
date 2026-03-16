package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/notwillk/swimlane/internal/ticket"
	"github.com/oklog/ulid/v2"
	"github.com/spf13/cobra"
)

// runCustomAction runs a substituted command string via sh -c.
func runCustomAction(cmd *cobra.Command, command string) error {
	c := exec.Command("sh", "-c", command)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func NewCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [title]",
		Short: "Create a new ticket (optionally with description from stdin)",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runCreate,
	}
	cmd.Flags().Bool("no-description", false, "do not read description from stdin")
	cmd.Flags().Bool("claim", false, "also assign to current user (SWIMLANE_USERNAME)")
	cmd.Flags().String("assign", "", "also assign to this user")
	cmd.Flags().Bool("start", false, "also set status to in-progress")
	cmd.Flags().Bool("activate", false, "also set ready to true")
	cmd.Flags().Bool("deactivate", false, "also set ready to false")
	return cmd
}

func runCreate(cmd *cobra.Command, args []string) error {
	cfg, err := getConfig(cmd)
	if err != nil {
		return err
	}
	// Optional custom action from config
	if cfg.Actions != nil {
		if a, ok := cfg.Actions["create"]; ok && a.Command != "" {
			title := ""
			if len(args) > 0 {
				title = args[0]
			}
			substituted := strings.ReplaceAll(a.Command, "{title}", title)
			return runCustomAction(cmd, substituted)
		}
	}

	title := ""
	if len(args) > 0 {
		title = args[0]
	}
	noDesc, _ := cmd.Flags().GetBool("no-description")

	var body []byte
	if !noDesc {
		fmt.Fprint(os.Stderr, "Enter description (Ctrl+D to end):\n")
		scanner := bufio.NewScanner(os.Stdin)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		body = []byte(strings.Join(lines, "\n"))
		if len(body) > 0 && body[len(body)-1] != '\n' {
			body = append(body, '\n')
		}
	}

	id := ulid.MustNew(ulid.Now(), nil)
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

	// Dogpiling: apply extra actions (conflict check)
	claim, _ := cmd.Flags().GetBool("claim")
	assign, _ := cmd.Flags().GetString("assign")
	start, _ := cmd.Flags().GetBool("start")
	activate, _ := cmd.Flags().GetBool("activate")
	deactivate, _ := cmd.Flags().GetBool("deactivate")
	if claim && assign != "" {
		return fmt.Errorf("cannot use both --claim and --assign")
	}
	if activate && deactivate {
		return fmt.Errorf("cannot use both --activate and --deactivate")
	}
	if claim {
		fm.Assignee = os.Getenv("SWIMLANE_USERNAME")
	} else if assign != "" {
		fm.Assignee = assign
	}
	if start {
		fm.Status = "in-progress"
	}
	if activate {
		fm.Ready = true
	} else if deactivate {
		fm.Ready = false
	}

	if err := ticket.WriteFrontmatterAndBody(path, fm, body); err != nil {
		return err
	}
	fmt.Println(path)
	return nil
}
