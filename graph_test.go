package graph

import "testing"

func TestUndirectedGraph_AddNodesAndEdges(t *testing.T) {
	t.Run("Node and edge count", func(t *testing.T) {
		// create an undirected graph
		g := NewUndirectedGraph[int]()
		// create three nodes
		u := Node[int]{1}
		v := Node[int]{2}
		w := Node[int]{3}
		x := Node[int]{3}
		y := Node[int]{3}

		// add u, v, and an edge between the two
		g.AddEdge(u, v, 10.0)
		// add w as a node, but without any edges
		g.AddNode(w)
		// add x and y from a list
		g.AddNodesFrom([]Node[int]{x, y})
		// add an edge from a list
		g.AddEdgesFrom([]Edge[int]{{x, y, 5.0}})

		// check node count
		n := g.NumberOfNodes()
		if n != 5 {
			t.Errorf("Expected 6 nodes, got %d", n)
		}

		// check edge count
		n = g.NumberOfEdges()
		if n != 4 {
			t.Errorf("Expected 4 edges, got %d", n)
		}
	})
}

func TestUndirectedGraph_RemoveNodesAndEdges(t *testing.T) {
	t.Run("Node and edge count", func(t *testing.T) {
		// create an undirected graph
		g := NewUndirectedGraph[int]()
		// create three nodes
		u := Node[int]{1}
		v := Node[int]{2}
		w := Node[int]{3}
		x := Node[int]{4}
		y := Node[int]{5}
		z := Node[int]{6}

		// add u, v, and an edge between the two
		g.AddEdge(u, v, 10.0)
		// add w, x, and an edge between the two
		g.AddEdge(w, x, 5.0)
		// add just y as a node from a list with some overlap
		g.AddNodesFrom([]Node[int]{w, x, y})

		// check node count
		n := g.NumberOfNodes()
		if n != 5 {
			t.Errorf("Expected 5 nodes, got %d", n)
		}

		// remove y as a node
		g.RemoveNode(y)
		// remove x and w from a list
		g.RemoveNodesFrom([]Node[int]{w, x})

		// check node count
		n = g.NumberOfNodes()
		if n != 2 {
			t.Errorf("Expected 2 nodes, got %d", n)
		}

		// add edges for (w, x) and (z, y)
		g.AddEdge(w, x, 5.0)
		g.AddEdge(y, z, 1.0)

		// check edge count
		n = g.NumberOfEdges()
		if n != 6 {
			t.Errorf("Expected 6 edges, got %d", n)
		}

		// remove edge between w and x
		g.RemoveEdge(w, x)
		// and remove y, z from a list
		g.RemoveEdgesFrom([]Edge[int]{{y, z, 1.0}})

		// check edge count
		n = g.NumberOfEdges()
		if n != 2 {
			t.Errorf("Expected 2 edges, got %d", n)
		}
	})
}

func TestUndirectedGraph_Clear(t *testing.T) {
	t.Run("Node and edge count", func(t *testing.T) {
		// create an undirected graph
		g := NewUndirectedGraph[int]()
		// create three nodes
		u := Node[int]{1}
		v := Node[int]{2}

		// add u, v, and an edge between the two
		g.AddEdge(u, v, 10.0)

		// and clear
		g.Clear()

		// check node count
		n := g.NumberOfNodes()
		if n != 0 {
			t.Errorf("Expected 0 nodes, got %d", n)
		}
		// check edge count
		n = g.NumberOfEdges()
		if n != 0 {
			t.Errorf("Expected 0 edges, got %d", n)
		}
	})
}
