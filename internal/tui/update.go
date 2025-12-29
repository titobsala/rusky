package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc", "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.items)-1 {
			m.cursor++
		}

	case "enter", " ":
		// Toggle completion status
		if len(m.items) > 0 && m.cursor < len(m.items) {
			item := &m.items[m.cursor]
			if item.IsCompleted() {
				item.Reopen()
			} else {
				item.Complete()
			}

			// Save changes to storage
			if err := m.manager.Save(m.items); err != nil {
				m.err = err
			}
		}
	}

	return m, nil
}
