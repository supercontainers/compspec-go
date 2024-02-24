package create

import (
	"fmt"
	"os"

	"github.com/compspec/compspec-go/pkg/types"
	ep "github.com/compspec/compspec-go/plugins/extractors"

	p "github.com/compspec/compspec-go/plugins"
)

// Artifact will create a compatibility artifact based on a request in YAML
// TODO likely want to refactor this into a proper create plugin
func Artifact(specname string, fields []string, saveto string, allowFail bool) error {

	// Cut out early if a spec not provided
	if specname == "" {
		return fmt.Errorf("a spec input -i/--input is required")
	}
	request, err := loadRequest(specname)
	if err != nil {
		return err
	}

	// Right now we only know about extractors, when we define subfields
	// we can further filter here.
	extractors := request.GetExtractors()
	plugins, err := ep.GetPlugins(extractors)
	if err != nil {
		return err
	}

	// Finally, add custom fields and extract metadata
	result, err := plugins.Extract(allowFail)
	if err != nil {
		return err
	}

	// Update with custom fields (either new or overwrite)
	result.AddCustomFields(fields)

	// The compspec returned is the populated Compatibility request!
	compspec, err := PopulateExtractors(&result, request)
	if err != nil {
		return err
	}

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
func PopulateExtractors(result *ep.Result, request *types.CompatibilityRequest) (*types.CompatibilityRequest, error) {

	// Every metadata attribute must be known under a schema
	schemas := request.Metadata.Schemas
	if len(schemas) == 0 {
		return nil, fmt.Errorf("the request must have one or more schemas")
	}
	for i, compat := range request.Compatibilities {

		// The compatibility section name is a schema, and must be defined
		url, ok := schemas[compat.Name]
		if !ok {
			return nil, fmt.Errorf("%s is missing a schema", compat.Name)
		}
		if url == "" {
			return nil, fmt.Errorf("%s has an empty schema", compat.Name)
		}

		for key, extractorKey := range compat.Attributes {

			// Get the extractor, section, and subfield from the extractor lookup key
			f, err := p.ParseField(extractorKey)
			if err != nil {
				fmt.Printf("warning: cannot parse %s: %s, setting to empty\n", key, extractorKey)
				compat.Attributes[key] = ""
				continue
			}

			// If we get here, we can parse it and look it up in our result metadata
			extractor, ok := result.Results[f.Extractor]
			if !ok {
				fmt.Printf("warning: extractor %s is unknown, setting to empty\n", f.Extractor)
				compat.Attributes[key] = ""
				continue
			}

			// Now get the section
			section, ok := extractor.Sections[f.Section]
			if !ok {
				fmt.Printf("warning: section %s.%s is unknown, setting to empty\n", f.Extractor, f.Section)
				compat.Attributes[key] = ""
				continue
			}

			// Now get the value!
			value, ok := section[f.Field]
			if !ok {
				fmt.Printf("warning: field %s.%s.%s is unknown, setting to empty\n", f.Extractor, f.Section, f.Field)
				compat.Attributes[key] = ""
				continue
			}

			// If we get here - we found it! Hooray!
			compat.Attributes[key] = value
		}

		// Update the compatibiity
		request.Compatibilities[i] = compat
	}

	return request, nil
}