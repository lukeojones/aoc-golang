package main

import (
	"github.com/lukeojones/aoc-golang/int/boilerplate"
	"github.com/lukeojones/aoc-golang/int/download"
	"log"
)

func main() {
	err := download.GetPuzzlePart(2023, 16)
	if err != nil {
		log.Fatal(err)
	}
	err = boilerplate.GenerateBoilerplate(2023, 16)
	if err != nil {
		log.Fatal(err)
	}
}
