package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

/*
Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
*/
func main() {
	input, err := readInput(2023, 4)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	cardCounts := make(map[int]int)
	lines := strings.Split(input, "\n")
	sumP1 := 0
	sumP2 := 0
	for cardNumber, line := range lines {
		content := strings.Split(line, ":")[1]
		sides := strings.Split(content, "|")
		winnersText, ticketText := sides[0], sides[1]
		winners := stringsToInts(strings.Split(winnersText, " "))
		tickets := stringsToInts(strings.Split(ticketText, " "))

		// Part 1
		matches := intersection(winners, tickets)
		numMatches := len(matches)
		fmt.Printf("Card %d has %d intersections: %v\n", cardNumber, numMatches, matches)
		sumP1 += powInt(2, numMatches-1)

		// Part 2
		cardCounts[cardNumber] = cardCounts[cardNumber] + 1
		for j := 0; j < numMatches; j++ {
			cardCounts[cardNumber+j+1] = cardCounts[cardNumber+j+1] + cardCounts[cardNumber]
		}
	}

	for card, count := range cardCounts {
		//Invariant - "Cards will never make you copy a card past the end of the table"
		if card > len(lines) {
			println(card, count)
			log.Panic("Card number is too high")
		}

		sumP2 += count
	}

	// Solution here
	fmt.Println("Solution:", sumP1)
	fmt.Println("Solution:", sumP2)
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

func stringsToInts(strs []string) []int {
	var ints []int
	for _, str := range strs {
		trimmed := strings.TrimSpace(str)
		num, err := strconv.Atoi(trimmed)
		if err != nil {
			continue
		}
		ints = append(ints, num)
	}
	return ints
}

func intersection(list1, list2 []int) []int {
	elements := make(map[int]bool)
	for _, item := range list1 {
		elements[item] = true
	}

	var intersect []int
	for _, item := range list2 {
		if _, found := elements[item]; found {
			intersect = append(intersect, item)
			delete(elements, item)
		}
	}
	return intersect
}

func powInt(base, exponent int) int {
	result := math.Pow(float64(base), float64(exponent))
	return int(result)
}
