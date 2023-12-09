package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	input, err := readInput(2023, 9)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	lines := strings.Split(input, "\n")
	var sequences [][]int
	for _, line := range lines {
		var nums []int
		tokens := strings.Split(line, " ")
		for _, numToken := range tokens {
			num, err := strconv.Atoi(numToken)
			if err != nil {
				log.Fatal(err)
			}
			nums = append(nums, num)
		}
		sequences = append(sequences, nums)
	}

	part1 := solve(sequences)
	fmt.Println("Part 1:", part1)

	//Part Two - Reverse and Solve
	for _, sequence := range sequences {
		slices.Reverse(sequence)
	}
	part2 := solve(sequences)
	fmt.Println("Part 2:", part2)
}

func solve(sequences [][]int) int {
	lastNum := 0
	for _, sequence := range sequences {
		//fmt.Println(sequence)
		complete := false
		var diffs []int
		for !complete && len(sequence) > 0 {
			complete = true
			diffs = append(diffs, sequence[len(sequence)-1])
			nextSeq := make([]int, len(sequence)-1)
			for i := 0; i < len(sequence)-1; i++ {
				nextSeq[i] = sequence[i+1] - sequence[i]
				if nextSeq[i] != 0 {
					complete = false
				}
			}
			sequence = nextSeq
		}
		for _, diff := range diffs {
			lastNum += diff
		}
	}
	return lastNum
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
