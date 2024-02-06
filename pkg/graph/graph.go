package graph

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/converged-computing/jsongraph-go/jsongraph/v2/graph"
	jgf "github.com/converged-computing/jsongraph-go/jsongraph/v2/graph"
	"github.com/supercontainers/compspec-go/pkg/utils"

	"github.com/scylladb/go-set/strset"
)

const (
	rootLabel = "compspec-root"
)

/*

Desired steps:

1. Load the schemas into a JSON Graph (called JGF).
2. Do a depth first search until we find the first match
3. Return the first match (greedy)

*/

type ImageMapping map[string]map[string]bool

// addNs adds the namespace to a node or edge label
func addNs(namespace, label string) string {
	return fmt.Sprintf("%s.%s", namespace, label)
}

type CompatibilityGraph struct {
	*jgf.JsonGraph

	// Lookup of schemas we've added already
	// This helps us ensure we add each only once
	Schemas map[string]bool

	// This isn't ideal, but this will work for now. This is a lookup of
	// container images (strings) to node labels.
	// TODO update to use set?
	Images ImageMapping `json:"imageMapping"`

	// Node labels to container images
	NodeLabels ImageMapping `json:"nodeLabels"`
}

// HasNode determines if the graph has a node, named by label
func (c *CompatibilityGraph) HasNode(name string) bool {
	_, ok := c.Graph.Nodes[name]
	return ok
}

// AddAttribute adds an attribute to the graph, and a reference to
// an image or application at each level.
func (c *CompatibilityGraph) AddAttribute(
	uri string,
	schemaName string,
	key string,
	value string,
) error {

	// Create a reference for our container image / application if we don't have it yet
	// I'm not sure which I'll need yet
	images, ok := c.Images[uri]
	if !ok {

		// Flat lookup from container image URI to node labels
		// where the application container is present
		images = map[string]bool{}

		// If this is the first time we've seen the image, add the root node
		// This says that all images we map to the graph are at the root
		// We only add the schema name given a feature below it is added (running this function)
		images[rootLabel] = true
		images[schemaName] = true
	}

	// Keep going until we hit the root
	// This would parse like:
	// hardware.gpu.available
	// hardware.gpu
	// hardware <- last node, parent is the schema node
	notAtRoot := true
	graphKey := key
	for notAtRoot {

		// Add the image identifier to the node
		nsGraphKey := addNs(schemaName, graphKey)
		images[nsGraphKey] = true

		// When we hit a node that doesn't have a ., it's the root
		// Note we use the graphKey to skip over the "." in the schema name
		if strings.Contains(graphKey, ".") {

			// Chop off last .<group> to the right
			parts := strings.Split(graphKey, ".")
			graphKey = strings.Join(parts[0:len(parts)-1], ".")

		} else {

			// If we are already at the last one, finish up
			notAtRoot = false
		}

	}

	// Final addition - we need to assign the image to the value
	// The last nsGraphKey is the fullpath
	nsValue := addNs(schemaName, fmt.Sprintf("%s.%s", key, value))
	valueParent := addNs(schemaName, key)

	// Create a node for the attribute, parent is the key
	attrNode := newNode(nsValue)
	edge := getEdge(valueParent, nsValue, "contains")
	c.Graph.Edges = append(c.Graph.Edges, edge)
	c.Graph.Nodes[nsValue] = *attrNode
	images[nsValue] = true

	// Update images for the uri
	c.Images[uri] = images

	// Add final URI to node label lookup.
	// This allows us to quickly get images for a label at the end
	labels, ok := c.NodeLabels[nsValue]
	if !ok {
		labels = map[string]bool{}
	}
	labels[uri] = true
	c.NodeLabels[nsValue] = labels
	return nil
}

// PrintMapping is a simple print function to show images and nodes mapped to
func (c *CompatibilityGraph) PrintMapping() error {

	fmt.Println(" -- Mapping for Images")
	for image, nodes := range c.Images {
		fmt.Printf("  image: %s\n", image)

		// sort by key just so it is prettier
		nodes := sortNodes(nodes)

		fmt.Println("  nodes:")
		for _, node := range nodes {
			fmt.Printf("  -  %s\n", node)
		}
	}
	return nil
}

// Match finds matching nodes in the graph
// This is actually very simple (and dumb) and given the data structure,
// we don't need to traverse anything. We can:
// 1. Get an exact lookup for a feature of interest to a node.
// 2. If this node doesn't exist, we cannot match - that feature is missing
// 3. If it exists, keep the set of images
// 4. Continue to get sets of images for all desired features
// 5. The intersection across those are the matches!
func (c *CompatibilityGraph) Match(fields []string) ([]string, error) {

	// No fields, all are matches!
	if len(fields) == 0 {
		fmt.Println("No field criteria provided, all images are matches.")

		matches := []string{}
		for uri, _ := range c.Images {
			matches = append(matches, uri)
		}
		return matches, nil
	}

	// Faux set
	matches := strset.New()
	started := false

	for _, field := range fields {
		if !strings.Contains(field, "=") {
			fmt.Printf("Field request %s is missing '=', skipping", field)
		}
		parts := strings.SplitN(field, "=", 2)
		key := parts[0]
		value := parts[1]
		nodeId := fmt.Sprintf("%s.%s", key, value)

		// Do we have the node images
		uris, ok := c.NodeLabels[nodeId]
		fmt.Println(uris)
		if !ok {
			return []string{}, fmt.Errorf("Field %s is not known and cannot be matched.", field)
		}

		// Cut out early if we don't have matches
		if len(uris) == 0 {
			return []string{}, fmt.Errorf("Field %s does not have any associated images, match not possible.", field)
		}

		// If we aren't started, just take the first set as the solution
		if !started {
			for uri := range uris {
				matches.Add(uri)
			}
			started = true
		} else {

			// Otherwise, we need to take the intersection
			contenders := strset.New()
			for uri := range uris {
				contenders.Add(uri)
			}
			matches = strset.Intersection(matches, contenders)

			// Intersection is empty, no solution
			if matches.IsEmpty() {
				return []string{}, fmt.Errorf("Adding field %s empties match set, no possible", field)
			}
		}
	}
	return matches.List(), nil
}

// Load graph from a cached file.
// This assumes you know what you are doing, meaning
// the schemas are not changing
func (c *CompatibilityGraph) LoadGraph(path string) (bool, error) {
	exists, err := utils.PathExists(path)

	// Some error, do not continue!
	if err != nil {
		return exists, err
	}
	// No error, but doesn't exist (will still load)
	if !exists {
		return exists, nil
	}

	// If we get here, load in!
	fd, err := os.Open(path)
	if err != nil {
		return exists, err
	}
	graph := graph.Graph{}
	jsonParser := json.NewDecoder(fd)
	err = jsonParser.Decode(&graph)
	if err != nil {
		return exists, err
	}

	// Go through the edges and add the root nodes
	c.Graph = graph
	return true, nil
}

// Save graph to a cached file
func (c *CompatibilityGraph) SaveGraph(path string) error {
	exists, err := utils.PathExists(path)
	if err != nil {
		return err
	}
	// Don't overwrite if exists
	if exists {
		fmt.Printf("Graph %s already exists, will not overwrite\n", path)
		return nil
	}
	content, err := json.MarshalIndent(c.Graph, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("Saving graph to %s\n", path)
	err = os.WriteFile(path, content, 0644)
	if err != nil {
		return err
	}
	return nil
}

// AddSchema by url to the graph
func (c *CompatibilityGraph) AddSchema(url string) error {

	// Have we added this schema before?
	_, ok := c.Schemas[url]
	if ok {
		return nil
	}
	c.Schemas[url] = true

	// Serialize the schema into JGF (it is version 2)
	jgf := &graph.JsonGraph{}
	err := utils.GetJsonUrl(url, jgf)
	if err != nil {
		return err
	}

	// Merge the graph, meaning we add it to our JsonGraph
	// Likely we should have a merge function upstream
	// This is the schema root node e.g, io.archspec
	root := newNode(jgf.Graph.Id)
	c.Graph.Nodes[jgf.Graph.Id] = *root

	fmt.Printf("Schema %s is being added to the graph\n", jgf.Graph.Id)

	// That needs to be added to the actual root
	edge := getEdge(rootLabel, jgf.Graph.Id, "contains")
	c.Graph.Edges = append(c.Graph.Edges, edge)

	// Add each node from our schema graph
	// If it's a top level node (no . to indicate nested)
	// Add an edge to the top root
	for nodeId, _ := range jgf.Graph.Nodes {

		// The node needs to be namespaced by the schema
		nsName := addNs(jgf.Graph.Id, nodeId)

		// Create a new node. This strips metadata for now for a leaner graph
		// this can be changed.
		nsNode := newNode(nsName)

		// Add the ids metadata for image ids
		c.Graph.Nodes[nsName] = *nsNode

		// This indicates a top level attribute in the schema namespace
		// Note we are using the new namespaced name
		if !strings.Contains(nodeId, ".") {
			edge := getEdge(jgf.Graph.Id, nsName, "contains")
			c.Graph.Edges = append(c.Graph.Edges, edge)
		}
	}

	// Now add the remainder of edges (little subtrees!)
	// We again need to add the namespace of the schema to avoid conflicts
	for _, edge := range jgf.Graph.Edges {
		source := addNs(jgf.Graph.Id, edge.Source)
		target := addNs(jgf.Graph.Id, edge.Target)
		edge := getEdge(source, target, "contains")
		c.Graph.Edges = append(c.Graph.Edges, edge)
	}
	return nil
}

// Init a new FlexGraph from a graphml filename
func NewGraph() (CompatibilityGraph, error) {

	schemas := map[string]bool{}
	images := map[string]map[string]bool{}
	labels := map[string]map[string]bool{}

	// prepare a graph to load targets into
	g := jgf.NewGraph()

	// An empty root off of which we will have schemas
	root := newNode(rootLabel)
	g.Graph.Nodes[rootLabel] = *root

	// Return the compatibility graph wrapping it
	cg := CompatibilityGraph{g, schemas, images, labels}
	return cg, nil
}
