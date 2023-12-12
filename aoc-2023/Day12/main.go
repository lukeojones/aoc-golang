package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var memo map[MemoKey]int

func main() {
	input, err := readInput(2023, 12)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	lines := strings.Split(input, "\n")
	ans1, _ := 0, 0
	for _, line := range lines {
		parts := strings.Split(line, " ")
		records, spans := strings.Split(parts[0], ""), parts[1]
		memo = make(map[MemoKey]int)
		ans1 += solve(records, strToInts(spans, ","))
	}
	fmt.Println("Solution Part 1:", ans1)

	//Part 2
	ans2 := 0
	for _, line := range lines {
		parts := strings.Split(line, " ")
		records, spans := strings.Split(parts[0], ""), strToInts(parts[1], ",")

		var largeRecords []string
		var largeSpans []int
		repeats := 5
		for i := 0; i < repeats; i++ {
			largeRecords = append(largeRecords, records...)
			largeSpans = append(largeSpans, spans...)
			if i < repeats-1 {
				largeRecords = append(largeRecords, "?")
			}
		}
		memo = make(map[MemoKey]int)
		ans2 += solve(largeRecords, largeSpans)
	}
	fmt.Println("Solution Part 2:", ans2)
}

func strToInts(str string, delim string) []int {
	var result []int
	for _, s := range strings.Split(str, delim) {
		num, err := strconv.Atoi(s)
		if err != nil {
			continue
		}
		result = append(result, num)
	}
	return result
}

func solve(records []string, spans []int) int {
	return f(0, 0, 0, records, spans)
}

type MemoKey struct {
	recordIndex, spanIndex, spanProgress int
}

func f(recordIndex, spanIndex, spanProgress int, records []string, spans []int) int {

	//Part 2 - Memoization
	if ans, exists := memo[MemoKey{recordIndex, spanIndex, spanProgress}]; exists {
		return ans
	}

	// Base Case - when we've reached the end of the records
	ans := 0
	if recordIndex == len(records) {
		if spanIndex == len(spans) && spanProgress == 0 {
			// All spans have been processed and nothing in progress
			ans = 1
			memo[MemoKey{recordIndex, spanIndex, spanProgress}] = ans
			return ans
		} else if spanIndex == len(spans)-1 && spanProgress == spans[spanIndex] {
			// One span is left but we're about to complete it
			ans = 1
			memo[MemoKey{recordIndex, spanIndex, spanProgress}] = ans
			return ans
		} else {
			ans = 0
			memo[MemoKey{recordIndex, spanIndex, spanProgress}] = ans
			return ans
		}
	}

	// Recursive Case - when we're still processing records
	switch records[recordIndex] {
	case ".":
		ans += handleOperational(recordIndex, spanIndex, spanProgress, records, spans)
	case "#":
		ans += handleDamaged(recordIndex, spanIndex, spanProgress, records, spans)
	case "?":
		// Handle Unknown Springs by forking down the operational and damaged paths
		ans += handleDamaged(recordIndex, spanIndex, spanProgress, records, spans)
		ans += handleOperational(recordIndex, spanIndex, spanProgress, records, spans)
	}
	memo[MemoKey{recordIndex, spanIndex, spanProgress}] = ans
	return ans
}

/**
 * Handle Damaged Springs: Increment span progress and record pointers
 */
func handleDamaged(recordIndex int, spanIndex int, spanProgress int, records []string, spans []int) int {
	return f(recordIndex+1, spanIndex, spanProgress+1, records, spans)
}

/*
  - Handle Operational Springs
    If we've not yet started a span, just move to the next record
    If we've completed a span, move to the next record, move to the next span and reset span progress
*/
func handleOperational(recordIndex int, spanIndex int, spanProgress int, records []string, spans []int) int {
	if spanProgress == 0 {
		return f(recordIndex+1, spanIndex, 0, records, spans)
	} else if spanProgress > 0 && spanIndex < len(spans) && spans[spanIndex] == spanProgress {
		return f(recordIndex+1, spanIndex+1, 0, records, spans)
	}
	return 0
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
