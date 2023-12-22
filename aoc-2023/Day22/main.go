package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
1,0,1~1,2,1
0,0,2~2,0,2
0,2,3~2,2,3
0,0,4~0,2,4
2,0,5~2,2,5
0,1,6~2,1,6
1,1,8~1,1,9
*/
func main() {
	input, err := readInput(2023, 22)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	var bricks []Brick

	//1,1,8~1,1,9
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		brickstr := strings.Split(line, "~")
		brick1 := strings.Split(brickstr[0], ",")
		brick2 := strings.Split(brickstr[1], ",")

		brick := Brick{
			p1: image.Point{X: toInt(brick1[0]), Y: toInt(brick1[1])},
			z1: toInt(brick1[2]),
			p2: image.Point{X: toInt(brick2[0]), Y: toInt(brick2[1])},
			z2: toInt(brick2[2]),
		}
		bricks = append(bricks, brick)

	}
	//fmt.Println(bricks)

	// Sort bricks by lowest z first
	sort.Slice(bricks, func(i, j int) bool {
		minZi := min(bricks[i].z1, bricks[i].z2)
		minZj := min(bricks[j].z1, bricks[j].z2)
		return minZi < minZj
	})

	var brickToBricksItSupports = make(map[Brick][]Brick)
	var bricksToBricksItIsSupportedBy = make(map[Brick][]Brick)

	// Now slide all bricks down until they hit the bottom or something else
	for bi, brick := range bricks {
		supported := false
		for (brick.z1 > 1 && brick.z2 > 1) && !supported {
			// Slide down
			brick.z1--
			brick.z2--
			bricks[bi] = brick

			// Check if we overlap with any other bricks
			needUndo := false
			for bu := bi - 1; bu >= 0; bu-- {
				under := bricks[bu]
				if bi == bu {
					continue
				}
				if doOverlap(brick, under) {
					// Undo the move and record dependency
					brickToBricksItSupports[under] = append(brickToBricksItSupports[under], brick)
					bricksToBricksItIsSupportedBy[brick] = append(bricksToBricksItIsSupportedBy[brick], under)
					needUndo = true
				}
			}

			if needUndo {
				brick.z1++
				brick.z2++
				bricks[bi] = brick
				supported = true
			}
		}
	}

	// Just for debugging
	sort.Slice(bricks, func(i, j int) bool {
		minZi := min(bricks[i].z1, bricks[i].z2)
		minZj := min(bricks[j].z1, bricks[j].z2)
		return minZi < minZj
	})

	ansP1 := 0

	for _, brickToRemove := range bricks {
		bricksItSupports := brickToBricksItSupports[brickToRemove]
		if bricksItSupports == nil || len(bricksItSupports) == 0 {
			ansP1++
			continue
		}

		allOtherBricksSupported := true
		for _, brickItSupports := range bricksItSupports {
			hasOtherSupport := false
			otherSupports := bricksToBricksItIsSupportedBy[brickItSupports]
			for _, otherSupport := range otherSupports {
				if otherSupport != brickToRemove {
					hasOtherSupport = true
				}
			}

			if !hasOtherSupport {
				allOtherBricksSupported = false
			}
		}

		if allOtherBricksSupported {
			ansP1++
		}
	}

	// Part 2
	ansP2 := 0
	for bi, brick := range bricks {
		println("Brick", bi)
		ansBrick := 0
		removed := make(map[Brick]bool)

		unprocessed := []Brick{brick}
		removed[brick] = true
		for len(unprocessed) > 0 {
			toDisintegrate := unprocessed[0]
			unprocessed = unprocessed[1:]
			bricksItSupports := brickToBricksItSupports[toDisintegrate]

			for _, brickItSupports := range bricksItSupports {
				hasOtherSupport := false
				otherSupports := bricksToBricksItIsSupportedBy[brickItSupports]
				for _, otherSupport := range otherSupports {
					if !removed[otherSupport] {
						hasOtherSupport = true
					}
				}

				if !hasOtherSupport && removed[brickItSupports] == false {
					ansBrick++
					removed[brickItSupports] = true
					unprocessed = append([]Brick{brickItSupports}, unprocessed...)
				}
			}
		}
		println(ansBrick)
		ansP2 += ansBrick
	}

	fmt.Println("Part 2:", ansP2) //1394 is too low

	// Solution here
	fmt.Println("Solution:", ansP1) //57008, 549 is too high
}

type Brick struct {
	p1, p2 image.Point
	z1, z2 int
}

func doOverlap(top, bottom Brick) bool {
	x1Min, x1Max := minMax(top.p1.X, top.p2.X)
	y1Min, y1Max := minMax(top.p1.Y, top.p2.Y)
	z1Min, z1Max := minMax(top.z1, top.z2)

	x2Min, x2Max := minMax(bottom.p1.X, bottom.p2.X)
	y2Min, y2Max := minMax(bottom.p1.Y, bottom.p2.Y)
	z2Min, z2Max := minMax(bottom.z1, bottom.z2)

	//zmin := min(top.z1, top.z2)dd

	// Check if one rectangle is on the left side of the other
	if x1Max < x2Min || x2Max < x1Min {
		return false
	}

	// Check if one rectangle is above the other
	if y1Max < y2Min || y2Max < y1Min {
		return false
	}

	if z1Max < z2Min || z2Max < z1Min {
		return false
	}

	// Rectangles overlap
	return true
}

func minMax(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Boilerplate: ReadInput reads the input file for the given year and day
func readInput(year, day int) (string, error) {
	filePath := fmt.Sprintf("aoc-%d/Day%d/input.test.txt", year, day)
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return string(file), nil
}
