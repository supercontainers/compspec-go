package plugins

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

// getPluginType returns a string to describe the plugin type
func getPluginType(p PluginRequest) string {

	if p.Plugin.IsCreator() && p.Plugin.IsExtractor() {
		return "extractor and creator"
	}
	if p.Plugin.IsExtractor() {
		return "extractor"
	}
	return "creator"
}

// List plugins available, print in a pretty table!
func (r *PluginsRequest) List() error {

	// Write out table with nodes
	t := table.NewWriter()
	t.SetTitle("Compatibility Plugins")
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"", "Type", "Name", "Section"})
	t.AppendSeparator()

	// keep count of plugins, total, and for each kind
	count := 0
	extractorCount := 0
	creatorCount := 0

	// Do creators first in one section (only a few)
	t.AppendSeparator()
	t.AppendRow(table.Row{"creation plugins", "", "", ""})

	// TODO add description column
	for _, p := range *r {

		if !p.Plugin.IsCreator() {
			continue
		}
		pluginType := getPluginType(p)

		// Creators don't have sections necessarily
		creatorCount += 1
		count += 1

		// Allow plugins to serve dual purposes
		// TODO what should sections be used for?
		t.AppendRow([]interface{}{"", pluginType, p.Name, ""})
	}

	// TODO add description column
	for _, p := range *r {

		if p.Plugin.IsExtractor() {
			extractorCount += 1
		}

		newPlugin := true
		pluginType := getPluginType(p)

		// Extractors are parsed by sections
		for _, section := range p.Plugin.Sections() {

			// Add the extractor plugin description only for first in the list
			if newPlugin {
				t.AppendSeparator()
				t.AppendRow(table.Row{p.Plugin.Description(), "", "", ""})
				newPlugin = false
			}
			count += 1

			// Allow plugins to serve dual purposes
			t.AppendRow([]interface{}{"", pluginType, p.Name, section})
		}
	}
	t.AppendSeparator()
	t.AppendFooter(table.Row{"Total", "", extractorCount + creatorCount, count})
	t.SetStyle(table.StyleColoredCyanWhiteOnBlack)
	t.Render()
	return nil
}
