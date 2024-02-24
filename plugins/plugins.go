package plugins

import (
	"strings"
)

// parseSections will return sections from the name string
// We could use regex here instead
func ParseSections(raw string) (string, []string) {

	sections := []string{}

	// If no opening bracker, it's just the name
	if !strings.Contains(raw, "[") {
		return raw, sections
	}
	// Get rid of last piece ]
	raw = strings.ReplaceAll(raw, "]", "")

	// Split into two pieces, the
	parts := strings.SplitN(raw, "[", 2)
	name, raw := parts[0], parts[1]

	// Split sections by comma
	sections = strings.Split(raw, ",")
	return name, sections
}
