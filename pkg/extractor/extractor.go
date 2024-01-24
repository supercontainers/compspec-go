package extractor

import "fmt"

// An Extractor interface has a single Extract function
// to return extractor data.
type Extractor interface {
	Name() string
	Extract(interface{}) (ExtractorData, error)
}

// ExtractorData is returned by an extractor
type ExtractorData struct {
	Sections map[string]ExtractorSection
}

// Dump extracted data to json
func (e *ExtractorData) Print() {
	for name, section := range e.Sections {
		fmt.Printf(" -- Section %s\n", name)
		for key, value := range section {
			fmt.Printf("   %s: %s\n", key, value)
		}
	}

}

// An extractor section corresponds to a named group of annotations
type ExtractorSection map[string]string

// Extractors is a lookup of registered extractors by name
type Extractors map[string]Extractor
