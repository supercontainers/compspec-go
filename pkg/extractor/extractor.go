package extractor

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

// An extractor section corresponds to a named group of annotations
type ExtractorSection map[string]string

// Extractors is a lookup of registered extractors by name
type Extractors map[string]Extractor
