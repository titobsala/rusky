package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete <id|index>",
	Short: "Mark a debt item as completed",
	Long:  "Mark a technical debt item as completed by providing its UUID or 1-based index.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]

		item, err := manager.Complete(identifier)
		if err != nil {
			return fmt.Errorf("failed to complete item: %w", err)
		}

		fmt.Printf("Completed debt item: %s\n", item.ID)
		fmt.Printf("Description: %s\n", item.Description)

		return nil
	},
}
