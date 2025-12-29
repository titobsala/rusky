package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tito-sala/rusky/internal/debt"
)

// Model represents the TUI state
type Model struct {
	manager  *debt.Manager
	items    []debt.DebtItem
	cursor   int
	width    int
	height   int
	err      error
	quitting bool
}

// NewModel creates a new TUI model
func NewModel(manager *debt.Manager) (*Model, error) {
	// Load items immediately
	items, err := manager.List()
	if err != nil {
		return nil, fmt.Errorf("failed to load items: %w", err)
	}

	return &Model{
		manager:  manager,
		items:    items,
		cursor:   0,
		quitting: false,
	}, nil
}

// Init initializes the model (required by Bubbletea)
func (m Model) Init() tea.Cmd {
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
