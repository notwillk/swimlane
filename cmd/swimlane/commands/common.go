package commands

import (
	"github.com/notwillk/swimlane/internal/config"
	"github.com/spf13/cobra"
)

// getConfig loads config using --config from the root command if set.
func getConfig(cmd *cobra.Command) (*config.Config, error) {
	configPath, _ := cmd.Flags().GetString("config")
	// Persistent flag is on root
	root := cmd.Root()
	if root != nil {
		configPath, _ = root.PersistentFlags().GetString("config")
	}
	return config.Load(configPath)
}
