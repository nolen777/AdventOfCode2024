package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func parseFile(fileName string) []int {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		textLine := scanner.Text()
		input := make([]int, 0, len(textLine))

		for _, b := range textLine {
			num, err := strconv.Atoi(string(b))

			if err != nil {
				fmt.Println(err)
				return input
			}
			input = append(input, num)
		}
		return input
	}
	log.Fatal("failed to read text")
	return []int{}
}

func makeDiskBlocks(format []int) []int {
	blocks := []int{}

	for idx, n := range format {
		var v int
		if idx%2 == 0 {
			v = idx / 2
		} else {
			v = -1 // represents empty
		}

		for i := 0; i < n; i++ {
			blocks = append(blocks, v)
		}
	}

	return blocks
}

func printBlocks(blocks []int) {
	for _, n := range blocks {
		if n == -1 {
			fmt.Print(".")
		} else {
			fmt.Print(string(rune(n + '0')))
		}
	}
	fmt.Println("")
}

func compress(blocks []int) {
	leftCursor := 0
	rightCursor := len(blocks) - 1

	for leftCursor < rightCursor {
		if blocks[leftCursor] != -1 {
			leftCursor++
		} else if blocks[rightCursor] == -1 {
			rightCursor--
		} else {
			blocks[leftCursor] = blocks[rightCursor]
			blocks[rightCursor] = -1
		}
	}
}

func checksum(blocks []int) int {
	sum := 0
	for idx, n := range blocks {
		if n == -1 {
			break
		}
		sum += idx * n
	}
	return sum
}

func main() {
	inputLine := parseFile("resources/Day9/sampleinput.txt")
	fmt.Println(inputLine)

	blocks := makeDiskBlocks(inputLine)
	printBlocks(blocks)

	compress(blocks)
	printBlocks(blocks)
	fmt.Println("sum: ", checksum(blocks))
}
