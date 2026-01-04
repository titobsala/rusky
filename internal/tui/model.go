package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tito-sala/rusky/internal/debt"
)

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
}

// buildVisualMapping creates a mapping from visual positions to array indices
// Visual order: open items first, then completed items
func (m *Model) buildVisualMapping() {
	openIndices := make([]int, 0, len(m.items))
	completedIndices := make([]int, 0, len(m.items))

	for i, item := range m.items {
		if item.IsCompleted() {
			completedIndices = append(completedIndices, i)
		} else {
			openIndices = append(openIndices, i)
		}
	}

	m.visualToArray = make([]int, 0, len(m.items))
	m.visualToArray = append(m.visualToArray, openIndices...)
	m.visualToArray = append(m.visualToArray, completedIndices...)
}

// getOpenCount returns the number of open items
func (m *Model) getOpenCount() int {
	count := 0
	for _, item := range m.items {
		if !item.IsCompleted() {
			count++
		}
	}
	return count
}

// NewModel creates a new TUI model
func NewModel(manager *debt.Manager) (*Model, error) {
	// Load items immediately
	items, err := manager.List()
	if err != nil {
		return nil, fmt.Errorf("failed to load items: %w", err)
	}

	m := &Model{
		manager:  manager,
		items:    items,
		cursor:   0,
		quitting: false,
	}

	// Build initial visual mapping
	m.buildVisualMapping()

	return m, nil
}

// Init initializes the model (required by Bubbletea)
func (m *Model) Init() tea.Cmd {
	return nil
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
