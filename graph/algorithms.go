package graph

import (
	"math"
	"slices"
)

// define a queue to work on - just a list of nodes
type Queue[K comparable] []Node[K]

// define a structure to repreesent the ordered path of nodes
// through a graph
type Path[K comparable] []Node[K]

// return types for path finding
type Distances[K comparable] map[Node[K]]float64
type Paths[K comparable] map[Node[K]]Node[K]

// implement a breadth-first search from a start node
// to a destination node. returns the path, and its length
func (g *graphData[K]) BFS(start Node[K]) (Distances[K], Paths[K]) {
	// create a queue
	queue := make(Queue[K], 1)
	// create a map to track which nodes have been explored
	visited := make(map[Node[K]]bool)

	// outputs
	distances := make(Distances[K])
	previous := make(Paths[K])

	// mark the starting node as explored
	visited[start] = true
	// distance for start
	distances[start] = 0
	// no parent for start
	previous[start] = start
	// seed the queue
	queue[0] = start

	// process while queue isn't empty
	for len(queue) > 0 {
		// pop the front of the queue
		current := queue[0]
		queue = queue[1:]

		// go through all the possible neighbors of the current node
		for neighbor := range g.Adjacencies[current] {
			// check if we've already been at this neighbor
			if _, explored := visited[neighbor]; !explored {
				visited[neighbor] = true
				previous[neighbor] = current
				distances[neighbor] = distances[current] + 1.0
				queue = append(queue, neighbor)
			}
		}
	}

	return distances, previous
}

func (g *graphData[K]) BFSTo(start, target Node[K]) (Path[K], int, float64) {
	// are we already there?
	if start == target {
		return Path[K]{start}, 1, 0.0
	}

	// run pathn finding
	distances, previous := g.BFS(start)

	// check if the target could in fact be reached
	if _, ok := previous[target]; !ok {
		// no, can't get to it, return empty path and zero length
		return Path[K]{}, 0, math.Inf(1)
	}

	// build the path from parent relationships
	path := make(Path[K], 1)
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
	return path, len(path), distances[target]
}

// calculate the shortest path from a given start to
// all other nodes. return the distances and previous
// nodes for each node in the graph
func (g *graphData[K]) Dijkstra(start Node[K]) (Distances[K], Paths[K]) {
	// initialize the queue and data structures to hold
	// the distances and prior nodes on the paths
	queue := make(Queue[K], 0)
	distances := make(Distances[K])
	previous := make(Paths[K])
	// for each node, set the distance to infinity and add
	// it to the queue
	for node := range g.Adjacencies {
		distances[node] = math.Inf(1)
		queue = append(queue, node)
	}
	// distance to the starting node is 0.0
	distances[start] = 0.0
	// can get to self
	previous[start] = start

	// process queue while it isn't empty
	for len(queue) > 0 {
		// find the node with the smallest distance still in the queue
		min_distance := math.Inf(1)
		min_index := 0
		for i := range queue {
			if distances[queue[i]] < min_distance {
				min_index = i
				min_distance = distances[queue[i]]
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
func (g *graphData[K]) DijkstraTo(start, target Node[K]) (Path[K], int, float64) {
	// calculate the graph distances and paths
	distances, previous := g.Dijkstra(start)

	// check that the target can be reached from the given start
	if _, ok := previous[target]; !ok {
		// it cannot
		return Path[K]{}, 0, math.Inf(1)
	}

	// build the path from parent relationships
	path := make(Path[K], 1)
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
