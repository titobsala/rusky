package debt

import (
	"fmt"
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
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	ID          string     `json:"id"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`

	// v0.2.0 fields for scanned items
	FilePath    *string `json:"file_path,omitempty"`    // Relative path to source file
	LineNumber  *int    `json:"line_number,omitempty"`  // Line number in file
	CommentType *string `json:"comment_type,omitempty"` // TODO, FIXME, HACK, etc.
	IsScanned   bool    `json:"is_scanned"`             // Auto-scanned vs manually added
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

// IsCodeReference returns true if this item references source code
func (d *DebtItem) IsCodeReference() bool {
	return d.FilePath != nil && d.LineNumber != nil
}

// GetLocation returns a human-readable location string (e.g., "src/main.go:123")
func (d *DebtItem) GetLocation() string {
	if !d.IsCodeReference() {
		return ""
	}
	return fmt.Sprintf("%s:%d", *d.FilePath, *d.LineNumber)
}

// Storage defines the interface for persisting debt items
type Storage interface {
	// Load retrieves all debt items from storage
	Load() ([]DebtItem, error)

	// Save persists all debt items to storage
	Save(items []DebtItem) error
}
