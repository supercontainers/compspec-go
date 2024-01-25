package plugins

import (
	"encoding/json"
	"fmt"

	"github.com/supercontainers/compspec-go/pkg/extractor"
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

// Do the extraction for a plugin request, meaning across a set of plugins
func (r *PluginsRequest) Extract() (Result, error) {

	// Prepare Result
	result := Result{}
	results := map[string]extractor.ExtractorData{}

	for _, p := range *r {
		r, err := p.Extractor.Extract(p.Sections)
		if err != nil {
			return result, fmt.Errorf("There was a kernel extraction error: %s\n", err)
		}
		results[p.Name] = r
	}
	result.Results = results
	return result, nil
}
