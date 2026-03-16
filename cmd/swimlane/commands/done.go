package commands

import (
	"fmt"
	"os"

	"github.com/notwillk/swimlane/internal/ticket"
	"github.com/spf13/cobra"
)

func NewDone() *cobra.Command {
	return &cobra.Command{
		Use:   "done [ulid]",
		Short: "Mark a ticket as complete by deleting its file",
		Args:  cobra.ExactArgs(1),
		RunE:  runDone,
	}
}

func runDone(cmd *cobra.Command, args []string) error {
	cfg, err := getConfig(cmd)
	if err != nil {
		return err
	}
	tickets, err := ticket.Discover(cfg)
	if err != nil {
		return err
	}
	ulidArg := args[0]
	var path string
	for _, t := range tickets {
		if t.ULID == ulidArg {
			path = t.Path
			break
		}
	}
	if path == "" {
		return fmt.Errorf("ticket not found: %s", ulidArg)
	}
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("delete ticket: %w", err)
	}
	return ticket.CheckParentsWhenSubtaskDone(cfg, ulidArg)
}
