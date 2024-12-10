package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func parseTopoMap(fileName string) [][]int {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	topoMap := make([][]int, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		textLine := scanner.Text()
		numLine := make([]int, 0, len(textLine))

		for _, b := range textLine {
			num, err := strconv.Atoi(string(b))

			if err != nil {
				fmt.Println(err)
				return topoMap
			}
			numLine = append(numLine, num)
		}
		topoMap = append(topoMap, numLine)
	}
	return topoMap
}

func makeReachableMap(topoMap [][]int, thRow int, thColumn int) [][]bool {
	reachableMap := make([][]bool, 0, len(topoMap))
	for i := 0; i < len(topoMap); i++ {
		reachableMap = append(reachableMap, make([]bool, len(topoMap[i])))
	}

	var setReachable func(row int, column int)
	setReachable = func(row int, column int) {
		if reachableMap[row][column] {
			return
		}
		reachableMap[row][column] = true
		value := topoMap[row][column]
		if row > 0 && topoMap[row-1][column]-value == 1 {
			setReachable(row-1, column)
		}
		if column > 0 && topoMap[row][column-1]-value == 1 {
			setReachable(row, column-1)
		}
		if row < len(topoMap)-1 && topoMap[row+1][column]-value == 1 {
			setReachable(row+1, column)
		}
		if column < len(topoMap[row])-1 && topoMap[row][column+1]-value == 1 {
			setReachable(row, column+1)
		}
	}

	setReachable(thRow, thColumn)
	return reachableMap
}

func makeTrailCountFrom(topoMap [][]int, trailheadRow int, trailheadColumn int) [][]int {
	trailCountMap := make([][]int, 0, len(topoMap))
	for i := 0; i < len(topoMap); i++ {
		trailCountMap = append(trailCountMap, make([]int, len(topoMap[i])))
	}

	var incrementTrailCount func(row int, column int)
	incrementTrailCount = func(row int, column int) {
		trailCountMap[row][column] += 1
		value := topoMap[row][column]
		if row > 0 && topoMap[row-1][column]-value == 1 {
			incrementTrailCount(row-1, column)
		}
		if column > 0 && topoMap[row][column-1]-value == 1 {
			incrementTrailCount(row, column-1)
		}
		if row < len(topoMap)-1 && topoMap[row+1][column]-value == 1 {
			incrementTrailCount(row+1, column)
		}
		if column < len(topoMap[row])-1 && topoMap[row][column+1]-value == 1 {
			incrementTrailCount(row, column+1)
		}
	}

	incrementTrailCount(trailheadRow, trailheadColumn)
	return trailCountMap
}

type coords struct {
	row    int
	column int
}

func main() {
	topoMap := parseTopoMap("resources/Day10/sampleinput.txt")
	fmt.Println(topoMap)

	trailheads := []coords{}
	peaks := []coords{}

	for r, row := range topoMap {
		for c, height := range row {
			if height == 0 {
				trailheads = append(trailheads, coords{row: r, column: c})
			}
			if height == 9 {
				peaks = append(peaks, coords{row: r, column: c})
			}
		}
	}

	// part 1
	totalScore := 0
	for _, th := range trailheads {
		reachable := makeReachableMap(topoMap, th.row, th.column)

		score := 0
		for _, p := range peaks {
			if reachable[p.row][p.column] {
				score++
			}
		}
		fmt.Println("(", th.row, ",", th.column, "): ", score)
		totalScore += score
	}

	fmt.Println(totalScore)

	totalRating := 0
	// part 2
	for _, th := range trailheads {
		trailCount := makeTrailCountFrom(topoMap, th.row, th.column)
		fmt.Println(trailCount)

		rating := 0
		for _, p := range peaks {
			rating += trailCount[p.row][p.column]
		}
		fmt.Println("(", th.row, ",", th.column, "): ", rating)
		totalRating += rating
	}
	fmt.Println("total rating: ", totalRating)
}
