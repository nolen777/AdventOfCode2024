package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func parseMap(filename string) [][]byte {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	mapBytes := [][]byte{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		textLine := scanner.Text()
		mapLine := make([]byte, 0)

		for _, r := range textLine {
			mapLine = append(mapLine, byte(r))
		}

		mapBytes = append(mapBytes, mapLine)
	}

	return mapBytes
}

type Cursor struct {
	row       int
	column    int
	direction int
}

type MapState struct {
	mapBytes [][]byte
	cursor   Cursor
	trail    []Cursor
}

func setupInitial(mapBytes [][]byte) MapState {
	var currentRow int
	var currentColumn int
	direction := 0

	for row, line := range mapBytes {
		for col, b := range line {
			switch b {
			case '^':
				currentRow = row
				currentColumn = col
				direction = 0
			case '>':
				currentRow = row
				currentColumn = col
				direction = 1
			case 'v':
				currentRow = row
				currentColumn = col
				direction = 2
			case '<':
				currentRow = row
				currentColumn = col
				direction = 3
			}
			fmt.Print(string(b))
		}
		fmt.Println("")
	}

	return MapState{mapBytes: mapBytes, cursor: Cursor{row: currentRow, column: currentColumn, direction: direction}}
}

func advanceOne(initial MapState) MapState {
	newCursor := initial.cursor
	newBytes := initial.mapBytes
	newTrail := make([]Cursor, len(initial.trail))
	copy(newTrail, initial.trail)

	switch initial.cursor.direction {
	case 0:
		newCursor.row -= 1
	case 1:
		newCursor.column += 1
	case 2:
		newCursor.row += 1
	case 3:
		newCursor.column -= 1
	}

	newBytes[initial.cursor.row][initial.cursor.column] = 'X'

	if positionIsValid(newBytes, newCursor) {
		destByte := initial.mapBytes[newCursor.row][newCursor.column]
		if destByte == '#' || destByte == 'O' {
			newTrail = append(newTrail, initial.cursor)
			newCursor.direction = (newCursor.direction + 1) % 4
			newCursor.row = initial.cursor.row
			newCursor.column = initial.cursor.column
		}
	}

	return MapState{mapBytes: newBytes, cursor: newCursor, trail: newTrail}
}

func countX(mapBytes [][]byte) int {
	total := 0
	for _, line := range mapBytes {
		for _, b := range line {
			if b == 'X' {
				total += 1
			}
		}
	}

	return total
}

func positionIsValid(mapBytes [][]byte, cursor Cursor) bool {
	return cursor.row >= 0 && cursor.column >= 0 && cursor.row < len(mapBytes) && cursor.column < len(mapBytes[0])
}

func printState(mapBytes [][]byte) {
	for _, line := range mapBytes {
		for _, b := range line {
			fmt.Print(string(b))
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func hasLoop(mapState MapState) bool {
	for _, oldCursor := range mapState.trail {
		if oldCursor == mapState.cursor {
			return true
		}
	}
	return false
}

func mapBytesCopy(mapBytes [][]byte) [][]byte {
	newBytes := [][]byte{}
	for _, line := range mapBytes {
		var lineCopy []byte
		lineCopy = append(lineCopy, line...)
		newBytes = append(newBytes, lineCopy)
	}
	return newBytes
}

func main() {
	mapBytes := parseMap("resources/Day6/input.txt")

	mapState := setupInitial(mapBytes)

	loopCount := 0
	for r := 0; r < len(mapBytes); r++ {
		for c := 0; c < len(mapBytes[r]); c++ {
			if r == mapState.cursor.row && c == mapState.cursor.column {
				continue
			}

			mapCopy := mapBytesCopy(mapState.mapBytes)
			mapCopy[r][c] = 'O'

			newMapState := MapState{mapBytes: mapCopy, cursor: mapState.cursor}

			for positionIsValid(newMapState.mapBytes, newMapState.cursor) {
				newMapState = advanceOne(newMapState)
				if hasLoop(newMapState) {
					fmt.Println("Found a loop!")
					loopCount++
					break
				}
			}
		}
	}

	fmt.Println("Total ", loopCount)
}
