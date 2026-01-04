package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/tito-sala/rusky/internal/debt"
	"github.com/tito-sala/rusky/internal/scanner"
	"github.com/tito-sala/rusky/internal/tui"
)

var (
	dryRun bool
	addAll bool
)

var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Scan codebase for technical debt markers",
	Long: `Scan your codebase for TODO, FIXME, HACK, XXX, BUG, and NOTE comments.

Examples:
  rusky scan              # Scan current directory
  rusky scan ./src        # Scan specific directory
  rusky scan --dry-run    # Preview without adding items
  rusky scan --add-all    # Add all items without confirmation`,
	Args: cobra.MaximumNArgs(1),
	RunE: runScan,
}

func runScan(cmd *cobra.Command, args []string) error {
	// Determine scan path
	scanPath := "."
	if len(args) > 0 {
		scanPath = args[0]
	}

	// Create scanner and scan
	s := scanner.NewScanner()
	results, err := s.Scan(scanPath)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	// Show summary
	fmt.Printf("\nScan complete: found %d items\n\n", len(results))

	if len(results) == 0 {
		fmt.Println("No technical debt markers found.")
		return nil
	}

	// Group by type for summary
	summary := groupByType(results)
	for typ, count := range summary {
		fmt.Printf("  %s: %d\n", typ, count)
	}
	fmt.Println()

	// Handle dry-run mode
	if dryRun {
		fmt.Println("Dry-run mode: items not added to .rusky.json")
		printPreview(results)
		return nil
	}

	// Handle add-all mode
	if addAll {
		return addAllResults(results)
	}

	// Interactive selection via TUI
	return tui.RunScanSelector(manager, results)
}

func groupByType(results []scanner.ScanResult) map[string]int {
	summary := make(map[string]int)
	for _, result := range results {
		summary[result.CommentType]++
	}
	return summary
}

func printPreview(results []scanner.ScanResult) {
	// Limit to first 10 for preview
	limit := 10
	if len(results) < limit {
		limit = len(results)
	}

	for i := 0; i < limit; i++ {
		r := results[i]
		fmt.Printf("  [%s] %s:%d - %s\n",
			r.CommentType, r.FilePath, r.LineNumber, r.Description)
	}

	if len(results) > limit {
		fmt.Printf("\n  ... and %d more items\n", len(results)-limit)
	}
}

func addAllResults(results []scanner.ScanResult) error {
	added := 0
	for _, result := range results {
		if err := addScanResult(result); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to add %s:%d: %v\n",
				result.FilePath, result.LineNumber, err)
			continue
		}
		added++
	}

	fmt.Printf("Added %d items to .rusky.json\n", added)
	return nil
}

func addScanResult(result scanner.ScanResult) error {
	items, err := manager.List()
	if err != nil {
		return err
	}

	// Create new debt item
	item := debt.DebtItem{
		ID:          uuid.New().String(),
		Description: result.Description,
		Status:      debt.StatusOpen,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		FilePath:    &result.FilePath,
		LineNumber:  &result.LineNumber,
		CommentType: &result.CommentType,
		IsScanned:   true,
	}

	items = append(items, item)
	return manager.Save(items)
}

func init() {
	scanCmd.Flags().BoolVar(&dryRun, "dry-run", false,
		"Preview scan results without adding items")
	scanCmd.Flags().BoolVar(&addAll, "add-all", false,
		"Add all found items without confirmation")
}
