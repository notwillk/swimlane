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
	root.AddCommand(commands.NewNew())
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
