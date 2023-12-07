package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var cardValues = map[rune]int{
	'A': 14,
	'K': 13,
	'Q': 12,
	//'J': 11,
	'J': 1,
	'T': 10,
	'9': 9,
	'8': 8,
	'7': 7,
	'6': 6,
	'5': 5,
	'4': 4,
	'3': 3,
	'2': 2,
}

type Hands []string

func (h Hands) Len() int {
	return len(h)
}
func (s Hands) Less(i, j int) bool {
	// Example custom logic: sort strings based on some custom rule
	// Modify this according to your specific needs
	scoreI := calcScore(s[i])
	scoreJ := calcScore(s[j])
	if scoreI == scoreJ {
		for k := 0; k < 5; k++ {
			rI := rune(s[i][k])
			rJ := rune(s[j][k])
			if rI != rJ {
				return cardValues[rI] < cardValues[rJ]
			}
		}
	}
	return scoreI < scoreJ
}
func (s Hands) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

/*
*
6.Five of a kind, where all five cards have the same label: AAAAA
5.Four of a kind, where four cards have the same label and one card has a different label: AA8AA
4.Full house, where three cards have the same label, and the remaining two cards share a different label: 23332
3.Three of a kind, where three cards have the same label, and the remaining two cards are each different from any other card in the hand: TTT98
2.Two pair, where two cards share one label, two other cards share a second label, and the remaining card has a third label: 23432
1.One pair, where two cards share one label, and the other three cards have a different label from the pair and each other: A23A4
0.High card, where all cards' labels are distinct: 23456
*/
func calcScore(hand string) int {
	runes := []rune(hand)
	counts := map[rune]int{}
	for _, r := range runes {
		counts[r]++
	}

	// Part 2
	//for int(\) counts['J'] > 0 {
	//
	//}
	wildcards := counts['J']
	for i := 0; i < wildcards; i++ {
		max := 0
		var kToIncrement rune
		for k, v := range counts {
			if k != 'J' && v > max {
				max = v
				kToIncrement = k
			}
		}
		counts[kToIncrement]++
		counts['J']--
	}

	score := 0
	hasPair := false
	hasTriple := false
	for _, v := range counts {
		if v == 5 {
			score = 6
			break
		}
		if v == 4 {
			score = 5
			break
		}
		if v == 3 {
			hasTriple = true
			score = 3
			if hasPair {
				score = 4
			}
			continue
		}
		if v == 2 {
			score = 1
			if hasTriple {
				score = 4
			} else if hasPair {
				score = 2
			}
			hasPair = true
		}
	}
	return score
}

func main() {
	input, err := readInput(2023, 7)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	hands := Hands{}
	bids := map[string]int{}
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		parts := strings.Split(line, " ")
		hand, bid := parts[0], parts[1]
		bids[hand], err = strconv.Atoi(bid)
		if err != nil {
			fmt.Printf("Failed to convert bid to int: %v\n", err)
			return
		}
		hands = append(hands, hand)
	}
	println("Length of hands:", len(hands))
	println("Length of bids:", len(bids))

	sort.Sort(hands)

	for i := 0; i < len(hands); i++ {
		println("Hand: ", hands[i], " = ", calcScore(hands[i]))
	}

	ans := 0
	for i := 0; i < len(hands); i++ {
		hand := hands[i]
		bid := bids[hand]
		ans += (i + 1) * bid
		//fmt.Println(hand, bid)

	}
	// Solution here
	fmt.Println("Solution:", ans)
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
