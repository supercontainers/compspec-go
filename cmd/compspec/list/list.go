package list

import (
	p "github.com/supercontainers/compspec-go/plugins"
)

// Run will list the extractor names and sections known
func Run(pluginNames []string) error {
	// parse [section,...,section] into named plugins and sections
	// return plugins
	plugins, err := p.GetPlugins(pluginNames)
	if err != nil {
		return err
	}
	// List plugin table
	return plugins.List()
}
