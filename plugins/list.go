package plugins

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

// List plugins available, print in a pretty table!
func (r *PluginsRequest) List() error {

	// Write out table with nodes
	t := table.NewWriter()
	t.SetTitle("Compatibility Plugins")
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"", "Type", "Name", "Section"})
	t.AppendSeparator()

	// keep count of plugins (just extractors for now)
	count := 0
	extractorCount := 0

	// TODO add description column
	for _, p := range *r {
		extractorCount += 1
		for i, section := range p.Extractor.Sections() {

			// Add the extractor plugin description only for first in the list
			if i == 0 {
				t.AppendSeparator()
				t.AppendRow(table.Row{p.Extractor.Description(), "", "", ""})
			}
			count += 1
			t.AppendRow([]interface{}{"", "extractor", p.Name, section})
		}
	}
	t.AppendSeparator()
	t.AppendFooter(table.Row{"Total", "", extractorCount, count})
	t.SetStyle(table.StyleColoredCyanWhiteOnBlack)
	t.Render()
	return nil
}
