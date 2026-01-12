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
		return 60 // ~3.0 seconds @ 20 FPS - smoother transitions
	case Scan:
		return 40 // ~2.0 seconds @ 20 FPS - keep as is
	case AlertMode:
		return 55 // ~2.75 seconds @ 20 FPS - smoother transitions
	default:
		return 40
	}
}
