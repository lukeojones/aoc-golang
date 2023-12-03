package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	input, err := readInput(2023, 3)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	lines := strings.Split(input, "\n")
	possibleGearLocations := make(map[string][]int)

	sumP1, sumP2 := 0, 0
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		re := regexp.MustCompile(`(\d+)`)
		numLocations := re.FindAllStringIndex(line, -1)
		for _, numLocation := range numLocations {
			isPart, adjChar, adjCharPos := findPartsAndAdjacentChar(lines, numLocation, i)
			partNum, _ := strconv.Atoi(line[numLocation[0]:numLocation[1]])
			if isPart { // Part 1
				sumP1 += partNum
			}
			if adjChar == '*' { // Part 2
				possibleGearLocations[adjCharPos] = append(possibleGearLocations[adjCharPos], partNum)
			}
		}
	}

	for possibleGearLocation := range possibleGearLocations {
		if len(possibleGearLocations[possibleGearLocation]) == 2 {
			sumP2 += possibleGearLocations[possibleGearLocation][0] * possibleGearLocations[possibleGearLocation][1]
		}
	}

	fmt.Println("Solution Part 1:", sumP1)
	fmt.Println("Solution Part 2:", sumP2)
}

func findPartsAndAdjacentChar(lines []string, numLocation []int, lineNum int) (bool, rune, string) {
	for col := numLocation[0]; col < numLocation[1]; col++ {
		if col > 0 && (!unicode.IsDigit(rune(lines[lineNum][col-1])) && lines[lineNum][col-1] != '.') {
			return true, rune(lines[lineNum][col-1]), fmt.Sprint(lineNum, col-1)
		}
		if col < len(lines[lineNum])-1 && (!unicode.IsDigit(rune(lines[lineNum][col+1])) && lines[lineNum][col+1] != '.') {
			return true, rune(lines[lineNum][col+1]), fmt.Sprint(lineNum, col+1)
		}
		if lineNum > 0 && (!unicode.IsDigit(rune(lines[lineNum-1][col])) && lines[lineNum-1][col] != '.') {
			return true, rune(lines[lineNum-1][col]), fmt.Sprint(lineNum-1, col)
		}
		if lineNum < len(lines)-1 && (!unicode.IsDigit(rune(lines[lineNum+1][col])) && lines[lineNum+1][col] != '.') {
			return true, rune(lines[lineNum+1][col]), fmt.Sprint(lineNum+1, col)
		}

		//do diagonals
		if lineNum > 0 && col > 0 && (!unicode.IsDigit(rune(lines[lineNum-1][col-1])) && lines[lineNum-1][col-1] != '.') {
			return true, rune(lines[lineNum-1][col-1]), fmt.Sprint(lineNum-1, col-1)
		}
		if lineNum > 0 && col < len(lines[lineNum])-1 && (!unicode.IsDigit(rune(lines[lineNum-1][col+1])) && lines[lineNum-1][col+1] != '.') {
			return true, rune(lines[lineNum-1][col+1]), fmt.Sprint(lineNum-1, col+1)
		}
		if col > 0 && lineNum < len(lines)-1 && (!unicode.IsDigit(rune(lines[lineNum+1][col-1])) && lines[lineNum+1][col-1] != '.') {
			return true, rune(lines[lineNum+1][col-1]), fmt.Sprint(lineNum+1, col-1)
		}
		if col < len(lines[lineNum])-1 && lineNum < len(lines)-1 && (!unicode.IsDigit(rune(lines[lineNum+1][col+1])) && lines[lineNum+1][col+1] != '.') {
			return true, rune(lines[lineNum+1][col+1]), fmt.Sprint(lineNum+1, col+1)
		}
	}
	return false, 'X', "XXX"
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
