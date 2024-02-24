package plugins

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

// List plugins available, print in a pretty table!
func List(ps []PluginInformation) error {

	// Write out table with nodes
	t := table.NewWriter()
	t.SetTitle("Compatibility Plugins")
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"", "Type", "Name", "Section"})
	t.AppendSeparator()

	// keep count of plugins (just extractors for now)
	count := 0
	pluginCount := 0

	// This will iterate across plugin types (e.g., extraction and converter)
	for _, p := range ps {
		pluginCount += 1

		// This iterates across plugins in the family
		for i, section := range p.GetSections() {

			// Add the extractor plugin description only for first in the list
			if i == 0 {
				t.AppendSeparator()
				t.AppendRow(table.Row{p.GetDescription(), "", "", ""})
			}

			count += 1
			t.AppendRow([]interface{}{"", p.GetType(), section.Name})
		}

	}
	t.AppendSeparator()
	t.AppendFooter(table.Row{"Total", "", pluginCount, count})
	t.SetStyle(table.StyleColoredCyanWhiteOnBlack)
	t.Render()
	return nil
}
