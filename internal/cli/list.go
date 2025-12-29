package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all technical debt items",
	Long:  "Display all technical debt items in a simple text format.",
	RunE: func(cmd *cobra.Command, args []string) error {
		items, err := manager.List()
		if err != nil {
			return fmt.Errorf("failed to list items: %w", err)
		}

		if len(items) == 0 {
			fmt.Println("No technical debt items found.")
			fmt.Println("Use 'rusky add <description>' to add your first item.")
			return nil
		}

		// Print header
		fmt.Printf("%-4s %-12s %-20s %s\n", "IDX", "STATUS", "CREATED", "DESCRIPTION")
		fmt.Println(strings.Repeat("-", 80))

		// Print items
		for i, item := range items {
			status := string(item.Status)
			if item.IsCompleted() {
				status = "✓ " + status
			}

			created := item.CreatedAt.Format("2006-01-02 15:04")

			fmt.Printf("%-4d %-12s %-20s %s\n", i+1, status, created, item.Description)
		}

		fmt.Printf("\nTotal: %d items\n", len(items))

		return nil
	},
}
