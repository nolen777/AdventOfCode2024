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

type TowelTree struct {
	present bool

	sub [5]*TowelTree
}

func (tt TowelTree) contains(t string) bool {
	if t == "" {
		return tt.present
	}
	sub := tt.sub[index(rune(t[0]))]
	if sub == nil {
		return false
	}
	return sub.contains(t[1:])
}

func recursiveFindTowelsForPattern(allTowels TowelTree, pattern string, usedTowels []string) ([]string, bool) {
	if pattern == "" {
		return usedTowels, true
	}

	for i := 1; i <= len(pattern); i++ {
		substr := pattern[:i]
		if allTowels.contains(substr) {
			newUsed := slices.Clone(usedTowels)
			newUsed = append(newUsed, substr)

			result, ok := recursiveFindTowelsForPattern(allTowels, pattern[i:], newUsed)

			if ok {
				return result, true
			}
		}
	}

	return usedTowels, false
}

func index(color rune) int {
	switch color {
	case 'w':
		return 0
	case 'u':
		return 1
	case 'b':
		return 2
	case 'r':
		return 3
	case 'g':
		return 4
	}
	log.Fatal("Unknown color", color)
	return -1
}

func makeTowelTree(towels []string) TowelTree {
	root := TowelTree{present: false}

	for _, towel := range towels {
		tt := &root

		for _, r := range towel {
			idx := index(r)

			if tt.sub[idx] == nil {
				tt.sub[idx] = &TowelTree{}
			}
			tt = tt.sub[idx]
		}
		tt.present = true
	}
	return root
}

func part1() {
	towels, patterns := parseTowelsAndPatterns("resources/Day19/sample.txt")

	towelTree := makeTowelTree(towels)

	successCount := 0
	for _, pattern := range patterns {
		usedTowels, ok := recursiveFindTowelsForPattern(towelTree, pattern, []string{})

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
