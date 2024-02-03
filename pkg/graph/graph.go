package graph

import (
	"fmt"
	"strings"

	"github.com/converged-computing/jsongraph-go/jsongraph/v2/graph"
	jgf "github.com/converged-computing/jsongraph-go/jsongraph/v2/graph"
	"github.com/supercontainers/compspec-go/pkg/utils"
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

	// I did this strategy because the metadata elements of the JGF are a list
	// This isn't ideal, but this will work for now. This is a lookup of node label
	// to container images (strings). This might be better to use integers for
	// larger graphs
	Images ImageMapping `json:"imageMapping"`
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

	// prepare a graph to load targets into
	g := jgf.NewGraph()

	// An empty root off of which we will have schemas
	root := newNode(rootLabel)
	g.Graph.Nodes[rootLabel] = *root

	// Return the compatibility graph wrapping it
	cg := CompatibilityGraph{g, schemas, images}
	return cg, nil
}
