package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/notwillk/swimlane/internal/config"
	"github.com/notwillk/swimlane/internal/ticket"
	"github.com/spf13/cobra"
)

func NewStatic() *cobra.Command {
	return &cobra.Command{
		Use:   "static",
		Short: "Validate config and ticket YAML/frontmatter (run all checks, report all failures)",
		RunE:  runStatic,
	}
}

func runStatic(cmd *cobra.Command, args []string) error {
	configPath, _ := cmd.Root().PersistentFlags().GetString("config")
	var cfg *config.Config
	var errs []string

	cfg, err := config.Load(configPath)
	if err != nil {
		errs = append(errs, fmt.Sprintf("config: %s", err.Error()))
		wd, _ := os.Getwd()
		cfg = &config.Config{
			Tickets:     config.DefaultTicketsGlob,
			DefaultPath: config.DefaultPath,
			ConfigDir:   wd,
			Default: config.Defaults{
				Priority: "p2",
				Ready:    true,
				Tags:     nil,
			},
		}
	}

	paths, err := ticket.GlobPaths(cfg)
	if err != nil {
		errs = append(errs, fmt.Sprintf("tickets glob: %s", err.Error()))
	} else {
		for _, path := range paths {
			_, err := ticket.ParseFile(path)
			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: %s", path, err.Error()))
			}
		}
	}

	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintln(os.Stderr, e)
		}
		fmt.Fprintf(os.Stderr, "\n%d failure(s)\n", len(errs))
		return errors.New("static analysis failed")
	}
	return nil
}
