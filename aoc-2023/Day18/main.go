package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"strconv"
	"strings"
)

type Instruction struct {
	dir    string
	length int
	colour string
}

func main() {
	input, err := readInput(2023, 18)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	lines := strings.Split(input, "\n")
	instructions := make([]Instruction, len(lines))

	for l, line := range lines {
		tokens := strings.Split(line, " ")
		dir := tokens[0]
		length, _ := strconv.Atoi(tokens[1])
		colour := tokens[2]
		instructions[l] = Instruction{dir, length, colour}
	}

	curr := image.Point{}
	var points []image.Point
	minX, minY := 0, 0
	for _, instruction := range instructions {
		for i := 0; i < instruction.length; i++ {
			p := image.Point{X: curr.X, Y: curr.Y}
			switch instruction.dir {
			case "U":
				p.Y -= 1
			case "D":
				p.Y += 1
			case "L":
				p.X -= 1
			case "R":
				p.X += 1
			}
			if p.X < minX {
				minX = p.X
			}
			if p.Y < minY {
				minY = p.Y
			}
			points = append(points, p)
			curr = p
		}
	}

	// Calculate Offsets
	offsetX := 0 - minX
	offsetY := 0 - minY

	maxX, maxY := 0, 0
	for p, _ := range points {
		points[p].X += offsetX
		points[p].Y += offsetY
		if points[p].X > maxX {
			maxX = points[p].X
		}
		if points[p].Y > maxY {
			maxY = points[p].Y
		}
	}

	// Create Grid
	grid := make([][]string, maxY+1)
	for r := 0; r <= maxY; r++ {
		grid[r] = make([]string, maxX+1)
		for c := 0; c <= maxX; c++ {
			if containsPoint(points, image.Point{X: c, Y: r}) {
				grid[r][c] = "#"
			} else {
				grid[r][c] = "."
			}
		}
	}

	// Expand Grid by 1 in each direction

	printGrid(grid)
	println("=============")
	grid = expandGrid(grid)
	printGrid(grid)

	// Part 1 - Flood Fill
	ansP1 := countLava(grid)
	//floodFill(grid, 0, 0)
	//printGrid(grid)

	// Solution here
	fmt.Println("Solution:", ansP1)
}

func expandGrid(grid [][]string) [][]string {
	rows := len(grid)
	cols := len(grid[0])
	newGrid := make([][]string, rows+2)

	// Add top border
	newGrid[0] = make([]string, cols+2)
	for c := range newGrid[0] {
		newGrid[0][c] = "."
	}

	// Copy the original grid and add left and right borders
	for i := range grid {
		newGrid[i+1] = make([]string, cols+2)
		newGrid[i+1][0] = "."
		copy(newGrid[i+1][1:], grid[i])
		newGrid[i+1][cols+1] = "."
	}

	// Add bottom border
	newGrid[rows+1] = make([]string, cols+2)
	for c := range newGrid[rows+1] {
		newGrid[0][c] = "."
	}

	return newGrid
}

func floodFill(grid [][]string, x, y int) {
	if x < 0 || x >= len(grid) || y < 0 || y >= len(grid[0]) || grid[x][y] != "." {
		return
	}
	grid[x][y] = "o" // Mark the filled cell with 'o'
	floodFill(grid, x+1, y)
	floodFill(grid, x-1, y)
	floodFill(grid, x, y+1)
	floodFill(grid, x, y-1)
}

func countLava(grid [][]string) int {
	floodFill(grid, 0, 0)

	area := 0
	for _, row := range grid {
		for _, cell := range row {
			if cell == "." || cell == "#" {
				area++
			}
		}
	}

	return area
}

func printGrid(grid [][]string) {
	for _, row := range grid {
		for _, col := range row {
			fmt.Print(col)
		}
		fmt.Println()
	}
}

func containsPoint(points []image.Point, p image.Point) bool {
	for _, point := range points {
		if point == p {
			return true
		}
	}
	return false
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
