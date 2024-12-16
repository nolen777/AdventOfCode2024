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

// Returns the point total and a "success" flag
func recursiveTry(m Map, points int, previousTurnDirection rune) (int, Map, bool) {
	currentTile := m.grid[m.position.row][m.position.column]
	if currentTile == End {
		return points, m, true
	}
	if currentTile == Wall {
		// Hit a wall, skip
		return points, m, false
	}
	if currentTile != Empty {
		// We've done this before; skip
		return points, m, false
	}
	if found && points > lowestPoints {
		return points, m, false
	}

	// move forward
	forwardMap := m
	forwardMap.grid[forwardMap.position.row][forwardMap.position.column] = m.direction
	switch m.direction {
	case North:
		forwardMap.position.row--
	case East:
		forwardMap.position.column++
	case South:
		forwardMap.position.row++
	case West:
		forwardMap.position.column--
	}
	forwardPoints, forwardMap, success := recursiveTry(forwardMap, points+forwardPoints, 0)
	if success {
		found = true
		if forwardPoints < lowestPoints {
			bestMap = copyMap(forwardMap)
			lowestPoints = forwardPoints
		}
	}
	m.grid[m.position.row][m.position.column] = Empty

	// Never turn twice in a row
	if previousTurnDirection == 0 && lowestPoints > points+turnPoints+forwardPoints {
		clockwise := clockwiseTurn(m.direction)
		counterClockwise := counterclockwiseTurn(m.direction)

		tryTurn := func(dir rune) {
			m.direction = dir
			pts, clMap, success := recursiveTry(m, points+turnPoints, 1)
			if success {
				found = true
				if pts < lowestPoints {
					bestMap = clMap
					lowestPoints = pts
				}
			}
		}

		if clockwise == North || clockwise == East {
			tryTurn(clockwise)
			tryTurn(counterClockwise)
		} else {
			tryTurn(counterClockwise)
			tryTurn(clockwise)
		}
	}

	return lowestPoints, bestMap, found
}

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
func progressiveFill(m Map) int {
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
		//fmt.Println("Changed ", changes, " tiles")
	}
	return minEndCost
}

func part1() {
	// Try spreading out costs
	startMap := parseMap("resources/Day16/input.txt")
	printMap(startMap)

	startMap = fillWalls(startMap)
	printMap(startMap)

	startMap.grid[startMap.position.row][startMap.position.column] = Empty

	progressiveFill(startMap)

	//	points, bestMap, success := recursiveTry(startMap, 0, 0)

	//if success {
	//	printMap(bestMap)
	//	fmt.Println(points)
	//} else {
	//	fmt.Println("Failed!")
	//}
}

func main() {
	part1()
}