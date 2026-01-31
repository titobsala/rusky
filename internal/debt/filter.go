package debt

import (
	"sort"
	"strings"
	"time"
)

// FilterOptions defines criteria for filtering debt items
type FilterOptions struct {
	Status       string   // "all", "open", "completed"
	DateRange    string   // "all", "today", "week", "month"
	PathPattern  string   // substring match, or "scanned" for all scanned items
	CommentTypes []string // Filter by comment types (e.g., ["TODO", "FIXME"])
}

// SortOptions defines criteria for sorting debt items
type SortOptions struct {
	Field     string // "status", "date", "path"
	Direction string // "asc", "desc"
}

// FilterAndSort applies filtering and sorting to a slice of debt items
// Returns a new slice, leaving the original unchanged
func FilterAndSort(items []DebtItem, filter FilterOptions, sortOpt SortOptions) []DebtItem {
	var result []DebtItem

	// 1. Filter
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	for _, item := range items {
		// Status Filter
		if filter.Status != "" && filter.Status != "all" {
			if string(item.Status) != filter.Status {
				continue
			}
		}

		// Date Filter
		if filter.DateRange != "" && filter.DateRange != "all" {
			itemDate := item.CreatedAt
			if filter.DateRange == "today" {
				if itemDate.Before(today) {
					continue
				}
			} else if filter.DateRange == "week" {
				// Last 7 days
				weekAgo := now.AddDate(0, 0, -7)
				if itemDate.Before(weekAgo) {
					continue
				}
			} else if filter.DateRange == "month" {
				// Last 30 days
				monthAgo := now.AddDate(0, 0, -30)
				if itemDate.Before(monthAgo) {
					continue
				}
			}
		}

		// Path Filter
		if filter.PathPattern != "" {
			if filter.PathPattern == "scanned" {
				if !item.IsScanned {
					continue
				}
			} else {
				// Substring match on FilePath
				if item.FilePath == nil || !strings.Contains(strings.ToLower(*item.FilePath), strings.ToLower(filter.PathPattern)) {
					continue
				}
			}
		}

		// Comment Type Filter
		if len(filter.CommentTypes) > 0 {
			if item.CommentType == nil {
				continue
			}
			// Check if item's comment type matches any of the filter types (case-insensitive)
			matched := false
			for _, filterType := range filter.CommentTypes {
				if strings.EqualFold(*item.CommentType, filterType) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		result = append(result, item)
	}

	// 2. Sort
	sort.SliceStable(result, func(i, j int) bool {
		a, b := result[i], result[j]

		// Direction logic: default asc means a < b
		// If desc, we invert result (except for specific logic if needed)
		less := false

		switch sortOpt.Field {
		case "date":
			// Date comparison
			if a.CreatedAt.Equal(b.CreatedAt) {
				return a.ID < b.ID // Stable fallback
			}
			less = a.CreatedAt.Before(b.CreatedAt)
		case "path":
			// Path comparison
			pathA := ""
			if a.FilePath != nil {
				pathA = *a.FilePath
			}
			pathB := ""
			if b.FilePath != nil {
				pathB = *b.FilePath
			}
			if pathA == pathB {
				return a.ID < b.ID
			}
			less = pathA < pathB
		case "status":
			// Status comparison: Open < Completed
			if a.Status == b.Status {
				// Secondary sort by date (newest first for better UX?)
				// Let's stick to stable sort order
				return a.CreatedAt.Before(b.CreatedAt)
			}
			// "open" > "completed" string-wise? 'o' > 'c'.
			// So "completed" comes before "open" alphabetically.
			// If we want Open first, we want "open" < "completed" (logic-wise).
			// Let's define: Open (0) < Completed (1)
			scoreA := 0
			if a.Status == StatusCompleted {
				scoreA = 1
			}
			scoreB := 0
			if b.Status == StatusCompleted {
				scoreB = 1
			}
			less = scoreA < scoreB
		default:
			// Default to date
			if a.CreatedAt.Equal(b.CreatedAt) {
				return a.ID < b.ID
			}
			less = a.CreatedAt.Before(b.CreatedAt)
		}

		if sortOpt.Direction == "desc" {
			return !less
		}
		return less
	})

	return result
}
