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

func part1() {
	towels, patterns := parseTowelsAndPatterns("resources/Day19/input.txt")

	towelSet := makeTowelSet(towels)
	fmt.Println("Towel set size: ", len(towelSet))

	successCount := 0
	for _, pattern := range patterns {
		ct := calcCounts(towelSet, pattern)
		fmt.Println(pattern, ct)
		if ct > 0 {
			successCount++
		}
	}
	fmt.Println("Total count:", successCount)

	// *** my original, horrible solution
	//singleChars := []rune{'w', 'u', 'b', 'r', 'g'}
	//absentSingleChars := []rune{}
	//for _, c := range singleChars {
	//	if !towelSet[string(c)] {
	//		absentSingleChars = append(absentSingleChars, c)
	//	}
	//}
	//
	//successCount := 0
	//for _, pattern := range patterns {
	//	ok := true
	//	for _, psc := range absentSingleChars {
	//		if !findChars(psc, towelSet, pattern) {
	//			fmt.Println(pattern, "is impossible for", string(psc))
	//			ok = false
	//			break
	//		}
	//	}
	//	if ok {
	//		successCount++
	//		fmt.Println(pattern, "found!")
	//	}
	//}
	//
	//fmt.Println("Success count: ", successCount)
}

var counts = map[string]int{"": 1}

func calcCounts(towelSet map[string]bool, pattern string) int {
	if ct, ok := counts[pattern]; ok {
		return ct
	}

	newCount := 0
	for sz := 1; sz <= len(pattern); sz++ {
		pre := pattern[:sz]

		if towelSet[pre] {
			suf := pattern[sz:]
			newCount += calcCounts(towelSet, suf)
		}
	}
	counts[pattern] = newCount
	return newCount
}

func part2() {
	towels, patterns := parseTowelsAndPatterns("resources/Day19/input.txt")

	towelSet := makeTowelSet(towels)
	fmt.Println("Towel set size: ", len(towelSet))

	totalCount := 0
	for _, pattern := range patterns {
		ct := calcCounts(towelSet, pattern)
		fmt.Println(pattern, ct)
		totalCount += ct
	}
	fmt.Println("Total count:", totalCount)
}

func main() {
	part2()
}
