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

func main() {
	left, right := SortedLists("resources/Day1/input.txt")

	entryCount := len(left)
	leftPos, rightPos := 0, 0
	lastLeft, lastLeftSum := 0, 0
	sum := 0

	for leftPos < entryCount && rightPos < entryCount {
		if left[leftPos] == lastLeft {
			sum += lastLeftSum
		} else {
			lastLeft = left[leftPos]
			lastLeftSum = 0
			for rightPos < entryCount && right[rightPos] < left[leftPos] {
				rightPos += 1
			}
			for rightPos < entryCount && right[rightPos] == left[leftPos] {
				lastLeftSum += lastLeft
				rightPos += 1
			}
			sum += lastLeftSum
		}

		leftPos += 1
	}

	println(sum)
}
