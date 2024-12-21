package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
)

const (
	// start positions
	Wall  = rune('#')
	Start = rune('S')
	End   = rune('E')
	Empty = rune('.')
)

type Coords struct {
	row    int
	column int
}

func parseMap(fileName string) ([][]rune, Coords, Coords) {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	startMap := make([][]rune, 0)

	scanner := bufio.NewScanner(file)
	start := Coords{}
	end := Coords{}

	rowIndex := 0
	for scanner.Scan() {
		line := scanner.Text()

		row := make([]rune, 0)
		for colIndex, r := range line {
			row = append(row, r)

			if r == Start {
				start = Coords{row: rowIndex, column: colIndex}
			} else if r == End {
				end = Coords{row: rowIndex, column: colIndex}
			}
		}
		startMap = append(startMap, row)
		rowIndex++
	}

	return startMap, start, end
}

func printMap(m [][]rune) {
	for _, row := range m {
		for _, t := range row {
			fmt.Print(string(t))
		}
		fmt.Println("")
	}
}

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

func CalculateCosts(from Coords, m [][]rune) [][]int {
	rowCount := len(m)
	columnCount := len(m[0])

	costs := make([][]int, 0, rowCount)
	oneRow := make([]int, 0, columnCount)
	for j := 0; j < columnCount; j++ {
		oneRow = append(oneRow, MaxInt)
	}
	for i := 0; i < rowCount; i++ {
		costs = append(costs, slices.Clone(oneRow))
	}

	costs[from.row][from.column] = 0

	isValid := func(row int, column int) bool {
		return row >= 0 && column >= 0 && row < rowCount && column < columnCount
	}

	changes := 1
	// Go through every tile repeatedly. On each tile, see if we've identified
	// a cheaper path to it
	for changes > 0 {
		changes = 0
		for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
			for colIndex := 0; colIndex < columnCount; colIndex++ {
				if m[rowIndex][colIndex] == Wall {
					continue
				}

				if isValid(rowIndex-1, colIndex) {
					if costs[rowIndex-1][colIndex] < MaxInt && costs[rowIndex-1][colIndex]+1 < costs[rowIndex][colIndex] {
						changes++
						costs[rowIndex][colIndex] = costs[rowIndex-1][colIndex] + 1
					}
				}
				if isValid(rowIndex+1, colIndex) {
					if costs[rowIndex+1][colIndex] < MaxInt && costs[rowIndex+1][colIndex]+1 < costs[rowIndex][colIndex] {
						changes++
						costs[rowIndex][colIndex] = costs[rowIndex+1][colIndex] + 1
					}
				}
				if isValid(rowIndex, colIndex-1) {
					if costs[rowIndex][colIndex-1] < MaxInt && costs[rowIndex][colIndex-1]+1 < costs[rowIndex][colIndex] {
						changes++
						costs[rowIndex][colIndex] = costs[rowIndex][colIndex-1] + 1
					}
				}
				if isValid(rowIndex, colIndex+1) {
					if costs[rowIndex][colIndex+1] < MaxInt && costs[rowIndex][colIndex+1]+1 < costs[rowIndex][colIndex] {
						changes++
						costs[rowIndex][colIndex] = costs[rowIndex][colIndex+1] + 1
					}
				}
			}
		}
	}

	return costs
}

func MinRemoved(nodes []Coords, costs [][]int) (Coords, []Coords) {
	minValue := MaxInt
	minIndex := -1
	for idx, tc := range nodes {
		if costs[tc.row][tc.column] < minValue {
			minIndex = idx
		}
	}

	tc := nodes[minIndex]
	nodes[minIndex] = nodes[len(nodes)-1]
	return tc, nodes[:len(nodes)-1]
}

func DijkstraCosts(from Coords, m [][]rune) [][]int {
	rowCount := len(m)
	columnCount := len(m[0])

	costs := make([][]int, 0, rowCount)
	oneRow := make([]int, 0, columnCount)
	for j := 0; j < columnCount; j++ {
		oneRow = append(oneRow, MaxInt)
	}
	for i := 0; i < rowCount; i++ {
		costs = append(costs, slices.Clone(oneRow))
	}

	costs[from.row][from.column] = 0

	isValid := func(row int, column int) bool {
		return row >= 0 && column >= 0 && row < rowCount && column < columnCount
	}

	unvisitedNodes := make([]Coords, 0, rowCount*columnCount)
	for rIdx, row := range m {
		for cIdx, _ := range row {
			if m[rIdx][cIdx] == Wall {
				continue
			}
			unvisitedNodes = append(unvisitedNodes, Coords{rIdx, cIdx})
		}
	}

	changes := 1
	// Go through every tile repeatedly. On each tile, see if we've identified
	// a cheaper path to it
	for changes > 0 {
		changes = 0

		node, newUnv := MinRemoved(unvisitedNodes, costs)

		rowIndex := node.row
		colIndex := node.column
		thisCost := costs[rowIndex][colIndex]

		checkDir := func(rd int, cd int) {
			if !isValid(rowIndex+rd, colIndex+cd) {
				return
			}
			if m[rowIndex+rd][colIndex+cd] == Wall {
				return
			}
			if thisCost+1 < costs[rowIndex+rd][colIndex+cd] {
				changes++
				costs[rowIndex+rd][colIndex+cd] = thisCost + 1
			}
		}

		checkDir(-1, 0)
		checkDir(1, 0)
		checkDir(0, -1)
		checkDir(0, 1)

		unvisitedNodes = newUnv
		if len(unvisitedNodes) < 1 {
			break
		}
	}

	return costs
}

type TileCost struct {
	row    int
	column int
	cost   int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type CheatPair struct {
	start Coords
	end   Coords
}

func naiveSearch(sortedCosts []TileCost, rawCostsToEnd [][]int, rawCostsToStart [][]int, start Coords, minSavings int, maxCheatDistance int) map[int]map[CheatPair]bool {
	initialCost := rawCostsToEnd[start.row][start.column]

	savingsCount := map[int]map[CheatPair]bool{}

	for highCostIdx := len(sortedCosts) - 1; highCostIdx >= 0; highCostIdx-- {
		highCost := sortedCosts[highCostIdx]

		rowIndex := highCost.row
		colIndex := highCost.column

		for _, tileCost := range sortedCosts[:highCostIdx] {
			dist := abs(rowIndex-tileCost.row) + abs(colIndex-tileCost.column)

			savings := highCost.cost - (tileCost.cost + dist)

			// since we're sorted, we can early-exit
			if savings < minSavings {
				break
			}

			if dist < 2 {
				continue
			}

			//right now only in straight lines
			//if rowIndex != tileCost.row && colIndex != tileCost.column {
			//	continue
			//}
			if dist > maxCheatDistance {
				continue
			}

			newCost := rawCostsToStart[rowIndex][colIndex] + dist + rawCostsToEnd[tileCost.row][tileCost.column]
			finalSavings := initialCost - newCost
			if finalSavings > 0 {
				if savingsCount[savings] == nil {
					savingsCount[savings] = map[CheatPair]bool{}
				}
				savingsCount[savings][CheatPair{start: Coords{rowIndex, colIndex}, end: Coords{tileCost.row, tileCost.column}}] = true
			}
		}
	}

	return savingsCount
}

func sortedTileCosts(rawCosts [][]int) []TileCost {
	tileCosts := make([]TileCost, 0, len(rawCosts)*len(rawCosts[0]))

	for rowIndex, row := range rawCosts {
		for colIndex, c := range row {
			if c == MaxInt {
				continue
			}
			tileCosts = append(tileCosts, TileCost{row: rowIndex, column: colIndex, cost: c})
		}
	}

	slices.SortFunc(tileCosts, func(a TileCost, b TileCost) int {
		return a.cost - b.cost
	})
	return tileCosts
}

func part1() {
	m, start, end := parseMap("resources/Day20/input.txt")
	printMap(m)
	fmt.Println(start)
	fmt.Println(end)

	dijkstraCostsToStart := DijkstraCosts(start, m)
	dijkstraCostsToEnd := DijkstraCosts(end, m)

	sortedTileCostsToEnd := sortedTileCosts(dijkstraCostsToEnd)

	fmt.Println("cost to end: ", dijkstraCostsToEnd[start.row][start.column])

	savingsCount := naiveSearch(sortedTileCostsToEnd, dijkstraCostsToEnd, dijkstraCostsToStart, start, 100, 2)

	cheatCount := 0
	for _, coords := range savingsCount {
		//		fmt.Println(sav, ": ", len(coords), coords)
		cheatCount += len(coords)
	}
	fmt.Println("Total cheats: ", cheatCount)
}

func part2() {
	m, start, end := parseMap("resources/Day20/input.txt")
	printMap(m)
	fmt.Println(start)
	fmt.Println(end)

	dijkstraCostsToStart := DijkstraCosts(start, m)
	dijkstraCostsToEnd := DijkstraCosts(end, m)
	sortedTileCostsToEnd := sortedTileCosts(dijkstraCostsToEnd)

	fmt.Println("cost to end: ", dijkstraCostsToEnd[start.row][start.column])

	savingsCount := naiveSearch(sortedTileCostsToEnd, dijkstraCostsToEnd, dijkstraCostsToStart, start, 100, 20)

	cheatCount := 0
	for _, coords := range savingsCount {
		//		fmt.Println(sav, ": ", len(coords), coords)
		cheatCount += len(coords)
	}
	fmt.Println("Total cheats: ", cheatCount)
}

func main() {
	part2()
}
