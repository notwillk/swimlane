package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewSkill() *cobra.Command {
	return &cobra.Command{
		Use:   "skill",
		Short: "Output a coding agent skill",
		RunE:  runSkill,
	}
}

func runSkill(cmd *cobra.Command, args []string) error {
	fmt.Fprint(os.Stdout, skillMarkdown)
	return nil
}
