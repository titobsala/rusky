package scanner

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ScanResult represents a found technical debt comment
type ScanResult struct {
	FilePath    string // Relative to scan root
	LineNumber  int    // 1-based line number
	CommentType string // TODO, FIXME, etc.
	Description string // Extracted comment text
	Line        string // Full line content for preview
}

// Scanner handles code scanning operations
type Scanner struct {
	patterns     []CommentPattern
	excludeDirs  []string
	excludeFiles []string
}

// NewScanner creates a new scanner with default patterns and exclusions
func NewScanner() *Scanner {
	return &Scanner{
		patterns: buildDefaultPatterns(),
		excludeDirs: []string{
			"node_modules", "vendor", ".git", ".hg", ".svn",
			"dist", "build", "target", "__pycache__", ".pytest_cache",
			".venv", "venv", "env", ".next", ".nuxt", "out",
		},
		excludeFiles: []string{
			".min.js", ".min.css", ".map",
		},
	}
}

// Scan scans a directory recursively for technical debt markers
func (s *Scanner) Scan(rootPath string) ([]ScanResult, error) {
	results := make([]ScanResult, 0)

	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	err = filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip excluded directories
		if d.IsDir() && s.shouldExcludeDir(d.Name()) {
			return filepath.SkipDir
		}

		// Skip excluded files
		if !d.IsDir() && s.shouldExcludeFile(d.Name()) {
			return nil
		}

		// Only scan text files
		if !d.IsDir() && s.isTextFile(d.Name()) {
			relPath, _ := filepath.Rel(absRoot, path)
			fileResults, err := s.scanFile(path, relPath)
			if err != nil {
				// Log error but continue scanning
				fmt.Fprintf(os.Stderr, "Warning: failed to scan %s: %v\n", relPath, err)
				return nil
			}
			results = append(results, fileResults...)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	return results, nil
}

// scanFile scans a single file for technical debt markers
func (s *Scanner) scanFile(absPath, relPath string) ([]ScanResult, error) {
	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	results := make([]ScanResult, 0)
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Check each pattern
		for _, pattern := range s.patterns {
			if matches := pattern.Pattern.FindStringSubmatch(line); matches != nil {
				// Extract description (capture group 1)
				description := strings.TrimSpace(matches[1])

				results = append(results, ScanResult{
					FilePath:    relPath,
					LineNumber:  lineNum,
					CommentType: pattern.Type,
					Description: description,
					Line:        strings.TrimSpace(line),
				})
				break // Only match one pattern per line
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
