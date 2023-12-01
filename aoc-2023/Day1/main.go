package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var mapDigits map[string]int = map[string]int{"one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9}

func main() {
	input, err := readInput(2023, 1)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	//Solution
	lines := strings.Split(input, "\n")
	sum := 0
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		sum = sum + calibrate(line)
	}
	print(sum)
}

func calibrate(row string) int {
	leftIndex, leftNum := -1, -1
	for firstIndex, char := range row {
		if unicode.IsDigit(char) {
			leftIndex = firstIndex
			leftNum, _ = strconv.Atoi(string(char))
			break
		}
	}

	rightIndex, rightNum := -1, -1
	for i := len(row) - 1; i >= 0; i-- {
		char := row[i]
		if unicode.IsDigit(rune(char)) {
			rightIndex = i
			rightNum, _ = strconv.Atoi(string(char))
			break
		}
	}

	// Part 2
	for str, _ := range mapDigits {
		firstIndexForStr := strings.Index(row, str)
		lastIndexForStr := strings.LastIndex(row, str)

		if firstIndexForStr != -1 && firstIndexForStr < leftIndex {
			leftIndex = firstIndexForStr
			leftNum = mapDigits[str]
		}
		if lastIndexForStr != -1 && lastIndexForStr > rightIndex {
			rightIndex = lastIndexForStr
			rightNum = mapDigits[str]
		}
	}

	return leftNum*10 + rightNum
}

// Boilerplate: ReadInput reads the input file for the given year and day
func readInput(year, day int) (string, error) {
	filePath := fmt.Sprintf("aoc-%d/Day%d/input.bak.txt", year, day)
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return string(file), nil
}
