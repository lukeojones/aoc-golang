package main

import (
	"container/heap"
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"strings"
)

type Direction int

var dirs = []image.Point{
	{0, -1}, // Up
	{1, 0},  // Right
	{0, 1},  // Down
	{-1, 0}, // Left
}
var reverseDirs = []Direction{
	Down,  // Up
	Left,  // Right
	Up,    // Down
	Right, // Left
}

const (
	Up Direction = iota
	Right
	Down
	Left
)

type Node struct {
	position   image.Point
	dist       int
	direction  Direction
	stepsInDir int
	index      int // Index in the heap
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].dist < pq[j].dist
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*Node)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*pq = old[0 : n-1]
	return node
}

func main() {
	input, err := readInput(2023, 17)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	lines := strings.Split(input, "\n")
	grid := make([][]int, len(lines))
	for i, line := range lines {
		grid[i] = make([]int, len(line))
		for j, char := range line {
			grid[i][j] = int(char - '0')
		}
	}

	ansP1 := calculateHeatLoss(grid, 0, 3)
	ansP2 := calculateHeatLoss(grid, 4, 10)

	// Solution here
	fmt.Println("Solution Part 1:", ansP1)
	fmt.Println("Solution Part 2:", ansP2)
}

func calculateHeatLoss(grid [][]int, minStepsToTake, maxStepsToTake int) int {
	openSet := &PriorityQueue{}
	heap.Init(openSet)
	startNode := &Node{position: image.Point{}, dist: 0, direction: Right, stepsInDir: 0}
	heap.Push(openSet, startNode)

	visited := make(map[image.Point]map[Direction]map[int]int)

	for openSet.Len() > 0 {
		currentNode := heap.Pop(openSet).(*Node)

		// Initialize nested maps if not already done
		if visited[currentNode.position] == nil {
			visited[currentNode.position] = make(map[Direction]map[int]int)
		}
		if visited[currentNode.position][currentNode.direction] == nil {
			visited[currentNode.position][currentNode.direction] = make(map[int]int)
		}

		_, exists := visited[currentNode.position][currentNode.direction][currentNode.stepsInDir]
		if exists {
			continue
		}

		visited[currentNode.position][currentNode.direction][currentNode.stepsInDir] = currentNode.dist

		for _, neighbor := range GetValidNeighbours(currentNode, grid, minStepsToTake, maxStepsToTake) {
			_, visited_ := visited[neighbor.position][neighbor.direction][neighbor.stepsInDir]
			if !visited_ {
				heap.Push(openSet, neighbor)
			}
		}
	}

	return GetHeatLossForCell(visited, image.Point{X: len(grid[0]) - 1, Y: len(grid) - 1})
}

func GetHeatLossForCell(visited map[image.Point]map[Direction]map[int]int, point image.Point) int {
	ans := math.MaxInt64
	for _, dirToDist := range visited[point] {
		for _, dist := range dirToDist {
			if dist < ans {
				ans = dist
			}
		}
	}
	return ans
}

func GetValidNeighbours(node *Node, grid [][]int, minStepsToTake, maxStepsToTake int) []*Node {
	neighbors := make([]*Node, 0)
	for di, d := range dirs {
		// Don't allow reverse
		if di == int(reverseDirs[node.direction]) {
			continue
		}

		newPos := node.position.Add(d)

		// Don't allow out of bounds
		if newPos.X < 0 || newPos.X >= len(grid[0]) || newPos.Y < 0 || newPos.Y >= len(grid) {
			continue
		}

		isNewDir := int(node.direction) != di
		isTargetReached := newPos.X == len(grid[0])-1 && newPos.Y == len(grid)-1
		isMinStepsTaken := node.stepsInDir >= minStepsToTake
		isMaxStepsTaken := node.stepsInDir >= maxStepsToTake

		// If can take more steps, increment cont steps and use same direction
		if !isNewDir && !isMaxStepsTaken {
			neighbors = append(neighbors, &Node{
				position:   newPos,
				dist:       node.dist + grid[newPos.Y][newPos.X],
				direction:  Direction(di),
				stepsInDir: node.stepsInDir + 1,
			})
			continue
		} else if (isNewDir || isTargetReached) && isMinStepsTaken {
			// If change in direction, reset continuous steps
			neighbors = append(neighbors, &Node{
				position:   newPos,
				dist:       node.dist + grid[newPos.Y][newPos.X],
				direction:  Direction(di),
				stepsInDir: 1,
			})
		}

	}
	return neighbors
}

// Boilerplate: ReadInput reads the input file for the given year and day
func readInput(year, day int) (string, error) {
	filePath := fmt.Sprintf("aoc-%d/Day%d/input.txt", year, day)
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return string(file), nil
}
