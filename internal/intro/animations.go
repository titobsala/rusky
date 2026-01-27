package intro

import (
	"strings"
)

// RenderAnimation is the main dispatcher that renders the GIF animation and progress bar
func RenderAnimation(animType AnimationType, frame int) string {
	var output strings.Builder
	
	width := 60
	maxFrames := GetMaxFrames(animType)
	
	// Calculate current GIF frame
	// We might have more or fewer GIF frames than the 6-second total (120 frames)
	// so we loop the GIF.
	gifFrameIdx := 0
	if len(preRenderedFrames) > 0 {
		gifFrameIdx = frame % len(preRenderedFrames)
	}

	// Calculate progress for the loading bar (0-100%)
	progress := (frame * 100) / maxFrames
	if progress > 100 {
		progress = 100
	}

	output.WriteString("\n\n")
	
	// Render the GIF frame
	if len(preRenderedFrames) > 0 {
		output.WriteString(preRenderedFrames[gifFrameIdx])
	} else {
		output.WriteString(RenderText("LOADING ANIMATION...", width))
	}
	
	output.WriteString("\n")
	
	// Render the gradient loading bar
	bar := RenderProgressBar(progress, 40)
	output.WriteString(RenderText(bar, width))
	
	output.WriteString("\n\n")
	if progress < 100 {
		output.WriteString(RenderText("INITIALIZING RUSKY PROTOCOLS...", width))
	} else {
		output.WriteString(RenderText("SYSTEM READY", width))
	}

	return output.String()
}