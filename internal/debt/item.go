package debt

import (
	"time"
)

// Status represents the current state of a debt item
type Status string

const (
	StatusOpen      Status = "open"
	StatusCompleted Status = "completed"
)

// DebtItem represents a single technical debt item
type DebtItem struct {
	ID          string     `json:"id"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// IsCompleted returns true if the debt item is marked as completed
func (d *DebtItem) IsCompleted() bool {
	return d.Status == StatusCompleted
}

// Complete marks the debt item as completed with the current timestamp
func (d *DebtItem) Complete() {
	now := time.Now()
	d.Status = StatusCompleted
	d.CompletedAt = &now
}

// Reopen marks a completed debt item as open again
func (d *DebtItem) Reopen() {
	d.Status = StatusOpen
	d.CompletedAt = nil
}

// Storage defines the interface for persisting debt items
type Storage interface {
	// Load retrieves all debt items from storage
	Load() ([]DebtItem, error)

	// Save persists all debt items to storage
	Save(items []DebtItem) error
}
