package create

import (
	"fmt"
	"os"

	"github.com/supercontainers/compspec-go/pkg/types"
	p "github.com/supercontainers/compspec-go/plugins"
	"sigs.k8s.io/yaml"
)

// loadRequest loads a Compatibility Request YAML into a struct
func loadRequest(filename string) (*types.CompatibilityRequest, error) {
	request := types.CompatibilityRequest{}
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return &request, err
	}

	err = yaml.Unmarshal(yamlFile, &request)
	if err != nil {
		return &request, err
	}
	return &request, nil
}

// Run will create a compatibility artifact based on a request in YAML
func Run(specname string, fields []string, saveto string) error {

	// Cut out early if a spec not provided
	if specname == "" {
		return fmt.Errorf("A spec input -i/--input is required")
	}
	request, err := loadRequest(specname)
	if err != nil {
		return err
	}

	// Right now we only know about extractors, when we define subfields
	// we can further filter here.
	extractors := request.GetExtractors()
	plugins, err := p.GetPlugins(extractors)
	if err != nil {
		return err
	}

	// Finally, add custom fields and extract metadata
	result, err := plugins.Extract()

	// Update with custom fields (either new or overwrite)
	result.AddCustomFields(fields)

	// The compspec returned is the populated Compatibility request!
	compspec, err := PopulateExtractors(&result, request)
	output, err := compspec.ToJson()
	if err != nil {
		return err
	}
	if saveto == "" {
		fmt.Println(string(output))
	} else {
		err = os.WriteFile(saveto, output, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadExtractors loads a compatibility result into a compatibility request
// After this we can save the populated thing into an artifact (json DUMP)
func PopulateExtractors(result *p.Result, request *types.CompatibilityRequest) (*types.CompatibilityRequest, error) {

	for i, compat := range request.Compatibilities {
		for key, extractorKey := range compat.Annotations {

			// Get the extractor, section, and subfield from the extractor lookup key
			f, err := p.ParseField(extractorKey)
			if err != nil {
				fmt.Printf("warning: cannot parse %s: %s, setting to empty\n", key, extractorKey)
				compat.Annotations[key] = ""
				continue
			}

			// If we get here, we can parse it and look it up in our result metadata
			extractor, ok := result.Results[f.Extractor]
			if !ok {
				fmt.Printf("warning: extractor %s is unknown, setting to empty\n", f.Extractor)
				compat.Annotations[key] = ""
				continue
			}

			// Now get the section
			section, ok := extractor.Sections[f.Section]
			if !ok {
				fmt.Printf("warning: section %s.%s is unknown, setting to empty\n", f.Extractor, f.Section)
				compat.Annotations[key] = ""
				continue
			}

			// Now get the value!
			value, ok := section[f.Field]
			if !ok {
				fmt.Printf("warning: field %s.%s.%s is unknown, setting to empty\n", f.Extractor, f.Section, f.Field)
				compat.Annotations[key] = ""
				continue
			}

			// If we get here - we found it! Hooray!
			compat.Annotations[key] = value
		}

		// Update the compatibiity
		request.Compatibilities[i] = compat
	}

	return request, nil
}