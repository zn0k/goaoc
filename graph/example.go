package graph

import (
	"fmt"
	"os"
	"strings"
)

// coordinates have X and Y components
type coordinate struct {
	X, Y int
}

// they implement Stringer for easy printing
func (c coordinate) String() string {
	return fmt.Sprintf("(%d, %d)", c.X, c.Y)
}

// directions for walking a grid are really just coordinates:
// {0, 1}, {0, -1}, {1, 0}, {-1, 0}
type Direction coordinate

// read in the maze grid and return an undirected graph as well as the start
// and end tile on the grid
func readLines(fname string, directions []Direction) (*UndirectedGraph[coordinate], Node[coordinate], Node[coordinate]) {
	buf, err := os.ReadFile(fname)
	if err != nil {
		panic(fmt.Sprintf("unable to open %s for reading", fname))
	}

	// initialize the start and end tiles, and the grid as 2d runes
	var start, target coordinate
	var grid [][]rune

	// walk the lines
	for y, line := range strings.Split(string(buf), "\n") {
		// walk each row
		var row []rune
		for x, c := range line {
			// check if we're on the start or end tile.
			// if so, record it, and then turn it into a normal tile
			if c == 'S' {
				start = coordinate{x, y}
				c = '.'
			}
			if c == 'T' {
				target = coordinate{x, y}
				c = '.'
			}
			// build the row
			row = append(row, c)
		}
		// build the grid from rows
		grid = append(grid, row)
	}

	// initialize a new graph
	g := NewUndirectedGraph[coordinate]()

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
					u := Node[coordinate]{coordinate{x, y}}
					v := Node[coordinate]{coordinate{new_x, new_y}}
					// add an edge between them, which also adds the nodes
					g.AddEdge(u, v, 1.0)
				}
			}
		}
	}

	return g, Node[coordinate]{start}, Node[coordinate]{target}
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
