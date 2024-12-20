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

	return UpdateCosts(m, costs, nil)
}

func UpdateCosts(m [][]rune, costs [][]int, shouldExit func(int, int) bool) [][]int {
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
						if shouldExit != nil && shouldExit(rowIndex, colIndex) {
							return costs
						}
					}
				}
				if isValid(rowIndex+1, colIndex) {
					if costs[rowIndex+1][colIndex] < MaxInt && costs[rowIndex+1][colIndex]+1 < costs[rowIndex][colIndex] {
						changes++
						costs[rowIndex][colIndex] = costs[rowIndex+1][colIndex] + 1
						if shouldExit != nil && shouldExit(rowIndex, colIndex) {
							return costs
						}
					}
				}
				if isValid(rowIndex, colIndex-1) {
					if costs[rowIndex][colIndex-1] < MaxInt && costs[rowIndex][colIndex-1]+1 < costs[rowIndex][colIndex] {
						changes++
						costs[rowIndex][colIndex] = costs[rowIndex][colIndex-1] + 1
						if shouldExit != nil && shouldExit(rowIndex, colIndex) {
							return costs
						}
					}
				}
				if isValid(rowIndex, colIndex+1) {
					if costs[rowIndex][colIndex+1] < MaxInt && costs[rowIndex][colIndex+1]+1 < costs[rowIndex][colIndex] {
						changes++
						costs[rowIndex][colIndex] = costs[rowIndex][colIndex+1] + 1
						if shouldExit != nil && shouldExit(rowIndex, colIndex) {
							return costs
						}
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

func naiveSearch(m [][]rune, costs [][]int, end Coords, maxCheatLength int, minSavings int) map[int][]Coords {
	rowCount := len(m)
	columnCount := len(m[0])

	initialCost := costs[end.row][end.column]

	savingsCount := map[int][]Coords{}

	isValid := func(r int, c int) bool {
		return r >= 0 && c >= 0 && r < rowCount && c < columnCount
	}

	for cheatStartRow, row := range costs {
		fmt.Println("Starting row ", cheatStartRow, " of ", rowCount)
		for cheatStartColumn, tileCost := range row {
			if tileCost == MaxInt {
				continue
			}

			for rd := 0; rd <= maxCheatLength; rd++ {
				if !isValid(cheatStartRow+rd, cheatStartColumn) {
					continue
				}
				for cd := 0; cd <= maxCheatLength-rd; cd++ {
					if rd == 0 && cd == 0 {
						continue
					}
					if !isValid(cheatStartRow+rd, cheatStartColumn+cd) {
						continue
					}
					if m[cheatStartRow+rd][cheatStartColumn+cd] == Wall {
						continue
					}
					currentCostToEnd := costs[cheatStartRow+rd][cheatStartColumn+cd]
					// don't bother if we wouldn't save the cost to cheat
					if currentCostToEnd <= tileCost+rd+cd {
						continue
					}

					cc := copyCosts(costs)
					cc[cheatStartRow+rd][cheatStartColumn+cd] = tileCost + rd + cd
					UpdateCosts(m, cc, func(r int, c int) bool {
						return initialCost-cc[end.row][end.column] >= minSavings
					})

					newCost := cc[end.row][end.column]
					savings := initialCost - newCost
					if savings >= minSavings {
						savingsCount[savings] = append(savingsCount[savings], Coords{cheatStartRow, cheatStartColumn})
					}
				}
			}
		}
	}

	return savingsCount
}

func part1() {
	m, start, end := parseMap("resources/Day20/input.txt")
	printMap(m)
	fmt.Println(start)
	fmt.Println(end)

	costs := CalculateCosts(start, m)
	fmt.Println("cost to end: ", costs[end.row][end.column])

	savingsCount := naiveSearch(m, costs, end, 2, 100)

	totalCheats := 0
	for sav, coords := range savingsCount {
		fmt.Println(sav, ": ", len(coords), coords)
		totalCheats += len(coords)
	}

	fmt.Println("Total cheats: ", totalCheats)
}

func part2() {
	m, start, end := parseMap("resources/Day20/sample.txt")
	printMap(m)
	fmt.Println(start)
	fmt.Println(end)

	costs := CalculateCosts(start, m)
	fmt.Println("cost to end: ", costs[end.row][end.column])

	savingsCount := naiveSearch(m, costs, end, 2, 1)

	totalCheats := 0
	for sav, coords := range savingsCount {
		fmt.Println(sav, ": ", len(coords), coords)
		totalCheats += len(coords)
	}

	fmt.Println("Total cheats: ", totalCheats)
}

func main() {
	//part1()
	part2()
}
