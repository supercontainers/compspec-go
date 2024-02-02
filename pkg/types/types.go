package types

import (
	"encoding/json"
	"strings"
)

// A compatibility spec is a
type CompatibilitySpec struct {
	Compatibilities map[string]CompatibilitySpec `json:"compatibilities"`
}
type CompatibiitySpec struct {
	Version    string     `json:"version"`
	Attributes Attributes `json:"attributes"`
}

type Attributes map[string]string

// A compatibility request is a mapping between a user preferences (some request to create a
// compatibility artifact) to a set of metadata attributes known by extractors.
// This is used for the create command

type CompatibilityRequest struct {
	Version         string                 `json:"version,omitempty"`
	Kind            string                 `json:"kind,omitempty"`
	Metadata        Metadata               `json:"metadata,omitempty"`
	Compatibilities []CompatibilityMapping `json:"compatibilities,omitempty"`
}
type Metadata struct {
	Name    string            `json:"name,omitempty"`
	Schemas map[string]string `json:"schemas,omitempty"`
}

// A compatibility mapping has one or more annotations that convert
// between extractor and compspec.json (the JsonSchema provided above)
type CompatibilityMapping struct {
	Name       string            `json:"name,omitempty"`
	Version    string            `json:"version,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// ToJson dumps our request to json for the artifact
func (r *CompatibilityRequest) ToJson() ([]byte, error) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return []byte{}, err
	}
	return b, err
}

// getExtractors parses the compatibility request for extractors needed
func (r *CompatibilityRequest) GetExtractors() []string {

	set := map[string]bool{}
	for _, compat := range r.Compatibilities {
		for _, request := range compat.Attributes {

			// The extractor name is the first field
			parts := strings.Split(request, ".")
			set[parts[0]] = true
		}
	}
	extractors := []string{}
	for name, _ := range set {
		extractors = append(extractors, name)
	}
	return extractors
}
