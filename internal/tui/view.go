package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tito-sala/rusky/internal/intro"
)

// View renders the TUI
func (m *Model) View() string {
	// Show intro animation first
	if m.introPhase == 0 {
		return intro.RenderAnimation(
			m.animationType,
			m.animationFrame,
		)
	}

	if m.quitting {
		return ""
	}

	var b strings.Builder

	header := titleStyle.Render("Rusky - Technical Debt Manager v0.2.2")
	b.WriteString(header)
	b.WriteString("\n\n")

	if m.err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		b.WriteString("\n\n")
	}

	if m.showDeleteConfirm {
		return m.renderDeleteConfirmation()
	}

	// Render Filter/Sort Status
	var filterStatus []string
	if m.filterOpts.Status != "all" {
		filterStatus = append(filterStatus, fmt.Sprintf("Status:%s", m.filterOpts.Status))
	}
	if m.filterOpts.DateRange != "all" {
		filterStatus = append(filterStatus, fmt.Sprintf("Date:%s", m.filterOpts.DateRange))
	}
	if m.filterOpts.PathPattern != "" {
		val := m.filterOpts.PathPattern
		if val == "scanned" {
			filterStatus = append(filterStatus, "Path:Scanned")
		} else {
			filterStatus = append(filterStatus, fmt.Sprintf("Path:%s", val))
		}
	}

	// Always show sort
	sortIcon := "↓" // asc
	if m.sortOpts.Direction == "desc" {
		sortIcon = "↑"
	}
	filterStatus = append(filterStatus, fmt.Sprintf("Sort:%s%s", m.sortOpts.Field, sortIcon))

	filterStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Italic(true)
	b.WriteString(filterStyle.Render(strings.Join(filterStatus, " | ")))
	b.WriteString("\n\n")

	// Render List
	if len(m.visualToArray) == 0 {
		b.WriteString(emptyStateStyle.Render("No items found matching current filters."))
		b.WriteString("\n")
	} else {
		for visualPos, arrayIndex := range m.visualToArray {
			b.WriteString(m.renderItem(arrayIndex, visualPos, visualPos+1))
		}
	}

	b.WriteString("\n")

	footer := statusBarStyle.Render("↑/↓: Nav | Enter: Toggle | d: Delete | f:Filter t:Time p:Path s:Sort o:Order | q: Quit")
	b.WriteString(footer)

	return b.String()
}

func (m *Model) renderDeleteConfirmation() string {
	if m.deleteTargetIndex < 0 || m.deleteTargetIndex >= len(m.items) {
		return ""
	}

	item := m.items[m.deleteTargetIndex]

	var b strings.Builder

	warningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true).
		Padding(1, 2)

	b.WriteString(warningStyle.Render("\u26a0 DELETE CONFIRMATION"))
	b.WriteString("\n\n")

	detailStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 4)

	b.WriteString(detailStyle.Render("You are about to permanently delete:"))
	b.WriteString("\n\n")

	itemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Padding(0, 6)

	itemText := item.Description
	if item.IsCodeReference() {
		itemText = fmt.Sprintf("%s [%s]", item.Description, item.GetLocation())
	}

	b.WriteString(itemStyle.Render(itemText))
	b.WriteString("\n\n")

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 4)

	b.WriteString(promptStyle.Render("This action cannot be undone."))
	b.WriteString("\n\n")

	actionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 4)

	b.WriteString(actionStyle.Render("Press 'y' to confirm, 'n' or 'Esc' to cancel"))

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

	text = fmt.Sprintf("%d. %s", displayIndex, item.Description)

	if item.IsCodeReference() {
		text = fmt.Sprintf("%d. %s [%s]", displayIndex, item.Description, item.GetLocation())
	}

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
