package download

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Config struct {
	SessionToken string `yaml:"sessionToken"`
}

const baseURL = "https://adventofcode.com"

func GetPuzzlePart(year, day int) error {
	configFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	fmt.Println("Downloading puzzle for year", year, "day", day, "part")
	url := fmt.Sprintf("%s/%d/day/%d/input", baseURL, year, day)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	request.AddCookie(&http.Cookie{Name: "session", Value: config.SessionToken})

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to download puzzle: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download puzzle: %s", response.Status)
	}

	dir := fmt.Sprintf("aoc-%d/Day%d", year, day)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	inputPath := filepath.Join(dir, "input.txt")
	file, err := os.Create(inputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	written, err := io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Downloaded %d bytes and written to %s\n", written, inputPath)

	// Create a placeholder for the test input too (as this is often useful)
	testInputPath := filepath.Join(dir, "input.test.txt")
	file, err = os.Create(testInputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.WriteString(file, "### This is the Test Input ###\n")
	if err != nil {
		return fmt.Errorf("failed to write test file: %w", err)
	}

	return nil
}
