package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
	size := len(m)

	isValid := func(x int, y int) bool {
		return x >= 0 && y >= 0 && x < size && y < size
	}

	costs := make([][]int, 0, size)
	for i := 0; i < size; i++ {
		line := make([]int, 0, size)
		for j := 0; j < size; j++ {
			line = append(line, MaxInt)
		}
		costs = append(costs, line)
	}

	costs[from.row][from.column] = 0

	changes := 1
	// Go through every tile repeatedly. On each tile, see if we've identified
	// a cheaper path to it
	for changes > 0 {
		changes = 0
		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				if m[x][y] == Wall {
					continue
				}

				if isValid(x-1, y) {
					if costs[x-1][y] < MaxInt && costs[x-1][y]+1 < costs[x][y] {
						changes++
						costs[x][y] = costs[x-1][y] + 1
					}
				}
				if isValid(x+1, y) {
					if costs[x+1][y] < MaxInt && costs[x+1][y]+1 < costs[x][y] {
						changes++
						costs[x][y] = costs[x+1][y] + 1
					}
				}
				if isValid(x, y-1) {
					if costs[x][y-1] < MaxInt && costs[x][y-1]+1 < costs[x][y] {
						changes++
						costs[x][y] = costs[x][y-1] + 1
					}
				}
				if isValid(x, y+1) {
					if costs[x][y+1] < MaxInt && costs[x][y+1]+1 < costs[x][y] {
						changes++
						costs[x][y] = costs[x][y+1] + 1
					}
				}
			}
		}
	}

	return costs
}

func part1() {
	m, start, end := parseMap("resources/Day20/sample.txt")
	printMap(m)
	fmt.Println(start)
	fmt.Println(end)

	costs := CalculateCosts(start, m)
	fmt.Println("cost to end: ", costs[end.row][end.column])
}

func part2() {
	//m, start, end := parseMap("resources/Day20/sample.txt")
	//_ = lines
}

func main() {
	part1()
}
