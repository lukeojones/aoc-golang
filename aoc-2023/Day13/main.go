package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	input, err := readInput(2023, 13)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	ansP1 := solve(input, 0)
	ansP2 := solve(input, 1)

	// Solution here
	fmt.Println("Solution Part 1:", ansP1)
	fmt.Println("Solution Part 2:", ansP2)
}

func solve(input string, smudgeTolerance int) int {
	patterns := strings.Split(input, "\n\n")
	rowsAbove, colsLeft := 0, 0
	for _, pattern := range patterns {
		lines := strings.Split(pattern, "\n")

		for r := 0; r < len(lines)-1; r++ {
			diff := 0
			for rr := 0; rr < len(lines); rr++ {
				up, down := r-rr, (r+rr)+1
				// Only compare differences if both are within the bounds of the pattern
				if 0 <= up && down < len(lines) {
					diff += compareRows(lines, up, down)
				}
			}
			if diff == smudgeTolerance {
				rowsAbove += r + 1 // Line of symmetry is the forward trailing edge of row r (so add 1)
			}
		}

		for c := 0; c < len(lines[0])-1; c++ {
			diff := 0
			for cc := 0; cc < len(lines[0]); cc++ {
				left, right := c-cc, c+1+cc

				// Only compare differences if both are within the bounds of the pattern
				if 0 <= left && right < len(lines[0]) {
					diff += compareColumns(lines, left, right)
				}
			}
			if diff == smudgeTolerance {
				colsLeft += c + 1 // Line of symmetry is the forward trailing edge of column c (so add 1)
			}
		}
	}
	return rowsAbove*100 + colsLeft
}

func compareColumns(lines []string, c1, c2 int) int {
	diff := 0
	for r := 0; r < len(lines); r++ {
		if lines[r][c1] != lines[r][c2] {
			diff++
		}
	}
	return diff
}

func compareRows(lines []string, r1, r2 int) int {
	diff := 0
	for c := 0; c < len(lines[r1]); c++ {
		if lines[r1][c] != lines[r2][c] {
			diff++
		}
	}
	return diff
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

//49952 too high
