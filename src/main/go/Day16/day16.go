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

type Coords struct {
	row    int
	column int
}

type Tile struct {
	status          rune
	northFacingCost int
	eastFacingCost  int
	southFacingCost int
	westFacingCost  int
}

func (t Tile) minCost() int {
	return min(t.northFacingCost, t.eastFacingCost, t.southFacingCost, t.westFacingCost)
}

type Map struct {
	grid      [][]Tile
	position  Coords
	direction rune // one of the four track directions above
}

func parseMap(fileName string) Map {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	startMap := Map{grid: make([][]Tile, 0), direction: East}

	scanner := bufio.NewScanner(file)

	rowIndex := 0
	for scanner.Scan() {
		line := scanner.Text()

		row := make([]Tile, 0)
		for colIndex, r := range line {
			newTile := Tile{status: r, northFacingCost: InitialCost, eastFacingCost: InitialCost, southFacingCost: InitialCost, westFacingCost: InitialCost}

			if r == Start {
				startMap.position = Coords{row: rowIndex, column: colIndex}
				newTile.eastFacingCost = 0
			}
			row = append(row, newTile)
		}
		startMap.grid = append(startMap.grid, row)
		rowIndex++
	}

	return startMap
}

func printMap(m Map) {
	for _, row := range m.grid {
		for _, t := range row {
			fmt.Print(string(t.status))
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

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)
const InitialCost = MaxInt - turnPoints - forwardPoints

const forwardPoints = 1
const turnPoints = 1000

func fillWalls(m Map) Map {
	changing := true

	for changing {
		changing = false
		for rowIndex, row := range m.grid {
			for columnIndex, t := range row {
				if t.status == Empty {
					wallCount := 0
					if m.grid[rowIndex][columnIndex-1].status == Wall {
						wallCount++
					}
					if m.grid[rowIndex][columnIndex+1].status == Wall {
						wallCount++
					}
					if m.grid[rowIndex-1][columnIndex].status == Wall {
						wallCount++
					}
					if m.grid[rowIndex+1][columnIndex].status == Wall {
						wallCount++
					}
					if wallCount >= 3 {
						changing = true
						m.grid[rowIndex][columnIndex].status = Wall
					}
				}
			}
		}
	}

	return m
}

func calcGrid(m Map, rowIndex int, columnIndex int) int {
	currentTile := m.grid[rowIndex][columnIndex]
	changes := 0

	// try by walking forward
	if m.grid[rowIndex-1][columnIndex].southFacingCost+forwardPoints < currentTile.southFacingCost {
		changes = 1
		m.grid[rowIndex][columnIndex].southFacingCost = m.grid[rowIndex-1][columnIndex].southFacingCost + forwardPoints
	}
	if m.grid[rowIndex][columnIndex-1].eastFacingCost+forwardPoints < currentTile.eastFacingCost {
		changes = 1
		m.grid[rowIndex][columnIndex].eastFacingCost = m.grid[rowIndex][columnIndex-1].eastFacingCost + forwardPoints
	}
	if m.grid[rowIndex+1][columnIndex].northFacingCost+forwardPoints < currentTile.northFacingCost {
		changes = 1
		m.grid[rowIndex][columnIndex].northFacingCost = m.grid[rowIndex+1][columnIndex].northFacingCost + forwardPoints
	}
	if m.grid[rowIndex][columnIndex+1].westFacingCost+forwardPoints < currentTile.westFacingCost {
		changes = 1
		m.grid[rowIndex][columnIndex].westFacingCost = m.grid[rowIndex][columnIndex+1].westFacingCost + forwardPoints
	}

	// try by turning
	if currentTile.northFacingCost+turnPoints < currentTile.eastFacingCost {
		changes = 1
		m.grid[rowIndex][columnIndex].eastFacingCost = currentTile.northFacingCost + turnPoints
	}
	if currentTile.northFacingCost+turnPoints < currentTile.westFacingCost {
		changes = 1
		m.grid[rowIndex][columnIndex].westFacingCost = currentTile.northFacingCost + turnPoints
	}
	if currentTile.eastFacingCost+turnPoints < currentTile.northFacingCost {
		changes = 1
		m.grid[rowIndex][columnIndex].northFacingCost = currentTile.eastFacingCost + turnPoints
	}
	if currentTile.eastFacingCost+turnPoints < currentTile.southFacingCost {
		changes = 1
		m.grid[rowIndex][columnIndex].southFacingCost = currentTile.eastFacingCost + turnPoints
	}
	if currentTile.southFacingCost+turnPoints < currentTile.eastFacingCost {
		changes = 1
		m.grid[rowIndex][columnIndex].eastFacingCost = currentTile.southFacingCost + turnPoints
	}
	if currentTile.southFacingCost+turnPoints < currentTile.westFacingCost {
		changes = 1
		m.grid[rowIndex][columnIndex].westFacingCost = currentTile.southFacingCost + turnPoints
	}
	if currentTile.westFacingCost+turnPoints < currentTile.northFacingCost {
		changes = 1
		m.grid[rowIndex][columnIndex].northFacingCost = currentTile.westFacingCost + turnPoints
	}
	if currentTile.westFacingCost+turnPoints < currentTile.southFacingCost {
		changes = 1
		m.grid[rowIndex][columnIndex].southFacingCost = currentTile.westFacingCost + turnPoints
	}

	return changes
}

// returns cost of End
func progressiveFill(m Map) Map {
	minEndCost := InitialCost

	var endPosition Coords
	changes := 1
	// Go through every tile repeatedly. On each tile, see if we've identified
	// a cheaper path to it
	for changes > 0 {
		changes = 0
		// go rows bottom to top
		for rowIndex := len(m.grid) - 2; rowIndex > 0; rowIndex-- {
			rowHasEntries := false
			row := m.grid[rowIndex]
			for columnIndex, t := range row {
				if t.status == Wall {
					continue
				}

				tileChanges := calcGrid(m, rowIndex, columnIndex)

				changes += tileChanges

				if tileChanges > 0 && t.status == End {
					endPosition.row = rowIndex
					endPosition.column = columnIndex
					//fmt.Println("Arrived!")
					minEndCost = min(
						minEndCost,
						m.grid[rowIndex][columnIndex].minCost())
					fmt.Println("New end cost: ", minEndCost)
				}

				if row[columnIndex].minCost() < InitialCost {
					rowHasEntries = true
				}
			}
			// Since we started from the bottom, if a row has no costs set, we can break -- nothing
			// above possibly can
			if !rowHasEntries {
				break
			}
		}
	}

	fmt.Println("Min cost: ", minEndCost)
	return recursiveMarkSeats(endPosition, minEndCost, m)
}

func recursiveMarkSeats(pos Coords, minCost int, m Map) Map {
	if m.grid[pos.row][pos.column].status == Wall {
		return m
	}

	if m.grid[pos.row][pos.column].northFacingCost == minCost {
		m.grid[pos.row][pos.column].status = Seat
		m = recursiveMarkSeats(Coords{row: pos.row + 1, column: pos.column}, minCost-forwardPoints, m)
	}
	if m.grid[pos.row][pos.column].eastFacingCost == minCost {
		m.grid[pos.row][pos.column].status = Seat
		m = recursiveMarkSeats(Coords{row: pos.row, column: pos.column - 1}, minCost-forwardPoints, m)
	}
	if m.grid[pos.row][pos.column].southFacingCost == minCost {
		m.grid[pos.row][pos.column].status = Seat
		m = recursiveMarkSeats(Coords{row: pos.row - 1, column: pos.column}, minCost-forwardPoints, m)
	}
	if m.grid[pos.row][pos.column].westFacingCost == minCost {
		m.grid[pos.row][pos.column].status = Seat
		m = recursiveMarkSeats(Coords{row: pos.row, column: pos.column + 1}, minCost-forwardPoints, m)
	}

	if m.grid[pos.row][pos.column].northFacingCost == minCost-turnPoints {
		m.grid[pos.row][pos.column].status = Seat
		m = recursiveMarkSeats(Coords{row: pos.row + 1, column: pos.column}, minCost-turnPoints-forwardPoints, m)
	}
	if m.grid[pos.row][pos.column].eastFacingCost == minCost-turnPoints {
		m.grid[pos.row][pos.column].status = Seat
		m = recursiveMarkSeats(Coords{row: pos.row, column: pos.column - 1}, minCost-turnPoints-forwardPoints, m)
	}
	if m.grid[pos.row][pos.column].southFacingCost == minCost-turnPoints {
		m.grid[pos.row][pos.column].status = Seat
		m = recursiveMarkSeats(Coords{row: pos.row - 1, column: pos.column}, minCost-turnPoints-forwardPoints, m)
	}
	if m.grid[pos.row][pos.column].westFacingCost == minCost-turnPoints {
		m.grid[pos.row][pos.column].status = Seat
		m = recursiveMarkSeats(Coords{row: pos.row, column: pos.column + 1}, minCost-turnPoints-forwardPoints, m)
	}
	return m
}

func part1() {
	startMap := parseMap("resources/Day16/input.txt")
	printMap(startMap)

	startMap = fillWalls(startMap)
	printMap(startMap)

	startMap.grid[startMap.position.row][startMap.position.column].status = Empty

	filledMap := progressiveFill(startMap)
	printMap(filledMap)

	seatCount := 0
	for _, row := range filledMap.grid {
		for _, b := range row {
			if b.status == Seat {
				seatCount++
			}
		}
	}
	fmt.Println("Seat count: ", seatCount)
}

func main() {
	part1()
}
