package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	input, err := readInput(2023, 14)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	grid := make([][]string, 0)
	lines := strings.Split(input, "\n")
	for _, lines := range lines {
		grid = append(grid, strings.Split(lines, ""))
	}

	loadsP1 := run(grid, 1, 1)
	loadsP2 := run(grid, 4, 2000)

	fmt.Println("Solution to Part 1:", loadsP1[0])
	fmt.Println("Solution to Part 2:", findNthTerm(loadsP2, 1_000_000_000, 10))
}

func findNthTerm(seq []int, n, tolerance int) int {
	counts := make(map[int][]int)
	for i, v := range seq {
		counts[v] = append(counts[v], i)
		if len(counts[v]) > tolerance {
			wavelength := counts[v][len(counts[v])-1] - counts[v][len(counts[v])-2]
			return seq[((n-(i+1))%wavelength)+i]
		}
	}
	return -1
}

func run(grid [][]string, rotations, cycles int) []int {
	grid = copy2DStringSlice(grid)
	var loads []int
	for i := 1; i <= cycles; i++ {
		if rotations == 0 {
			roll(grid)
			loads = append(loads, calcLoad(grid))
			continue
		}
		for r := 0; r < rotations; r++ {
			roll(grid)
			grid = rotateCW(grid)
		}
		loads = append(loads, calcLoad(grid))
	}
	return loads
}

func copy2DStringSlice(src [][]string) [][]string {
	dst := make([][]string, len(src))
	for i := range src {
		dst[i] = make([]string, len(src[i]))
		copy(dst[i], src[i])
	}
	return dst
}

func calcLoad(grid [][]string) int {
	ans := 0
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[0]); c++ {
			if grid[r][c] == "O" {
				ans += len(grid) - r
			}
		}
	}
	return ans
}

func rotateCW(grid [][]string) [][]string {
	newGrid := make([][]string, len(grid[0]))
	for r := 0; r < len(grid[0]); r++ {
		newGrid[r] = make([]string, len(grid))
		for c := 0; c < len(grid); c++ {
			newGrid[r][c] = grid[len(grid)-1-c][r]
		}
	}
	return newGrid
}

func roll(grid [][]string) {
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[r]); c++ {
			if grid[r][c] == "O" {
				rd := r // destination row
				for rr := r - 1; rr >= 0; rr-- {
					if grid[rr][c] == "." {
						rd = rr
						continue
					}
					break
				}
				if rd != r {
					grid[rd][c] = grid[r][c]
					grid[r][c] = "."
				}
			}
		}
	}
}

func printGrid(grid [][]string, desc string) {
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
