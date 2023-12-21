package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"strings"
)

var LEFT = image.Point{X: -1}
var RIGHT = image.Point{X: 1}
var UP = image.Point{Y: -1}
var DOWN = image.Point{Y: 1}
var DIRS = []image.Point{LEFT, RIGHT, UP, DOWN}

func main() {
	input, err := readInput(2023, 21)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	lines := strings.Split(input, "\n")
	grid := make([][]string, len(lines))
	start := image.Point{X: 0, Y: 0}
	for i, line := range lines {
		grid[i] = strings.Split(line, "")
		for j, c := range grid[i] {
			if c == "S" {
				start.X = j
				start.Y = i
			}
		}
	}

	println("Grid size:", len(grid), "x", len(grid[0]), " = ", len(grid)*len(grid[0]))

	// do bfs to find all possible paths to '.' from 'S'
	ansP1 := bfs(grid, start, 64)

	//Part 2 - Sequence Spotting for 1x, 2x and 3x step ranges sizes
	innerWidth := 65
	steps := 26501365
	ans65 := bfs(grid, start, innerWidth)
	ans186 := bfs(grid, start, 1*len(grid)+innerWidth)
	ans327 := bfs(grid, start, 2*len(grid)+innerWidth)
	a, b, c := fitQuadratic(ans65, ans186, ans327)
	n := steps / len(grid)
	ansP2 := a*n*n + b*n + c

	fmt.Println("Solution Part 1:", ansP1)
	fmt.Println("Solution Part 2:", ansP2)
}

type PathPoint struct {
	Pos    image.Point
	length int
}

func bfs(grid [][]string, start image.Point, steps int) int {
	unprocessed := []PathPoint{{start, 0}}
	plots := make(map[image.Point]PathPoint, 17161)
	seen := make(map[PathPoint]bool)
	for len(unprocessed) > 0 {
		curr := unprocessed[0]
		unprocessed = unprocessed[1:]

		_, exists := seen[curr]
		if exists {
			continue
		}

		if getTileElement(grid, curr.Pos.X, curr.Pos.Y) != "#" && curr.length == steps {
			_, exists := plots[curr.Pos]
			if !exists {
				plots[curr.Pos] = curr
			} else {
				continue
			}
		}

		if curr.length >= steps {
			continue
		}

		for _, dir := range DIRS {
			next := curr.Pos.Add(dir)
			if getTileElement(grid, next.X, next.Y) == "#" {
				continue
			}
			unprocessed = append(unprocessed, PathPoint{next, curr.length + 1})
		}

		seen[curr] = true
	}

	return len(plots)
}

func getTileElement(grid [][]string, x, y int) string {
	rows := len(grid)
	cols := len(grid[0])

	if x < 0 {
		x = (x%cols + cols) % cols
	} else {
		x = x % cols
	}

	if y < 0 {
		y = (y%rows + rows) % rows
	} else {
		y = y % rows
	}

	return grid[y][x]
}

func fitQuadratic(v0, v1, v2 int) (int, int, int) {
	a := (v0 - 2*v1 + v2) / 2
	b := (-3*v0 + 4*v1 - v2) / 2
	c := v0
	return a, b, c
}

func PrintGrid(grid [][]string, desc string) {
	fmt.Println("=======", desc, "========")
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[0]); c++ {
			fmt.Print(grid[r][c])
		}
		fmt.Println()
	}
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
