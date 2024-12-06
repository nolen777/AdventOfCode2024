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

type MapState struct {
	mapBytes      [][]byte
	currentRow    int
	currentColumn int
	direction     int
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

	return MapState{mapBytes: mapBytes, currentRow: currentRow, currentColumn: currentColumn, direction: direction}
}

func advanceOne(initial MapState) MapState {
	newRow := initial.currentRow
	newColumn := initial.currentColumn
	newDirection := initial.direction
	newBytes := initial.mapBytes

	switch initial.direction {
	case 0:
		newRow -= 1
	case 1:
		newColumn += 1
	case 2:
		newRow += 1
	case 3:
		newColumn -= 1
	}

	newBytes[initial.currentRow][initial.currentColumn] = 'X'

	if newRow >= 0 && newColumn >= 0 && newRow < len(newBytes) && newColumn < len(newBytes[0]) {
		if initial.mapBytes[newRow][newColumn] == '#' {
			newDirection = (newDirection + 1) % 4
			newRow = initial.currentRow
			newColumn = initial.currentColumn
		}
	}

	return MapState{mapBytes: newBytes, currentRow: newRow, currentColumn: newColumn, direction: newDirection}
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

func main() {
	mapBytes := parseMap("resources/Day6/sampleinput.txt")

	mapState := setupInitial(mapBytes)

	for mapState.currentRow >= 0 && mapState.currentRow < len(mapState.mapBytes) && mapState.currentColumn >= 0 && mapState.currentColumn < len(mapState.mapBytes[0]) {
		mapState = advanceOne(mapState)

		for _, line := range mapState.mapBytes {
			for _, b := range line {
				fmt.Print(string(b))
			}
			fmt.Println("")
		}
		fmt.Println("")
	}

	fmt.Println("Total ", countX(mapState.mapBytes))
}
