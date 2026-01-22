package cli

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"github.com/tito-sala/rusky/internal/debt"
	"github.com/tito-sala/rusky/internal/styles"
)

var (
	statusFilter string
	dateFilter   string
	pathFilter   string
	sortField    string
	sortOrder    string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List technical debt items with optional filtering and sorting",
	Long:  "Display technical debt items in a formatted table. Filter by status, date, or path, and sort the results.",
	RunE: func(cmd *cobra.Command, args []string) error {
		items, err := manager.List()
		if err != nil {
			return fmt.Errorf("failed to list items: %w", err)
		}

		// Validate status filter (basic check, though FilterAndSort handles it safely)
		validStatuses := map[string]bool{"all": true, "open": true, "completed": true}
		if !validStatuses[statusFilter] {
			return fmt.Errorf("invalid status '%s': must be 'all', 'open', or 'completed'", statusFilter)
		}

		// Build options
		filterOpts := debt.FilterOptions{
			Status:      statusFilter,
			DateRange:   dateFilter,
			PathPattern: pathFilter,
		}

		sortOpts := debt.SortOptions{
			Field:     sortField,
			Direction: sortOrder,
		}

		// Apply filters and sort
		filteredItems := debt.FilterAndSort(items, filterOpts, sortOpts)

		if len(filteredItems) == 0 {
			emptyMsg := styles.EmptyStateStyle.Render("No technical debt items found matching your criteria.")
			hint := styles.EmptyStateStyle.Render("Try adjusting your filters or use 'rusky add <description>' to add a new item.")
			fmt.Println(emptyMsg)
			fmt.Println(hint)
			return nil
		}

		// Check if any items have file locations
		hasLocations := false
		for _, item := range filteredItems {
			if item.FilePath != nil {
				hasLocations = true
				break
			}
		}

		// Create table headers
		headers := []string{"IDX", "STATUS", "DESCRIPTION"}
		if hasLocations {
			headers = append(headers, "LOCATION")
		}

		// Create and configure the table
		t := table.New().
			Headers(headers...).
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
			row := []string{
				fmt.Sprintf("%d", i+1),
				styles.GetStatusSymbol(item),
				item.Description,
			}

			if hasLocations {
				location := ""
				if item.FilePath != nil {
					location = item.GetLocation()
				}
				row = append(row, location)
			}

			t.Row(row...)
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
	listCmd.Flags().StringVar(&statusFilter, "status", "all", "Filter by status: open, completed, or all")
	listCmd.Flags().StringVar(&dateFilter, "date", "all", "Filter by creation date: today, week, month, or all")
	listCmd.Flags().StringVar(&pathFilter, "path", "", "Filter by file path (substring match or 'scanned')")
	listCmd.Flags().StringVar(&sortField, "sort", "status", "Sort field: status, date, path")
	listCmd.Flags().StringVar(&sortOrder, "order", "asc", "Sort order: asc, desc")
}
