package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
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

func makeDiskBlocks(format []int) ([]filePosition, []emptyBlockPosition) {
	filePositions := []filePosition{}
	emptyBlockPositions := []emptyBlockPosition{}
	currentPosition := 0

	for idx, n := range format {
		var v int
		if idx%2 == 0 {
			v = idx / 2
			filePositions = append(filePositions, filePosition{fileId: v, startIndex: currentPosition, length: n})
		} else {
			v = -1 // represents empty
			emptyBlockPositions = append(emptyBlockPositions, emptyBlockPosition{startIndex: currentPosition, length: n})
		}

		currentPosition += n
	}

	return filePositions, emptyBlockPositions
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

type filePosition struct {
	fileId     int
	startIndex int
	length     int
}

type emptyBlockPosition struct {
	startIndex int
	length     int
}

func shiftOneFile(files []filePosition, emptyBlocks []emptyBlockPosition, filteredEmpties []emptyBlockPosition) []emptyBlockPosition {
	if len(files) == 0 {
		return filteredEmpties
	}

	nextFilteredEmpty := []emptyBlockPosition{}
	lastEntry := files[len(files)-1]

	var newLastEntry filePosition

	// the nonfiltered list
	for idx, emptyBlock := range emptyBlocks {
		if emptyBlock.length >= lastEntry.length {
			newLastEntry = lastEntry
			newLastEntry.startIndex = emptyBlock.startIndex
			files[len(files)-1] = newLastEntry
			emptyBlock.length -= lastEntry.length
			emptyBlock.startIndex += lastEntry.length

			emptyBlocks[idx] = emptyBlock

			break
		}
	}

	// the filtered list
	for idx, nonzeroEmptyBlock := range filteredEmpties {
		if nonzeroEmptyBlock.length >= lastEntry.length {
			lastEntry.startIndex = nonzeroEmptyBlock.startIndex

			if lastEntry != newLastEntry {
				fmt.Println("Found a discrepancy! ", idx)
			}
			//files[len(files)-1] = lastEntry
			nonzeroEmptyBlock.length -= lastEntry.length
			nonzeroEmptyBlock.startIndex += lastEntry.length

			if nonzeroEmptyBlock.length == 0 {
				nextFilteredEmpty = append(nextFilteredEmpty, filteredEmpties[:idx]...)
				nextFilteredEmpty = append(nextFilteredEmpty, filteredEmpties[idx+1:]...)
			} else {
				filteredEmpties[idx] = nonzeroEmptyBlock
				nextFilteredEmpty = filteredEmpties
			}
			break
		}
	}
	return shiftOneFile(files[:len(files)-1], emptyBlocks, nextFilteredEmpty)
}

func newChecksum(filePositions []filePosition) int {
	sum := 0

	for _, fPos := range filePositions {
		for i := 0; i < fPos.length; i++ {
			sum += fPos.fileId * (i + fPos.startIndex)
		}
	}
	return sum
}

func main() {
	inputLine := parseFile("resources/Day9/input.txt")
	//fmt.Println(inputLine)

	filePositions, emptyBlockPositions := makeDiskBlocks(inputLine)
	//printBlocks(blocks)
	//fmt.Println(filePositions)
	//fmt.Println(emptyBlockPositions)
	//fmt.Println("")

	filteredEmpties := make([]emptyBlockPosition, 0, len(emptyBlockPositions))
	filteredEmpties = append(filteredEmpties, emptyBlockPositions...)

	emptyBlockPositions = shiftOneFile(filePositions, emptyBlockPositions, filteredEmpties)
	sort.Slice(filePositions, func(i int, j int) bool {
		return filePositions[i].startIndex < filePositions[j].startIndex
	})
	sort.Slice(emptyBlockPositions, func(i int, j int) bool {
		return emptyBlockPositions[i].startIndex < emptyBlockPositions[j].startIndex
	})
	fmt.Println(filePositions)
	fmt.Println(emptyBlockPositions)
	//fmt.Println(emptyBlockPositions)
	fmt.Println("new sum: ", newChecksum(filePositions))
	//for fileId := len(blocks) - 1; fileId >= 0; fileId-- {
	//	shiftOneFile(blocks, fileId)
	//	printBlocks(blocks)
	//}
	//compress(blocks)
	//printBlocks(blocks)
	//fmt.Println("sum: ", checksum(blocks))
}
