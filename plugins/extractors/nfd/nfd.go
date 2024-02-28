package nfd

import (
	"fmt"

	source "github.com/converged-computing/nfd-source/source"

	// Note that "fake" is removed from here
	_ "github.com/converged-computing/nfd-source/source/cpu"
	_ "github.com/converged-computing/nfd-source/source/kernel"
	_ "github.com/converged-computing/nfd-source/source/local"
	_ "github.com/converged-computing/nfd-source/source/memory"
	_ "github.com/converged-computing/nfd-source/source/network"
	_ "github.com/converged-computing/nfd-source/source/pci"
	_ "github.com/converged-computing/nfd-source/source/storage"
	_ "github.com/converged-computing/nfd-source/source/system"
	_ "github.com/converged-computing/nfd-source/source/usb"

	"github.com/compspec/compspec-go/pkg/plugin"
	"github.com/compspec/compspec-go/pkg/utils"
)

const (
	ExtractorName        = "nfd"
	ExtractorDescription = "node feature discovery"

	// Each of these corresponds to a source
	CPUSection = "cpu"

	// TODO can we do a check that this is desired / enabled before running?
	KernelSection  = "kernel"
	LocalSection   = "local"
	MemorySection  = "memory"
	NetworkSection = "network"
	PCISection     = "pci"
	StorageSection = "storage"
	SystemSection  = "system"
	USBSection     = "usb"
)

var (
	validSections = []string{
		CPUSection,
		KernelSection,
		LocalSection,
		MemorySection,
		NetworkSection,
		PCISection,
		StorageSection,
		SystemSection,
		USBSection,
	}
)

// NFDExtractor is an extractor for node feature discovery
type NFDExtractor struct {
	sections []string
}

func (e NFDExtractor) Name() string {
	return ExtractorName
}

func (e NFDExtractor) Sections() []string {
	return e.sections
}

func (e NFDExtractor) Description() string {
	return ExtractorDescription
}

func (e NFDExtractor) Create(plugin.PluginOptions) error { return nil }
func (e NFDExtractor) IsCreator() bool                   { return false }
func (e NFDExtractor) IsExtractor() bool                 { return true }

// Validate ensures that the sections provided are in the list we know
func (e NFDExtractor) Validate() bool {
	invalids, valid := utils.StringArrayIsSubset(e.sections, validSections)
	for _, invalid := range invalids {
		fmt.Printf("Sections %s is not known for extractor plugin %s\n", invalid, e.Name())
	}
	return valid
}

// Extract returns system metadata, for a set of named sections
func (e NFDExtractor) Extract(allowFail bool) (plugin.PluginData, error) {

	sections := map[string]plugin.PluginSection{}
	data := plugin.PluginData{}

	// Get all registered feature sources
	sources := source.GetAllFeatureSources()

	// Only extract the sections we asked for
	for _, name := range e.sections {
		discovery, ok := sources[name]

		// This should not happen
		if !ok {
			fmt.Printf("%s is not a known feature source\n", name)
			continue
		}
		err := discovery.Discover()
		if err != nil {
			fmt.Printf("Issue discovering features for %s\n", discovery.Name())
			continue
		}

		// Create a new section for the <name> group
		// For each of the below, "fs" is a feature set
		// AttributeFeatureSet
		section := plugin.PluginSection{}
		features := discovery.GetFeatures()
		for k, fs := range features.Attributes {
			for fName, feature := range fs.Elements {
				uid := fmt.Sprintf("%s.%s", k, fName)
				section[uid] = feature
			}
		}

		// FlagFeatureSet
		// Note that the second value to feature is v1alpha.Nil
		// I think this acts as a flag, double check
		for k, fs := range features.Flags {
			for feature, _ := range fs.Elements {
				uid := fmt.Sprintf("%s.%s", k, feature)
				section[uid] = "true"
			}
		}

		// InstanceFeatureSet
		for k, fs := range features.Instances {
			for idx, feature := range fs.Elements {
				for fName, attr := range feature.Attributes {
					uid := fmt.Sprintf("%s.%d.%s", k, idx, fName)
					section[uid] = attr
				}
			}
		}
		sections[name] = section
	}
	data.Sections = sections
	return data, nil
}

// NewPlugin validates and returns a new kernel plugin
func NewPlugin(sections []string) (plugin.PluginInterface, error) {
	if len(sections) == 0 {
		sections = validSections
	}
	e := NFDExtractor{sections: sections}
	if !e.Validate() {
		return nil, fmt.Errorf("plugin %s is not valid", e.Name())
	}
	return e, nil
}
