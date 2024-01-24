package kernel

import (
	"github.com/supercontainers/compspec-go/pkg/extractor"
)

const (
	KernelExtractorName  = "KernelExtractor"
	KernelBootSection    = "kernel.boot"
	KernelConfigSection  = "kernel.config"
	KernelModulesSection = "kernel.config"
)

type KernelExtractor struct{}

func (c KernelExtractor) Name() string {
	return KernelExtractorName
}

// Extract returns kernel metadata
// TODO eventually the user could select which sections they want
func (c KernelExtractor) Extract(interface{}) (extractor.ExtractorData, error) {
	sections := map[string]extractor.ExtractorSection{}
	data := extractor.ExtractorData{}

	// Add kernel boot parameters
	section, err := getKernelBootParams()
	if err != nil {
		return data, err
	}
	sections[KernelBootSection] = section

	// Add kernel config parameters
	section, err = getKernelBootConfig()
	if err != nil {
		return data, err
	}
	sections[KernelConfigSection] = section

	// Add kernel modules section
	section, err = getKernelModules()
	if err != nil {
		return data, err
	}
	sections[KernelModulesSection] = section
	data.Sections = sections
	return data, nil
}
