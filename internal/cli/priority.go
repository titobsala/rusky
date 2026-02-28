package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tito-sala/rusky/internal/debt"
	"github.com/tito-sala/rusky/internal/styles"
)

var priorityCmd = &cobra.Command{
	Use:   "priority <id> <level>",
	Short: "Set the priority of a debt item",
	Long:  "Set the priority level (none, low, medium, high, critical) for a debt item.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		priorityStr := args[1]

		validPriorities := map[string]debt.Priority{
			"none":     debt.PriorityNone,
			"low":      debt.PriorityLow,
			"medium":   debt.PriorityMedium,
			"high":     debt.PriorityHigh,
			"critical": debt.PriorityCritical,
		}

		priority, ok := validPriorities[priorityStr]
		if !ok {
			return fmt.Errorf("invalid priority '%s': must be none, low, medium, high, or critical", priorityStr)
		}

		item, err := manager.SetPriority(id, priority)
		if err != nil {
			return err
		}

		fmt.Println(styles.SuccessStyle.Render("Priority updated:"))
		fmt.Printf("  %s\n", item.Description)
		fmt.Printf("  Priority: %s\n", styles.GetPriorityLabel(item.GetPriority()))

		return nil
	},
}
