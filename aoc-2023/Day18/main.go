package main

import (
	"fmt"
	"image"
	"log"
	"math"
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
	instructions := parseInstructions(lines)
	ansP1 := solve(instructions)

	instructions = parseHexInstructions(lines)
	ansP2 := solve(instructions)

	fmt.Println("Solution to Part 1:", ansP1)
	fmt.Println("Solution to Part 2:", ansP2)
}

func parseHexInstructions(lines []string) []Instruction {
	instructions := make([]Instruction, len(lines))
	directions := []string{"R", "D", "L", "U"}
	for l, line := range lines {
		tokens := strings.Split(line, " ")
		hexCode := tokens[2]
		length := hexToDec(hexCode[2:7])
		instructions[l] = Instruction{directions[toInt(hexCode[7:8])], int(length), hexCode}
	}
	return instructions
}

func parseInstructions(lines []string) []Instruction {
	instructions := make([]Instruction, len(lines))
	for l, line := range lines {
		tokens := strings.Split(line, " ")
		instructions[l] = Instruction{tokens[0], toInt(tokens[1]), tokens[2]}
	}
	return instructions
}

func solve(instructions []Instruction) int {
	curr := image.Point{}
	var points []image.Point
	for _, instruction := range instructions {
		p := image.Point{X: curr.X, Y: curr.Y}
		switch instruction.dir {
		case "U":
			p.Y -= instruction.length
		case "D":
			p.Y += instruction.length
		case "L":
			p.X -= instruction.length
		case "R":
			p.X += instruction.length
		}
		points = append(points, p)
		curr = p
	}

	return calculatePolygonArea(points)
}

func calculatePolygonArea(vertices []image.Point) int {
	area := 0
	perimeter := 0

	for i := 0; i < len(vertices); i++ {
		j := (i + 1) % len(vertices)
		area += vertices[i].X * vertices[j].Y
		area -= vertices[j].X * vertices[i].Y

		// Calculate perimeter
		if vertices[i].X == vertices[j].X {
			perimeter += int(math.Abs(float64(vertices[j].Y - vertices[i].Y)))
		} else {
			perimeter += int(math.Abs(float64(vertices[j].X - vertices[i].X)))
		}
	}

	// Area of the polygon
	polygonArea := int(math.Abs(float64(area)) / 2.0)
	return polygonArea + (perimeter / 2) + 1
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func hexToDec(s string) int {
	i, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		panic(err)
	}
	return int(i)
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
