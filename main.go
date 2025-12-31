package main

import (
	"fmt"
	"iter"
	"maps"
	"math"
	"os"
	"slices"
	"strings"
)

// coordinates have X and Y components
type Coordinate struct {
	X, Y int
}

// they implement Stringer for easy printing
func (c Coordinate) String() string {
	return fmt.Sprintf("(%d, %d)", c.X, c.Y)
}

// directions for walking a grid are really just coordinates:
// {0, 1}, {0, -1}, {1, 0}, {-1, 0}
type Direction Coordinate

// define an interface for an abstract graph that can have nodes
// and edges added to it, can have them removed, and return iterators
// over the nodes and edges of the graph
type Graph interface {
	AddNode(n Node)
	AddEdge(u, v Node, w float64)
	RemoveNode(n Node)
	RemoveEdge(u, v Node)
	Nodes() iter.Seq[Node]
	Edges() iter.Seq[Edge]
}

// nodes are identified by their grid coordinates
type Node struct {
	ID Coordinate
}

// edges are identified by the nodes they connect,
// and the weight of the connection
type Edge struct {
	u, v   Node
	weight float64
}

// an adjacency is defined by the other end point and
// the edge weight
type Adjancency struct {
	v      Node
	weight float64
}

// generic data structure for a graph. it's a simple lookup
// table for graphs and list of graphs with the weight associated
// with the edge between the two keys
type graphData struct {
	Adjacencies map[Node]map[Node]float64
}

// function to add a node to the graph
func (g *graphData) AddNode(n Node) {
	// does the node already exist in the graph?
	if _, ok := g.Adjacencies[n]; !ok {
		// no, add it with no adjacencies
		g.Adjacencies[n] = make(map[Node]float64)
	}
}

// function to remove a node from the graph
func (g *graphData) RemoveNode(n Node) {
	// remove all adjancencies to the node
	for node := range g.Adjacencies {
		delete(g.Adjacencies[node], n)
	}
	// remove adjacencies from the node, and with that its record
	delete(g.Adjacencies, n)
}

// function to retrieve an iterator over the nodes of the graph
func (g *graphData) Nodes() iter.Seq[Node] {
	return maps.Keys(g.Adjacencies)
}

// function to retrieve a list of edges from a graph
func (g *graphData) Edges() iter.Seq[Edge] {
	// create the iterator
	return func(yield func(Edge) bool) {
		// walk the nodes
		for u := range g.Adjacencies {
			// walk the node's adjacencies
			for v, w := range g.Adjacencies[u] {
				// create the edge
				edge := Edge{u: u, v: v, weight: w}
				// and yield it
				if !yield(edge) {
					return
				}
			}
		}
	}
}

// define a queue to work on - just a list of nodes
type Queue []Node

// define a structure to repreesent the ordered path of nodes
// through a graph
type Path []Node

// implement a breadth-first search from a start node
// to a destination node. returns the path, and its length
func (g *graphData) BFS(start, target Node) (Path, int) {
	// if we're already there...
	if start == target {
		return Path{}, 0
	}

	// create a queue
	queue := make(Queue, 1)
	// create a map to track which nodes have been explored
	visited := make(map[Node]bool)

	// mark the starting node as explored
	visited[start] = true
	// seed the queue
	queue[0] = start

	// initialize the path by keeping track of the prior step to each node
	previous := make(map[Node]Node)

	// process while queue isn't empty
	for len(queue) > 0 {
		// pop the front of the queue
		current := queue[0]
		queue = queue[1:]

		// check if we're at the target
		if current == target {
			break
		}
		// go through all the possible neighbors of the current node
		for neighbor, _ := range g.Adjacencies[current] {
			// check if we've already been at this neighbor
			if _, explored := visited[neighbor]; !explored {
				visited[neighbor] = true
				previous[neighbor] = current
				queue = append(queue, neighbor)
			}
		}
	}

	// build the path from parent relationships
	path := make(Path, 1)
	// walk back from the target
	path[0] = target
	current := target
	for current != start {
		step := previous[current]
		current = previous[current]
		path = append(path, step)
	}
	// and reverse it
	slices.Reverse(path)

	// return the path and its length
	return path, len(path)
}

// calculate the shortest path from a given start to
// all other nodes. return the distances and previous
// nodes for each node in the graph
func (g *graphData) Dijkstra(start Node) (map[Node]float64, map[Node]Node) {
	// initialize the queue and data structures to hold
	// the distances and prior nodes on the paths
	queue := make(Queue, 0)
	distances := make(map[Node]float64)
	previous := make(map[Node]Node)
	// for each node, set the distance to infinity and add
	// it to the queue
	for node := range g.Adjacencies {
		distances[node] = math.Inf(1)
		queue = append(queue, node)
	}
	// distance to the starting node is 0.0
	distances[start] = 0.0

	// process queue while it isn't empty
	for len(queue) > 0 {
		// find the node with the smallest distance still in the queue
		min_distance := math.Inf(1)
		min_index := 0
		for i := range queue {
			if distances[queue[i]] < min_distance {
				min_index = i
			}
		}
		// fetch it, and remove it from the queue
		current := queue[min_index]
		queue = slices.Delete(queue, min_index, min_index+1)

		// go through all the possible neighbors of the current node
		for neighbor, weight := range g.Adjacencies[current] {
			// calculate the distance from this node to the neighbor
			// by adding the weight of the edge
			alternative := distances[current] + weight
			// is that a cheaper way to the neighbor?
			if alternative < distances[neighbor] {
				// yes. update its distance and set this node to be
				// on the path to it
				distances[neighbor] = alternative
				previous[neighbor] = current
			}
		}
	}

	return distances, previous
}

// calculate the shortest path from a given node to a given node
// returns the path, the length of the path, and the distance cost
func (g *graphData) DijkstraTo(start, target Node) (Path, int, float64) {
	// calculate the graph distances and paths
	distances, previous := g.Dijkstra(start)

	// check that the target can be reached from the given start
	if _, ok := previous[target]; !ok {
		// it cannot
		return Path{}, 0, 0.0
	}

	// build the path from parent relationships
	path := make(Path, 1)
	// walk back from the target
	path[0] = target
	current := target
	for current != start {
		step := previous[current]
		current = previous[current]
		path = append(path, step)
	}
	// and reverse it
	slices.Reverse(path)

	return path, len(path), distances[target]
}

// helper to create an empty new graphData structure
func newGraphData() graphData {
	return graphData{
		Adjacencies: make(map[Node]map[Node]float64),
	}
}

// UndirectedGraph inherits from graphData
type UndirectedGraph struct {
	graphData
}

// helper to generate a new UndirectedGraph
func NewUndirectedGraph() *UndirectedGraph {
	return &UndirectedGraph{
		graphData: newGraphData(),
	}
}

// adding new edges to an undirected graphs adds
// them both ways, from u to v and from v to u
func (g *UndirectedGraph) AddEdge(u, v Node, w float64) {
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

// remove an edge from an undirected graph
// this removes the edge both ways
func (g *UndirectedGraph) RemoveEdge(u, v Node) {
	delete(g.Adjacencies[u], v)
	delete(g.Adjacencies[v], u)
}

// DirectedGraph also inherits from graphData
type DirectedGraph struct {
	graphData
}

func NewDirectedGraph() *DirectedGraph {
	return &DirectedGraph{
		graphData: newGraphData(),
	}
}

// directed graphs add edges only in the indicated direction
func (g *DirectedGraph) AddEdge(u, v Node, w float64) {
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

// remove an edge from a directed graph
func (g *DirectedGraph) RemoveEdge(u, v Node) {
	delete(g.Adjacencies[u], v)
}

// read in the maze grid and return an undirected graph as well as the start
// and end tile on the grid
func readLines(fname string, directions []Direction) (*UndirectedGraph, Node, Node) {
	buf, err := os.ReadFile(fname)
	if err != nil {
		panic(fmt.Sprintf("unable to open %s for reading", fname))
	}

	// initialize the start and end tiles, and the grid as 2d runes
	start := Node{}
	target := Node{}
	var grid [][]rune

	// walk the lines
	for y, line := range strings.Split(string(buf), "\n") {
		// walk each row
		var row []rune
		for x, c := range line {
			// check if we're on the start or end tile.
			// if so, record it, and then turn it into a normal tile
			if c == 'S' {
				start = Node{Coordinate{x, y}}
				c = '.'
			}
			if c == 'T' {
				target = Node{Coordinate{x, y}}
				c = '.'
			}
			// build the row
			row = append(row, c)
		}
		// build the grid from rows
		grid = append(grid, row)
	}

	// initialize a new graph
	g := NewUndirectedGraph()

	// walk the grid
	height, width := len(grid), len(grid[0])
	for y, row := range grid {
		for x, c := range row {
			// on a wall, this isn't a valid node
			if c != '.' {
				continue
			}
			// on a walkable tile. explore its neighbors
			for _, d := range directions {
				// calculate the neighbor coordinates
				new_x, new_y := x+d.X, y+d.Y
				// are they within the grid?
				if new_x < 0 || new_x >= width || new_y < 0 || new_y >= height {
					// no, outside the grid
					continue
				}
				// is the neighbor walkable?
				if grid[new_y][new_x] == '.' {
					// yes. create a node for the current position
					// and a node for the neighbor
					u := Node{Coordinate{x, y}}
					v := Node{Coordinate{new_x, new_y}}
					// add an edge between them, which also adds the nodes
					g.AddEdge(u, v, 1.0)
				}
			}
		}
	}

	return g, start, target
}

func main() {
	// this grid is only walkable in cardinal directions
	ds := []Direction{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
	}
	// read in the grid
	g, s, t := readLines("input.txt", ds)

	// run a BFS
	path, length := g.BFS(s, t)
	fmt.Printf("path=%v, length=%d\n", path, length)

	// run dijkstra
	path, length, _ = g.DijkstraTo(s, t)
	fmt.Printf("path=%v, length=%d\n", path, length)
}
