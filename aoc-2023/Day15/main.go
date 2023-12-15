package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, err := readInput(2023, 15)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	steps := strings.Split(input, ",")

	// Part 1
	ansP1 := 0
	for _, step := range steps {
		ansP1 += hash(step)
	}

	// Part 2
	var boxes = make(map[int][]Lens)
	for _, step := range steps {
		parts := strings.FieldsFunc(step, func(c rune) bool {
			return c == '-' || c == '='
		})

		op := step[strings.Index(step, parts[0])+len(parts[0])]
		label := parts[0]
		boxId := hash(label)
		if op == '=' {
			val, _ := strconv.Atoi(parts[1])
			lenses, exists := boxes[boxId]
			if !exists {
				boxes[boxId] = []Lens{Lens{label, val}}
			} else {
				_, i := findLens(lenses, label)
				if i > -1 {
					lenses[i].focalLength = val
				} else {
					boxes[boxId] = append(lenses, Lens{label, val})
				}
			}
		} else if op == '-' {
			_, i := findLens(boxes[boxId], label)
			if i > -1 {
				//remove lens
				boxes[boxId] = append(boxes[boxId][:i], boxes[boxId][i+1:]...)
			}
		}
	}

	ansP2 := 0
	for boxId, lenses := range boxes {
		for lensNum, lens := range lenses {
			ansP2 += (boxId + 1) * (lensNum + 1) * lens.focalLength
		}
	}

	// Solution here
	fmt.Println("Solution Part 1:", ansP1)
	fmt.Println("Solution Part 2:", ansP2)
}

func findLens(lenses []Lens, label string) (Lens, int) {
	for i, lens := range lenses {
		if lens.label == label {
			return lens, i
		}
	}
	return Lens{}, -1
}

type Lens struct {
	label       string
	focalLength int
}

/*
* Hash of Box Label => BoxId
 */
func hash(s string) int {
	curr := 0
	for c := 0; c < len(s); c++ {
		curr += int(s[c])
		curr *= 17
		curr %= 256
	}
	return curr
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
