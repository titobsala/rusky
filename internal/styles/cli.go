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
				Foreground(lipgloss.Color(ColorGray))

	// EmptyStateStyle for empty state messages
	EmptyStateStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorGray)).
			Italic(true)
)

// GetStatusSymbol returns a styled status symbol for the given debt item
func GetStatusSymbol(item debt.DebtItem) string {
	if item.IsCompleted() {
		return CompletedStatusStyle.Render(SymbolCompleted)
	}
	return OpenStatusStyle.Render(SymbolOpen)
}
