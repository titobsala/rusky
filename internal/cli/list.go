package cli

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"github.com/tito-sala/rusky/internal/debt"
	"github.com/tito-sala/rusky/internal/styles"
)

var statusFilter string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List technical debt items with optional filtering",
	Long:  "Display technical debt items in a formatted table. Use --status to filter by item status.",
	RunE: func(cmd *cobra.Command, args []string) error {
		items, err := manager.List()
		if err != nil {
			return fmt.Errorf("failed to list items: %w", err)
		}

		// Validate status filter
		validStatuses := map[string]bool{"all": true, "open": true, "completed": true}
		if !validStatuses[statusFilter] {
			return fmt.Errorf("invalid status '%s': must be 'all', 'open', or 'completed'", statusFilter)
		}

		// Apply status filter
		filteredItems := items
		if statusFilter != "all" {
			filteredItems = []debt.DebtItem{}
			targetStatus := debt.StatusOpen
			if statusFilter == "completed" {
				targetStatus = debt.StatusCompleted
			}
			for _, item := range items {
				if item.Status == targetStatus {
					filteredItems = append(filteredItems, item)
				}
			}
		}

		if len(filteredItems) == 0 {
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
		for i, item := range filteredItems {
			t.Row(
				fmt.Sprintf("%d", i+1),
				styles.GetStatusSymbol(item),
				item.Description,
			)
		}

		// Calculate statistics
		openCount := 0
		completedCount := 0
		for _, item := range filteredItems {
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
				len(filteredItems), openCount, completedCount),
		))

		return nil
	},
}

func init() {
	listCmd.Flags().StringVar(&statusFilter, "status", "all",
		"Filter by status: open, completed, or all")
}
