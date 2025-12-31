package graph

// DirectedGraph also inherits from graphData
type DirectedGraph[K comparable] struct {
	graphData[K]
}

// constructor
func NewDirectedGraph[K comparable]() *DirectedGraph[K] {
	return &DirectedGraph[K]{
		graphData: newGraphData[K](),
	}
}

// directed graphs add edges only in the indicated direction
func (g *DirectedGraph[K]) AddEdge(u, v Node[K], w float64) {
	// add nodes to graph if they don't exist yet
	if _, ok := g.Adjacencies[u]; !ok {
		g.AddNode(u)
	}
	if _, ok := g.Adjacencies[v]; !ok {
		g.AddNode(v)
	}

	// add the edge and adjancency
	g.Adjacencies[u][v] = w
}

// add from an iter of edges
func (g *DirectedGraph[K]) AddEdgesFrom(es []Edge[K]) {
	for _, e := range es {
		g.AddEdge(e.u, e.v, e.weight)
	}
}

// remove an edge from a directed graph
func (g *DirectedGraph[K]) RemoveEdge(u, v Node[K]) {
	delete(g.Adjacencies[u], v)
}

// remove edges from an undirected graph using an iter as the source
func (g *DirectedGraph[K]) RemoveEdgesFrom(es []Edge[K]) {
	for _, e := range es {
		delete(g.Adjacencies[e.u], e.v)
	}
}
