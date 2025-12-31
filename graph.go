package graph

import (
	"iter"
	"maps"
	"slices"
)

// nodes can be identified by anything that can be used as a key in a map
type Node[K comparable] struct {
	ID K
}

// edges are identified by the nodes they connect,
// and the weight of the connection
type Edge[K comparable] struct {
	u, v   Node[K]
	weight float64
}

// an adjacency is defined by the other end point and
// the edge weight
type Adjancency[K comparable] struct {
	v      Node[K]
	weight float64
}

// define an interface for an abstract graph
type Graph[K comparable] interface {
	AddNode(n Node[K])
	AddEdge(u, v Node[K], w float64)
	AddNodesFrom(ns []Node[K])
	AddEdgesFrom(es []Edge[K])
	RemoveNode(n Node[K])
	RemoveEdge(u, v Node[K])
	RemoveNodesFrom(ns []Node[K])
	RemoveEdgesFrom(es []Edge[K])
	Nodes() iter.Seq[Node[K]]
	Edges() iter.Seq[Edge[K]]
	Clear()
	NumberOfNodes() int
	NumberOfEdges() int
	Successors(n Node[K]) iter.Seq[Node[K]]
	Predecessors(n Node[K]) iter.Seq[Node[K]]
	Neighbors(n Node[K]) iter.Seq[Node[K]]
	InDegree(n Node[K]) int
	OutDegree(n Node[K]) int
	Degree(n Node[K]) int
}

// generic data structure for a graph. it's a simple lookup
// table for graphs and list of graphs with the weight associated
// with the edge between the two keys
type graphData[K comparable] struct {
	Adjacencies map[Node[K]]map[Node[K]]float64
}

// function to add a node to the graph
func (g *graphData[K]) AddNode(n Node[K]) {
	// does the node already exist in the graph?
	if _, ok := g.Adjacencies[n]; !ok {
		// no, add it with no adjacencies
		g.Adjacencies[n] = make(map[Node[K]]float64)
	}
}

// functions to add nodes to the graph from some iter
func (g *graphData[K]) AddNodesFrom(ns []Node[K]) {
	for _, n := range ns {
		g.AddNode(n)
	}
}

// function to remove a node from the graph
func (g *graphData[K]) RemoveNode(n Node[K]) {
	// remove all adjancencies to the node
	for node := range g.Adjacencies {
		delete(g.Adjacencies[node], n)
	}
	// remove adjacencies from the node, and with that its record
	delete(g.Adjacencies, n)
}

// function to remove ndoes from the graph sourced from some iter
func (g *graphData[K]) RemoveNodesFrom(ns []Node[K]) {
	for _, n := range ns {
		g.RemoveNode(n)
	}
}

// function to retrieve an iterator over the nodes of the graph
func (g *graphData[K]) Nodes() iter.Seq[Node[K]] {
	return maps.Keys(g.Adjacencies)
}

// function to retrieve a list of edges from a graph
func (g *graphData[K]) Edges() iter.Seq[Edge[K]] {
	// create the iterator
	return func(yield func(Edge[K]) bool) {
		// walk the nodes
		for u := range g.Adjacencies {
			// walk the node's adjacencies
			for v, w := range g.Adjacencies[u] {
				// create the edge
				edge := Edge[K]{u: u, v: v, weight: w}
				// and yield it
				if !yield(edge) {
					return
				}
			}
		}
	}
}

// function to reset a graph by clearing its edges and nodes
func (g *graphData[K]) Clear() {
	clear(g.Adjacencies)
}

// function to return the number of nodes in the graph
func (g *graphData[K]) NumberOfNodes() int {
	return len(slices.Collect(g.Nodes()))
}

// function to return the number of edges in the graph
func (g *graphData[K]) NumberOfEdges() int {
	return len(slices.Collect(g.Edges()))
}

// function to return the successors of a node in the graph
func (g *graphData[K]) Successors(n Node[K]) iter.Seq[Node[K]] {
	// create the iterator
	return func(yield func(Node[K]) bool) {
		// walk the neighbors of the node
		for succ := range g.Adjacencies[n] {
			// and yield it
			if !yield(succ) {
				return
			}
		}
	}
}

// function to return the predecessors of a node in the graph
func (g *graphData[K]) Predecessors(n Node[K]) iter.Seq[Node[K]] {
	// create the iterator
	return func(yield func(Node[K]) bool) {
		// walk all nodes
		for node := range g.Adjacencies {
			// walk the node's neighbors
			for neigh := range g.Adjacencies[node] {
				// is it the node we are looking for?
				if neigh == n {
					// yield it
					if !yield(neigh) {
						return
					}
				}
			}
		}
	}
}

// functions to return the in-degree, out-degree, and its sum
func (g *graphData[K]) InDegree(n Node[K]) int {
	return len(slices.Collect(g.Predecessors(n)))
}

func (g *graphData[K]) OutDegree(n Node[K]) int {
	return len(slices.Collect(g.Successors(n)))
}

func (g *graphData[K]) Degree(n Node[K]) int {
	return g.InDegree(n) + g.OutDegree(n)
}

// function to return all the neighbors of a node in the graph
func (g *graphData[K]) Neighbors(n Node[K]) iter.Seq[Node[K]] {
	// create the iterator
	return func(yield func(Node[K]) bool) {
		// combine the successors and predecessors
		for n := range g.Successors(n) {
			if !yield(n) {
				return
			}
		}
		for n := range g.Predecessors(n) {
			if !yield(n) {
				return
			}
		}
	}
}

// helper to create an empty new graphData structure
func newGraphData[K comparable]() graphData[K] {
	return graphData[K]{
		Adjacencies: make(map[Node[K]]map[Node[K]]float64),
	}
}
