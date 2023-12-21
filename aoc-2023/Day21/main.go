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
	ans := bfs(grid, start, 64)

	PrintGrid(grid, "Input")

	// Solution here
	fmt.Println("Solution:", ans)
}

type PathPoint struct {
	Pos    image.Point
	length int
}

func bfs(grid [][]string, start image.Point, steps int) int {
	unprocessed := []PathPoint{{start, 0}}
	plots := make(map[image.Point]PathPoint, 17161)
	seen := make(map[PathPoint]bool)
	//var paths []PathPoint
	for len(unprocessed) > 0 {
		curr := unprocessed[0]
		unprocessed = unprocessed[1:]

		_, exists := seen[curr]
		if exists {
			continue
		}

		if (grid[curr.Pos.Y][curr.Pos.X] == "." || grid[curr.Pos.Y][curr.Pos.X] == "S") && curr.length == steps {
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
			if next.X < 0 || next.X >= len(grid[0]) || next.Y < 0 || next.Y >= len(grid) {
				continue
			}
			if grid[next.Y][next.X] == "#" {
				continue
			}
			unprocessed = append(unprocessed, PathPoint{next, curr.length + 1})
		}

		seen[curr] = true
	}

	return len(plots)
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
