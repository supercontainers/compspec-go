package create

import (
	"os"

	"github.com/compspec/compspec-go/pkg/types"
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
