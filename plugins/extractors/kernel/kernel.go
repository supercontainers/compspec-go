package kernel

import (
	"fmt"

	"github.com/compspec/compspec-go/pkg/plugin"
	"github.com/compspec/compspec-go/pkg/utils"
)

const (
	ExtractorName        = "kernel"
	ExtractorDescription = "generic kernel extractor"
	KernelBootSection    = "boot"
	KernelConfigSection  = "config"
	KernelModulesSection = "modules"
)

var (
	validSections = []string{KernelBootSection, KernelConfigSection, KernelModulesSection}
)

type KernelExtractor struct {
	sections []string
}

func (e KernelExtractor) Description() string {
	return ExtractorDescription
}

func (e KernelExtractor) Create(plugin.PluginOptions) error { return nil }
func (e KernelExtractor) IsCreator() bool                   { return false }
func (e KernelExtractor) IsExtractor() bool                 { return true }

func (e KernelExtractor) Sections() []string {
	return e.sections
}

func (c KernelExtractor) Name() string {
	return ExtractorName
}

// Validate ensures that the sections provided are in the list we know
// This is implemented on the level of the plugin, assuming each
// plugin might have custom logic to do this.
func (c KernelExtractor) Validate() bool {
	invalids, valid := utils.StringArrayIsSubset(c.sections, validSections)
	for _, invalid := range invalids {
		fmt.Printf("Sections %s is not known for extractor plugin %s\n", invalid, c.Name())
	}
	return valid
}

// Extract returns kernel metadata, for a set of named sections
// TODO eventually the user could select which sections they want
func (c KernelExtractor) Extract(allowFail bool) (plugin.PluginData, error) {

	sections := map[string]plugin.PluginSection{}
	data := plugin.PluginData{}

	// Only extract the sections we asked for
	for _, name := range c.sections {

		// Boot!
		if name == KernelBootSection {
			section, err := getKernelBootParams()
			if err != nil && !allowFail {
				return data, err
			}
			sections[KernelBootSection] = section
		}

		// Kernel full config file
		if name == KernelConfigSection {
			section, err := getKernelBootConfig()
			if err != nil && !allowFail {
				return data, err
			}
			sections[KernelConfigSection] = section
		}

		// Kernel full config file
		if name == KernelModulesSection {
			section, err := getKernelModules()
			if err != nil && !allowFail {
				return data, err
			}
			sections[KernelModulesSection] = section
		}
	}
	data.Sections = sections
	return data, nil
}

// NewPlugin validates and returns a new kernel plugins
func NewPlugin(sections []string) (plugin.PluginInterface, error) {
	if len(sections) == 0 {
		sections = validSections
	}
	e := KernelExtractor{sections: sections}
	if !e.Validate() {
		return nil, fmt.Errorf("plugin %s is not valid", e.Name())
	}
	return e, nil
}
