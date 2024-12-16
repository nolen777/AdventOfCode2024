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
				r = Empty
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

	found := false
	lowestPoints := MaxInt
	var bestMap Map

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

func part1() {
	startMap := parseMap("resources/Day16/sample.txt")
	printMap(startMap)

	points, bestMap, success := recursiveTry(startMap, 0, 0)

	if success {
		printMap(bestMap)
		fmt.Println(points)
	} else {
		fmt.Println("Failed!")
	}
}

func main() {
	part1()
}
