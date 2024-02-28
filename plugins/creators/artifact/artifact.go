package artifact

import (
	"fmt"
	"os"

	"github.com/compspec/compspec-go/pkg/plugin"
	"github.com/compspec/compspec-go/pkg/types"
	p "github.com/compspec/compspec-go/plugins"
	"sigs.k8s.io/yaml"
)

const (
	CreatorName        = "artifact"
	CreatorDescription = "describe an application or environment"
)

type ArtifactCreator struct{}

func (c ArtifactCreator) Description() string {
	return CreatorDescription
}

func (c ArtifactCreator) Name() string {
	return CreatorName
}

func (c ArtifactCreator) Sections() []string {
	return []string{}
}

func (c ArtifactCreator) Extract(allowFail bool) (plugin.PluginData, error) {
	return plugin.PluginData{}, nil
}
func (c ArtifactCreator) IsCreator() bool   { return true }
func (c ArtifactCreator) IsExtractor() bool { return false }

// Create generates the desired output
func (c ArtifactCreator) Create(options plugin.PluginOptions) error {

	// unwrap options (we can be sure they are at least provided)
	specname := options.StrOpts["specname"]
	saveto := options.StrOpts["saveto"]
	fields := options.ListOpts["fields"]

	// This is uber janky. We could use interfaces
	// But I just feel so lazy right now
	allowFail := options.BoolOpts["allowFail"]

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
	plugins, err := p.GetPlugins(extractors)
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

// LoadExtractors loads a compatibility result into a compatibility request
// After this we can save the populated thing into an artifact (json DUMP)
func PopulateExtractors(result *plugin.Result, request *types.CompatibilityRequest) (*types.CompatibilityRequest, error) {

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
			f, err := plugin.ParseField(extractorKey)
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

func (c ArtifactCreator) Validate() bool {
	return true
}

// NewPlugin creates a new ArtifactCreator
func NewPlugin() (plugin.PluginInterface, error) {
	c := ArtifactCreator{}
	return c, nil
}
