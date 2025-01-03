package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func ReadFile(filename string) [][]rune {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var runes = make([][]rune, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		textLine := scanner.Text()
		runeLine := make([]rune, 0)

		for _, r := range textLine {
			runeLine = append(runeLine, r)
		}

		runes = append(runes, runeLine)
	}

	return runes
}

// part 1
func hasMatch(runes [][]rune, str string, i int, j int) int {
	matchCount := 0
	// horizontal forward
	if j+len(str) <= len(runes[i]) {
		success := true
		for x, r := range str {
			if runes[i][j+x] != r {
				success = false
				break
			}
		}
		if success {
			matchCount++
		}
	}
	// horizontal backward
	if j-len(str) >= -1 {
		success := true
		for x, r := range str {
			if runes[i][j-x] != r {
				success = false
				break
			}
		}
		if success {
			matchCount++
		}
	}
	// vertical down
	if i+len(str) <= len(runes) {
		success := true
		for x, r := range str {
			if runes[i+x][j] != r {
				success = false
				break
			}
		}
		if success {
			matchCount++
		}
	}
	// vertical up
	if i-len(str) >= -1 {
		success := true
		for x, r := range str {
			if runes[i-x][j] != r {
				success = false
				break
			}
		}
		if success {
			matchCount++
		}
	}
	// Down-right
	if i+len(str) <= len(runes) && j+len(str) <= len(runes[i]) {
		success := true
		for x, r := range str {
			if runes[i+x][j+x] != r {
				success = false
				break
			}
		}
		if success {
			matchCount++
		}
	}
	// Down-left
	if i+len(str) <= len(runes) && j-len(str) >= -1 {
		success := true
		for x, r := range str {
			if runes[i+x][j-x] != r {
				success = false
				break
			}
		}
		if success {
			matchCount++
		}
	}
	// Up-right
	if i-len(str) >= -1 && j+len(str) <= len(runes[i]) {
		success := true
		for x, r := range str {
			if runes[i-x][j+x] != r {
				success = false
				break
			}
		}
		if success {
			matchCount++
		}
	}
	// Up-left
	if i-len(str) >= -1 && j-len(str) >= -1 {
		success := true
		for x, r := range str {
			if runes[i-x][j-x] != r {
				success = false
				break
			}
		}
		if success {
			matchCount++
		}
	}

	if matchCount > 0 {
		fmt.Println("Found ", matchCount, " matches for ", i, j)
	}
	return matchCount
}

func count(runes [][]rune, str string) int {
	c := 0
	for i := 0; i < len(runes); i++ {
		for j := 0; j < len(runes[i]); j++ {
			c += hasMatch(runes, str, i, j)
		}
	}

	return c
}

// part 2
func hasXMatch(runes [][]rune, str string, i int, j int) int {
	if i+2 >= len(runes) || j+2 >= len(runes[i]) {
		return 0
	}

	if runes[i+1][j+1] != 'A' {
		return 0
	}

	downRight := ((runes[i][j] == 'M' && runes[i+2][j+2] == 'S') ||
		(runes[i][j] == 'S' && runes[i+2][j+2] == 'M'))

	upRight := ((runes[i+2][j] == 'M' && runes[i][j+2] == 'S') ||
		(runes[i+2][j] == 'S' && runes[i][j+2] == 'M'))

	if upRight && downRight {
		return 1
	}
	return 0
}

func count2(runes [][]rune, str string) int {
	c := 0
	for i := 0; i < len(runes); i++ {
		for j := 0; j < len(runes[i]); j++ {
			c += hasXMatch(runes, str, i, j)
		}
	}

	return c
}

func main() {
	runes := ReadFile("resources/Day4/input.txt")

	fmt.Println(count2(runes, "MAS"))
}
