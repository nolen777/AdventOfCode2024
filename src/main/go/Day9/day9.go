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

func makeDiskBlocks(format []int) ([]int, []emptyBlockPosition) {
	blocks := []int{}
	emptyBlockPositions := []emptyBlockPosition{}

	for idx, n := range format {
		var v int
		if idx%2 == 0 {
			v = idx / 2
		} else {
			v = -1 // represents empty
			emptyBlockPositions = append(emptyBlockPositions, emptyBlockPosition{startIndex: len(blocks), length: n})
		}

		for i := 0; i < n; i++ {
			blocks = append(blocks, v)
		}
	}

	return blocks, emptyBlockPositions
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
			continue
		}
		sum += idx * n
	}
	return sum
}

type emptyBlockPosition struct {
	startIndex int
	length     int
}

func shiftOneFile(blocks []int, currentFileId int) {
	totalBlockCount := len(blocks)

	startIndex := 0
	fileCount := 0
	// How long a contiguous block do we need?
	for idx, n := range blocks {
		if n == currentFileId {
			startIndex = idx
			fileCount = 0

			for fileCount+idx < totalBlockCount && blocks[idx+fileCount] == n {
				fileCount++
			}
			break
		}
	}

	// Find a block that size
	for idx, n := range blocks[:startIndex] {
		if n == -1 {
			emptyCount := 0

			for emptyCount+idx < startIndex && blocks[idx+emptyCount] == -1 {
				emptyCount++

				if emptyCount >= fileCount {
					// We can move
					for i := idx; i < idx+fileCount; i++ {
						blocks[i] = currentFileId
					}
					for i := startIndex; i < startIndex+fileCount; i++ {
						blocks[i] = -1
					}

					return
				}
			}
		}
	}
}

func main() {
	inputLine := parseFile("resources/Day9/input.txt")
	fmt.Println(inputLine)

	blocks, emptyBlockPositions := makeDiskBlocks(inputLine)
	printBlocks(blocks)
	fmt.Println(emptyBlockPositions)

	for fileId := len(blocks) - 1; fileId >= 0; fileId-- {
		shiftOneFile(blocks, fileId)
		//printBlocks(blocks)
	}
	//compress(blocks)
	//printBlocks(blocks)
	fmt.Println("sum: ", checksum(blocks))
}
