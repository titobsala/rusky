package intro

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderAnimation is the main dispatcher that renders the appropriate animation based on type and frame
func RenderAnimation(animType AnimationType, frame int) string {
	switch animType {
	case Lightning:
		return renderLightning(frame)
	case Scan:
		return renderScan(frame)
	case AlertMode:
		return renderAlertMode(frame)
	default:
		return renderLightning(frame)
	}
}

// renderLightning implements the "Guard Dog at Computer" animation with loading bar
// Phase 1 (0-15): Dog at computer with empty bar, dark purple
// Phase 2 (16-30): Progress bar fills 0-100%, gradient to brighter purple
// Phase 3 (31-33): Complete bar with "READY TO GUARD" message
// Phase 4 (34-42): Glitch transition with purple morphing
// Phase 5 (43+): Bright purple with glow, "SYSTEM ARMED" status
func renderLightning(frame int) string {
	var output strings.Builder

	// Phase 1: Initial state with empty progress bar
	if frame <= 15 {
		dog := ColorizeByFrame(DogAtComputer, frame, 15)
		bar := RenderProgressBar(0)
		title := RenderText("RUSKY SECURITY SYSTEM", white)

		output.WriteString("\n\n")
		output.WriteString(title)
		output.WriteString("\n\n")
		output.WriteString(dog)
		output.WriteString("\n")
		output.WriteString(lipgloss.NewStyle().Width(40).Align(lipgloss.Center).Render(bar))

		// Phase 2: Progress bar filling
	} else if frame <= 30 {
		// Calculate progress (0-100%)
		progress := ((frame - 16) * 100) / 14

		// Gradient effect as progress increases
		dog := ColorizeByFrame(DogAtComputer, frame-16, 14)
		bar := RenderProgressBar(progress)
		title := RenderText("INITIALIZING...", midPurple)

		output.WriteString("\n\n")
		output.WriteString(title)
		output.WriteString("\n\n")
		output.WriteString(dog)
		output.WriteString("\n")
		output.WriteString(lipgloss.NewStyle().Width(40).Align(lipgloss.Center).Render(bar))

		// Phase 3: Complete with ready message
	} else if frame <= 33 {
		dog := ApplyBrightPurple(DogAtComputer)
		bar := RenderProgressBar(100)
		title := RenderText("READY TO GUARD", brightPurple)

		output.WriteString("\n\n")
		output.WriteString(title)
		output.WriteString("\n\n")
		output.WriteString(dog)
		output.WriteString("\n")
		output.WriteString(lipgloss.NewStyle().Width(40).Align(lipgloss.Center).Render(bar))

		// Phase 4: Glitch transition
	} else if frame <= 42 {
		// Oscillate between colors for glitch effect
		var dog string
		if frame%2 == 0 {
			dog = ApplyBrightPurple(DogAtComputer)
		} else {
			dog = ApplyMidPurple(DogAtComputer)
		}

		title := RenderText("SYSTEM CHECK...", white)

		output.WriteString("\n\n")
		output.WriteString(title)
		output.WriteString("\n\n")
		output.WriteString(dog)

		// Phase 5: Final glow state
	} else {
		dog := ApplyGlow(DogAtComputer)
		title := RenderText("✓ SYSTEM ARMED ✓", brightPurple)
		subtitle := RenderText("Technical Debt Detection Active", gray)

		output.WriteString("\n\n")
		output.WriteString(title)
		output.WriteString("\n")
		output.WriteString(subtitle)
		output.WriteString("\n\n")
		output.WriteString(dog)
	}

	return output.String()
}

// renderScan implements the "Dog Sniffing Out Debt" line-by-line scanning animation
// Progressive scanning with percentage counter
func renderScan(frame int) string {
	var output strings.Builder

	// Calculate scan progress
	totalLines := len(DogSniffing)
	scannedLines := min((frame*totalLines)/GetMaxFrames(Scan), totalLines)

	// Calculate percentage
	percent := min((frame*100)/GetMaxFrames(Scan), 100)

	title := RenderText(fmt.Sprintf("DEBT SCAN: %d%%", percent), white)

	output.WriteString("\n\n")
	output.WriteString(title)
	output.WriteString("\n\n")

	// Render dog with progressive color application
	for i, line := range DogSniffing {
		var styledLine string
		if i < scannedLines {
			// Already scanned - bright purple
			styledLine = lipgloss.NewStyle().Foreground(brightPurple).Render(line)
		} else if i == scannedLines {
			// Currently scanning - white highlight
			styledLine = lipgloss.NewStyle().Foreground(white).Bold(true).Render(line)
		} else {
			// Not yet scanned - dark purple
			styledLine = lipgloss.NewStyle().Foreground(darkPurple).Render(line)
		}
		output.WriteString(styledLine)
		output.WriteString("\n")
	}

	// Add scanning indicator
	if percent < 100 {
		indicator := RenderText(">>> SNIFFING FOR TECHNICAL DEBT <<<", midPurple)
		output.WriteString("\n")
		output.WriteString(indicator)
	} else {
		complete := RenderText("✓ SCAN COMPLETE ✓", brightPurple)
		output.WriteString("\n")
		output.WriteString(complete)
	}

	return output.String()
}

// renderAlertMode implements the "Sleeping to Guard" transition animation
// Phase 1 (0-10): Dog sleeping, dark purple with zzz
// Phase 2 (11-15): Alert waves, ears perk up
// Phase 3 (16-25): Dog stands up, eyes open
// Phase 4 (26-30): Guard position
// Phase 5 (31-40): Purple gradient sweep, "ON DUTY"
// Phase 6 (41-45): Final glow, ready state
func renderAlertMode(frame int) string {
	var output strings.Builder

	// Phase 1: Sleeping
	if frame <= 10 {
		dog := ApplyDarkPurple(DogSleeping)
		title := RenderText("[ STANDBY MODE ]", gray)

		output.WriteString("\n\n")
		output.WriteString(title)
		output.WriteString("\n\n")
		output.WriteString(dog)

		// Phase 2: Alert sound
	} else if frame <= 15 {
		dog := ApplyMidPurple(DogSleeping)
		alert := "))) ALERT ((("
		alertText := lipgloss.NewStyle().
			Foreground(white).
			Bold(true).
			Width(40).
			Align(lipgloss.Center).
			Render(alert)

		title := RenderText("[ ALERT DETECTED ]", midPurple)

		output.WriteString("\n\n")
		output.WriteString(alertText)
		output.WriteString("\n")
		output.WriteString(title)
		output.WriteString("\n\n")
		output.WriteString(dog)

		// Phase 3: Standing up, alert posture
	} else if frame <= 25 {
		dog := ColorizeByFrame(DogAlert, frame-16, 9)
		title := RenderText("[ ACTIVATING ]", midPurple)

		output.WriteString("\n\n")
		output.WriteString(title)
		output.WriteString("\n\n")
		output.WriteString(dog)

		// Phase 4: Guard position
	} else if frame <= 30 {
		dog := ApplyBrightPurple(DogGuarding)
		title := RenderText("[ READY ]", brightPurple)

		output.WriteString("\n\n")
		output.WriteString(title)
		output.WriteString("\n\n")
		output.WriteString(dog)

		// Phase 5: Purple gradient sweep with "ON DUTY"
	} else if frame <= 40 {
		dog := ColorizeByFrame(DogGuarding, frame-31, 9)
		title := RenderText("ON DUTY", brightPurple)
		subtitle := RenderText("Guarding Code Quality", white)

		output.WriteString("\n\n")
		output.WriteString(title)
		output.WriteString("\n")
		output.WriteString(subtitle)
		output.WriteString("\n\n")
		output.WriteString(dog)

		// Phase 6: Final glow
	} else {
		dog := ApplyGlow(DogGuarding)
		title := RenderText("✓ GUARD MODE ACTIVE ✓", brightPurple)
		subtitle := RenderText("Technical Debt Monitoring Enabled", gray)

		output.WriteString("\n\n")
		output.WriteString(title)
		output.WriteString("\n")
		output.WriteString(subtitle)
		output.WriteString("\n\n")
		output.WriteString(dog)
	}

	return output.String()
}
