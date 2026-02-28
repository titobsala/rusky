package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/tito-sala/rusky/internal/debt"
)

// Color palette matching the TUI theme
const (
	ColorPurple = "#7D56F4"
	ColorWhite  = "#FFFFFF"
	ColorGray   = "#888888"
	ColorGreen  = "#00AA00"
	ColorRed    = "#FF4444"
	ColorOrange = "#FF8800"
	ColorYellow = "#FFDD00"
	ColorBlue   = "#4488FF"
)

// Status symbols
const (
	SymbolOpen      = "●"
	SymbolCompleted = "✓"
)

var (
	// TitleStyle for the main title
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(ColorPurple)).
			MarginBottom(1)

	// TableBorderStyle for table borders
	TableBorderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorGray))

	// HeaderStyle for table headers
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(ColorWhite))

	// FooterStyle for the summary footer
	FooterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorWhite)).
			Background(lipgloss.Color(ColorPurple)).
			Padding(0, 1).
			MarginTop(1)

	// OpenStatusStyle for open items
	OpenStatusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorWhite))

	// CompletedStatusStyle for completed items
	CompletedStatusStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorGreen)).
				Bold(true)

	// EmptyStateStyle for empty state messages
	EmptyStateStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorGray)).
			Italic(true)

	// SuccessStyle for success messages
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorGreen)).
			Bold(true)

	// ErrorStyle for error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorRed)).
			Bold(true)
)

// GetStatusSymbol returns a styled status symbol for the given debt item
func GetStatusSymbol(item debt.DebtItem) string {
	if item.IsCompleted() {
		return CompletedStatusStyle.Render(SymbolCompleted)
	}
	return OpenStatusStyle.Render(SymbolOpen)
}

// priorityColors maps priority levels to colors
var priorityColors = map[debt.Priority]lipgloss.Color{
	debt.PriorityCritical: lipgloss.Color(ColorRed),
	debt.PriorityHigh:     lipgloss.Color(ColorOrange),
	debt.PriorityMedium:   lipgloss.Color(ColorYellow),
	debt.PriorityLow:      lipgloss.Color(ColorBlue),
	debt.PriorityNone:     lipgloss.Color(ColorGray),
}

// priorityLabels maps priority levels to display labels
var priorityLabels = map[debt.Priority]string{
	debt.PriorityCritical: "critical",
	debt.PriorityHigh:     "high",
	debt.PriorityMedium:   "medium",
	debt.PriorityLow:      "low",
	debt.PriorityNone:     "none",
}

// GetPriorityLabel returns a human-readable priority label
func GetPriorityLabel(priority debt.Priority) string {
	if label, ok := priorityLabels[priority]; ok {
		return label
	}
	return "none"
}

// GetPriorityStyle returns a styled priority label with color
func GetPriorityStyle(priority debt.Priority) string {
	color, ok := priorityColors[priority]
	if !ok {
		color = lipgloss.Color(ColorGray)
	}
	style := lipgloss.NewStyle().
		Foreground(color).
		Bold(true)
	return style.Render(GetPriorityLabel(priority))
}
