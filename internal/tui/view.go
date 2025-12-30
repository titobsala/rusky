package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the TUI
func (m *Model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Header
	header := titleStyle.Render("Rusky - Technical Debt Manager v0.1.0")
	b.WriteString(header)
	b.WriteString("\n\n")

	// Show error if present
	if m.err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		b.WriteString("\n\n")
	}

	// Get count of open items
	openCount := m.getOpenCount()

	// Render open items section
	b.WriteString(subtitleStyle.Render(fmt.Sprintf("Open Items (%d)", openCount)))
	b.WriteString("\n")

	if openCount == 0 {
		b.WriteString(emptyStateStyle.Render("No open items"))
		b.WriteString("\n")
	} else {
		for visualPos := 0; visualPos < openCount; visualPos++ {
			arrayIndex := m.visualToArray[visualPos]
			b.WriteString(m.renderItem(arrayIndex, visualPos, visualPos+1))
		}
	}

	b.WriteString("\n")

	// Render completed items section
	completedCount := len(m.visualToArray) - openCount
	b.WriteString(subtitleStyle.Render(fmt.Sprintf("Completed Items (%d)", completedCount)))
	b.WriteString("\n")

	if completedCount == 0 {
		b.WriteString(emptyStateStyle.Render("No completed items"))
		b.WriteString("\n")
	} else {
		for visualPos := openCount; visualPos < len(m.visualToArray); visualPos++ {
			arrayIndex := m.visualToArray[visualPos]
			displayIndex := visualPos - openCount + 1
			b.WriteString(m.renderItem(arrayIndex, visualPos, displayIndex))
		}
	}

	b.WriteString("\n")

	// Footer with help
	footer := statusBarStyle.Render("↑/↓: Navigate | Enter/Space: Toggle Complete | q/Esc: Quit")
	b.WriteString(footer)

	return b.String()
}

// renderItem renders a single debt item
// arrayIndex: index in m.items array
// visualPos: position in visual display (0-based)
// displayIndex: number shown to user (1-based, resets for each section)
func (m *Model) renderItem(arrayIndex, visualPos, displayIndex int) string {
	item := m.items[arrayIndex]
	isCurrent := m.cursor == visualPos

	var prefix, text string

	if isCurrent {
		prefix = cursor + " "
	} else {
		prefix = "  "
	}

	// Format the item text
	text = fmt.Sprintf("%d. %s", displayIndex, item.Description)

	// Apply styling based on state
	var style lipgloss.Style
	if item.IsCompleted() {
		if isCurrent {
			style = selectedCompletedItemStyle
			text = prefix + text + " ✓"
		} else {
			style = completedItemStyle
			text = "  " + text + " ✓"
		}
	} else {
		if isCurrent {
			style = selectedItemStyle
			text = prefix + text
		} else {
			style = normalItemStyle
			text = "  " + text
		}
	}

	return style.Render(text) + "\n"
}
