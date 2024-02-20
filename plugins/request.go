package plugins

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/compspec/compspec-go/pkg/extractor"
)

// A plugin request has a Name and sections
type PluginRequest struct {
	Name      string
	Sections  []string
	Extractor extractor.Extractor
	// TODO add checker here eventually too.
}

type PluginsRequest []PluginRequest

// A Result wraps named extractor data, just for easy dumping to json
type Result struct {
	Results map[string]extractor.ExtractorData `json:"extractors,omitempty"`
}

// ToJson serializes a result to json
func (r *Result) ToJson() ([]byte, error) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return []byte{}, err
	}
	return b, err
}

// Print prints the result to the terminal
func (r *Result) Print() {
	for name, result := range r.Results {
		fmt.Printf(" --Result for %s\n", name)
		result.Print()
	}
}

// AddCustomFields adds or updates an existing result with
// custom metadata, either new or to overwrite
func (r *Result) AddCustomFields(fields []string) {

	for _, field := range fields {
		if !strings.Contains(field, "=") {
			fmt.Printf("warning: field %s does not contain an '=', skipping\n", field)
			continue
		}
		parts := strings.Split(field, "=")
		if len(parts) < 2 {
			fmt.Printf("warning: field %s has an empty value, skipping\n", field)
			continue
		}

		// No reason the value cannot have additional =
		field = parts[0]
		value := strings.Join(parts[1:], "=")

		// Get the extractor, section, and subfield from the field
		f, err := ParseField(field)
		if err != nil {
			fmt.Printf(err.Error(), field)
			continue
		}

		// Is the extractor name in the result?
		_, ok := r.Results[f.Extractor]
		if !ok {
			sections := extractor.Sections{}
			r.Results[f.Extractor] = extractor.ExtractorData{Sections: sections}
		}
		data := r.Results[f.Extractor]

		// Is the section name in the extractor data?
		_, ok = data.Sections[f.Section]
		if !ok {
			data.Sections[f.Section] = extractor.ExtractorSection{}
		}
		section := data.Sections[f.Section]
		section[f.Field] = value

		// Wrap it back up!
		data.Sections[f.Section] = section
		r.Results[f.Extractor] = data
	}
}

// Do the extraction for a plugin request, meaning across a set of plugins
func (r *PluginsRequest) Extract(allowFail bool) (Result, error) {

	// Prepare Result
	result := Result{}
	results := map[string]extractor.ExtractorData{}

	for _, p := range *r {
		r, err := p.Extractor.Extract(p.Sections)

		// We can allow failure
		if err != nil && !allowFail {
			return result, fmt.Errorf("There was an extraction error for %s: %s\n", p.Name, err)
		} else if err != nil && allowFail {
			fmt.Printf("Allowing failure - ignoring extraction error for %s: %s\n", p.Name, err)
		}
		results[p.Name] = r
	}
	result.Results = results
	return result, nil
}
