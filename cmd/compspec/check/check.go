package check

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/supercontainers/compspec-go/pkg/graph"
	"github.com/supercontainers/compspec-go/pkg/oras"
	"github.com/supercontainers/compspec-go/pkg/types"
	"sigs.k8s.io/yaml"
)

var (
	defaultMediaType = "application/org.supercontainers.compspec"
)

// loadManifest loads the manifest into a ManifestList
func loadManifest(filename string) (*types.ManifestList, error) {
	m := types.ManifestList{}
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return &m, err
	}

	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		return &m, err
	}
	return &m, nil
}

// Run will check a manifest list of artifacts against a host machine
// For now, the host machine parameters will be provided as flags
func Run(manifestFile string, hostFields []string, mediaType string) error {

	// Default media type if one not provided
	if mediaType == "" {
		mediaType = defaultMediaType
	}

	// Cut out early if a spec not provided
	if manifestFile == "" {
		return fmt.Errorf("A manifest file input -i/--input is required")
	}
	manifestList, err := loadManifest(manifestFile)
	if err != nil {
		return err
	}
	fmt.Println(manifestList)

	// Prepare a graph with our compspec schemas added
	g, err := graph.NewGraph()
	if err != nil {
		return err
	}

	// Load the compatibility specs into a lookup by image
	// This assumes we allow one image per compability spec, not sure
	// if there is a use case to have an image twice with two (sounds weird)
	lookup := map[string]types.CompatibilityRequest{}
	for _, item := range manifestList.Images {

		compspec, err := oras.LoadArtifact(item.Artifact, mediaType)
		if err != nil {
			fmt.Printf("warning, there was an issue loading the artifact for %s, skipping\n", item.Name)
		}
		lookup[item.Name] = compspec

		// Add schemas to the graph
		for _, schema := range compspec.Metadata.Schemas {
			fmt.Printf("Adding schema to graph %s\n", schema)
			err = g.AddSchema(schema)
			if err != nil {
				return err
			}
		}

		// When all schemas are added to a compatibility spec, we can walk graph to add metadata attributes
		// Each compspec has a list of compatibilities
		for _, compat := range compspec.Compatibilities {

			// If we don't have the root node, no go
			if !g.HasNode(compat.Name) {
				return fmt.Errorf("Schema root node %s is missing from the graph, missing from schemas.", compat.Name)
			}

			// Now add each attribute. Each attribute turns into a child node of the attribute
			// and if we are missing an attribute (meaning it isn't defined in the schema)
			// that is an error! item.Name is the container image to link to each node
			// in the traversal.
			for key, value := range compat.Attributes {
				err = g.AddAttribute(item.Name, compat.Name, key, value)
			}
		}
	}
	err = g.PrintMapping()
	if err != nil {
		return err
	}

	// Print the graph
	printme, _ := json.MarshalIndent(g.Graph, "", "\t")
	fmt.Println(string(printme))

	// TODO we will take this set of requests, load them into a graph,
	// and then query the graph based on user preferences (the host fields)
	// that are provided that describe the host we want to match compatibility with
	return nil
}
