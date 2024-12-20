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

	return UpdateCosts(m, costs)
}

func UpdateCosts(m [][]rune, costs [][]int) [][]int {
	rowCount := len(costs)
	columnCount := len(costs[0])

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

func copyCosts(c [][]int) [][]int {
	cc := make([][]int, 0, len(c))

	for _, row := range c {
		cc = append(cc, slices.Clone(row))
	}
	return cc
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

func naiveSearch(m [][]rune, sortedCosts []TileCost, rawCosts [][]int, start Coords, minSavings int) map[int][]Coords {
	//rowCount := len(m)
	//columnCount := len(m[0])

	initialCost := rawCosts[start.row][start.column]

	savingsCount := map[int][]Coords{}

	maxDiff := 2

	for rowIndex, row := range rawCosts {
		for colIndex, rawCost := range row {
			if m[rowIndex][colIndex] == Wall {
				continue
			}

			for _, tileCost := range sortedCosts {
				dist := abs(rowIndex-tileCost.row) + abs(colIndex-tileCost.column)

				savings := rawCost - (tileCost.cost + dist)

				// since we're sorted, we can early-exit
				if savings <= 0 {
					break
				}

				if dist < 2 {
					continue
				}

				// right now only in straight lines
				if rowIndex != tileCost.row && colIndex != tileCost.column {
					continue
				}
				if dist > maxDiff {
					continue
				}
				cc := copyCosts(rawCosts)
				cc[rowIndex][colIndex] = tileCost.cost + dist
				UpdateCosts(m, cc)

				newCost := cc[start.row][start.column]
				finalSavings := initialCost - newCost
				if finalSavings > 0 {
					savingsCount[savings] = append(savingsCount[savings], Coords{rowIndex, colIndex})
				}
			}
		}
	}

	return savingsCount
}

func part1() {
	m, start, end := parseMap("resources/Day20/sample.txt")
	printMap(m)
	fmt.Println(start)
	fmt.Println(end)

	rawCosts := CalculateCosts(end, m)
	tileCosts := make([]TileCost, 0, len(m)*len(m[0]))
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

	fmt.Println("cost to end: ", rawCosts[start.row][start.column])

	savingsCount := naiveSearch(m, tileCosts, rawCosts, start, 1)

	for sav, coords := range savingsCount {
		fmt.Println(sav, ": ", len(coords), coords)
	}
}

func part2() {
	//m, start, end := parseMap("resources/Day20/sample.txt")
	//_ = lines
}

func main() {
	part1()
}
