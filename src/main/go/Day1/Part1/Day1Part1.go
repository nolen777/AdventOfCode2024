package main

import (
	"bufio"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func SortedLists(filename string) ([]int, []int) {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var left []int
	var right []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		entries := strings.Split(line, "   ")
		leftEntry, err := strconv.Atoi(entries[0])
		if err != nil {
			log.Fatal(err)
		}
		rightEntry, err := strconv.Atoi(entries[1])

		if err != nil {
			log.Fatal(err)
		}

		left = append(left, leftEntry)
		right = append(right, rightEntry)
	}

	slices.Sort(left)
	slices.Sort(right)

	return left, right
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	left, right := SortedLists("resources/Day1/input.txt")

	entryCount := len(left)
	pos := 0
	sum := 0

	for pos < entryCount {
		sum += abs(left[pos] - right[pos])
		pos += 1
	}

	println(sum)
}
