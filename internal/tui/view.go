package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the TUI
func (m Model) View() string {
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

	// Separate items into open and completed
	var openItems, completedItems []int
	for i, item := range m.items {
		if item.IsCompleted() {
			completedItems = append(completedItems, i)
		} else {
			openItems = append(openItems, i)
		}
	}

	// Render open items
	b.WriteString(subtitleStyle.Render(fmt.Sprintf("Open Items (%d)", len(openItems))))
	b.WriteString("\n")

	if len(openItems) == 0 {
		b.WriteString(emptyStateStyle.Render("No open items"))
		b.WriteString("\n")
	} else {
		for idx, i := range openItems {
			b.WriteString(m.renderItem(i, idx+1))
		}
	}

	b.WriteString("\n")

	// Render completed items
	b.WriteString(subtitleStyle.Render(fmt.Sprintf("Completed Items (%d)", len(completedItems))))
	b.WriteString("\n")

	if len(completedItems) == 0 {
		b.WriteString(emptyStateStyle.Render("No completed items"))
		b.WriteString("\n")
	} else {
		for idx, i := range completedItems {
			b.WriteString(m.renderItem(i, len(openItems)+idx+1))
		}
	}

	b.WriteString("\n")

	// Footer with help
	footer := statusBarStyle.Render("↑/↓: Navigate | Enter/Space: Toggle Complete | q/Esc: Quit")
	b.WriteString(footer)

	return b.String()
}

// renderItem renders a single debt item
func (m Model) renderItem(itemIndex, displayIndex int) string {
	item := m.items[itemIndex]
	isCurrent := m.cursor == itemIndex

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
