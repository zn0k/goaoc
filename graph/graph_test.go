package graph

import (
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

func TestNewNode(t *testing.T) {
	t.Run("Create new nodes", func(t *testing.T) {
		// create a new graph
		g:= NewUndirectedGraph[int]()
		// get the node
		node := g.NewNode(2)
		// check its ID
		if node.ID != 2 {
			t.Errorf("Expected node ID to be 2, got %v", node.ID)
		}
	})
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

		// check for node existence
		ok := g.HasNode(w)
		if !ok {
			t.Errorf("Expected graph to have node w, got %t", ok)
		}

		// check for edge existence
		ok = g.HasEdge(u, v)
		if !ok {
			t.Errorf("Expected graph to have edge between u and v, got %t", ok)
		}

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

		// add a duplicate edge, with a different weight
		g.AddEdge(u, v, 20.0)
		// edge count should still be the same
		// check edge count
		n = g.NumberOfEdges()
		if n != 4 {
			t.Errorf("Expected 4 edges after adding duplicate, got %d", n)
		}

		// and the weight should be the new value
		if g.Adjacencies[u][v] != 20.0 {
			t.Errorf("Expected new edge weight %f, got %f", 20.0, g.Adjacencies[u][v])
		}

		// add a duplicate node
		g.AddNode(u)
		// node count should remain unchanged
		n = g.NumberOfNodes()
		if n != 5 {
			t.Errorf("Expected 5 nodes, got %d", n)
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

		// remove a node that doesn't exist
		// should not panic, node count should remain the same
		g.RemoveNode(z)
		n = g.NumberOfNodes()
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

		// remove an edge that doesn't exist, shouldn't affect edge count
		g.RemoveEdge(z, u)
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
		expected := []Node[int]{v, w}
		ns := g.Neighbors(u)
		if !slices.Contains(ns, v) && slices.Contains(ns, w) {
			t.Errorf("Neighbors function expected %v, got %v", expected, ns)
		}
		ns = g.Successors(u)
		if !slices.Contains(ns, v) && slices.Contains(ns, w) {
			t.Errorf("Successors function expected %v, got %v", expected, ns)
		}
		ns = g.Predecessors(u)
		if !slices.Contains(ns, v) && slices.Contains(ns, w) {
			t.Errorf("Predecessors function expected %v, got %v", expected, ns)
		}
	})
}

func TestUndirectedGraph_Loop(t *testing.T) {
	t.Run("Undirected graph self loop", func(t *testing.T) {
		// create an undirected graph
		g := NewUndirectedGraph[int]()
		u, _, _, _, _, _ := getNodes()

		// add an edge from a node to itself
		g.AddEdge(u, u, 1.0)
		// that should result in one edges, and a degree of 1
		if len(g.Adjacencies[u]) != 1 {
			t.Errorf("Self loop for directed graph expected 1 adjancency, got %d", len(g.Adjacencies[u]))
		}
		if g.Degree(u) != 1 {
			t.Errorf("Self loop for directed graph degree of 1, got %d", g.Degree(u))
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

		// check for node existence
		ok := g.HasNode(w)
		if !ok {
			t.Errorf("Expected graph to have node w, got %t", ok)
		}

		// check for edge existence
		ok = g.HasEdge(u, v)
		if !ok {
			t.Errorf("Expected graph to have edge between u and v, got %t", ok)
		}

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

		// add a duplicate edge, with a different weight
		g.AddEdge(u, v, 20.0)
		// edge count should still be the same
		// check edge count
		n = g.NumberOfEdges()
		if n != 2 {
			t.Errorf("Expected 2 edges after adding duplicate, got %d", n)
		}

		// and the weight should be the new value
		if g.Adjacencies[u][v] != 20.0 {
			t.Errorf("Expected new edge weight %f, got %f", 20.0, g.Adjacencies[u][v])
		}

		// add a duplicate node
		g.AddNode(u)
		// node count should remain unchanged
		n = g.NumberOfNodes()
		if n != 5 {
			t.Errorf("Expected 5 nodes, got %d", n)
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

		// remove a node that doesn't exist
		// should not panic, node count should remain the same
		g.RemoveNode(z)
		n = g.NumberOfNodes()
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

		// remove an edge that doesn't exist, shouldn't affect edge count
		g.RemoveEdge(z, u)
		n = g.NumberOfEdges()
		if n != 3 {
			t.Errorf("Expected 6 edges, got %d", n)
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

		// remove node v
		g.RemoveNode(v)
		// check edge count, should be 0
		n = g.NumberOfEdges()
		if n != 0 {
			t.Errorf("Expected 0 edges, got %d", n)
		}
		// check that u's adjacency list is empty
		if len(g.Adjacencies[u]) != 0 {
			t.Errorf("Expected u's adjacency list to be empty, got length %d", len(g.Adjacencies[u]))
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
		expected := []Node[int]{v, w}
		ns := g.Neighbors(u)
		if !slices.Contains(ns, v) && slices.Contains(ns, w) {
			t.Errorf("Neighbors function expected %v, got %v", expected, ns)
		}
		ns = g.Successors(u)
		if !slices.Contains(ns, v) && slices.Contains(ns, w) {
			t.Errorf("Successors function expected %v, got %v", expected, ns)
		}
		n = len(g.Predecessors(u))
		if n != 0 {
			t.Errorf("Length of predecessors for u was expected to be 0, got %d", n)
		}
		n = len(g.Predecessors(v))
		if n != 1 {
			t.Errorf("Length of predecessors for v was expected to be 1, got %d", n)
		}
	})
}

func TestGraphCopy(t *testing.T) {
	// set up a small graph
	g := NewUndirectedGraph[string]()
	u := Node[string]{ID: "A"}
	v := Node[string]{ID: "B"}
	g.AddEdge(u, v, 1.0)

	// make a copy of it
	h := g.Copy()

	// test structural consistency
	if len(h.Adjacencies) != len(g.Adjacencies) {
		t.Errorf("Size mismatch: expected %d nodes, got %d", len(g.Adjacencies), len(h.Adjacencies))
	}

	for oldNode, oldNeighbors := range g.Adjacencies {
		newNode := Node[string]{ID: oldNode.ID}
		newNeighbors, exists := h.Adjacencies[newNode]
		if !exists {
			t.Fatalf("Node %v missing in copied graph", oldNode.ID)
		}

		for oldNeighbor, oldWeight := range oldNeighbors {
			newNeighbor := Node[string]{ID: oldNeighbor.ID}
			newWeight, weightExists := newNeighbors[newNeighbor]
			if !weightExists || newWeight != oldWeight {
				t.Errorf("Edge mismatch for node %s: expected weight %f, got %f",
					oldNode.ID, oldWeight, newWeight)
			}
		}
	}

	// test referential integrity
	// in the copy, node A that connects to B must be the same
	// node A that B connects to
	// first, retrieve them
	var copyNodeA, copyNodeB Node[string]
	for n := range h.Adjacencies {
		if n.ID == "A" {
			copyNodeA = n
		}
		if n.ID == "B" {
			copyNodeB = n
		}
	}
	// now check
	foundA := false
	for neighbor := range h.Adjacencies[copyNodeB] {
		if neighbor == copyNodeA {
			foundA = true
		}
	}
	if !foundA {
		t.Error("Referential integrity failed")
	}

	// verify that changing the copy doesn't affect the original
	// change the weight in the original
	g.Adjacencies[u][v] = 10.0
	if h.Adjacencies[u][v] == 10.0 {
		t.Error("Deep independence failed")
	}
}
