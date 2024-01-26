package plugins

import (
	"fmt"
	"strings"
)

// Field holds (previously) flattened metadata for an <extractor>.<section>.<field>
type Field struct {
	Extractor string
	Section   string
	Field     string
}

// parseField parses a flattened <extractor>.<section>.<field> into components
func ParseField(field string) (Field, error) {

	f := Field{}
	// Get the extractor from the field
	parts := strings.Split(field, ".")

	// We need at least an extractor name, section, and value
	if len(parts) < 3 {
		return f, fmt.Errorf("warning: field %s value needs to have at least <extractor>.<section>.<field>\n", field)

	}
	f.Extractor = parts[0]
	f.Section = parts[1]
	f.Field = strings.Join(parts[2:], ".")
	return f, nil
}
