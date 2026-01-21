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
