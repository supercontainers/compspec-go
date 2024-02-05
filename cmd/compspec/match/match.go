package match

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/supercontainers/compspec-go/pkg/graph"
	"github.com/supercontainers/compspec-go/pkg/oras"
	"github.com/supercontainers/compspec-go/pkg/types"
	"github.com/supercontainers/compspec-go/pkg/utils"
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
func Run(
	manifestFile string,
	hostFields []string,
	mediaType string,
	cachePath string,
	printMapping bool,
	printGraph bool,
	allowFail bool,
	checkArtifacts bool,
	randomize bool,
	single bool,
) error {

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

	// Prepare a graph with our compspec schemas added
	g, err := graph.NewGraph()
	if err != nil {
		return err
	}

	// Keep check of artifacts missing
	missing := 0

	// Load the compatibility specs into a lookup by image
	// This assumes we allow one image per compability spec, not sure
	// if there is a use case to have an image twice with two (sounds weird)
	lookup := map[string]types.CompatibilityRequest{}
	for _, item := range manifestList.Images {

		compspec, err := oras.LoadArtifact(item.Artifact, mediaType, cachePath)
		if err != nil {
			fmt.Printf("warning, there was an issue loading the artifact for %s: %s, skipping\n", item.Name, err)

			// Don't allow any artifacts to be missing, unless we are just checking
			if !allowFail && !checkArtifacts {
				return err
			}
			missing += 1
		}
		// If just checking artifacts, continue
		if checkArtifacts {
			continue
		}
		lookup[item.Name] = compspec

		// Add schemas to the graph
		for _, schema := range compspec.Metadata.Schemas {
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

	if checkArtifacts {
		fmt.Printf("Checking artifacts complete. There were %d artifacts missing.\n", missing)
		return nil
	}

	// We only want to print the mapping and exit
	if printMapping {
		err = g.PrintMapping()
		if err != nil {
			return err
		}
		return nil
	}

	// Print the graph
	if printGraph {
		toprint, _ := json.MarshalIndent(g.Graph, "", "\t")
		fmt.Println(string(toprint))
		return nil
	}

	// Perform the match to the desired host
	matches, err := g.Match(hostFields)
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		fmt.Println("There was no match. Try changig your constaints.")
		return nil
	}

	// Do we want to randomize (randomly shuffle the list)?
	if randomize {
		matches = utils.RandomSort(matches)
	}

	// Do we want just one match?
	if single && len(matches) > 0 {
		matches = []string{matches[0]}
	}

	fmt.Println(" --- Found matches ---")
	for _, match := range matches {
		fmt.Println(match)
	}

	return nil
}
