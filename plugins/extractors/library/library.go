package library

import (
	"fmt"

	"github.com/compspec/compspec-go/pkg/plugin"
	"github.com/compspec/compspec-go/pkg/utils"
)

const (
	ExtractorName        = "library"
	ExtractorDescription = "generic library extractor"
	MPISection           = "mpi"
)

var (
	validSections = []string{MPISection}
)

type LibraryExtractor struct {
	sections []string
}

func (e LibraryExtractor) Name() string {
	return ExtractorName
}

func (e LibraryExtractor) Sections() []string {
	return e.sections
}

func (e LibraryExtractor) Description() string {
	return ExtractorDescription
}

func (e LibraryExtractor) Create(options map[string]string) error {
	return nil
}

func (e LibraryExtractor) IsCreator() bool   { return false }
func (e LibraryExtractor) IsExtractor() bool { return true }

// Validate ensures that the sections provided are in the list we know
func (e LibraryExtractor) Validate() bool {
	invalids, valid := utils.StringArrayIsSubset(e.sections, validSections)
	for _, invalid := range invalids {
		fmt.Printf("Sections %s is not known for extractor plugin %s\n", invalid, e.Name())
	}
	return valid
}

// Extract returns library metadata, for a set of named sections
func (e LibraryExtractor) Extract(interface{}) (plugin.PluginData, error) {

	sections := map[string]plugin.PluginSection{}
	data := plugin.PluginData{}

	// Only extract the sections we asked for
	for _, name := range e.sections {
		if name == MPISection {
			section, err := getMPIInformation()
			if err != nil {
				return data, err
			}
			sections[MPISection] = section
		}
	}
	data.Sections = sections
	return data, nil
}

// NewPlugin validates and returns a new plugin
func NewPlugin(sections []string) (plugin.PluginInterface, error) {
	if len(sections) == 0 {
		sections = validSections
	}
	e := LibraryExtractor{sections: sections}
	if !e.Validate() {
		return nil, fmt.Errorf("plugin %s is not valid", e.Name())
	}
	return e, nil
}
