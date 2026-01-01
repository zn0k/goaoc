package graph

import (
	"iter"
	"slices"
	"testing"
)

func getNodes() (Node[int], Node[int], Node[int], Node[int], Node[int], Node[int]) {
	u := Node[int]{1}
	v := Node[int]{2}
	w := Node[int]{3}
	x := Node[int]{4}
	y := Node[int]{5}
	z := Node[int]{6}

	return u, v, w, x, y, z
}

func TestUndirectedGraph_AddNodesAndEdges(t *testing.T) {
	t.Run("Undirected graph adding nodes and edges", func(t *testing.T) {
		// create an undirected graph
		g := NewUndirectedGraph[int]()

		u, v, w, x, y, _ := getNodes()

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
			t.Errorf("Expected 5 nodes, got %d", n)
		}

		// check edge count
		n = g.NumberOfEdges()
		if n != 4 {
			t.Errorf("Expected 4 edges, got %d", n)
		}
	})
}

func TestUndirectedGraph_RemoveNodesAndEdges(t *testing.T) {
	t.Run("Undirected graph removing nodes and edges", func(t *testing.T) {
		// create an undirected graph
		g := NewUndirectedGraph[int]()

		u, v, w, x, y, z := getNodes()

		// add u, v, and an edge between the two
		g.AddEdge(u, v, 10.0)
		// add w, x, and an edge between the two
		g.AddEdge(w, x, 5.0)
		// add just y as a node from a list with some overlap
		g.AddNode(y)

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
	t.Run("Undirected node clearing", func(t *testing.T) {
		// create an undirected graph
		g := NewUndirectedGraph[int]()

		u, v, _, _, _, _ := getNodes()

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

func TestUndirectedGraph_Neighbors(t *testing.T) {
	t.Run("Undirected node neighbor counts", func(t *testing.T) {
		// create an undirected graph
		g := NewUndirectedGraph[int]()

		u, v, w, x, y, z := getNodes()

		g.AddNodesFrom([]Node[int]{u, v, w, x, y, z})

		// check that u doesn't have neighbors
		n := g.Degree(u)
		if n != 0 {
			t.Errorf("Expected degree of 0, got %d", n)
		}

		// add edges from u to v and w
		g.AddEdge(u, v, 1.0)
		g.AddEdge(u, w, 1.0)

		// test again for degree
		n = g.Degree(u)
		if n != 2 {
			t.Errorf("Expected degree of 2, got %d", n)
		}
		// test in-degree and out-degree
		n = g.InDegree(u)
		if n != 2 {
			t.Errorf("Expected in-degree of 2, got %d", n)
		}
		n = g.OutDegree(u)
		if n != 2 {
			t.Errorf("Expected out-degree of 2, got %d", n)
		}

		n = g.Degree(v)
		if n != 1 {
			t.Errorf("Expected degree of 1, got %d", n)
		}
		// test in-degree and out-degree
		n = g.InDegree(v)
		if n != 1 {
			t.Errorf("Expected in-degree of 1, got %d", n)
		}
		n = g.OutDegree(v)
		if n != 1 {
			t.Errorf("Expected out-degree of 1, got %d", n)
		}

		// check for the two specific neighbors
		// sucessors and predecessors should be the same
		check := func(ns iter.Seq[Node[int]]) bool {
			first := false
			second := false
			for n := range ns {
				if n == v {
					first = true
				}
				if n == w {
					second = true
				}
			}
			return first && second
		}
		expected := []Node[int]{v, w}
		ns := g.Neighbors(u)
		if !check(ns) {
			t.Errorf("Neighbors function expected %v, got %v", expected, ns)
		}
		ns = g.Successors(u)
		if !check(ns) {
			t.Errorf("Successors function expected %v, got %v", expected, ns)
		}
		ns = g.Predecessors(u)
		if !check(ns) {
			t.Errorf("Predecessors function expected %v, got %v", expected, ns)
		}
	})
}

func TestDirectedGraph_AddNodesAndEdges(t *testing.T) {
	t.Run("Directed graph adding nodes and edges", func(t *testing.T) {
		// create a directed graph
		g := NewDirectedGraph[int]()

		u, v, w, x, y, _ := getNodes()

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
			t.Errorf("Expected 5 nodes, got %d", n)
		}

		// check edge count
		n = g.NumberOfEdges()
		if n != 2 {
			t.Errorf("Expected 2 edges, got %d", n)
		}
	})
}

func TestDirectedGraph_RemoveNodesAndEdges(t *testing.T) {
	t.Run("Directed graph removing nodes and edges", func(t *testing.T) {
		// create a directed graph
		g := NewDirectedGraph[int]()

		u, v, w, x, y, z := getNodes()

		// add u, v, and an edge between the two
		g.AddEdge(u, v, 10.0)
		// add w, x, and an edge between the two
		g.AddEdge(w, x, 5.0)
		// add just y as a node from a list with some overlap
		g.AddNode(y)

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
		if n != 3 {
			t.Errorf("Expected 3 edges, got %d", n)
		}

		// remove edge between w and x
		g.RemoveEdge(w, x)
		// and remove y, z from a list
		g.RemoveEdgesFrom([]Edge[int]{{y, z, 1.0}})

		// check edge count
		n = g.NumberOfEdges()
		if n != 1 {
			t.Errorf("Expected 1 edge, got %d", n)
		}
	})
}

func TestDirectedGraph_Clear(t *testing.T) {
	t.Run("Directed node clearing", func(t *testing.T) {
		// create a directed graph
		g := NewDirectedGraph[int]()

		u, v, _, _, _, _ := getNodes()

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

func TestDirectedGraph_Neighbors(t *testing.T) {
	t.Run("Directed node neighbor counts", func(t *testing.T) {
		// create a directed graph
		g := NewDirectedGraph[int]()

		u, v, w, x, y, z := getNodes()

		g.AddNodesFrom([]Node[int]{u, v, w, x, y, z})

		// check that u doesn't have neighbors
		n := g.Degree(u)
		if n != 0 {
			t.Errorf("Expected degree of 0, got %d", n)
		}

		// add edges from u to v and w
		g.AddEdge(u, v, 1.0)
		g.AddEdge(u, w, 1.0)

		// test again for degree
		n = g.Degree(u)
		if n != 2 {
			t.Errorf("Expected degree of 2, got %d", n)
		}
		// test in-degree and out-degree
		n = g.InDegree(u)
		if n != 0 {
			t.Errorf("Expected in-degree of 0, got %d", n)
		}
		n = g.OutDegree(u)
		if n != 2 {
			t.Errorf("Expected out-degree of 2, got %d", n)
		}

		n = g.Degree(v)
		if n != 1 {
			t.Errorf("Expected degree of 1, got %d", n)
		}
		// test in-degree and out-degree
		n = g.InDegree(v)
		if n != 1 {
			t.Errorf("Expected in-degree of 1, got %d", n)
		}
		n = g.OutDegree(v)
		if n != 0 {
			t.Errorf("Expected out-degree of 0, got %d", n)
		}

		// check for the two specific neighbors
		// sucessors should be the same
		check := func(ns iter.Seq[Node[int]]) bool {
			first := false
			second := false
			for n := range ns {
				if n == v {
					first = true
				}
				if n == w {
					second = true
				}
			}
			return first && second
		}
		expected := []Node[int]{v, w}
		ns := g.Neighbors(u)
		if !check(ns) {
			t.Errorf("Neighbors function expected %v, got %v", expected, ns)
		}
		ns = g.Successors(u)
		if !check(ns) {
			t.Errorf("Successors function expected %v, got %v", expected, ns)
		}
		n = len(slices.Collect(g.Predecessors(u)))
		if n != 0 {
			t.Errorf("Length of predecessors for u was expected to be 0, got %d", n)
		}
		n = len(slices.Collect(g.Predecessors(v)))
		if n != 1 {
			t.Errorf("Length of predecessors for v was expected to be 1, got %d", n)
		}
	})
}
