package intro

// AnimationType represents the type of intro animation to display
type AnimationType int

const (
	// Lightning represents the "Guard Dog at Computer" animation with loading bar
	Lightning AnimationType = iota
	// Scan represents the "Dog Sniffing Out Debt" line-by-line scanning animation
	Scan
	// AlertMode represents the "Sleeping to Guard" transition animation
	AlertMode
)

// GetMaxFrames returns the total number of frames for a given animation type
func GetMaxFrames(animType AnimationType) int {
	switch animType {
	case Lightning:
		return 50 // ~2.5 seconds @ 20 FPS
	case Scan:
		return 40 // ~2.0 seconds @ 20 FPS
	case AlertMode:
		return 45 // ~2.25 seconds @ 20 FPS
	default:
		return 40
	}
}
