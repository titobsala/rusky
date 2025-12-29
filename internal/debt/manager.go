package debt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// Manager handles business logic for debt operations
type Manager struct {
	storage Storage
}

// NewManager creates a new Manager instance
func NewManager(storage Storage) *Manager {
	return &Manager{storage: storage}
}

// Add creates a new debt item with the given description
func (m *Manager) Add(description string) (*DebtItem, error) {
	if description == "" {
		return nil, fmt.Errorf("description cannot be empty")
	}

	// Load existing items
	items, err := m.storage.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load items: %w", err)
	}

	// Create new item
	item := DebtItem{
		ID:          uuid.New().String(),
		Description: description,
		Status:      StatusOpen,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}

	// Append and save
	items = append(items, item)
	if err := m.storage.Save(items); err != nil {
		return nil, fmt.Errorf("failed to save items: %w", err)
	}

	return &item, nil
}

// Complete marks a debt item as completed
// The identifier can be either a UUID or a 1-based index
func (m *Manager) Complete(identifier string) (*DebtItem, error) {
	items, err := m.storage.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load items: %w", err)
	}

	// Try to find item by UUID or index
	index := m.findItemIndex(items, identifier)
	if index == -1 {
		return nil, fmt.Errorf("item not found: %s", identifier)
	}

	// Mark as completed
	items[index].Complete()

	// Save
	if err := m.storage.Save(items); err != nil {
		return nil, fmt.Errorf("failed to save items: %w", err)
	}

	return &items[index], nil
}

// List returns all debt items
func (m *Manager) List() ([]DebtItem, error) {
	items, err := m.storage.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load items: %w", err)
	}
	return items, nil
}

// Save persists the given items to storage
func (m *Manager) Save(items []DebtItem) error {
	if err := m.storage.Save(items); err != nil {
		return fmt.Errorf("failed to save items: %w", err)
	}
	return nil
}

// findItemIndex finds an item by UUID or 1-based index
// Returns -1 if not found
func (m *Manager) findItemIndex(items []DebtItem, identifier string) int {
	// Try UUID match first
	for i, item := range items {
		if item.ID == identifier {
			return i
		}
	}

	// Try index match (1-based)
	if idx, err := strconv.Atoi(identifier); err == nil {
		if idx > 0 && idx <= len(items) {
			return idx - 1
		}
	}

	return -1
}
