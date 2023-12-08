package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type NodePair struct {
	Left  string
	Right string
}

func main() {
	input, err := readInput(2023, 8)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	parts := strings.Split(input, "\n\n")
	directions := []rune(parts[0])
	nodesMap := make(map[string]NodePair)
	var ghostStarts []string
	for _, line := range strings.Split(parts[1], "\n") {
		sides := strings.Split(line, " = ")
		pairText := strings.Split(sides[1], ", ")
		pair := NodePair{Left: strings.Trim(pairText[0], "()"), Right: strings.Trim(pairText[1], "()")}
		nodesMap[sides[0]] = pair
		if sides[0][2] == 'A' {
			ghostStarts = append(ghostStarts, sides[0])
		}
	}

	var i = 0
	var ghostsFirstZ = make([]int, len(ghostStarts))
	for !allEndInZ(ghostStarts) {
		dir := directions[i%len(directions)]
		for g, ghostStart := range ghostStarts {
			curr := ghostStart
			pair := nodesMap[curr]
			if dir == 'L' {
				curr = pair.Left
			} else if dir == 'R' {
				curr = pair.Right
			}
			ghostStarts[g] = curr
			if curr[2] == 'Z' && ghostsFirstZ[g] == 0 {
				ghostsFirstZ[g] = i + 1
			}
		}
		i++
	}

	ans := lcmArr(ghostsFirstZ)
	fmt.Println("Solution:", ans)
}

func allEndInZ(ghostPositions []string) bool {
	for _, pos := range ghostPositions {
		if pos[2] != 'Z' {
			return false
		}
	}
	return true
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

func lcmArr(nums []int) int {
	var res int = 1
	for _, num := range nums {
		res = lcm(res, num)
	}
	return res
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
