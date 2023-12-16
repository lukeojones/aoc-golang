package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"strings"
)

func main() {
	input, err := readInput(2023, 16)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}
	lines := strings.Split(input, "\n")

	visited := make([][]string, len(lines))
	grid := make([][]string, len(lines))
	for i, line := range lines {
		grid[i] = strings.Split(line, "")
		visited[i] = strings.Split(strings.Repeat(".", len(line)), "")
	}
	PrintGrid(grid, "Input")
	PrintGrid(visited, "Visited")

	// Solution here
	fmt.Println("Solution:", "?")
}

/*
/, \, |, -
*/
func GetBeams(dir image.Point, tile rune) []image.Point {
	switch tile {
	case '/':
		switch dir {
		case UP:
			return []image.Point{RIGHT}
		case RIGHT:
			return []image.Point{UP}
		case DOWN:
			return []image.Point{LEFT}
		case LEFT:
			return []image.Point{DOWN}
		}
	case '\\':
		switch dir {
		case UP:
			return []image.Point{LEFT}
		case RIGHT:
			return []image.Point{DOWN}
		case DOWN:
			return []image.Point{RIGHT}
		case LEFT:
			return []image.Point{UP}
		}
	case '|':
		switch dir {
		case UP:
			return []image.Point{UP}
		case RIGHT:
			return []image.Point{UP, DOWN}
		case DOWN:
			return []image.Point{DOWN}
		case LEFT:
			return []image.Point{UP, DOWN}
		}
	case '-':
		switch dir {
		case UP:
			return []image.Point{LEFT, RIGHT}
		case RIGHT:
			return []image.Point{RIGHT}
		case DOWN:
			return []image.Point{LEFT, RIGHT}
		case LEFT:
			return []image.Point{LEFT}
		}
	}
	return nil
}

var LEFT = image.Point{X: -1}
var RIGHT = image.Point{X: 1}
var UP = image.Point{Y: -1}
var DOWN = image.Point{Y: 1}

// Boilerplate: ReadInput reads the input file for the given year and day
func readInput(year, day int) (string, error) {
	filePath := fmt.Sprintf("aoc-%d/Day%d/input.test.txt", year, day)
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return string(file), nil
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
