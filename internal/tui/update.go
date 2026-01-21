package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model state
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		// Handle intro animation tick
		if m.introPhase == 0 {
			m.animationFrame++
			if m.animationFrame >= m.animationMaxFrames {
				m.introPhase = 1 // Transition to main UI
			}
			return m, tick()
		}
		return m, nil

	case tea.KeyMsg:
		// Skip intro animation on any key press
		if m.introPhase == 0 {
			m.introPhase = 1
			return m, nil
		}
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
	// If showing delete confirmation dialog, handle only y/n/esc
	if m.showDeleteConfirm {
		switch msg.String() {
		case "y", "Y":
			return m.confirmDelete()
		case "n", "N", "esc":
			m.resetDeleteConfirmation()
			return m, nil
		}
		// Ignore all other keys when in confirmation mode
		return m, nil
	}

	// Normal key handling
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

	case "d":
		// Initiate delete confirmation
		if len(m.visualToArray) > 0 && m.cursor < len(m.visualToArray) {
			arrayIndex := m.visualToArray[m.cursor]
			m.showDeleteConfirm = true
			m.deleteTargetIndex = arrayIndex
		}

	case "f":
		// Cycle status filter
		switch m.filterOpts.Status {
		case "all":
			m.filterOpts.Status = "open"
		case "open":
			m.filterOpts.Status = "completed"
		case "completed":
			m.filterOpts.Status = "all"
		default:
			m.filterOpts.Status = "all"
		}
		m.buildVisualMapping()
		m.cursor = 0

	case "t":
		// Cycle date filter
		switch m.filterOpts.DateRange {
		case "all":
			m.filterOpts.DateRange = "today"
		case "today":
			m.filterOpts.DateRange = "week"
		case "week":
			m.filterOpts.DateRange = "month"
		case "month":
			m.filterOpts.DateRange = "all"
		default:
			m.filterOpts.DateRange = "all"
		}
		m.buildVisualMapping()
		m.cursor = 0

	case "p":
		// Toggle path filter (scanned only)
		if m.filterOpts.PathPattern == "" {
			m.filterOpts.PathPattern = "scanned"
		} else {
			m.filterOpts.PathPattern = ""
		}
		m.buildVisualMapping()
		m.cursor = 0

	case "s":
		// Cycle sort field
		switch m.sortOpts.Field {
		case "status":
			m.sortOpts.Field = "date"
		case "date":
			m.sortOpts.Field = "path"
		case "path":
			m.sortOpts.Field = "status"
		default:
			m.sortOpts.Field = "status"
		}
		m.buildVisualMapping()
		m.cursor = 0

	case "o":
		// Cycle sort direction
		if m.sortOpts.Direction == "asc" {
			m.sortOpts.Direction = "desc"
		} else {
			m.sortOpts.Direction = "asc"
		}
		m.buildVisualMapping()
		m.cursor = 0

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

// confirmDelete executes the delete operation and updates the model
func (m *Model) confirmDelete() (tea.Model, tea.Cmd) {
	if m.deleteTargetIndex < 0 || m.deleteTargetIndex >= len(m.items) {
		m.resetDeleteConfirmation()
		return m, nil
	}

	// Get the item ID for deletion
	itemID := m.items[m.deleteTargetIndex].ID

	// Delete via manager
	_, err := m.manager.Delete(itemID)
	if err != nil {
		m.err = err
		m.resetDeleteConfirmation()
		return m, nil
	}

	// Reload items from storage
	items, err := m.manager.List()
	if err != nil {
		m.err = err
		m.resetDeleteConfirmation()
		return m, nil
	}

	m.items = items

	// Rebuild visual mapping
	m.buildVisualMapping()

	// Adjust cursor position after deletion
	if len(m.visualToArray) == 0 {
		// No items left
		m.cursor = 0
	} else if m.cursor >= len(m.visualToArray) {
		// Cursor was on last item, move up
		m.cursor = len(m.visualToArray) - 1
	}
	// Otherwise, cursor stays at same position (next item takes its place)

	m.resetDeleteConfirmation()
	return m, nil
}
