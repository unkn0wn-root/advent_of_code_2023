package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"log"
	"os"
)

// represent the orientation of a plane.
type Direction int

const (
	Vertical Direction = iota
	Horizontal
	Undecided
)

const Infinity = 1 << 30

// represent a 2D coordinate.
type Point struct {
	X, Y int
}

// represent a node in the graph.
type Node struct {
	Point
	Direction Direction
	Visited   bool
	HeatLoss  int

	CalculatedHeatLoss int
	Total              int
	Index              int
}

// represent the overall graph structure.
type Graph struct {
	Nodes  []Node
	Width  int
	Height int
}

// implement the heap.Interface and holds Nodes.
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int {
    return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
    return pq[i].Total < pq[j].Total
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Node)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1
	*pq = old[0 : n-1]

    return item
}
func (pq *PriorityQueue) update(item *Node, priority int) {
	heap.Fix(pq, item.Index)
}

// read the input file and returns the 2D grid.
func parseInput(filename string) ([][]int, error) {
	input, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	input = bytes.TrimSpace(input)
	lines := bytes.Split(input, []byte("\n"))
	grid := make([][]int, len(lines))
	for i := range lines {
		grid[i] = make([]int, len(lines[i]))

		for j, ch := range lines[i] {
			n := int(ch) - '0'
			grid[i][j] = n
		}
	}

	return grid, nil
}

// return the neighboring nodes for the given node.
func (g *Graph) getNeighbors(node *Node, minSteps, maxSteps int) []*Node {
	neighbors := make([]*Node, 0, 6)

	if node.Direction == Horizontal || node.Direction == Undecided {
		for heatLoss, dy := 0, 1; dy <= maxSteps; dy++ {
			neighbor := g.getNodeByCoords(node.X, node.Y+dy, Vertical)
			if neighbor != nil {
				heatLoss += neighbor.HeatLoss
				if dy >= minSteps {
					neighbor.CalculatedHeatLoss = heatLoss
					neighbors = append(neighbors, neighbor)
				}
			}
		}
		for heatLoss, dy := 0, 1; dy <= maxSteps; dy++ {
			neighbor := g.getNodeByCoords(node.X, node.Y-dy, Vertical)
			if neighbor != nil {
				heatLoss += neighbor.HeatLoss
				if dy >= minSteps {
					neighbor.CalculatedHeatLoss = heatLoss
					neighbors = append(neighbors, neighbor)
				}
			}
		}
	}

	if node.Direction == Vertical || node.Direction == Undecided {
		for heatLoss, dx := 0, 1; dx <= maxSteps; dx++ {
			neighbor := g.getNodeByCoords(node.X+dx, node.Y, Horizontal)
			if neighbor != nil {
				heatLoss += neighbor.HeatLoss
				if dx >= minSteps {
					neighbor.CalculatedHeatLoss = heatLoss
					neighbors = append(neighbors, neighbor)
				}
			}
		}
		for heatLoss, dx := 0, 1; dx <= maxSteps; dx++ {
			neighbor := g.getNodeByCoords(node.X-dx, node.Y, Horizontal)
			if neighbor != nil {
				heatLoss += neighbor.HeatLoss
				if dx >= minSteps {
					neighbor.CalculatedHeatLoss = heatLoss
					neighbors = append(neighbors, neighbor)
				}
			}
		}
	}

	return neighbors
}

// return the node at the specified coordinates.
func (g *Graph) getNodeByCoords(x, y int, direction Direction) *Node {
	if x < 0 || y < 0 || y >= g.Height || x >= g.Width {
		return nil
	}

	return &g.Nodes[y*2*g.Width+x*2+int(direction)]
}

// initialize the graph based on the input grid.
func createGraph(grid [][]int) Graph {
	graph := Graph{}
	nodes := make([]Node, 0, len(grid)*len(grid)*2)
	graph.Height = len(grid)
	graph.Width = len(grid[0])

	for y := range grid {
		for x := range grid[y] {
			nodes = append(nodes, Node{
				Point:     Point{X: x, Y: y},
				Direction: Vertical,
				Total:     Infinity,
				HeatLoss:  grid[y][x],
			})
			nodes = append(nodes, Node{
				Point:     Point{X: x, Y: y},
				Direction: Horizontal,
				Total:     Infinity,
				HeatLoss:  grid[y][x],
			})
		}
	}
	graph.Nodes = nodes

    return graph
}

// use Dijkstra's algorithm to find the shortest path.
func findShortestPath(grid [][]int, minSteps, maxSteps int) int {
	graph := createGraph(grid)
	nodes := graph.Nodes

	nodes[0].Total = 0
	nodes[0].Direction = Undecided

	priorityQueue := make(PriorityQueue, len(nodes))
	for i := range nodes {
		nodes[i].Index = i
		priorityQueue[i] = &nodes[i]
	}
	heap.Init(&priorityQueue)

	var currentNode *Node
	var endNode = &nodes[len(nodes)-1]
	for {
		currentNode = heap.Pop(&priorityQueue).(*Node)

		if currentNode.X == endNode.X && currentNode.Y == endNode.Y {
			break
		}

		currentNode.Visited = true

		for _, neighbor := range graph.getNeighbors(currentNode, minSteps, maxSteps) {
			if currentNode.Total+neighbor.CalculatedHeatLoss < neighbor.Total {
				neighbor.Total = currentNode.Total + neighbor.CalculatedHeatLoss
				priorityQueue.update(neighbor, neighbor.Total)
			}
		}
	}

	return currentNode.Total
}

func main() {
	grid, err := parseInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	part1 := findShortestPath(grid, 1, 3)
	part2 := findShortestPath(grid, 4, 10)
    fmt.Println("Part 1:", part1)
	fmt.Println("Part 2:", part2)
}
