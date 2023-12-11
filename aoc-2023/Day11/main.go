package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"strings"
)

func main() {
	input, err := readInput(2023, 11)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	lines := strings.Split(input, "\n")
	rows := make([][]rune, len(lines))

	var emptyRows []int
	var emptyCols []int
	var galaxies []image.Point

	for i, line := range lines {
		rows[i] = []rune(line)
	}
	for r, row := range rows {
		emptyRow := true
		for c := range row {
			if row[c] == '#' {
				galaxies = append(galaxies, image.Point{X: c, Y: r})
				emptyRow = false
			}
		}
		if emptyRow {
			emptyRows = append(emptyRows, r)
		}
	}

	for c := 0; c < len(rows[0]); c++ {
		emptyCol := true
		for r := range rows {
			if rows[r][c] == '#' {
				emptyCol = false
			}
		}
		if emptyCol {
			emptyCols = append(emptyCols, c)
		}
	}

	ans := 0
	for _, g := range galaxies {
		//println("Galaxy:", galaxies[g].X, ",", galaxies[g].Y)
		for _, o := range galaxies {
			ans += int(math.Abs(float64(g.X-o.X))) + int(math.Abs(float64(g.Y-o.Y)))
			for _, col := range emptyCols {
				if (int(math.Min(float64(g.X), float64(o.X))) < col) && (int(math.Max(float64(g.X), float64(o.X))) > col) {
					ans += 1_000_000 - 1
				}
			}
			for _, row := range emptyRows {
				if (int(math.Min(float64(g.Y), float64(o.Y))) < row) && (int(math.Max(float64(g.Y), float64(o.Y))) > row) {
					ans += 1_000_000 - 1
				}
			}
		}
	}

	// Solution here
	fmt.Println("Solution:", ans/2)
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
