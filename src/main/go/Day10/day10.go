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

	//for r, row := range topoMap {
	//	for c, height := range row {
	//		if height == 0 {
	//			reachable := makeReachableMap(topoMap, r, c)
	//			fmt.Println("(", r, ",", c, "):")
	//
	//			fmt.Println(reachable)
	//		}
	//	}
	//}
}
