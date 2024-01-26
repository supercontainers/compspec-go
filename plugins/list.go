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

	// TODO add description column
	for _, p := range *r {
		for _, section := range p.Extractor.Sections() {
			count += 1
			t.AppendRow([]interface{}{"", "extractor", p.Name, section})
		}
	}
	t.AppendSeparator()
	t.AppendFooter(table.Row{"Total", count, "", ""})
	t.SetStyle(table.StyleColoredCyanWhiteOnBlack)
	t.Render()
	return nil
}
