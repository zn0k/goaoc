package graph

// UndirectedGraph inherits from graphData
type UndirectedGraph[K comparable] struct {
	graphData[K]
}

// helper to generate a new UndirectedGraph
func NewUndirectedGraph[K comparable]() *UndirectedGraph[K] {
	return &UndirectedGraph[K]{
		graphData: newGraphData[K](),
	}
}

// adding new edges to an undirected graphs adds
// them both ways, from u to v and from v to u
func (g *UndirectedGraph[K]) AddEdge(u, v Node[K], w float64) {
	// add nodes to graph if they don't exist yet
	if _, ok := g.Adjacencies[u]; !ok {
		g.AddNode(u)
	}
	if _, ok := g.Adjacencies[v]; !ok {
		g.AddNode(v)
	}

	// add the edges and adjacencies both ways
	g.Adjacencies[u][v] = w
	g.Adjacencies[v][u] = w
}

// add from an iter of edges
func (g *UndirectedGraph[K]) AddEdgesFrom(es []Edge[K]) {
	for _, e := range es {
		g.AddEdge(e.u, e.v, e.weight)
	}
}

// remove an edge from an undirected graph
// this removes the edge both ways
func (g *UndirectedGraph[K]) RemoveEdge(u, v Node[K]) {
	delete(g.Adjacencies[u], v)
	delete(g.Adjacencies[v], u)
}

// remove edges from an undirected graph using an iter as the source
func (g *UndirectedGraph[K]) RemoveEdgesFrom(es []Edge[K]) {
	for _, e := range es {
		delete(g.Adjacencies[e.u], e.v)
		delete(g.Adjacencies[e.v], e.u)
	}
}

// override Neighbors, Predecessors, and Degrees for UndirectedGraph
// Neighbors and Predecessors are all the same as Successors, so make
// the former not double count and the latter cheaper to implement
func (g *UndirectedGraph[K]) Neighbors(n Node[K]) []Node[K] {
	return g.Successors(n)
}

func (g *UndirectedGraph[K]) Predecessors(n Node[K]) []Node[K] {
	return g.Successors(n)
}

// and Degrees is just the number of neighbors
func (g *UndirectedGraph[K]) Degree(n Node[K]) int {
	return len(g.Neighbors(n))
}
