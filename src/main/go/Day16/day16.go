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

	// track
	North = rune('^')
	East  = rune('>')
	South = rune('v')
	West  = rune('<')

	// seats
	Seat = rune('O')
)

func clockwiseTurn(dir rune) rune {
	switch dir {
	case North:
		return East
	case East:
		return South
	case South:
		return West
	case West:
		return North
	}
	log.Fatal("invalid direction")
	return 0
}

func counterclockwiseTurn(dir rune) rune {
	switch dir {
	case North:
		return West
	case East:
		return North
	case South:
		return East
	case West:
		return South
	}
	log.Fatal("invalid direction")
	return 0
}

type Coords struct {
	row    int
	column int
}

type Map struct {
	grid      [][]rune
	position  Coords
	direction rune // one of the four track directions above
}

func parseMap(fileName string) Map {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	startMap := Map{grid: make([][]rune, 0), direction: East}

	scanner := bufio.NewScanner(file)

	rowIndex := 0
	for scanner.Scan() {
		line := scanner.Text()

		row := make([]rune, 0)
		for colIndex, r := range line {
			if r == Start {
				startMap.position = Coords{row: rowIndex, column: colIndex}
			}
			row = append(row, r)

		}
		startMap.grid = append(startMap.grid, row)
		rowIndex++
	}

	return startMap
}

func printMap(m Map) {
	for _, row := range m.grid {
		for _, r := range row {
			fmt.Print(string(r))
		}
		fmt.Println("")
	}
	fmt.Println("Position: row ", m.position.row, ", column ", m.position.column)
	fmt.Print("Direction: ")
	switch m.direction {
	case North:
		fmt.Println("North")
	case East:
		fmt.Println("East")
	case South:
		fmt.Println("South")
	case West:
		fmt.Println("West")
	}
}

func copyMap(m Map) Map {
	mc := Map{grid: make([][]rune, 0, len(m.grid)), position: m.position, direction: m.direction}

	for _, row := range m.grid {
		rowCopy := make([]rune, 0, len(row))
		rowCopy = append(rowCopy, row...)
		mc.grid = append(mc.grid, rowCopy)
	}
	return mc
}

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

const forwardPoints = 1
const turnPoints = 1000

var found bool = false
var lowestPoints int = MaxInt
var bestMap Map

func fillWalls(m Map) Map {
	changing := true

	for changing {
		changing = false
		for rowIndex, row := range m.grid {
			for columnIndex, r := range row {
				if r == Empty {
					wallCount := 0
					if m.grid[rowIndex][columnIndex-1] == Wall {
						wallCount++
					}
					if m.grid[rowIndex][columnIndex+1] == Wall {
						wallCount++
					}
					if m.grid[rowIndex-1][columnIndex] == Wall {
						wallCount++
					}
					if m.grid[rowIndex+1][columnIndex] == Wall {
						wallCount++
					}
					if wallCount >= 3 {
						changing = true
						m.grid[rowIndex][columnIndex] = Wall
					}
				}
			}
		}
	}

	return m
}

func calcGrid(oneGrid [][]int, clock [][]int, counter [][]int, rowIndex int, columnIndex int, rowDiff int, colDiff int) int {
	currentCost := oneGrid[rowIndex][columnIndex]
	changes := 0
	// try by walking forward
	if oneGrid[rowIndex+rowDiff][columnIndex+colDiff]+forwardPoints < currentCost {
		changes = 1
		oneGrid[rowIndex][columnIndex] = oneGrid[rowIndex+rowDiff][columnIndex+colDiff] + forwardPoints
	}
	// try by turning
	if clock[rowIndex][columnIndex]+turnPoints < currentCost {
		changes = 1
		oneGrid[rowIndex][columnIndex] = clock[rowIndex][columnIndex] + turnPoints
	}
	if counter[rowIndex][columnIndex]+turnPoints < currentCost {
		changes = 1
		oneGrid[rowIndex][columnIndex] = counter[rowIndex][columnIndex] + turnPoints
	}
	return changes
}

// returns cost of End
func progressiveFill(m Map) Map {
	initialCost := MaxInt - turnPoints - forwardPoints

	makeInitial := func() [][]int {
		g := make([][]int, 0, len(m.grid))
		for _, row := range m.grid {
			gRow := make([]int, 0, len(row))
			for _, _ = range row {
				gRow = append(gRow, initialCost)
			}
			g = append(g, gRow)
		}
		return g
	}

	// initialize costs for grid
	northFacing := makeInitial()
	eastFacing := makeInitial()
	southFacing := makeInitial()
	westFacing := makeInitial()

	eastFacing[m.position.row][m.position.column] = 0

	minEndCost := initialCost

	var endPosition Coords
	changes := 1
	// Go through every tile repeatedly. On each tile, see if we've identified
	// a cheaper path to it
	for changes > 0 {
		changes = 0
		// go rows bottom to top
		for rowIndex := len(m.grid) - 2; rowIndex > 0; rowIndex-- {
			row := m.grid[rowIndex]
			for columnIndex, b := range row {
				if b == Wall {
					continue
				}

				tileChanges := 0
				tileChanges += calcGrid(northFacing, westFacing, eastFacing, rowIndex, columnIndex, 1, 0)
				tileChanges += calcGrid(eastFacing, northFacing, southFacing, rowIndex, columnIndex, 0, -1)
				tileChanges += calcGrid(southFacing, westFacing, eastFacing, rowIndex, columnIndex, -1, 0)
				tileChanges += calcGrid(westFacing, northFacing, southFacing, rowIndex, columnIndex, 0, 1)

				changes += tileChanges

				if tileChanges > 0 && b == End {
					endPosition.row = rowIndex
					endPosition.column = columnIndex
					//fmt.Println("Arrived!")
					minEndCost = min(
						minEndCost,
						northFacing[rowIndex][columnIndex],
						eastFacing[rowIndex][columnIndex],
						southFacing[rowIndex][columnIndex],
						westFacing[rowIndex][columnIndex])
					fmt.Println("New end cost: ", minEndCost)
				}
			}
		}
	}

	return recursiveMarkSeats(endPosition, minEndCost, m, northFacing, eastFacing, southFacing, westFacing)
}

func recursiveMarkSeats(pos Coords, minCost int, m Map, nf [][]int, ef [][]int, sf [][]int, wf [][]int) Map {
	if m.grid[pos.row][pos.column] == Wall {
		return m
	}

	if nf[pos.row][pos.column] == minCost {
		m.grid[pos.row][pos.column] = Seat
		m = recursiveMarkSeats(Coords{row: pos.row + 1, column: pos.column}, minCost-forwardPoints, m, nf, ef, sf, wf)
	}
	if ef[pos.row][pos.column] == minCost {
		m.grid[pos.row][pos.column] = Seat
		m = recursiveMarkSeats(Coords{row: pos.row, column: pos.column - 1}, minCost-forwardPoints, m, nf, ef, sf, wf)
	}
	if sf[pos.row][pos.column] == minCost {
		m.grid[pos.row][pos.column] = Seat
		m = recursiveMarkSeats(Coords{row: pos.row - 1, column: pos.column}, minCost-forwardPoints, m, nf, ef, sf, wf)
	}
	if wf[pos.row][pos.column] == minCost {
		m.grid[pos.row][pos.column] = Seat
		m = recursiveMarkSeats(Coords{row: pos.row, column: pos.column + 1}, minCost-forwardPoints, m, nf, ef, sf, wf)
	}

	if nf[pos.row][pos.column] == minCost-turnPoints {
		m.grid[pos.row][pos.column] = Seat
		m = recursiveMarkSeats(Coords{row: pos.row + 1, column: pos.column}, minCost-turnPoints-forwardPoints, m, nf, ef, sf, wf)
	}
	if ef[pos.row][pos.column] == minCost-turnPoints {
		m.grid[pos.row][pos.column] = Seat
		m = recursiveMarkSeats(Coords{row: pos.row, column: pos.column - 1}, minCost-turnPoints-forwardPoints, m, nf, ef, sf, wf)
	}
	if sf[pos.row][pos.column] == minCost-turnPoints {
		m.grid[pos.row][pos.column] = Seat
		m = recursiveMarkSeats(Coords{row: pos.row - 1, column: pos.column}, minCost-turnPoints-forwardPoints, m, nf, ef, sf, wf)
	}
	if wf[pos.row][pos.column] == minCost-turnPoints {
		m.grid[pos.row][pos.column] = Seat
		m = recursiveMarkSeats(Coords{row: pos.row, column: pos.column + 1}, minCost-turnPoints-forwardPoints, m, nf, ef, sf, wf)
	}
	return m
}

func part1() {
	// Try spreading out costs
	startMap := parseMap("resources/Day16/input.txt")
	printMap(startMap)

	startMap = fillWalls(startMap)
	printMap(startMap)

	startMap.grid[startMap.position.row][startMap.position.column] = Empty

	filledMap := progressiveFill(startMap)
	printMap(filledMap)

	seatCount := 0
	for _, row := range filledMap.grid {
		for _, b := range row {
			if b == Seat {
				seatCount++
			}
		}
	}
	fmt.Println("Seat count: ", seatCount)
}

func main() {
	part1()
}
