package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, err := readInput(2023, 6)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	lines := strings.Split(input, "\n")
	times := strings.Fields(lines[0])[1:]
	dists := strings.Fields(lines[1])[1:]

	// Part 1
	ans := 1
	for i := 0; i < len(times); i++ {
		time, err := strconv.Atoi(times[i])
		dist, err := strconv.Atoi(dists[i])
		if err != nil {
			fmt.Printf("Failed to convert time to int: %v\n", err)
			return
		}
		ans *= calculateWaysToBeat(time, dist)
	}

	// Part 2

	fmt.Println("Solution:", ans)
}

func calculateWaysToBeat(time, record int) int {
	waysToBeat := 0
	for h := 0; h <= time; h++ {
		d := h * (time - h)
		if d > record {
			waysToBeat++
		}
	}
	return waysToBeat
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
