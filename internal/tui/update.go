package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model state
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc", "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.visualToArray)-1 {
			m.cursor++
		}

	case "enter", " ":
		// Toggle completion status
		if len(m.visualToArray) > 0 && m.cursor < len(m.visualToArray) {
			// Map visual position to array index
			arrayIndex := m.visualToArray[m.cursor]
			item := &m.items[arrayIndex]

			if item.IsCompleted() {
				item.Reopen()
			} else {
				item.Complete()
			}

			// Save changes to storage
			if err := m.manager.Save(m.items); err != nil {
				m.err = err
			} else {
				// Rebuild visual mapping after status change
				m.buildVisualMapping()

				// Keep cursor on the same item (it moved to new visual position)
				for visualPos, arrIdx := range m.visualToArray {
					if arrIdx == arrayIndex {
						m.cursor = visualPos
						break
					}
				}
			}
		}
	}

	return m, nil
}
