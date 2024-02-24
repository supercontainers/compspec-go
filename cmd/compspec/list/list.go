package list

import (
	"github.com/compspec/compspec-go/plugins/extractors"

	p "github.com/compspec/compspec-go/plugins"
)

// Run will list the extractor names and sections known
func Run(pluginNames []string) error {

	// parse [section,...,section] into named plugins and sections
	// return plugins
	plugins, err := extractors.GetPlugins(pluginNames)
	if err != nil {
		return err
	}
	// Convert to plugin information
	info := []p.PluginInformation{}
	for _, p := range plugins {
		info = append(info, &p)
	}
	// List plugin table
	return p.List(info)
}
