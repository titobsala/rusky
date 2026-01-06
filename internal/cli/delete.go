package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tito-sala/rusky/internal/debt"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <id|index>",
	Short: "Permanently delete a debt item",
	Long:  "Permanently delete a technical debt item by providing its UUID or 1-based index. This action cannot be undone.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]

		items, err := manager.List()
		if err != nil {
			return fmt.Errorf("failed to load items: %w", err)
		}

		var itemToDelete *debt.DebtItem
		for i := range items {
			if items[i].ID == identifier {
				itemToDelete = &items[i]
				break
			}
		}

		if itemToDelete == nil {
			if idx, err := strconv.Atoi(identifier); err == nil {
				if idx > 0 && idx <= len(items) {
					itemToDelete = &items[idx-1]
				}
			}
		}

		if itemToDelete == nil {
			return fmt.Errorf("item not found: %s", identifier)
		}

		fmt.Printf("\u26a0 WARNING: You are about to permanently delete:\n\n")
		fmt.Printf("  %s", itemToDelete.Description)
		if itemToDelete.IsCodeReference() {
			fmt.Printf(" [%s]", itemToDelete.GetLocation())
		}
		fmt.Printf("\n\n")
		fmt.Printf("This action cannot be undone.\n")
		fmt.Printf("Are you sure? (y/N): ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Delete cancelled.")
			return nil
		}

		item, err := manager.Delete(identifier)
		if err != nil {
			return fmt.Errorf("failed to delete item: %w", err)
		}

		fmt.Printf("\n\u2713 Deleted debt item: %s\n", item.ID)
		fmt.Printf("Description: %s\n", item.Description)

		return nil
	},
}
