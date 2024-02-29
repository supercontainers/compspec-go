package system

import (
	"fmt"

	"github.com/compspec/compspec-go/pkg/plugin"
	"github.com/compspec/compspec-go/pkg/utils"
)

const (
	ExtractorName        = "system"
	ExtractorDescription = "generic system extractor"

	// Just cores, etc.
	CPUSection       = "cpu"
	ProcessorSection = "processor"
	ArchSection      = "arch"
	OsSection        = "os"
	MemorySection    = "memory"
	HWLOC            = "hwloc"
)

var (
	validSections = []string{ProcessorSection, OsSection, ArchSection, MemorySection, CPUSection, HWLOC}
)

type SystemExtractor struct {
	sections []string
}

func (e SystemExtractor) Name() string {
	return ExtractorName
}

func (e SystemExtractor) Description() string {
	return ExtractorDescription
}

func (e SystemExtractor) Sections() []string {
	return e.sections
}

func (e SystemExtractor) Create(plugin.PluginOptions) error { return nil }
func (e SystemExtractor) IsCreator() bool                   { return false }
func (e SystemExtractor) IsExtractor() bool                 { return true }

// Validate ensures that the sections provided are in the list we know
func (e SystemExtractor) Validate() bool {
	invalids, valid := utils.StringArrayIsSubset(e.sections, validSections)
	for _, invalid := range invalids {
		fmt.Printf("Sections %s is not known for extractor plugin %s\n", invalid, e.Name())
	}
	return valid
}

// Extract returns system metadata, for a set of named sections
func (e SystemExtractor) Extract(allowFail bool) (plugin.PluginData, error) {

	sections := map[string]plugin.PluginSection{}
	data := plugin.PluginData{}

	// Only extract the sections we asked for
	for _, name := range e.sections {
		if name == ProcessorSection {
			section, err := getProcessorInformation()
			if err != nil && !allowFail {
				return data, err
			}
			sections[ProcessorSection] = section
		}
		if name == OsSection {
			section, err := getOsInformation()
			if err != nil && !allowFail {
				return data, err
			}
			sections[OsSection] = section
		}
		if name == HWLOC {
			section, err := getHWLOC()
			if err != nil {
				return data, err
			}
			sections[HWLOC] = section
		}

		if name == CPUSection {
			section, err := getCPUInformation()
			if err != nil && !allowFail {
				return data, err
			}
			sections[CPUSection] = section
		}
		if name == ArchSection {
			section, err := getArchInformation()
			if err != nil && !allowFail {
				return data, err
			}
			sections[ArchSection] = section
		}

		if name == MemorySection {
			section, err := getMemoryInformation()
			if err != nil && !allowFail {
				return data, err
			}
			sections[MemorySection] = section
		}

	}
	data.Sections = sections
	return data, nil
}

// NewPlugin validates and returns a new kernel plugin
func NewPlugin(sections []string) (plugin.PluginInterface, error) {
	if len(sections) == 0 {
		sections = validSections
	}
	e := SystemExtractor{sections: sections}
	if !e.Validate() {
		return nil, fmt.Errorf("plugin %s is not valid", e.Name())
	}
	return e, nil
}
