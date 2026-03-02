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

// Delete permanently removes a debt item from storage
// The identifier can be either a UUID or a 1-based index
func (m *Manager) Delete(identifier string) (*DebtItem, error) {
	items, err := m.storage.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load items: %w", err)
	}

	// Try to find item by UUID or index
	index := m.findItemIndex(items, identifier)
	if index == -1 {
		return nil, fmt.Errorf("item not found: %s", identifier)
	}

	// Store the item for return value
	deletedItem := items[index]

	// Remove item from slice
	items = append(items[:index], items[index+1:]...)

	// Save
	if err := m.storage.Save(items); err != nil {
		return nil, fmt.Errorf("failed to save items: %w", err)
	}

	return &deletedItem, nil
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

// FindByLocation finds a debt item by file path and line number
// Returns nil if not found
func (m *Manager) FindByLocation(filePath string, lineNumber int) (*DebtItem, error) {
	items, err := m.storage.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load items: %w", err)
	}

	for i := range items {
		item := &items[i]
		if item.FilePath != nil && item.LineNumber != nil {
			if *item.FilePath == filePath && *item.LineNumber == lineNumber {
				return item, nil
			}
		}
	}

	return nil, nil
}

// UpdateDescription updates the description of a debt item by ID
func (m *Manager) UpdateDescription(id string, newDescription string) error {
	items, err := m.storage.Load()
	if err != nil {
		return fmt.Errorf("failed to load items: %w", err)
	}

	for i := range items {
		if items[i].ID == id {
			items[i].Description = newDescription
			if err := m.storage.Save(items); err != nil {
				return fmt.Errorf("failed to save items: %w", err)
			}
			return nil
		}
	}

	return fmt.Errorf("item not found: %s", id)
}

// SetPriority updates the priority of a debt item by ID
func (m *Manager) SetPriority(id string, priority Priority) (*DebtItem, error) {
	items, err := m.storage.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load items: %w", err)
	}

	index := m.findItemIndex(items, id)
	if index == -1 {
		return nil, fmt.Errorf("item not found: %s", id)
	}

	items[index].Priority = priority

	if err := m.storage.Save(items); err != nil {
		return nil, fmt.Errorf("failed to save items: %w", err)
	}

	return &items[index], nil
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
