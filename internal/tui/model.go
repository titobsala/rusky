package tui

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tito-sala/rusky/internal/debt"
	"github.com/tito-sala/rusky/internal/intro"
)

// tickMsg is used for animation frame advancement
type tickMsg time.Time

// Model represents the TUI state
type Model struct {
	items         []debt.DebtItem
	err           error
	manager       *debt.Manager
	cursor        int   // Visual position (0-based)
	visualToArray []int // Maps visual index to array index
	width         int
	height        int
	quitting      bool

	// Filter and Sort state
	filterOpts debt.FilterOptions
	sortOpts   debt.SortOptions

	// Delete confirmation state
	showDeleteConfirm bool
	deleteTargetIndex int // Array index of item to delete

	// Intro animation state
	introPhase         int // 0 = show animation, 1 = show main UI
	animationType      intro.AnimationType
	animationFrame     int // current frame number (within current animation)
	animationMaxFrames int // total frames for current animation
}

// buildVisualMapping creates a mapping from visual positions to array indices
// Applies current filters and sort options
func (m *Model) buildVisualMapping() {
	var indices []int

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 1. Filter
	for i, item := range m.items {
		// Status Filter
		if m.filterOpts.Status != "" && m.filterOpts.Status != "all" {
			if string(item.Status) != m.filterOpts.Status {
				continue
			}
		}

		// Date Filter
		if m.filterOpts.DateRange != "" && m.filterOpts.DateRange != "all" {
			itemDate := item.CreatedAt
			if m.filterOpts.DateRange == "today" {
				if itemDate.Before(today) {
					continue
				}
			} else if m.filterOpts.DateRange == "week" {
				weekAgo := now.AddDate(0, 0, -7)
				if itemDate.Before(weekAgo) {
					continue
				}
			} else if m.filterOpts.DateRange == "month" {
				monthAgo := now.AddDate(0, 0, -30)
				if itemDate.Before(monthAgo) {
					continue
				}
			}
		}

		// Path Filter
		if m.filterOpts.PathPattern != "" {
			if m.filterOpts.PathPattern == "scanned" {
				if !item.IsScanned {
					continue
				}
			} else {
				if item.FilePath == nil || !strings.Contains(strings.ToLower(*item.FilePath), strings.ToLower(m.filterOpts.PathPattern)) {
					continue
				}
			}
		}

		// Comment Type Filter
		if len(m.filterOpts.CommentTypes) > 0 {
			if item.CommentType == nil {
				continue
			}
			// Check if item's comment type matches any of the filter types (case-insensitive)
			matched := false
			for _, filterType := range m.filterOpts.CommentTypes {
				if strings.EqualFold(*item.CommentType, filterType) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		// Priority Filter
		if m.filterOpts.Priority != "" && m.filterOpts.Priority != "all" {
			itemPriority := item.GetPriority()
			if m.filterOpts.Priority == "none" {
				if itemPriority != debt.PriorityNone {
					continue
				}
			} else if !strings.EqualFold(string(itemPriority), m.filterOpts.Priority) {
				continue
			}
		}

		indices = append(indices, i)
	}

	// 2. Sort
	sort.SliceStable(indices, func(i, j int) bool {
		idxA := indices[i]
		idxB := indices[j]
		a := m.items[idxA]
		b := m.items[idxB]

		less := false

		switch m.sortOpts.Field {
		case "date":
			if a.CreatedAt.Equal(b.CreatedAt) {
				return a.ID < b.ID
			}
			less = a.CreatedAt.Before(b.CreatedAt)
		case "path":
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
			if a.Status == b.Status {
				// Secondary sort by priority (higher first)
				scoreA := priorityScoreTUI(a.GetPriority())
				scoreB := priorityScoreTUI(b.GetPriority())
				if scoreA == scoreB {
					return a.CreatedAt.Before(b.CreatedAt)
				}
				less = scoreA > scoreB
			}
			// Open (0) < Completed (1)
			scoreA := 0
			if a.Status == debt.StatusCompleted {
				scoreA = 1
			}
			scoreB := 0
			if b.Status == debt.StatusCompleted {
				scoreB = 1
			}
			less = scoreA < scoreB
		case "priority":
			scoreA := priorityScoreTUI(a.GetPriority())
			scoreB := priorityScoreTUI(b.GetPriority())
			if scoreA == scoreB {
				return a.CreatedAt.Before(b.CreatedAt)
			}
			less = scoreA > scoreB
		default:
			// Default to date
			if a.CreatedAt.Equal(b.CreatedAt) {
				return a.ID < b.ID
			}
			less = a.CreatedAt.Before(b.CreatedAt)
		}

		if m.sortOpts.Direction == "desc" {
			return !less
		}
		return less
	})

	m.visualToArray = indices
}

// resetDeleteConfirmation clears the delete confirmation dialog state
func (m *Model) resetDeleteConfirmation() {
	m.showDeleteConfirm = false
	m.deleteTargetIndex = -1
}

// getFilterIndex returns the current position (1-based) and total count for a given filter
// Returns (currentIndex, total) or (0, 0) if filter is not active
func (m *Model) getStatusFilterIndex() (int, int) {
	statuses := []string{"all", "open", "completed"}
	for i, s := range statuses {
		if s == m.filterOpts.Status {
			return i + 1, len(statuses)
		}
	}
	return 1, len(statuses)
}

func (m *Model) getDateFilterIndex() (int, int) {
	dates := []string{"all", "today", "week", "month"}
	for i, d := range dates {
		if d == m.filterOpts.DateRange {
			return i + 1, len(dates)
		}
	}
	return 1, len(dates)
}

func (m *Model) getPathFilterIndex() (int, int) {
	if m.filterOpts.PathPattern == "" {
		return 1, 2
	}
	return 2, 2
}

func typesMatch(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	seen := make(map[string]bool)
	for _, t := range a {
		seen[t] = true
	}
	for _, t := range b {
		if !seen[t] {
			return false
		}
	}
	return true
}

// priorityScoreTUI returns a numeric score for sorting priorities in TUI
func priorityScoreTUI(p debt.Priority) int {
	switch p {
	case debt.PriorityCritical:
		return 4
	case debt.PriorityHigh:
		return 3
	case debt.PriorityMedium:
		return 2
	case debt.PriorityLow:
		return 1
	default:
		return 0
	}
}

func (m *Model) getCommentTypeFilterIndex() (int, int) {
	commentTypes := [][]string{
		{},                // all (no filter)
		{"TODO"},          // TODO only
		{"FIXME"},         // FIXME only
		{"HACK"},          // HACK only
		{"TODO", "FIXME"}, // TODO + FIXME combined
	}

	for i, types := range commentTypes {
		if typesMatch(types, m.filterOpts.CommentTypes) {
			return i + 1, len(commentTypes)
		}
	}
	return 1, len(commentTypes)
}

func (m *Model) getPriorityFilterIndex() (int, int) {
	priorities := []string{"all", "none", "low", "medium", "high", "critical"}
	for i, p := range priorities {
		if p == m.filterOpts.Priority {
			return i + 1, len(priorities)
		}
	}
	return 1, len(priorities)
}

func (m *Model) getSortFilterIndex() (int, int) {
	fields := []string{"status", "date", "path", "priority"}
	for i, f := range fields {
		if f == m.sortOpts.Field {
			return i + 1, len(fields)
		}
	}
	return 1, len(fields)
}

// NewModel creates a new TUI model
func NewModel(manager *debt.Manager) (*Model, error) {
	// Load items immediately
	items, err := manager.List()
	if err != nil {
		return nil, fmt.Errorf("failed to load items: %w", err)
	}

	// Seed random for any future randomized behaviors (kept local and deterministic-safe)
	rand.Seed(time.Now().UnixNano())

	// Use the main GIF-based intro animation
	animType := intro.Main

	m := &Model{
		manager:            manager,
		items:              items,
		cursor:             0,
		quitting:           false,
		introPhase:         0, // Start with intro animation
		animationType:      animType,
		animationFrame:     0,
		animationMaxFrames: intro.GetMaxFrames(animType),
		filterOpts: debt.FilterOptions{
			Status:    "all",
			DateRange: "all",
		},
		sortOpts: debt.SortOptions{
			Field:     "status",
			Direction: "asc",
		},
	}

	// Build initial visual mapping
	m.buildVisualMapping()

	return m, nil
}

// tick returns a command that sends a tick message after 50ms (20 FPS)
func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Init initializes the model (required by Bubbletea)
func (m *Model) Init() tea.Cmd {
	return tick()
}

// Run launches the TUI
func Run(manager *debt.Manager) error {
	model, err := NewModel(manager)
	if err != nil {
		return err
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}
