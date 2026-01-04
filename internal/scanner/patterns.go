package scanner

import "regexp"

// CommentPattern represents a technical debt marker pattern
type CommentPattern struct {
	Type    string         // TODO, FIXME, HACK, XXX, BUG, NOTE
	Pattern *regexp.Regexp // Compiled regex pattern
}

// buildDefaultPatterns creates regex patterns for all comment styles
func buildDefaultPatterns() []CommentPattern {
	markers := []string{"TODO", "FIXME", "HACK", "XXX", "BUG", "NOTE"}
	patterns := make([]CommentPattern, 0, len(markers)*4)

	for _, marker := range markers {
		// Pattern 1: // TODO: description (C-style single line)
		patterns = append(patterns, CommentPattern{
			Type:    marker,
			Pattern: regexp.MustCompile(`//\s*` + marker + `[:\s]+(.+)`),
		})

		// Pattern 2: # TODO: description (Shell-style)
		patterns = append(patterns, CommentPattern{
			Type:    marker,
			Pattern: regexp.MustCompile(`#\s*` + marker + `[:\s]+(.+)`),
		})

		// Pattern 3: /* TODO: description */ or /* TODO: description (C-style block)
		patterns = append(patterns, CommentPattern{
			Type:    marker,
			Pattern: regexp.MustCompile(`/\*\s*` + marker + `[:\s]+(.+?)(?:\s*\*/|$)`),
		})

		// Pattern 4: <!-- TODO: description --> (HTML/XML)
		patterns = append(patterns, CommentPattern{
			Type:    marker,
			Pattern: regexp.MustCompile(`<!--\s*` + marker + `[:\s]+(.+?)\s*-->`),
		})
	}

	return patterns
}
