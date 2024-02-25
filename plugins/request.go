package plugins

import (
	"fmt"

	"github.com/compspec/compspec-go/pkg/plugin"
	pg "github.com/compspec/compspec-go/pkg/plugin"
)

// A plugin request has a Name and sections
type PluginRequest struct {
	Name     string
	Sections []string
	Plugin   pg.PluginInterface
}

type PluginsRequest []PluginRequest

// Do the extraction for a plugin request, meaning across a set of plugins
func (r *PluginsRequest) Extract(allowFail bool) (pg.Result, error) {

	// Prepare Result
	result := pg.Result{}
	results := map[string]pg.PluginData{}

	for _, p := range *r {

		// Skip plugins that don't define extraction
		if !p.Plugin.IsExtractor() {
			continue
		}
		r, err := p.Plugin.Extract(p.Sections)

		// We can allow failure
		if err != nil && !allowFail {
			return result, fmt.Errorf("there was an extraction error for %s: %s", p.Name, err)
		} else if err != nil && allowFail {
			fmt.Printf("Allowing failure - ignoring extraction error for %s: %s\n", p.Name, err)
		}
		results[p.Name] = r
	}
	result.Results = results
	return result, nil
}

// Do creation
func (r *PluginsRequest) Create() (pg.Result, error) {

	// Prepare Result
	result := pg.Result{}

	for _, p := range *r {

		// Skip plugins that don't define extraction
		if !p.Plugin.IsCreator() {
			continue
		}
		options := plugin.PluginOptions{}
		err := p.Plugin.Create(options)
		if err != nil {
			return result, err
		}

	}
	return result, nil
}
