package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Interval struct {
	min   int
	max   int
	valid bool
}

func NewInterval(min, max int) Interval {
	return Interval{min, max, true}
}
func EmptyInterval() Interval {
	return Interval{-1, -1, false}
}

type NutrientTransform struct {
	interval Interval
	offset   int
}

var maps = map[string][]NutrientTransform{}

func main() {
	input, err := readInput(2023, 5)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	parts := strings.Split(input, "\n\n")

	seeds := strings.Fields(strings.TrimPrefix(parts[0], "seeds: "))

	locations := make([]int, 0)
	for _, seed := range seeds {
		locations = append(locations, toInt(seed))
	}

	var blockNames []string
	for i := 1; i < len(parts); i++ {
		block := strings.Split(parts[i], "\n")
		blockName := strings.TrimSuffix(block[0], " map:")
		blockNames = append(blockNames, blockName)
		for t := 1; t < len(block); t++ {
			transform := strings.Split(block[t], " ")
			dest, source, length := toInt(transform[0]), toInt(transform[1]), toInt(transform[2])
			interval := NewInterval(source, source+length)
			maps[blockName] = append(maps[blockName], NutrientTransform{interval, dest - source})
		}
		sort.Slice(maps[blockName], func(i, j int) bool {
			return maps[blockName][i].interval.min < maps[blockName][j].interval.min
		})
	}

	partOne(locations, blockNames)
	partTwo(locations, blockNames)
}

func partOne(locations []int, blockNames []string) {
	var destinations = make([]int, len(locations))
	for seed, seedTemp := range locations {
		for _, bn := range blockNames {
			seedTemp = applySimpleTransform(seedTemp, bn)
		}
		destinations[seed] = seedTemp
	}
	ans := minArray(destinations)
	println("Solution Part 1:", ans)
}

func partTwo(locations []int, blockNames []string) {
	var seedIntervals []Interval
	for i := 0; i < len(locations)-1; i = i + 2 {
		seedMin, seedRange := locations[i], locations[i+1]
		interval := NewInterval(seedMin, seedMin+seedRange)
		seedIntervals = append(seedIntervals, interval)
	}

	for _, blockName := range blockNames {
		result := make([]Interval, 0)
		for si := 0; si < len(seedIntervals); si++ {
			seedInterval := seedIntervals[si]
			processed := applyBlock(seedInterval, maps[blockName], blockName)
			result = append(result, processed...)
		}
		seedIntervals = result
	}

	ans := math.MaxInt
	for _, interval := range seedIntervals {
		if interval.min < ans {
			ans = interval.min
		}
	}
	println("Solution Part 2:", ans)
}

func applyBlock(seedInterval Interval, nutrientTransformBlock []NutrientTransform, blockName string) []Interval {
	//println("Processing ", len(nutrientTransformBlock), " for ", blockName)
	mapped := make([]Interval, 0)        //those where a nutrientTransformName has been found
	unmapped := []Interval{seedInterval} //those where a nutrientTransformName has been found - this gets very big

	for t := 0; t < len(nutrientTransformBlock); t++ {
		m, u := applyTransform(nutrientTransformBlock[t], unmapped)
		mapped = append(mapped, m...)
		unmapped = u
		if len(unmapped) == 0 {
			return mapped
		}
	}
	return append(mapped, unmapped...)
}

func applyTransform(t NutrientTransform, seedIntervals []Interval) (mapped, unmapped []Interval) {
	mapped = make([]Interval, 0)
	unmapped = make([]Interval, 0)
	for _, seedInterval := range seedIntervals {
		left, right, overlap := overlap(seedInterval, t.interval)
		if overlap.valid {
			shifted := NewInterval(overlap.min+t.offset, overlap.max+t.offset)
			mapped = append(mapped, shifted)
		}
		if left.valid {
			unmapped = append(unmapped, left)
		}
		if right.valid {
			unmapped = append(unmapped, right)
		}
	}
	return mapped, unmapped
}

func overlap(seedInterval Interval, nutrientInterval Interval) (left, right, overlap Interval) {
	if seedInterval.max <= nutrientInterval.min {
		return seedInterval, EmptyInterval(), EmptyInterval()
	}

	if seedInterval.min >= nutrientInterval.max {
		return EmptyInterval(), seedInterval, EmptyInterval()
	}

	if seedInterval.min < nutrientInterval.min {
		left = NewInterval(seedInterval.min, nutrientInterval.min)
	}

	if seedInterval.max > nutrientInterval.max {
		right = NewInterval(nutrientInterval.max, seedInterval.max)
	}

	overlapLow := maxInt(seedInterval.min, nutrientInterval.min)
	overlapHigh := minInt(seedInterval.max, nutrientInterval.max)
	overlap = NewInterval(overlapLow, overlapHigh)

	return left, right, overlap
}

func applySimpleTransform(input int, mapping string) int {
	for _, Range := range maps[mapping] {
		if input >= Range.interval.min && input <= Range.interval.max {
			return input + Range.offset
		}
	}
	return input
}

// Utility Stuff
func maxInt(a, b int) int {
	if a < b {
		return b
	}
	return a
}
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func minArray(arr []int) int {
	min := math.MaxInt
	for _, v := range arr {
		if v < min {
			min = v
		}
	}
	return min
}
func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
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
