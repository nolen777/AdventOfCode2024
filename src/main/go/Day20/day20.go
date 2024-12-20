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

type CheatPair struct {
	start Coords
	end   Coords
}

func copyMap(m [][]rune) [][]rune {
	mc := make([][]rune, 0, len(m))

	for _, row := range m {
		mc = append(mc, slices.Clone(row))
	}
	return mc
}

func copyCosts(c [][]int) [][]int {
	cc := make([][]int, 0, len(c))

	for _, row := range c {
		cc = append(cc, slices.Clone(row))
	}
	return cc
}

func naiveSearch(m [][]rune, costs [][]int, start Coords, end Coords) map[int][]Coords {
	rowCount := len(m)
	columnCount := len(m[0])

	initialCost := costs[end.row][end.column]

	savingsCount := map[int][]Coords{}

	isValid := func(r int, c int) bool {
		return r >= 0 && c >= 0 && r < rowCount && c < columnCount
	}

	for rowIndex, row := range costs {
		for colIndex, tileCost := range row {
			if tileCost == MaxInt {
				continue
			}

			checkOneDir := func(rd int, cd int) {
				if m[rowIndex+rd][colIndex+cd] == Wall || m[rowIndex+2*rd][colIndex+2*cd] == Wall {
					// First see if we can jump directly to exit
					if rowIndex+2*rd == end.row && colIndex+2*cd == end.column {
						newCost := costs[rowIndex][colIndex] + 2
						savings := initialCost - newCost
						if savings > 0 {
							fmt.Println("Found one going ", rd, cd, ", directly to the end!, that saved ", savings)
							fmt.Println(rowIndex, colIndex)
							savingsCount[savings] = append(savingsCount[savings], Coords{rowIndex, colIndex})
						}
						return
					}

					// Otherwise, need to be going to a valid empty space 3 away
					if !isValid(rowIndex+2*rd, colIndex+2*cd) {
						return
					}
					if m[rowIndex+3*rd][colIndex+3*cd] == Wall {
						return
					}
					currentCostToEnd := costs[rowIndex+3*rd][colIndex+3*cd]
					// Costs 3 to jump, so don't bother if we wouldn't save that
					if currentCostToEnd <= 3+tileCost {
						return
					}

					cc := copyCosts(costs)
					cc[rowIndex+3*rd][colIndex+3*cd] = tileCost + 3
					UpdateCosts(m, cc)

					newCost := cc[end.row][end.column]
					savings := initialCost - newCost
					if savings > 0 {
						fmt.Println("Found one going ", rd, cd, " that saved ", savings)
						fmt.Println(rowIndex, colIndex)
						savingsCount[savings] = append(savingsCount[savings], Coords{rowIndex, colIndex})
					}
				}
			}

			// up
			checkOneDir(-1, 0)
			checkOneDir(0, -1)
			checkOneDir(1, 0)
			checkOneDir(0, 1)
		}
	}

	return savingsCount
}

func part1() {
	m, start, end := parseMap("resources/Day20/sample.txt")
	printMap(m)
	fmt.Println(start)
	fmt.Println(end)

	costs := CalculateCosts(start, m)
	fmt.Println("cost to end: ", costs[end.row][end.column])

	savingsCount := naiveSearch(m, costs, start, end)

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
