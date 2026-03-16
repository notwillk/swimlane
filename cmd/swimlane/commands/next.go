package commands

import (
	"fmt"
	"os"

	"github.com/notwillk/swimlane/internal/filter"
	"github.com/notwillk/swimlane/internal/graph"
	"github.com/notwillk/swimlane/internal/ticket"
	"github.com/spf13/cobra"
)

func NewNext() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "next",
		Short: "Return the next ticket to implement",
		RunE:  runNext,
	}
	cmd.Flags().StringSlice("priority", nil, "filter by priority (prefix with ! to negate)")
	cmd.Flags().StringSlice("tag", nil, "filter by tag (prefix with ! to negate)")
	cmd.Flags().StringSlice("status", nil, "filter by status (prefix with ! to negate)")
	cmd.Flags().StringSlice("ready", nil, "filter by ready (prefix with ! to negate)")
	return cmd
}

func runNext(cmd *cobra.Command, args []string) error {
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
	g := graph.Build(tickets)
	path := g.Next(f)
	if path == "" {
		return fmt.Errorf("no next ticket")
	}
	fmt.Fprintln(os.Stdout, path)
	return nil
}
