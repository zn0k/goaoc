package graph

import (
	"bufio"
	"fmt"
	"maps"
	"os"
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
	HasNode(n Node[K]) bool
	HasEdge(u, v Node[K]) bool
	RemoveNode(n Node[K])
	RemoveEdge(u, v Node[K])
	RemoveNodesFrom(ns []Node[K])
	RemoveEdgesFrom(es []Edge[K])
	Nodes() []Node[K]
	Edges() []Edge[K]
	Clear()
	NumberOfNodes() int
	NumberOfEdges() int
	NewNode(obj K) Node[K]
	Successors(n Node[K]) []Node[K]
	Predecessors(n Node[K]) []Node[K]
	Neighbors(n Node[K]) []Node[K]
	InDegree(n Node[K]) int
	OutDegree(n Node[K]) int
	Degree(n Node[K]) int
	Copy() graphData[K]
}

// generic data structure for a graph. it's a simple lookup
// table for graphs and list of graphs with the weight associated
// with the edge between the two keys
type graphData[K comparable] struct {
	Adjacencies map[Node[K]]map[Node[K]]float64
}

// function to wrap a new node
func (g *graphData[K]) NewNode(obj K) Node[K] {
	return Node[K]{ID: obj}
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

// function to check whether the graph has a node
func (g *graphData[K]) HasNode(u Node[K]) bool {
	_, ok := g.Adjacencies[u]
	return ok
}

// function to check whether the grpah has an edge
func (g *graphData[K]) HasEdge(u, v Node[K]) bool {
	_, hasU := g.Adjacencies[u]
	if !hasU {
		return false
	}
	_, hasV := g.Adjacencies[u][v]
	return hasV
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
func (g *graphData[K]) Nodes() []Node[K] {
	return slices.Collect(maps.Keys(g.Adjacencies))
}

// function to retrieve a list of edges from a graph
func (g *graphData[K]) Edges() []Edge[K] {
	edges := make([]Edge[K], 0)
	for u := range g.Adjacencies {
		// walk the node's adjacencies
		for v, w := range g.Adjacencies[u] {
			// create the edge
			edges = append(edges, Edge[K]{u: u, v: v, weight: w})
		}
	}
	return edges
}

// function to reset a graph by clearing its edges and nodes
func (g *graphData[K]) Clear() {
	clear(g.Adjacencies)
}

// function to return the number of nodes in the graph
func (g *graphData[K]) NumberOfNodes() int {
	return len(g.Nodes())
}

// function to return the number of edges in the graph
func (g *graphData[K]) NumberOfEdges() int {
	return len(g.Edges())
}

// function to return the successors of a node in the graph
func (g *graphData[K]) Successors(n Node[K]) []Node[K] {
	return slices.Collect(maps.Keys(g.Adjacencies[n]))
}

// function to return the predecessors of a node in the graph
func (g *graphData[K]) Predecessors(n Node[K]) []Node[K] {
	predecessors := make([]Node[K], 0)
	// walk all nodes
	for node := range g.Adjacencies {
		// walk the node's neighbors
		for neigh := range g.Adjacencies[node] {
			// is it the node we are looking for?
			if neigh == n {
				predecessors = append(predecessors, node)
			}
		}
	}
	return predecessors
}

// functions to return the in-degree, out-degree, and its sum
func (g *graphData[K]) InDegree(n Node[K]) int {
	return len(g.Predecessors(n))
}

func (g *graphData[K]) OutDegree(n Node[K]) int {
	return len(g.Successors(n))
}

func (g *graphData[K]) Degree(n Node[K]) int {
	return g.InDegree(n) + g.OutDegree(n)
}

// function to return all the neighbors of a node in the graph
func (g *graphData[K]) Neighbors(n Node[K]) []Node[K] {
	return append(g.Successors(n), g.Predecessors(n)...)
}

// function to deep copy a graph
func (g *graphData[K]) Copy() *graphData[K] {
	// create new graph
	newG := newGraphData[K]()
	// registry for copied nodes
	nodesMap := make(map[K]Node[K])
	// function to either retrieve a copy of a node, or create it
	getOrCreate := func(id K) Node[K] {
		if node, ok := nodesMap[id]; ok {
			return node
		}
		newNode := Node[K]{ID: id}
		nodesMap[id] = newNode
		return newNode
	}

	for n, neighbors := range g.Adjacencies {
		newNode := getOrCreate(n.ID)
		if _, ok := newG.Adjacencies[newNode]; !ok {
			newG.Adjacencies[newNode] = make(map[Node[K]]float64)
		}
		for nei, weight := range neighbors {
			newNeighbor := getOrCreate(nei.ID)
			newG.Adjacencies[newNode][newNeighbor] = weight
		}
	}
	return &newG
}

// helper to create an empty new graphData structure
func newGraphData[K comparable]() graphData[K] {
	return graphData[K]{
		Adjacencies: make(map[Node[K]]map[Node[K]]float64),
	}
}

// function to export the edge list into a given file
// this can usually be imported by other graphing libraries
func (g *graphData[K]) ExportEdgeList(fname string) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	for _, e := range g.Edges() {
		fmt.Fprintf(writer, "'%v' '%v'\n", e.u.ID, e.v.ID)
	}
	return writer.Flush()
}
