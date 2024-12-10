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

func adjacents(topoMap [][]int, c coords) []coords {
	adj := make([]coords, 0, 4)

	if c.row > 0 {
		adj = append(adj, coords{row: c.row - 1, column: c.column})
	}
	if c.column > 0 {
		adj = append(adj, coords{row: c.row, column: c.column - 1})
	}
	if c.row < len(topoMap)-1 {
		adj = append(adj, coords{row: c.row + 1, column: c.column})
	}
	if c.column < len(topoMap[c.row])-1 {
		adj = append(adj, coords{row: c.row, column: c.column + 1})
	}

	return adj
}

func gentleUpAdjacents(topoMap [][]int, c coords) []coords {
	value := topoMap[c.row][c.column]
	allAdj := adjacents(topoMap, c)
	upAdj := make([]coords, 0, len(allAdj))

	for _, a := range allAdj {
		if topoMap[a.row][a.column]-value == 1 {
			upAdj = append(upAdj, a)
		}
	}

	return upAdj
}

func makeReachableMap(topoMap [][]int, thCoords coords) [][]bool {
	reachableMap := make([][]bool, 0, len(topoMap))
	for i := 0; i < len(topoMap); i++ {
		reachableMap = append(reachableMap, make([]bool, len(topoMap[i])))
	}

	var setReachable func(c coords)
	setReachable = func(c coords) {
		if reachableMap[c.row][c.column] {
			return
		}
		reachableMap[c.row][c.column] = true

		for _, a := range gentleUpAdjacents(topoMap, c) {
			setReachable(a)
		}
	}

	setReachable(thCoords)
	return reachableMap
}

func makeTrailCountFrom(topoMap [][]int, trailhead coords) [][]int {
	trailCountMap := make([][]int, 0, len(topoMap))
	for i := 0; i < len(topoMap); i++ {
		trailCountMap = append(trailCountMap, make([]int, len(topoMap[i])))
	}

	var incrementTrailCount func(c coords)
	incrementTrailCount = func(c coords) {
		trailCountMap[c.row][c.column] += 1

		for _, a := range gentleUpAdjacents(topoMap, c) {
			incrementTrailCount(a)
		}
	}

	incrementTrailCount(trailhead)
	return trailCountMap
}

type coords struct {
	row    int
	column int
}

func main() {
	topoMap := parseTopoMap("resources/Day10/input.txt")
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

	fmt.Println("***        ***")
	fmt.Println("*** PART 1 ***")
	fmt.Println("***        ***")
	totalScore := 0
	for _, th := range trailheads {
		reachable := makeReachableMap(topoMap, th)

		score := 0
		for _, p := range peaks {
			if reachable[p.row][p.column] {
				score++
			}
		}
		fmt.Println("(", th.row, ",", th.column, "): ", score)
		totalScore += score
	}

	fmt.Println("Total score: ", totalScore)
	fmt.Println("\n\n")

	fmt.Println("***        ***")
	fmt.Println("*** PART 2 ***")
	fmt.Println("***        ***")
	totalRating := 0
	// part 2
	for _, th := range trailheads {
		trailCount := makeTrailCountFrom(topoMap, th)

		rating := 0
		for _, p := range peaks {
			rating += trailCount[p.row][p.column]
		}
		fmt.Println("(", th.row, ",", th.column, "): ", rating)
		totalRating += rating
	}
	fmt.Println("total rating: ", totalRating)
}
