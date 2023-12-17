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

	//visited := make([][]string, len(lines))
	//memo := make(map[image.Point][]image.Point)
	grid := make([][]string, len(lines))
	for i, line := range lines {
		grid[i] = strings.Split(line, "")
		//visited[i] = strings.Split(strings.Repeat(".", len(line)), "")
	}
	//PrintGrid(grid, "Input")
	//PrintGrid(visited, "Visited")

	ansP1 := solve([]Beam{{Pos: image.Point{}, Dir: RIGHT}}, grid)

	// Part 1
	//ansP1 := len(memo)
	ans := 0
	for r := 0; r < len(grid); r++ {
		startingBeams := []Beam{{Pos: image.Point{X: 0, Y: r}, Dir: RIGHT}}
		temp := solve(startingBeams, grid)
		if temp > ans {
			ans = temp
		}

		startingBeams = []Beam{{Pos: image.Point{X: len(grid[0]) - 1, Y: r}, Dir: LEFT}}
		temp = solve(startingBeams, grid)
		if temp > ans {
			ans = temp
		}
	}
	for c := 0; c < len(grid[0]); c++ {
		startingBeams := []Beam{{Pos: image.Point{X: c, Y: 0}, Dir: DOWN}}
		temp := solve(startingBeams, grid)
		if temp > ans {
			ans = temp
		}

		startingBeams = []Beam{{Pos: image.Point{X: c, Y: len(grid) - 1}, Dir: UP}}
		temp = solve(startingBeams, grid)
		if temp > ans {
			ans = temp
		}

	}

	// Solution here
	fmt.Println("Solution:", ansP1)
	fmt.Println("Solution:", ans)
}

func solve(startingBeams []Beam, grid [][]string) int {
	memo := make(map[image.Point][]image.Point)
	//beams := []Beam{{Pos: image.Point{}, Dir: RIGHT}} // Start with single beam going right
	beams := []Beam{}
	beams = append(beams, startingBeams...)
	i := 0
	for len(beams) > 0 {
		beam := beams[0]
		beams = beams[1:]

		dirs, exists := memo[beam.Pos]
		if exists && containsDir(dirs, beam.Dir) {
			continue
		}

		// Stop at edges
		if beam.Pos.X < 0 || beam.Pos.X >= len(grid[0]) || beam.Pos.Y < 0 || beam.Pos.Y >= len(grid) {
			continue
		}

		// Continue at dots
		if grid[beam.Pos.Y][beam.Pos.X] == "." {
			//visited[beam.Pos.Y][beam.Pos.X] = "#"
			memo[beam.Pos] = append(memo[beam.Pos], beam.Dir)
			beam.Pos = beam.Pos.Add(beam.Dir)
			beams = append(beams, beam)
			continue
		}

		//visited[beam.Pos.Y][beam.Pos.X] = "#"
		memo[beam.Pos] = append(memo[beam.Pos], beam.Dir)
		changes := GetBeams(beam.Dir, rune(grid[beam.Pos.Y][beam.Pos.X][0]))
		for _, change := range changes {
			newBeam := Beam{Pos: beam.Pos.Add(change), Dir: change}
			beams = append(beams, newBeam)
		}
		//PrintGrid(visited, fmt.Sprint("Visited after iteration: ", i))
		i++
	}
	return len(memo)
}

type Beam struct {
	Pos image.Point
	Dir image.Point
}

func containsDir(dirs []image.Point, dir image.Point) bool {
	for _, d := range dirs {
		if d == dir {
			return true
		}
	}
	return false
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
	filePath := fmt.Sprintf("aoc-%d/Day%d/input.txt", year, day)
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
