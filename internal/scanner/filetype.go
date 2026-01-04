package scanner

import (
	"path/filepath"
	"strings"
)

// Common text file extensions
var textExtensions = map[string]bool{
	// Programming languages
	".go":    true,
	".js":    true,
	".ts":    true,
	".tsx":   true,
	".jsx":   true,
	".py":    true,
	".rb":    true,
	".java":  true,
	".c":     true,
	".cpp":   true,
	".cc":    true,
	".cxx":   true,
	".h":     true,
	".hpp":   true,
	".hxx":   true,
	".cs":    true,
	".rs":    true,
	".swift": true,
	".kt":    true,
	".kts":   true,
	".scala": true,
	".php":   true,
	".sh":    true,
	".bash":  true,
	".zsh":   true,
	".fish":  true,

	// Web
	".html":   true,
	".htm":    true,
	".css":    true,
	".scss":   true,
	".sass":   true,
	".less":   true,
	".vue":    true,
	".svelte": true,

	// Config/Data
	".yaml": true,
	".yml":  true,
	".json": true,
	".toml": true,
	".xml":  true,
	".md":   true,
	".txt":  true,
	".sql":  true,
	".env":  true,
}

// isTextFile checks if the filename has a text file extension
func (s *Scanner) isTextFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return textExtensions[ext]
}

// shouldExcludeDir checks if a directory should be excluded from scanning
func (s *Scanner) shouldExcludeDir(dirname string) bool {
	for _, excluded := range s.excludeDirs {
		if dirname == excluded {
			return true
		}
	}
	return false
}

// shouldExcludeFile checks if a file should be excluded from scanning
func (s *Scanner) shouldExcludeFile(filename string) bool {
	for _, suffix := range s.excludeFiles {
		if strings.HasSuffix(filename, suffix) {
			return true
		}
	}
	return false
}
