package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <description>",
	Short: "Add a new technical debt item",
	Long:  "Add a new technical debt item with the provided description.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Join all arguments to support descriptions with spaces
		description := strings.Join(args, " ")

		item, err := manager.Add(description)
		if err != nil {
			return fmt.Errorf("failed to add item: %w", err)
		}

		fmt.Printf("Added debt item: %s\n", item.ID)
		fmt.Printf("Description: %s\n", item.Description)

		return nil
	},
}
