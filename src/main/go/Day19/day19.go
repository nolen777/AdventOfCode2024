package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

func parseTowelsAndPatterns(fileName string) ([]string, []string) {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// first parse available towels
	scanner.Scan()
	line := scanner.Text()
	towels := strings.Split(line, ", ")

	scanner.Scan()
	line = scanner.Text()

	patterns := []string{}
	for scanner.Scan() {
		pattern := scanner.Text()
		patterns = append(patterns, pattern)
	}
	return towels, patterns
}

func recursiveFindTowelsForPattern(allTowels []string, pattern string, usedTowels []string) ([]string, bool) {
	if pattern == "" {
		return usedTowels, true
	}

	for _, towel := range allTowels {
		if strings.HasPrefix(pattern, towel) {
			newUsed := slices.Clone(usedTowels)
			newUsed = append(newUsed, towel)

			result, ok := recursiveFindTowelsForPattern(allTowels, pattern[len(towel):], newUsed)

			if ok {
				return result, true
			}
		}
	}

	return usedTowels, false
}

func part1() {
	towels, patterns := parseTowelsAndPatterns("resources/Day19/sample.txt")

	successCount := 0
	for _, pattern := range patterns {
		usedTowels, ok := recursiveFindTowelsForPattern(towels, pattern, []string{})

		if ok {
			successCount++
			fmt.Println("Success for ", pattern, " with ", usedTowels)
		} else {
			fmt.Println(pattern, " is impossible")
		}
	}

	fmt.Println("Success count: ", successCount)
}

func part2() {
	towels, patterns := parseTowelsAndPatterns("resources/Day19/sample.txt")
	_ = towels
	_ = patterns
}

func main() {
	part1()
}
