package debt

import (
	"testing"
	"time"
)

func TestFilterAndSort(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	lastWeek := now.AddDate(0, 0, -8)

	items := []DebtItem{
		{ID: "1", Description: "Task 1", Status: StatusOpen, CreatedAt: now, FilePath: strPtr("src/main.go"), IsScanned: true},
		{ID: "2", Description: "Task 2", Status: StatusCompleted, CreatedAt: yesterday, FilePath: nil, IsScanned: false},
		{ID: "3", Description: "Task 3", Status: StatusOpen, CreatedAt: lastWeek, FilePath: strPtr("readme.md"), IsScanned: true},
	}

	tests := []struct {
		name      string
		filter    FilterOptions
		sort      SortOptions
		wantCount int
		check     func([]DebtItem) bool
	}{
		{
			name:      "Filter All Open",
			filter:    FilterOptions{Status: "open"},
			sort:      SortOptions{},
			wantCount: 2,
			check: func(res []DebtItem) bool {
				return res[0].ID == "3" && res[1].ID == "1"
			},
		},
		{
			name:      "Filter Completed",
			filter:    FilterOptions{Status: "completed"},
			sort:      SortOptions{},
			wantCount: 1,
			check: func(res []DebtItem) bool {
				return res[0].ID == "2"
			},
		},
		{
			name:      "Filter Date Today",
			filter:    FilterOptions{DateRange: "today"},
			sort:      SortOptions{},
			wantCount: 1,
			check: func(res []DebtItem) bool {
				return res[0].ID == "1"
			},
		},
		{
			name:      "Filter Path Scanned",
			filter:    FilterOptions{PathPattern: "scanned"},
			sort:      SortOptions{},
			wantCount: 2,
			check: func(res []DebtItem) bool {
				for _, item := range res {
					if !item.IsScanned {
						return false
					}
				}
				return true
			},
		},
		{
			name:      "Sort Date Desc",
			filter:    FilterOptions{},
			sort:      SortOptions{Field: "date", Direction: "desc"},
			wantCount: 3,
			check: func(res []DebtItem) bool {
				return res[0].ID == "1" && res[1].ID == "2" && res[2].ID == "3"
			},
		},
		{
			name:      "Sort Date Asc",
			filter:    FilterOptions{},
			sort:      SortOptions{Field: "date", Direction: "asc"},
			wantCount: 3,
			check: func(res []DebtItem) bool {
				return res[0].ID == "3" && res[1].ID == "2" && res[2].ID == "1"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterAndSort(items, tt.filter, tt.sort)
			if len(got) != tt.wantCount {
				t.Errorf("FilterAndSort() count = %v, want %v", len(got), tt.wantCount)
			}
			if tt.check != nil && !tt.check(got) {
				t.Errorf("FilterAndSort() check failed")
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}

func TestCommentTypeFilter(t *testing.T) {
	now := time.Now()

	items := []DebtItem{
		{ID: "1", Description: "Fix bug", Status: StatusOpen, CreatedAt: now, CommentType: strPtr("TODO"), FilePath: strPtr("src/main.go"), IsScanned: true},
		{ID: "2", Description: "Fix urgent", Status: StatusOpen, CreatedAt: now, CommentType: strPtr("FIXME"), FilePath: strPtr("src/app.go"), IsScanned: true},
		{ID: "3", Description: "Quick fix", Status: StatusOpen, CreatedAt: now, CommentType: strPtr("HACK"), FilePath: strPtr("src/util.go"), IsScanned: true},
		{ID: "4", Description: "Investigate", Status: StatusOpen, CreatedAt: now, CommentType: strPtr("XXX"), FilePath: strPtr("src/test.go"), IsScanned: true},
		{ID: "5", Description: "Critical bug", Status: StatusOpen, CreatedAt: now, CommentType: strPtr("BUG"), FilePath: strPtr("src/core.go"), IsScanned: true},
		{ID: "6", Description: "Remember this", Status: StatusOpen, CreatedAt: now, CommentType: strPtr("NOTE"), FilePath: strPtr("src/api.go"), IsScanned: true},
		{ID: "7", Description: "Manual entry", Status: StatusOpen, CreatedAt: now, CommentType: nil, FilePath: nil, IsScanned: false},
	}

	tests := []struct {
		name      string
		filter    FilterOptions
		wantCount int
		wantIDs   []string
	}{
		{
			name:      "Filter Single Type - TODO",
			filter:    FilterOptions{CommentTypes: []string{"TODO"}},
			wantCount: 1,
			wantIDs:   []string{"1"},
		},
		{
			name:      "Filter Single Type - FIXME",
			filter:    FilterOptions{CommentTypes: []string{"FIXME"}},
			wantCount: 1,
			wantIDs:   []string{"2"},
		},
		{
			name:      "Filter Multiple Types - TODO and FIXME",
			filter:    FilterOptions{CommentTypes: []string{"TODO", "FIXME"}},
			wantCount: 2,
			wantIDs:   []string{"1", "2"},
		},
		{
			name:      "Filter Multiple Types - HACK, XXX, BUG",
			filter:    FilterOptions{CommentTypes: []string{"HACK", "XXX", "BUG"}},
			wantCount: 3,
			wantIDs:   []string{"3", "4", "5"},
		},
		{
			name:      "Filter Case Insensitive - todo (lowercase)",
			filter:    FilterOptions{CommentTypes: []string{"todo"}},
			wantCount: 1,
			wantIDs:   []string{"1"},
		},
		{
			name:      "Filter Case Insensitive - FiXmE (mixed case)",
			filter:    FilterOptions{CommentTypes: []string{"FiXmE"}},
			wantCount: 1,
			wantIDs:   []string{"2"},
		},
		{
			name:      "Filter Excludes Manual Entries (nil CommentType)",
			filter:    FilterOptions{CommentTypes: []string{"TODO", "FIXME", "HACK"}},
			wantCount: 3,
			wantIDs:   []string{"1", "2", "3"},
		},
		{
			name:      "No Filter - All Items",
			filter:    FilterOptions{CommentTypes: []string{}},
			wantCount: 7,
			wantIDs:   []string{"1", "2", "3", "4", "5", "6", "7"},
		},
		{
			name:      "Combine with Status Filter - Open TODO",
			filter:    FilterOptions{Status: "open", CommentTypes: []string{"TODO"}},
			wantCount: 1,
			wantIDs:   []string{"1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterAndSort(items, tt.filter, SortOptions{Field: "date", Direction: "asc"})

			if len(got) != tt.wantCount {
				t.Errorf("FilterAndSort() count = %v, want %v", len(got), tt.wantCount)
			}

			// Verify correct items were returned
			gotIDs := make([]string, len(got))
			for i, item := range got {
				gotIDs[i] = item.ID
			}

			if len(gotIDs) != len(tt.wantIDs) {
				t.Errorf("FilterAndSort() returned IDs %v, want %v", gotIDs, tt.wantIDs)
				return
			}

			for i := range gotIDs {
				found := false
				for _, wantID := range tt.wantIDs {
					if gotIDs[i] == wantID {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("FilterAndSort() returned unexpected ID %v", gotIDs[i])
				}
			}
		})
	}
}
