package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tito-sala/rusky/internal/debt"
)

const defaultFilePath = ".rusky.json"

// Version constants for .rusky.json format
const (
	Version010     = "0.1.0"
	Version020     = "0.2.0"
	CurrentVersion = Version020
)

// fileFormat represents the structure of the .rusky.json file
type fileFormat struct {
	Version string          `json:"version"`
	Items   []debt.DebtItem `json:"items"`
}

// JSONStorage implements the Storage interface using a JSON file
type JSONStorage struct {
	filepath string
}

// NewJSONStorage creates a new JSONStorage instance
// If filepath is empty, it uses the default ".rusky.json"
func NewJSONStorage(filepath string) *JSONStorage {
	if filepath == "" {
		filepath = defaultFilePath
	}
	return &JSONStorage{filepath: filepath}
}

// Load reads and parses the JSON file, returning all debt items
// Returns an empty slice if the file doesn't exist (new project)
func (s *JSONStorage) Load() ([]debt.DebtItem, error) {
	// Check if file exists
	if _, err := os.Stat(s.filepath); os.IsNotExist(err) {
		return []debt.DebtItem{}, nil
	}

	// Read file
	data, err := os.ReadFile(s.filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON
	var fileData fileFormat
	if err := json.Unmarshal(data, &fileData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Migrate from older versions if needed
	if fileData.Version == Version010 {
		// v0.1.0 items don't have the new fields, but Go will initialize them properly:
		// - FilePath, LineNumber, CommentType will be nil
		// - IsScanned will be false (zero value)
		// No explicit migration needed, fields are already correct
	}

	return fileData.Items, nil
}

// Save writes all debt items to the JSON file
func (s *JSONStorage) Save(items []debt.DebtItem) error {
	// Prepare data structure
	fileData := fileFormat{
		Version: CurrentVersion,
		Items:   items,
	}

	// Marshal to JSON with indentation for readability
	data, err := json.MarshalIndent(fileData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file with appropriate permissions
	if err := os.WriteFile(s.filepath, data, 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
