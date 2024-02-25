package graph

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/compspec/compspec-go/pkg/utils"
	"github.com/converged-computing/jsongraph-go/jsongraph/metadata"
	"github.com/converged-computing/jsongraph-go/jsongraph/v2/graph"
	jgf "github.com/converged-computing/jsongraph-go/jsongraph/v2/graph"
)

// A ClusterGraph is meant to be a plain (flux oriented) JGF to describe a cluster (nodes)
type ClusterGraph struct {
	*jgf.JsonGraph

	Name string

	// Top level counter for node labels (JGF v2) that maps to ids (JGF v1)
	nodeCounter int32

	// Counters for specific resource types (e.g., rack, node)
	resourceCounters map[string]int32
}

// HasNode determines if the graph has a node, named by label
func (c *ClusterGraph) HasNode(name string) bool {
	_, ok := c.Graph.Nodes[name]
	return ok
}

// Save graph to a cached file
func (c *ClusterGraph) SaveGraph(path string) error {
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

// Path gets a new path
func getNodePath(root, subpath string) string {
	if subpath == "" {
		return fmt.Sprintf("/%s", root)
	}
	return fmt.Sprintf("/%s/%s", root, subpath)
}

// AddNode adds a node to the graph
// g.AddNode("rack", 1, false, "", root)
func (c *ClusterGraph) AddNode(
	resource string,
	name string,
	size int32,
	exclusive bool,
	unit string,
) *graph.Node {
	node := c.getNode(resource, name, size, exclusive, unit)
	c.Graph.Nodes[*node.Label] = *node
	return node
}

// Add an edge from source to dest with some relationship
func (c *ClusterGraph) AddEdge(source, dest graph.Node, relation string) {
	edge := getEdge(*source.Label, *dest.Label, relation)
	c.Graph.Edges = append(c.Graph.Edges, edge)
}

// getNode is a private shared function that can also be used to generate the root!
func (c *ClusterGraph) getNode(
	resource string,
	name string,
	size int32,
	exclusive bool,
	unit string,
) *graph.Node {

	// Get the identifier for the resource type
	counter, ok := c.resourceCounters[resource]
	if !ok {
		counter = 0
	}

	// The current count in the graph (global)
	count := c.nodeCounter

	// The id in the metadata is the counter for that resource type
	resourceCounter := fmt.Sprintf("%d", counter)

	// The resource name is the type + the resource counter
	resourceName := fmt.Sprintf("%s%d", name, counter)

	// New Metadata with expected fluxion data
	m := metadata.Metadata{}
	m.AddElement("type", resource)
	m.AddElement("basename", name)
	m.AddElement("id", resourceCounter)
	m.AddElement("name", resourceName)

	// uniq_id should be the same as the label, but as an integer
	m.AddElement("uniq_id", count)
	m.AddElement("rank", -1)
	m.AddElement("exclusive", exclusive)
	m.AddElement("unit", unit)
	m.AddElement("size", size)
	m.AddElement("paths", map[string]string{"containment": getNodePath(name, "")})

	// Update the resource counter
	counter += 1
	c.resourceCounters[resource] = counter

	// Update the global counter
	c.nodeCounter += 1

	// Assemble the node!
	// Label for v2 will be identifier "id" for JGF v1
	label := fmt.Sprintf("%d", count)
	node := graph.Node{Label: &label, Metadata: m}
	return &node
}

// Init a new FlexGraph from a graphml filename
// The cluster root is slightly different so we don't use getNode here
func NewClusterGraph(name string) (ClusterGraph, error) {

	// prepare a graph to load targets into
	g := jgf.NewGraph()

	clusterName := fmt.Sprintf("%s0", name)

	// New Metadata with expected fluxion data
	m := metadata.Metadata{}
	m.AddElement("type", "cluster")
	m.AddElement("basename", name)
	m.AddElement("name", clusterName)
	m.AddElement("id", 0)
	m.AddElement("uniq_id", 0)
	m.AddElement("rank", -1)
	m.AddElement("exclusive", false)
	m.AddElement("unit", "")
	m.AddElement("size", 1)
	m.AddElement("paths", map[string]string{"containment": getNodePath(name, "")})

	// Root cluster node
	label := "0"
	node := graph.Node{Label: &label, Metadata: m}

	// Set the root node
	g.Graph.Nodes[label] = node

	// Create a new cluster!
	// Start counting at 1 - index 0 is the cluster root
	resourceCounters := map[string]int32{"cluster": int32(1)}
	cluster := ClusterGraph{g, name, 1, resourceCounters}

	return cluster, nil
}
