package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Title style for the header
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	// Subtitle style for section headers
	subtitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 0, 0, 2)

	// Selected item style (cursor is on this item)
	selectedItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		PaddingLeft(2)

	// Normal open item style
	normalItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		PaddingLeft(4)

	// Completed item style
	completedItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Strikethrough(true).
		PaddingLeft(4)

	// Selected completed item style
	selectedCompletedItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Strikethrough(true).
		Bold(true).
		PaddingLeft(2)

	// Status bar style (footer with help text)
	statusBarStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	// Empty state message style
	emptyStateStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true).
		PaddingLeft(4)

	// Cursor symbol
	cursor = "❯"
)
