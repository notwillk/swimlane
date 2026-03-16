package commands

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/notwillk/swimlane/internal/filter"
	"github.com/notwillk/swimlane/internal/ticket"
	"github.com/spf13/cobra"
)

func NewLS() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List tickets",
		RunE:  runLS,
	}
	cmd.Flags().StringSlice("priority", nil, "filter by priority (prefix with ! to negate)")
	cmd.Flags().StringSlice("tag", nil, "filter by tag (prefix with ! to negate)")
	cmd.Flags().StringSlice("status", nil, "filter by status (prefix with ! to negate)")
	cmd.Flags().StringSlice("ready", nil, "filter by ready (true/false, prefix with ! to negate)")
	cmd.Flags().Bool("mine", false, "filter to tickets assigned to current user (SWIMLANE_USERNAME)")
	cmd.Flags().Bool("csv", false, "output as CSV")
	cmd.Flags().Bool("json", false, "output as JSON")
	return cmd
}

func runLS(cmd *cobra.Command, args []string) error {
	cfg, err := getConfig(cmd)
	if err != nil {
		return err
	}
	tickets, err := ticket.Discover(cfg)
	if err != nil {
		return err
	}
	pri, _ := cmd.Flags().GetStringSlice("priority")
	tag, _ := cmd.Flags().GetStringSlice("tag")
	status, _ := cmd.Flags().GetStringSlice("status")
	ready, _ := cmd.Flags().GetStringSlice("ready")
	f := filter.FromFlags(pri, tag, status, ready)
	tickets = filter.Apply(tickets, f)
	if mine, _ := cmd.Flags().GetBool("mine"); mine {
		currentUser := os.Getenv("SWIMLANE_USERNAME")
		var filtered []*ticket.Ticket
		for _, t := range tickets {
			if t.Assignee == currentUser {
				filtered = append(filtered, t)
			}
		}
		tickets = filtered
	}

	csvOut, _ := cmd.Flags().GetBool("csv")
	jsonOut, _ := cmd.Flags().GetBool("json")

	if jsonOut {
		return writeLSJSON(os.Stdout, tickets)
	}
	if csvOut {
		return writeLSCSV(os.Stdout, tickets)
	}
	return writeLSTable(os.Stdout, tickets)
}

func writeLSTable(w *os.File, tickets []*ticket.Ticket) error {
	for _, t := range tickets {
		// P1 todo   01J9T8ZK1BC5A9JH56T9Y9M1DX implement-login-api
		displayTitle := t.Title
		if displayTitle == "" {
			displayTitle = ticketSlugFromPath(t.Path)
		}
		fmt.Fprintf(w, "P%s %s   %s %s\n", strings.TrimPrefix(t.Priority, "p"), t.Status, t.ULID, displayTitle)
	}
	return nil
}

func ticketSlugFromPath(path string) string {
	base := path
	for i := len(base) - 1; i >= 0; i-- {
		if base[i] == '/' || base[i] == '\\' {
			base = base[i+1:]
			break
		}
	}
	if len(base) > 30 && base[26] == '-' && len(base) > 4 && base[len(base)-3:] == ".md" {
		return base[27 : len(base)-3]
	}
	return base
}

func writeLSCSV(w *os.File, tickets []*ticket.Ticket) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"ulid", "title", "priority", "status", "ready", "assignee", "tags", "blocked_by", "subtasks", "path"}); err != nil {
		return err
	}
	for _, t := range tickets {
		title := t.Title
		if title == "" {
			title = ticketSlugFromPath(t.Path)
		}
		readyStr := "false"
		if t.Ready {
			readyStr = "true"
		}
		if err := cw.Write([]string{
			t.ULID,
			title,
			t.Priority,
			t.Status,
			readyStr,
			t.Assignee,
			strings.Join(t.Tags, "|"),
			strings.Join(t.BlockedBy, "|"),
			strings.Join(t.Subtasks, "|"),
			t.Path,
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func writeLSJSON(w *os.File, tickets []*ticket.Ticket) error {
	type row struct {
		ULID      string   `json:"ulid"`
		Title     string   `json:"title"`
		Priority  string   `json:"priority"`
		Status    string   `json:"status"`
		Ready     bool     `json:"ready"`
		Assignee  string   `json:"assignee"`
		Tags      []string `json:"tags"`
		BlockedBy []string `json:"blocked_by"`
		Subtasks  []string `json:"subtasks"`
		Path      string   `json:"path"`
	}
	var rows []row
	for _, t := range tickets {
		title := t.Title
		if title == "" {
			title = ticketSlugFromPath(t.Path)
		}
		rows = append(rows, row{
			ULID: t.ULID, Title: title, Priority: t.Priority, Status: t.Status,
			Ready: t.Ready, Assignee: t.Assignee, Tags: t.Tags, BlockedBy: t.BlockedBy, Subtasks: t.Subtasks, Path: t.Path,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
