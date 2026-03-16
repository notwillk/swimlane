package main

import (
	"os"

	"github.com/notwillk/swimlane/cmd/swimlane/commands"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "swimlane",
		Short: "CLI kanban-style task system built on Markdown files with YAML frontmatter",
	}
	root.PersistentFlags().String("config", "", "path to config file")

	root.AddCommand(commands.NewLS())
	root.AddCommand(commands.NewCreate())
	root.AddCommand(commands.NewAssign())
	root.AddCommand(commands.NewClaim())
	root.AddCommand(commands.NewUnclaim())
	root.AddCommand(commands.NewStart())
	root.AddCommand(commands.NewStop())
	root.AddCommand(commands.NewComplete())
	root.AddCommand(commands.NewDelete())
	root.AddCommand(commands.NewActivate())
	root.AddCommand(commands.NewDeactivate())
	root.AddCommand(commands.NewNext())
	root.AddCommand(commands.NewDone())
	root.AddCommand(commands.NewSchemaJSON())
	root.AddCommand(commands.NewStatic())
	root.AddCommand(commands.NewCompletion())
	root.AddCommand(commands.NewSkill())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
