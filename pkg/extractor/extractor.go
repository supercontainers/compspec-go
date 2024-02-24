package extractor

import (
	"encoding/json"
	"fmt"
)

// An Extractor interface has:
//
//	an Extract function to return extractor data across sections
//	a validate function to typically check that the plugin is valid
type Extractor interface {
	Name() string
	Description() string
	Extract(interface{}) (ExtractorData, error)
	Validate() bool
	Sections() []string
	// GetSection(string) ExtractorData
}

// ExtractorData is returned by an extractor
type ExtractorData struct {
	Sections Sections `json:"sections,omitempty"`
}
type Sections map[string]ExtractorSection

// Print extractor data to the console
func (e *ExtractorData) Print() {
	for name, section := range e.Sections {
		fmt.Printf(" -- Section %s\n", name)
		for key, value := range section {
			fmt.Printf("   %s: %s\n", key, value)
		}
	}

}

// ToJson serializes to json
func (e *ExtractorData) ToJson() (string, error) {
	b, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), err
}

// An extractor section corresponds to a named group of attributes
type ExtractorSection map[string]string

// Extractors is a lookup of registered extractors by name
type Extractors map[string]Extractor
