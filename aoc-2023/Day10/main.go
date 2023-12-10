package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"strings"
)

type Dir struct {
	dr int // delta row - -1 = up, 0 = same, 1 = down
	dc int // delta column - -1 = left, 0 = same, 1 = right
}

var dirs = []Dir{
	{-1, 0}, // up
	{1, 0},  // down
	{0, -1}, // left
	{0, 1},  // right
}

func main() {
	input, err := readInput(2023, 10)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	grid := make([][]string, 0)
	rows := strings.Split(input, "\n")
	for i := 0; i < len(rows); i++ {
		grid = append(grid, strings.Split(rows[i], ""))
	}

	var dir Dir
	var pipe string
	var pr, pc int
	var sr, sc int
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[r]); c++ {
			if grid[r][c] == "S" {
				//fmt.Println("Found S at", r, c)
				pr = r
				sr = pr
				pc = c
				sc = pc
				pipe, dir = IdentifyStartDirection(grid, r, c)
				//fmt.Println("Start direction:", pipe, dir.dr, dir.dc)
				grid[pr][pc] = pipe
			}
		}
	}

	pipeline := make([]image.Point, 0)
	dist := 0
	for {
		pipeline = append(pipeline, image.Point{X: pc, Y: pr})
		dist++
		pr = pr + dir.dr
		pc = pc + dir.dc
		if grid[pr][pc] == "L" {
			//find entry dir to L and modify
			if dir.dr == 1 {
				dir = Dir{0, 1}
			} else if dir.dc == -1 {
				dir = Dir{-1, 0}
			}
		}
		if grid[pr][pc] == "J" {
			//find entry dir to J and modify
			if dir.dr == 1 {
				dir = Dir{0, -1}
			} else if dir.dc == 1 {
				dir = Dir{-1, 0}
			}
		}
		if grid[pr][pc] == "7" {
			//find entry dir to 7 and modify
			if dir.dr == -1 {
				dir = Dir{0, -1}
			} else if dir.dc == 1 {
				dir = Dir{1, 0}
			}
		}
		if grid[pr][pc] == "F" {
			//find entry dir to F and modify
			if dir.dr == -1 {
				dir = Dir{0, 1}
			} else if dir.dc == -1 {
				dir = Dir{1, 0}
			}
		}
		if grid[pr][pc] == "|" || grid[pr][pc] == "-" {
			//continue
		}
		if pr == sr && pc == sc {
			break
		}
	}
	//Part 1
	ans1 := dist / 2
	fmt.Println("Solution Part 1: ", ans1)

	// Part 2
	// Continuous Area via Shoelace algorithm (https://www.101computing.net/the-shoelace-algorithm/)
	cArea := polygonContinuousArea(pipeline)

	// Inner Area via re-arranging Pick's theorem (https://en.wikipedia.org/wiki/Pick%27s_theorem)
	ans2 := cArea - ans1 + 1
	fmt.Println("Solution Part 2:", ans2)
}

/*
def polygonArea(vertices):

	#A function to apply the Shoelace algorithm
	numberOfVertices = len(vertices)
	sum1 = 0
	sum2 = 0

	for i in range(0,numberOfVertices-1):
	  sum1 = sum1 + vertices[i][0] *  vertices[i+1][1]
	  sum2 = sum2 + vertices[i][1] *  vertices[i+1][0]

	#Add xn.y1
	sum1 = sum1 + vertices[numberOfVertices-1][0]*vertices[0][1]
	#Add x1.yn
	sum2 = sum2 + vertices[0][0]*vertices[numberOfVertices-1][1]

	area = abs(sum1 - sum2) / 2
	return area
*/
func polygonContinuousArea(vertices []image.Point) int {
	sum1 := 0
	sum2 := 0
	for i := 0; i < len(vertices)-1; i++ {
		sum1 = sum1 + vertices[i].X*vertices[i+1].Y
		sum2 = sum2 + vertices[i].Y*vertices[i+1].X
	}
	sum1 = sum1 + vertices[len(vertices)-1].X*vertices[0].Y
	sum2 = sum2 + vertices[0].X*vertices[len(vertices)-1].Y
	area := int(math.Abs(float64(sum1)-float64(sum2))) / 2
	return area
}

/*
*
| is a vertical pipe connecting north and south.
- is a horizontal pipe connecting east and west.
L is a 90-degree bend connecting north and east.
J is a 90-degree bend connecting north and west.
7 is a 90-degree bend connecting south and west.
F is a 90-degree bend connecting south and east.
*/
func IdentifyStartDirection(grid [][]string, r, c int) (string, Dir) {
	//Possible UP
	pup := r > 0 && (grid[r-1][c] == "|" || grid[r-1][c] == "7" || grid[r-1][c] == "F")
	//Possible DOWN
	pdown := r < len(grid) && (grid[r+1][c] == "|" || grid[r+1][c] == "L" || grid[r+1][c] == "J")
	//Possible LEFT
	pleft := c > 0 && (grid[r][c-1] == "-" || grid[r][c-1] == "F" || grid[r][c-1] == "L")
	//Possible RIGHT
	pright := c < len(grid[r]) && (grid[r][c+1] == "-" || grid[r][c+1] == "7" || grid[r][c+1] == "J")

	if pup && pdown {
		return "|", Dir{1, 0}
		// "S" is a vertical pipe connecting north and south : |
	} else if pleft && pright {
		return "-", Dir{0, 1}
		// "S" is a horizontal pipe connecting east and west : -
	} else if pup && pright {
		return "L", Dir{0, 1}
		// "S" is a 90-degree bend connecting north and east : L
	} else if pup && pleft {
		return "J", Dir{-1, 0}
		// "S" is a 90-degree bend connecting north and west : J
	} else if pdown && pleft {
		return "7", Dir{0, -1}
		// "S" is a 90-degree bend connecting south and west : 7
	} else if pdown && pright {
		return "F", Dir{1, 0}
		// "S" is a 90-degree bend connecting south and east : F
	}

	panic("Could not identify start direction")
	return "X", Dir{0, 0}
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
