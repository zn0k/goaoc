package graph

import (
	"math"
	"testing"
)

func TestBFS(t *testing.T) {
	// create an undirected graph
	g := NewUndirectedGraph[int]()
	u, v, w, x, y, z := getNodes()

	// create a 5-node line graph
	g.AddEdge(u, v, 1.0)
	g.AddEdge(v, w, 1.0)
	g.AddEdge(w, x, 1.0)
	g.AddEdge(x, y, 1.0)

	// add an unreachable node
	g.AddNode(z)

	t.Run("BFS algorithm unreachable node", func(t *testing.T) {
		// check that z is unreachable
		path_to_z, zl := g.BFS(u, z)
		if len(path_to_z) != 0 || zl != 0 {
			t.Errorf("BFS expected zero length path to z, got %d and %d", len(path_to_z), zl)
		}

	})

	t.Run("BFS algorithm shortest path", func(t *testing.T) {
		// check correct path from u to y
		path_to_y, yl := g.BFS(u, y)
		if len(path_to_y) != 5 || yl != 5 {
			t.Errorf("BFS expected shortest length path to y over 5 nodes, got %d and %d", len(path_to_y), yl)
		}
	})

	// add a cycle
	g.AddEdge(y, u, 1.0)

	t.Run("BFS algorithm with cyclical graph", func(t *testing.T) {
		// check correct path from u to y with a cycle present
		// path is now also shorter
		path_to_y, yl := g.BFS(u, y)
		if len(path_to_y) != 2 || yl != 2 {
			t.Errorf("BFS expected shortest length path to y over 2 nodes, got %d and %d", len(path_to_y), yl)
		}
	})
}

func TestDijkstra(t *testing.T) {
	// create an undirected graph
	g := NewUndirectedGraph[int]()
	u, v, w, x, y, z := getNodes()

	// create a 5-node cyclical graph
	// all weights are 1.0
	g.AddEdge(u, v, 1.0)
	g.AddEdge(v, w, 1.0)
	g.AddEdge(w, x, 1.0)
	g.AddEdge(x, y, 1.0)

	// add an unreachable node
	g.AddNode(z)

	t.Run("Dijkstra algorithm unreachable node", func(t *testing.T) {
		// check that z is unreachable
		path_to_z, zl, cost := g.DijkstraTo(u, z)
		if len(path_to_z) != 0 || zl != 0 {
			t.Errorf("Dijkstra expected zero length path to z, got %d and %d", len(path_to_z), zl)
		}
		if cost != math.Inf(1) {
			t.Errorf("Dijkstra expected infinite cost to z, got %v", cost)
		}

	})

	t.Run("Dijktra algorithm shortest path", func(t *testing.T) {
		// check correct path from u to y
		path_to_y, yl, cost := g.DijkstraTo(u, y)
		if len(path_to_y) != 5 || yl != 5 {
			t.Errorf("Dijkstra expected path from u to y to be 5, got %d and %d", len(path_to_y), yl)
		}
		if cost != 4.0 {
			t.Errorf("Dijkstra expected cost to y to be 5.0, got %f", cost)
		}
	})

	t.Run("Dijktra algorithm cost to self", func(t *testing.T) {
		// check correct path from u to y
		path_to_u, ul, cost := g.DijkstraTo(u, u)
		if len(path_to_u) != 1 || ul != 1 {
			t.Errorf("Dijkstra expected length path to u from u to be 1, got %d and %d", len(path_to_u), ul)
		}
		if cost != 0.0 {
			t.Errorf("Dijkstra expected cost from u to u to be 0.0, got %f", cost)
		}
	})

	// add a cycle
	g.AddEdge(y, u, 1.0)

	t.Run("Dijktra algorithm shortest path with cycle", func(t *testing.T) {
		// check correct path from u to y
		path_to_y, yl, cost := g.DijkstraTo(u, y)
		if len(path_to_y) != 2 || yl != 2 {
			t.Errorf("Dijkstra expected path from u to y to be 2, got %d and %d", len(path_to_y), yl)
		}
		if cost != 1.0 {
			t.Errorf("Dijkstra expected cost to y to be 1.0, got %f", cost)
		}
	})

	// add a direct path from u to w with cost 1.0
	// check that it will be taken
	g.AddEdge(u, w, 1.0)
	t.Run("Dijktra algorithm weighted shortest path", func(t *testing.T) {
		// check correct path from u to y
		path_to_w, wl, cost := g.DijkstraTo(u, w)
		if len(path_to_w) != 2 || wl != 2 {
			t.Errorf("Dijkstra expected path from u to w to be 2, got %d and %d", len(path_to_w), wl)
		}
		if cost != 1.0 {
			t.Errorf("Dijkstra expected cost to y to be 1.0, got %f", cost)
		}
	})

}
