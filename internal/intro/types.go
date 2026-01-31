package intro

// AnimationType represents the type of intro animation to display
type AnimationType int

const (
	// Main represents the high-resolution GIF animation
	Main AnimationType = iota
)

// GetMaxFrames returns the total number of frames for a given animation type
func GetMaxFrames(animType AnimationType) int {
	return 80 // 4 seconds @ 20 FPS
}
