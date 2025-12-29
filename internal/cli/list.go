package cli

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"github.com/tito-sala/rusky/internal/styles"
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
			emptyMsg := styles.EmptyStateStyle.Render("No technical debt items found.")
			hint := styles.EmptyStateStyle.Render("Use 'rusky add <description>' to add your first item.")
			fmt.Println(emptyMsg)
			fmt.Println(hint)
			return nil
		}

		// Create and configure the table
		t := table.New().
			Headers("IDX", "STATUS", "DESCRIPTION").
			Border(lipgloss.RoundedBorder()).
			BorderStyle(styles.TableBorderStyle).
			StyleFunc(func(row, col int) lipgloss.Style {
				// Header row styling
				if row == 0 {
					return styles.HeaderStyle
				}
				return lipgloss.NewStyle()
			})

		// Add rows to the table
		for i, item := range items {
			t.Row(
				fmt.Sprintf("%d", i+1),
				styles.GetStatusSymbol(item),
				item.Description,
			)
		}

		// Calculate statistics
		openCount := 0
		completedCount := 0
		for _, item := range items {
			if item.IsCompleted() {
				completedCount++
			} else {
				openCount++
			}
		}

		// Print styled output
		fmt.Println(styles.TitleStyle.Render("Technical Debt Items"))
		fmt.Println(t.String())
		fmt.Println(styles.FooterStyle.Render(
			fmt.Sprintf("Total: %d items • %d open • %d completed",
				len(items), openCount, completedCount),
		))

		return nil
	},
}
