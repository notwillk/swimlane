package commands

import (
	"fmt"
	"os"

	"github.com/notwillk/swimlane/internal/ticket"
	"github.com/spf13/cobra"
)

// pathByULID returns the ticket path for the given ULID, or "" if not found.
func pathByULID(cmd *cobra.Command, ulid string) (string, error) {
	cfg, err := getConfig(cmd)
	if err != nil {
		return "", err
	}
	tickets, err := ticket.Discover(cfg)
	if err != nil {
		return "", err
	}
	for _, t := range tickets {
		if t.ULID == ulid {
			return t.Path, nil
		}
	}
	return "", nil
}

func NewAssign() *cobra.Command {
	return &cobra.Command{
		Use:   "assign <ulid> <user>",
		Short: "Assign a ticket to a user",
		Args:  cobra.ExactArgs(2),
		RunE:  runAssign,
	}
}

func runAssign(cmd *cobra.Command, args []string) error {
	path, err := pathByULID(cmd, args[0])
	if err != nil || path == "" {
		return fmt.Errorf("ticket not found: %s", args[0])
	}
	fm, body, err := ticket.ReadFrontmatterAndBody(path)
	if err != nil {
		return err
	}
	fm.Assignee = args[1]
	return ticket.WriteFrontmatterAndBody(path, fm, body)
}

func NewClaim() *cobra.Command {
	return &cobra.Command{
		Use:   "claim <ulid>",
		Short: "Assign the ticket to the current user (SWIMLANE_USERNAME)",
		Args:  cobra.ExactArgs(1),
		RunE:  runClaim,
	}
}

func runClaim(cmd *cobra.Command, args []string) error {
	user := os.Getenv("SWIMLANE_USERNAME")
	if user == "" {
		return fmt.Errorf("SWIMLANE_USERNAME is not set")
	}
	path, err := pathByULID(cmd, args[0])
	if err != nil || path == "" {
		return fmt.Errorf("ticket not found: %s", args[0])
	}
	fm, body, err := ticket.ReadFrontmatterAndBody(path)
	if err != nil {
		return err
	}
	fm.Assignee = user
	return ticket.WriteFrontmatterAndBody(path, fm, body)
}

func NewUnclaim() *cobra.Command {
	return &cobra.Command{
		Use:   "unclaim <ulid>",
		Short: "Unassign the ticket",
		Args:  cobra.ExactArgs(1),
		RunE:  runUnclaim,
	}
}

func runUnclaim(cmd *cobra.Command, args []string) error {
	path, err := pathByULID(cmd, args[0])
	if err != nil || path == "" {
		return fmt.Errorf("ticket not found: %s", args[0])
	}
	fm, body, err := ticket.ReadFrontmatterAndBody(path)
	if err != nil {
		return err
	}
	fm.Assignee = ""
	return ticket.WriteFrontmatterAndBody(path, fm, body)
}

func setStatus(cmd *cobra.Command, ulid, status string) error {
	path, err := pathByULID(cmd, ulid)
	if err != nil || path == "" {
		return fmt.Errorf("ticket not found: %s", ulid)
	}
	fm, body, err := ticket.ReadFrontmatterAndBody(path)
	if err != nil {
		return err
	}
	fm.Status = status
	return ticket.WriteFrontmatterAndBody(path, fm, body)
}

func NewStart() *cobra.Command {
	return &cobra.Command{
		Use:   "start <ulid>",
		Short: "Move ticket to in-progress",
		Args:  cobra.ExactArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error { return setStatus(cmd, args[0], "in-progress") },
	}
}

func NewStop() *cobra.Command {
	return &cobra.Command{
		Use:   "stop <ulid>",
		Short: "Move ticket to todo",
		Args:  cobra.ExactArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error { return setStatus(cmd, args[0], "todo") },
	}
}

func NewComplete() *cobra.Command {
	return &cobra.Command{
		Use:   "complete <ulid>",
		Short: "Move ticket to done (update status only; does not delete file)",
		Args:  cobra.ExactArgs(1),
		RunE:  runComplete,
	}
}
func runComplete(cmd *cobra.Command, args []string) error {
	if err := setStatus(cmd, args[0], "done"); err != nil {
		return err
	}
	cfg, err := getConfig(cmd)
	if err != nil {
		return err
	}
	return ticket.CheckParentsWhenSubtaskDone(cfg, args[0])
}

func NewDelete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <ulid>",
		Short: "Delete the ticket file",
		Args:  cobra.ExactArgs(1),
		RunE:  runDelete,
	}
}

func runDelete(cmd *cobra.Command, args []string) error {
	path, err := pathByULID(cmd, args[0])
	if err != nil || path == "" {
		return fmt.Errorf("ticket not found: %s", args[0])
	}
	return os.Remove(path)
}

func setReady(cmd *cobra.Command, ulid string, ready bool) error {
	path, err := pathByULID(cmd, ulid)
	if err != nil || path == "" {
		return fmt.Errorf("ticket not found: %s", ulid)
	}
	fm, body, err := ticket.ReadFrontmatterAndBody(path)
	if err != nil {
		return err
	}
	fm.Ready = ready
	return ticket.WriteFrontmatterAndBody(path, fm, body)
}

func NewActivate() *cobra.Command {
	return &cobra.Command{
		Use:   "activate <ulid>",
		Short: "Set ready to true",
		Args:  cobra.ExactArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error { return setReady(cmd, args[0], true) },
	}
}

func NewDeactivate() *cobra.Command {
	return &cobra.Command{
		Use:   "deactivate <ulid>",
		Short: "Set ready to false",
		Args:  cobra.ExactArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error { return setReady(cmd, args[0], false) },
	}
}
