package extractors

import (
	"strings"

	"github.com/compspec/compspec-go/plugins"
	"github.com/compspec/compspec-go/plugins/extractors/kernel"
	"github.com/compspec/compspec-go/plugins/extractors/library"
	"github.com/compspec/compspec-go/plugins/extractors/nfd"
	"github.com/compspec/compspec-go/plugins/extractors/system"
)

// Add new plugin names here. They should correspond with the package name, then NewPlugin()
var (
	KernelExtractor  = "kernel"
	SystemExtractor  = "system"
	LibraryExtractor = "library"
	NFDExtractor     = "nfd"
	pluginNames      = []string{KernelExtractor, SystemExtractor, LibraryExtractor, NFDExtractor}
)

// Get plugins parses a request and returns a list of plugins
// We honor the order that the plugins and sections are provided in
func GetPlugins(names []string) (PluginsRequest, error) {

	if len(names) == 0 {
		names = pluginNames
	}

	request := PluginsRequest{}

	// Prepare an extractor for each, and validate the requested sections
	// TODO: this could also be done with an init -> Register pattern
	for _, name := range names {

		// If we are given a list of section names, parse.
		name, sections := plugins.ParseSections(name)

		if strings.HasPrefix(name, KernelExtractor) {
			p, err := kernel.NewPlugin(sections)
			if err != nil {
				return request, err
			}
			// Save the name, the instantiated interface, and sections
			pr := PluginRequest{Name: name, Extractor: p, Sections: sections}
			request = append(request, pr)
		}

		if strings.HasPrefix(name, NFDExtractor) {
			p, err := nfd.NewPlugin(sections)
			if err != nil {
				return request, err
			}
			// Save the name, the instantiated interface, and sections
			pr := PluginRequest{Name: name, Extractor: p, Sections: sections}
			request = append(request, pr)
		}

		if strings.HasPrefix(name, SystemExtractor) {
			p, err := system.NewPlugin(sections)
			if err != nil {
				return request, err
			}
			// Save the name, the instantiated interface, and sections
			pr := PluginRequest{Name: name, Extractor: p, Sections: sections}
			request = append(request, pr)
		}

		if strings.HasPrefix(name, LibraryExtractor) {
			p, err := library.NewPlugin(sections)
			if err != nil {
				return request, err
			}
			// Save the name, the instantiated interface, and sections
			pr := PluginRequest{Name: name, Extractor: p, Sections: sections}
			request = append(request, pr)
		}
	}
	return request, nil
}
