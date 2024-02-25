package plugin

import (
	"encoding/json"
	"fmt"
)

// A Plugin interface can define any of the following:
//
//	an Extract function to return extractor data across sections
//	a validate function to typically check that the plugin is valid
//	a Creation interface that can use extractor data to generate something new
type PluginInterface interface {
	Name() string
	Description() string

	// This is probably a dumb way to do it, but it works
	IsExtractor() bool
	IsCreator() bool

	// Extractors
	Extract(interface{}) (PluginData, error)
	Validate() bool
	Sections() []string

	// Creators take a map of named options
	Create(map[string]string) error
}

// ExtractorData is returned by an extractor
type PluginData struct {
	Sections Sections `json:"sections,omitempty"`
}
type Sections map[string]PluginSection

// Print extractor data to the console
func (e *PluginData) Print() {
	for name, section := range e.Sections {
		fmt.Printf(" -- Section %s\n", name)
		for key, value := range section {
			fmt.Printf("   %s: %s\n", key, value)
		}
	}

}

// ToJson serializes to json
func (e *PluginData) ToJson() (string, error) {
	b, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), err
}

// An extractor section corresponds to a named group of attributes
type PluginSection map[string]string

// Extractors is a lookup of registered extractors by name
type Plugins map[string]PluginInterface
