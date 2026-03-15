package commands

import (
	"os"

	"github.com/spf13/cobra"
)

func NewCompletion() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "completion [bash|zsh|fish]",
		Short:     "Generate shell completions",
		Args:      cobra.ExactValidArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish"},
		RunE:      runCompletion,
	}
	return cmd
}

func runCompletion(cmd *cobra.Command, args []string) error {
	root := cmd.Root()
	shell := args[0]
	switch shell {
	case "bash":
		return root.GenBashCompletion(os.Stdout)
	case "zsh":
		return root.GenZshCompletion(os.Stdout)
	case "fish":
		return root.GenFishCompletion(os.Stdout, true)
	default:
		return nil
	}
}
