package types

// A compatibility spec is a
type CompatibilitySpec struct {
	Compatibilities map[string]CompatibilitySpec `json:"compatibilities"`
}
type CompatibiitySpec struct {
	Version     string      `json:"version"`
	Annotations Annotations `json:"annotations"`
}

type Annotations map[string]string
