package cli

import (
	"github.com/spf13/cobra"
	"github.com/tito-sala/rusky/internal/debt"
	"github.com/tito-sala/rusky/internal/storage"
	"github.com/tito-sala/rusky/internal/tui"
)

const version = "0.1.0"

var (
	manager *debt.Manager
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "rusky",
	Short:   "Rusky - Technical Debt Manager",
	Long:    "A simple, language-agnostic TUI/CLI tool for tracking technical debt in your projects.",
	Version: version,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Launch TUI when run without subcommands
		return tui.Run(manager)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Initialize storage and manager
	store := storage.NewJSONStorage("")
	manager = debt.NewManager(store)

	// Add subcommands
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(completeCmd)
	rootCmd.AddCommand(listCmd)
}
