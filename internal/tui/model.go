package tui

import (
	"fmt"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tito-sala/rusky/internal/debt"
	"github.com/tito-sala/rusky/internal/intro"
)

// tickMsg is used for animation frame advancement
type tickMsg time.Time

// Model represents the TUI state
type Model struct {
	items         []debt.DebtItem
	err           error
	manager       *debt.Manager
	cursor        int   // Visual position (0-based)
	visualToArray []int // Maps visual index to array index
	width         int
	height        int
	quitting      bool

	// Delete confirmation state
	showDeleteConfirm bool
	deleteTargetIndex int // Array index of item to delete

	// Intro animation state
	introPhase         int // 0 = show animation, 1 = show main UI
	animationType      intro.AnimationType
	animationFrame     int // current frame number (within current animation)
	animationMaxFrames int // total frames for current animation
}

// buildVisualMapping creates a mapping from visual positions to array indices
// Visual order: open items first, then completed items
func (m *Model) buildVisualMapping() {
	openIndices := make([]int, 0, len(m.items))
	completedIndices := make([]int, 0, len(m.items))

	for i, item := range m.items {
		if item.IsCompleted() {
			completedIndices = append(completedIndices, i)
		} else {
			openIndices = append(openIndices, i)
		}
	}

	m.visualToArray = make([]int, 0, len(m.items))
	m.visualToArray = append(m.visualToArray, openIndices...)
	m.visualToArray = append(m.visualToArray, completedIndices...)
}

// getOpenCount returns the number of open items
func (m *Model) getOpenCount() int {
	count := 0
	for _, item := range m.items {
		if !item.IsCompleted() {
			count++
		}
	}
	return count
}

// resetDeleteConfirmation clears the delete confirmation dialog state
func (m *Model) resetDeleteConfirmation() {
	m.showDeleteConfirm = false
	m.deleteTargetIndex = -1
}

// NewModel creates a new TUI model
func NewModel(manager *debt.Manager) (*Model, error) {
	// Load items immediately
	items, err := manager.List()
	if err != nil {
		return nil, fmt.Errorf("failed to load items: %w", err)
	}

	// Seed random for any future randomized behaviors (kept local and deterministic-safe)
	rand.Seed(time.Now().UnixNano())

	// Pick exactly one intro animation at random (Lightning, Scan, AlertMode).
	animType := intro.AnimationType(rand.Intn(3))

	m := &Model{
		manager:            manager,
		items:              items,
		cursor:             0,
		quitting:           false,
		introPhase:         0, // Start with intro animation
		animationType:      animType,
		animationFrame:     0,
		animationMaxFrames: intro.GetMaxFrames(animType),
	}

	// Build initial visual mapping
	m.buildVisualMapping()

	return m, nil
}

// tick returns a command that sends a tick message after 50ms (20 FPS)
func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Init initializes the model (required by Bubbletea)
func (m *Model) Init() tea.Cmd {
	return tick()
}

// Run launches the TUI
func Run(manager *debt.Manager) error {
	model, err := NewModel(manager)
	if err != nil {
		return err
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}
