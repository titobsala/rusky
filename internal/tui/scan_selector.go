package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/tito-sala/rusky/internal/debt"
	"github.com/tito-sala/rusky/internal/scanner"
)

// ScanSelectorModel handles scan result selection
type ScanSelectorModel struct {
	results  []scanner.ScanResult
	selected map[int]bool // Index -> selected
	cursor   int
	manager  *debt.Manager
	width    int
	height   int
	quitting bool
	saved    bool
}

// NewScanSelectorModel creates a new scan selector model
func NewScanSelectorModel(manager *debt.Manager, results []scanner.ScanResult) *ScanSelectorModel {
	// All items selected by default
	selected := make(map[int]bool)
	for i := range results {
		selected[i] = true
	}

	return &ScanSelectorModel{
		results:  results,
		selected: selected,
		cursor:   0,
		manager:  manager,
		quitting: false,
		saved:    false,
	}
}

// Init initializes the model (required by Bubbletea)
func (m *ScanSelectorModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m *ScanSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.results)-1 {
				m.cursor++
			}

		case " ":
			// Toggle selection
			m.selected[m.cursor] = !m.selected[m.cursor]

		case "a":
			// Select all
			for i := range m.results {
				m.selected[i] = true
			}

		case "n":
			// Deselect all
			for i := range m.results {
				m.selected[i] = false
			}

		case "enter":
			// Save selected items
			return m, m.saveSelectedItems()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case savedMsg:
		m.saved = true
		m.quitting = true
		return m, tea.Quit
	}

	return m, nil
}

// View renders the TUI
func (m *ScanSelectorModel) View() string {
	if m.quitting && m.saved {
		selectedCount := 0
		for _, sel := range m.selected {
			if sel {
				selectedCount++
			}
		}
		return fmt.Sprintf("Added %d items to .rusky.json\n", selectedCount)
	}

	if m.quitting {
		return "Cancelled\n"
	}

	var b strings.Builder

	// Header
	b.WriteString(titleStyle.Render("Select Items to Add"))
	b.WriteString("\n\n")

	// Items (with scrolling for many results)
	visibleStart, visibleEnd := m.getVisibleRange()

	for i := visibleStart; i < visibleEnd; i++ {
		result := m.results[i]

		// Checkbox
		checkbox := "[ ]"
		if m.selected[i] {
			checkbox = "[✓]"
		}

		// Cursor
		cursorSymbol := "  "
		if i == m.cursor {
			cursorSymbol = cursor + " "
		}

		// Format line
		line := fmt.Sprintf("%s%s [%s] %s:%d - %s",
			cursorSymbol, checkbox, result.CommentType,
			result.FilePath, result.LineNumber,
			truncate(result.Description, 60))

		// Style
		var style lipgloss.Style
		if i == m.cursor {
			style = selectedItemStyle
		} else {
			style = normalItemStyle
		}

		b.WriteString(style.Render(line))
		b.WriteString("\n")
	}

	// Show scroll indicator if needed
	if visibleEnd < len(m.results) {
		scrollInfo := fmt.Sprintf("... %d more items (scroll down)", len(m.results)-visibleEnd)
		b.WriteString(emptyStateStyle.Render(scrollInfo))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Footer
	selectedCount := 0
	for _, sel := range m.selected {
		if sel {
			selectedCount++
		}
	}

	footer := fmt.Sprintf(
		"↑/↓: Navigate | Space: Toggle | a: Select All | n: Deselect All | Enter: Add (%d selected) | q: Cancel",
		selectedCount)
	b.WriteString(statusBarStyle.Render(footer))

	return b.String()
}

// getVisibleRange returns the range of items to display (pagination)
func (m *ScanSelectorModel) getVisibleRange() (int, int) {
	maxVisible := m.height - 10 // Reserve space for header/footer
	if maxVisible < 5 {
		maxVisible = 5
	}
	if m.height == 0 {
		// Default height if not set
		maxVisible = 20
	}

	if len(m.results) <= maxVisible {
		return 0, len(m.results)
	}

	// Center cursor in view
	start := m.cursor - maxVisible/2
	if start < 0 {
		start = 0
	}

	end := start + maxVisible
	if end > len(m.results) {
		end = len(m.results)
		start = end - maxVisible
		if start < 0 {
			start = 0
		}
	}

	return start, end
}

type savedMsg struct{}

func (m *ScanSelectorModel) saveSelectedItems() tea.Cmd {
	return func() tea.Msg {
		// Load existing items
		items, err := m.manager.List()
		if err != nil {
			// Log error but continue
			fmt.Fprintf(os.Stderr, "Error loading items: %v\n", err)
			return savedMsg{}
		}

		// Add selected scan results
		for i, result := range m.results {
			if !m.selected[i] {
				continue
			}

			item := debt.DebtItem{
				ID:          uuid.New().String(),
				Description: result.Description,
				Status:      debt.StatusOpen,
				CreatedAt:   time.Now(),
				FilePath:    &result.FilePath,
				LineNumber:  &result.LineNumber,
				CommentType: &result.CommentType,
				IsScanned:   true,
			}

			items = append(items, item)
		}

		// Save all items
		if err := m.manager.Save(items); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving items: %v\n", err)
		}

		return savedMsg{}
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// RunScanSelector launches the scan selector TUI
func RunScanSelector(manager *debt.Manager, results []scanner.ScanResult) error {
	model := NewScanSelectorModel(manager, results)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}
