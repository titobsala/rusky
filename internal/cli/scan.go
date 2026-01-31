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
	dryRun    bool
	addAll    bool
	forceAdd  bool
	updateAll bool
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

	if len(results) == 0 {
		fmt.Println("\nNo technical debt markers found.")
		return nil
	}

	// Detect duplicates (unless --force-add is set)
	var resultsWithStatus []scanner.ScanResultWithStatus
	newCount := 0
	duplicateCount := 0
	changedCount := 0

	if forceAdd {
		// Convert all to new results
		for _, r := range results {
			resultsWithStatus = append(resultsWithStatus, scanner.ScanResultWithStatus{
				ScanResult: r,
				Status:     scanner.StatusNew,
			})
		}
		newCount = len(results)
	} else {
		// Check for duplicates
		for _, r := range results {
			existing, err := manager.FindByLocation(r.FilePath, r.LineNumber)
			if err != nil {
				return fmt.Errorf("failed to check for duplicates: %w", err)
			}

			if existing == nil {
				// New item
				resultsWithStatus = append(resultsWithStatus, scanner.ScanResultWithStatus{
					ScanResult: r,
					Status:     scanner.StatusNew,
				})
				newCount++
			} else if existing.Description != r.Description {
				// Changed description
				resultsWithStatus = append(resultsWithStatus, scanner.ScanResultWithStatus{
					ScanResult:   r,
					Status:       scanner.StatusChanged,
					ExistingID:   existing.ID,
					ExistingDesc: existing.Description,
				})
				changedCount++
			} else {
				// Duplicate (no change)
				resultsWithStatus = append(resultsWithStatus, scanner.ScanResultWithStatus{
					ScanResult: r,
					Status:     scanner.StatusDuplicate,
					ExistingID: existing.ID,
				})
				duplicateCount++
			}
		}
	}

	// Print summary
	fmt.Printf("\nFound %d technical debt markers across %d files\n", len(results), countUniqueFiles(results))
	if !forceAdd {
		if duplicateCount > 0 {
			fmt.Printf("%d duplicates skipped (no changes)\n", duplicateCount)
		}
		if changedCount > 0 {
			fmt.Printf("%d items with changed descriptions\n", changedCount)
		}
		if newCount > 0 {
			fmt.Printf("%d new items found\n", newCount)
		}
	}
	fmt.Println()

	summary := groupByType(results)
	for typ, count := range summary {
		fmt.Printf("  %s: %d\n", typ, count)
	}
	fmt.Println()

	if dryRun {
		fmt.Println("Dry-run mode: items not added to .rusky.json")
		printPreview(results)
		return nil
	}

	if addAll || forceAdd {
		return addAllResultsWithStatus(resultsWithStatus)
	}

	return tui.RunScanSelector(manager, results)
}

func countUniqueFiles(results []scanner.ScanResult) int {
	files := make(map[string]bool)
	for _, r := range results {
		files[r.FilePath] = true
	}
	return len(files)
}

func groupByType(results []scanner.ScanResult) map[string]int {
	summary := make(map[string]int)
	for _, result := range results {
		summary[result.CommentType]++
	}
	return summary
}

func printPreview(results []scanner.ScanResult) {
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

func addAllResultsWithStatus(results []scanner.ScanResultWithStatus) error {
	added := 0
	updated := 0
	skipped := 0

	for _, result := range results {
		switch result.Status {
		case scanner.StatusNew:
			if err := addScanResult(result.ScanResult); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to add %s:%d: %v\n",
					result.FilePath, result.LineNumber, err)
				continue
			}
			added++

		case scanner.StatusChanged:
			if updateAll {
				// Auto-update
				if err := manager.UpdateDescription(result.ExistingID, result.Description); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to update %s:%d: %v\n",
						result.FilePath, result.LineNumber, err)
					continue
				}
				updated++
			} else {
				// Prompt for confirmation
				fmt.Printf("\n[CHANGED] %s:%d\n", result.FilePath, result.LineNumber)
				fmt.Printf("  Old: %s\n", result.ExistingDesc)
				fmt.Printf("  New: %s\n", result.Description)
				fmt.Print("Update? (y/n): ")

				var response string
				fmt.Scanln(&response)
				if response == "y" || response == "Y" {
					if err := manager.UpdateDescription(result.ExistingID, result.Description); err != nil {
						fmt.Fprintf(os.Stderr, "Warning: failed to update: %v\n", err)
						continue
					}
					updated++
				} else {
					skipped++
				}
			}

		case scanner.StatusDuplicate:
			// Skip silently
			skipped++
		}
	}

	fmt.Printf("\nSummary:\n")
	if added > 0 {
		fmt.Printf("  Added: %d new items\n", added)
	}
	if updated > 0 {
		fmt.Printf("  Updated: %d changed descriptions\n", updated)
	}
	if skipped > 0 {
		fmt.Printf("  Skipped: %d duplicates\n", skipped)
	}

	return nil
}

func addScanResult(result scanner.ScanResult) error {
	items, err := manager.List()
	if err != nil {
		return err
	}

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
	scanCmd.Flags().BoolVar(&forceAdd, "force-add", false,
		"Force add all items, ignoring duplicates")
	scanCmd.Flags().BoolVar(&updateAll, "update-all", false,
		"Automatically update all changed descriptions without prompting")
}
