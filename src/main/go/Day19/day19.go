package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	maxTowelLength := 0
	for _, t := range towels {
		maxTowelLength = max(maxTowelLength, len(t))
	}
	fmt.Println("Max towel length: ", maxTowelLength)

	scanner.Scan()
	line = scanner.Text()

	patterns := []string{}
	for scanner.Scan() {
		pattern := scanner.Text()
		patterns = append(patterns, pattern)
	}
	return towels, patterns
}

func recursiveFindTowelsForPattern(allTowels map[string]bool, pattern string, usedTowels []string) ([]string, bool) {
	if pattern == "" {
		return usedTowels, true
	}

	for i := min(len(pattern), 8); i > 0; i-- {
		//	for i := 1; i <= len(pattern); i++ {
		substr := pattern[:i]
		if allTowels[substr] {
			usedTowels = append(usedTowels, substr)

			result, ok := recursiveFindTowelsForPattern(allTowels, pattern[i:], usedTowels)

			if ok {
				return result, true
			} else {
				usedTowels = usedTowels[:len(usedTowels)-1]
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

func makeTowelSet(towels []string) map[string]bool {
	m := make(map[string]bool, len(towels))
	for _, t1 := range towels {
		m[t1] = true
		//for _, t2 := range towels {
		//	m[t1+t2] = true
		//}
	}
	return m
}

func findChars(c rune, allTowels map[string]bool, pattern string) bool {
	i := 0
	for i < len(pattern) {
		l := rune(pattern[i])
		found := false
		if l != c {
			i++
			continue
		}
		//outer:
		// start 7 steps behind
		for j := max(0, i-7); j <= i; j++ {
			for sz := min(8, len(pattern)-j); sz >= i-j+1; sz-- {
				//	for sz := i - j + 1; sz <= min(8, len(pattern)-j); sz++ {
				toFind := pattern[j : j+sz]
				if allTowels[toFind] {
					newPattern := pattern[j+sz:]
					if findChars(c, allTowels, newPattern) {
						return true
					}
				}
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func part1() {
	towels, patterns := parseTowelsAndPatterns("resources/Day19/input.txt")

	towelSet := makeTowelSet(towels)
	fmt.Println("Towel set size: ", len(towelSet))

	singleChars := []rune{'w', 'u', 'b', 'r', 'g'}
	absentSingleChars := []rune{}
	for _, c := range singleChars {
		if !towelSet[string(c)] {
			absentSingleChars = append(absentSingleChars, c)
		}
	}

	successCount := 0
	for _, pattern := range patterns {
		ok := true
		for _, psc := range absentSingleChars {
			if !findChars(psc, towelSet, pattern) {
				fmt.Println(pattern, "is impossible for", string(psc))
				ok = false
				break
			}
		}
		if ok {
			successCount++
			fmt.Println(pattern, "found!")
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
