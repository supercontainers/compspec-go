package graph

import (
	"sort"

	"github.com/converged-computing/jsongraph-go/jsongraph/metadata"
	"github.com/converged-computing/jsongraph-go/jsongraph/v2/graph"
)

// newNode generates a new node with a lookup for image ids
// This is unecessary, but added if this process is eventually more complex
func newNode(label string) *graph.Node {

	m := metadata.Metadata{}

	// We will use this to map container image ids to each node
	node := graph.Node{Label: &label, Metadata: m}
	return &node

}

// sortNodes so they are pretty!
func sortNodes(nodes map[string]bool) []string {

	keys := make([]string, 0, len(nodes))
	for node := range nodes {
		keys = append(keys, node)
	}
	sort.Strings(keys)
	return keys
}
