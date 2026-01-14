package intro

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Color definitions for purple gradient theme
var (
	darkPurple   = lipgloss.Color("#4A3A6B")
	midPurple    = lipgloss.Color("#6647A3")
	brightPurple = lipgloss.Color("#7D56F4")
	white        = lipgloss.Color("#FFFFFF")
	gray         = lipgloss.Color("#888888")
)

// STYLE 3: THE PEEKING DOG
// Only the top half of the dog is visible. Very stable animation.

var DogAtComputerRaw = `
      [ LOADING... ]
      
       / \__
      (    -\___
      /         O
     /   (_____/
    /_____/   U`

// DogAtComputerBarkRaw is a subtle alternate frame (same layout) used to simulate barking/talking.
var DogAtComputerBarkRaw = `
      [ LOADING... ]

       / \__
      (    @\___
      /         O
     /   (_____/
    /_____/   U`

var DogSleepingRaw = `
       z Z z
      
       / \__
      (    -\___
      /         O
     /   (_____/
    /_____/   U`

var DogAlertRaw = `
       !     !
       
       / \__
      (    O\___
      /         O
     /   (_____/
    /_____/   U`

var DogGuardingRaw = `
      [ SECURITY ]
      
       / \__
      (    ▼\___
      /    |    O
     /   (_____/
    /_____/   U`

var DogSniffingRaw = `
      
       / \__
      (    >\___
      /   *sniff*
     /   (_____/
    /_____/   U`

func rawToLines(raw string) []string {
	raw = strings.Trim(raw, "\n")
	return strings.Split(raw, "\n")
}

// DogAtComputer represents a guard dog sitting at a computer workstation
var DogAtComputer = rawToLines(DogAtComputerRaw)

// DogAtComputerBark is an alternate frame for subtle barking/talking animation.
var DogAtComputerBark = rawToLines(DogAtComputerBarkRaw)

// DogSleeping represents a dog resting/sleeping
var DogSleeping = rawToLines(DogSleepingRaw)

// DogAlert represents a dog with perked ears in alert state
var DogAlert = rawToLines(DogAlertRaw)

// DogGuarding represents a dog in sitting guard position
var DogGuarding = rawToLines(DogGuardingRaw)

// DogSniffing represents a dog with nose down, sniffing
var DogSniffing = rawToLines(DogSniffingRaw)

// Helper functions to apply colors to ASCII art

// ApplyDarkPurple applies dark purple color to the entire ASCII art
func ApplyDarkPurple(art []string) string {
	style := lipgloss.NewStyle().Foreground(darkPurple)
	return style.Render(strings.Join(art, "\n"))
}

// ApplyMidPurple applies mid-range purple color to the entire ASCII art
func ApplyMidPurple(art []string) string {
	style := lipgloss.NewStyle().Foreground(midPurple)
	return style.Render(strings.Join(art, "\n"))
}

// ApplyBrightPurple applies bright purple color (Rusky's primary color) to the entire ASCII art
func ApplyBrightPurple(art []string) string {
	style := lipgloss.NewStyle().Foreground(brightPurple)
	return style.Render(strings.Join(art, "\n"))
}

// ApplyWhite applies white color to the entire ASCII art
func ApplyWhite(art []string) string {
	style := lipgloss.NewStyle().Foreground(white)
	return style.Render(strings.Join(art, "\n"))
}

// ApplyGray applies gray color to the entire ASCII art
func ApplyGray(art []string) string {
	style := lipgloss.NewStyle().Foreground(gray)
	return style.Render(strings.Join(art, "\n"))
}

// ApplyGlow applies a glow effect by rendering with bold and bright purple
func ApplyGlow(art []string) string {
	style := lipgloss.NewStyle().
		Foreground(brightPurple).
		Bold(true)
	return style.Render(strings.Join(art, "\n"))
}

// ColorizeByFrame applies color based on animation frame for gradient effects
func ColorizeByFrame(art []string, frame int, maxFrames int) string {
	// Calculate gradient position (0.0 to 1.0)
	progress := float64(frame) / float64(maxFrames)

	// Apply color based on progress through gradient
	if progress < 0.33 {
		return ApplyDarkPurple(art)
	} else if progress < 0.66 {
		return ApplyMidPurple(art)
	}
	return ApplyBrightPurple(art)
}

// RenderProgressBar creates a loading bar with specified percentage
func RenderProgressBar(percent int) string {
	barWidth := 10
	filled := (percent * barWidth) / 100
	empty := barWidth - filled

	var bar strings.Builder
	bar.WriteString("[")
	for range filled {
		bar.WriteString("█")
	}
	for range empty {
		bar.WriteString("░")
	}
	bar.WriteString("]")

	style := lipgloss.NewStyle().Foreground(brightPurple)
	percentText := lipgloss.NewStyle().Foreground(white).Render(fmt.Sprintf("%d%% ", percent))
	return style.Render(bar.String()) + " " + percentText
}

// RenderText renders centered text with specified color
func RenderText(text string, color lipgloss.Color) string {
	style := lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Width(40).
		Align(lipgloss.Center)
	return style.Render(text)
}
