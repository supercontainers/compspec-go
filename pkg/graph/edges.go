package graph

import (
	"github.com/converged-computing/jsongraph-go/jsongraph/v2/graph"
)

// Get an edge with a specific relationship (typically "contains" or "in")
func getEdge(source string, dest string, relation string) graph.Edge {
	return graph.Edge{Source: source, Target: dest, Relation: relation}
}
