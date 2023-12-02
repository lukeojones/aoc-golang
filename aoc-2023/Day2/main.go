package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red
Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red
Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green
*/
func main() {
	input, err := readInput(2023, 2)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}
	lines := strings.Split(input, "\n")

	sumP1, sumP2 := 0, 0
	for i, game := range lines {
		if len(game) == 0 {
			continue
		}
		fmt.Println(game)

		//Part 1
		possible := true

		//Part 2
		maxRed, maxGreen, maxBlue := 0, 0, 0

		draws := strings.Split(game, ";")
		for _, draw := range draws {
			red, green, blue := countCubeColours(draw)

			//Part 1
			if red > 12 || green > 13 || blue > 14 {
				possible = false
			}

			//Part 2
			if red > maxRed {
				maxRed = red
			}
			if green > maxGreen {
				maxGreen = green
			}
			if blue > maxBlue {
				maxBlue = blue
			}
		}

		if possible {
			sumP1 += i + 1
		}
		sumP2 += maxRed * maxGreen * maxBlue
	}

	// Solution here
	fmt.Println("Solution Part 1:", sumP1)
	fmt.Println("Solution Part 2:", sumP2)
}

func countCubeColours(draw string) (int, int, int) {
	re := regexp.MustCompile(`(\d+) (red|green|blue)`)
	red, green, blue := 0, 0, 0
	matches := re.FindAllStringSubmatch(draw, -1)
	for _, match := range matches {
		fmt.Println(match)
		count, _ := strconv.Atoi(match[1])
		switch match[2] {
		case "red":
			red += count
		case "green":
			green += count
		case "blue":
			blue += count
		}
	}
	return red, green, blue
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
