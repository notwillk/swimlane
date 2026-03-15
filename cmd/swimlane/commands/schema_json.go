package commands

import (
	"fmt"
	"os"

	"github.com/notwillk/swimlane/internal/schema"
	"github.com/spf13/cobra"
)

func NewSchemaJSON() *cobra.Command {
	return &cobra.Command{
		Use:       "schema-json [config|ticket]",
		Short:     "Print JSON schema for config or ticket",
		Args:      cobra.ExactValidArgs(1),
		ValidArgs: []string{"config", "ticket"},
		RunE:      runSchemaJSON,
	}
}

func runSchemaJSON(cmd *cobra.Command, args []string) error {
	kind := args[0]
	var out string
	switch kind {
	case "config":
		out = schema.Config
	case "ticket":
		out = schema.Ticket
	default:
		return fmt.Errorf("unknown schema type: %s", kind)
	}
	fmt.Fprint(os.Stdout, out)
	return nil
}
