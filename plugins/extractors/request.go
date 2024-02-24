package extractors

import (
	"fmt"

	"github.com/compspec/compspec-go/pkg/extractor"
	"github.com/compspec/compspec-go/plugins"
)

// A plugin request has a Name and sections
type PluginRequest struct {
	Name      string
	Sections  []string
	Extractor extractor.Extractor
}

// These functions make it possible to use the PluginRequest as a PluginInformation interface
func (p *PluginRequest) GetName() string {
	return p.Name
}
func (p *PluginRequest) GetType() string {
	return "extractor"
}
func (p *PluginRequest) GetDescription() string {
	return p.Extractor.Description()
}
func (p *PluginRequest) GetSections() []plugins.PluginSection {
	sections := make([]plugins.PluginSection, len(p.Extractor.Sections()))

	for _, section := range p.Extractor.Sections() {
		newSection := plugins.PluginSection{Name: section}
		sections = append(sections, newSection)
	}
	return sections
}

type PluginsRequest []PluginRequest

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
