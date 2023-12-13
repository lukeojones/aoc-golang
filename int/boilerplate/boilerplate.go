package boilerplate

import (
	"fmt"
	"os"
	"path/filepath"
)

func GenerateBoilerplate(year, day int) error {
	dirPath := fmt.Sprintf("aoc-%d/Day%d", year, day)
	filePath := filepath.Join(dirPath, "main.go")

	// Check if the directory exists, if not create it
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file %s already exists", filePath)
	}

	// Create the boilerplate file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Write the boilerplate code to the file
	boilerplateCode := fmt.Sprintf(`
package main

import (
	"fmt"
	"strings"
)

func main() {
	input, err := readInput(%d, %d)
	if err != nil {
		fmt.Printf("Failed to read input file: %%v\n", err)
		return
	}

	// Solution here
	fmt.Println("Solution:", input)
}

// Boilerplate: ReadInput reads the input file for the given year and day
func readInput(year, day int) (string, error) {
	filePath := fmt.Sprintf("aoc-%%d/Day%%d/input.txt", year, day)
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return string(file), nil
}
`, year, day)

	if _, err := file.WriteString(boilerplateCode); err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	fmt.Printf("Boilerplate created at %s\n", filePath)
	return nil
}
