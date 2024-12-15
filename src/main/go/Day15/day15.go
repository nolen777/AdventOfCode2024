package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	Wall     = rune('#')
	Empty    = rune('.')
	Box      = rune('O')
	BoxLeft  = rune('[')
	BoxRight = rune(']')
	Robot    = rune('@')
)

const (
	Up    = rune('^')
	Right = rune('>')
	Left  = rune('<')
	Down  = rune('v')
)

type WarehouseMap struct {
	grid                  [][]rune
	robotRow              int
	robotColumn           int
	remainingInstructions []rune
}

func parseMap(fileName string) WarehouseMap {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	wm := WarehouseMap{}

	scanner := bufio.NewScanner(file)

	rowIndex := 0
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			break
		}
		row := []rune{}
		for columnIndex, b := range line {
			row = append(row, b)
			if b == Robot {
				wm.robotRow = rowIndex
				wm.robotColumn = columnIndex
			}
		}

		wm.grid = append(wm.grid, row)
		rowIndex++
	}

	for scanner.Scan() {
		line := scanner.Text()
		for _, b := range line {
			wm.remainingInstructions = append(wm.remainingInstructions, b)
		}
	}
	return wm
}

func parseWideMap(fileName string) WarehouseMap {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	wm := WarehouseMap{}

	scanner := bufio.NewScanner(file)

	rowIndex := 0
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			break
		}
		row := []rune{}
		for columnIndex, b := range line {
			switch b {
			case Wall:
				row = append(row, Wall, Wall)
			case Box:
				row = append(row, BoxLeft, BoxRight)
			case Empty:
				row = append(row, Empty, Empty)
			case Robot:
				row = append(row, Robot, Empty)
				wm.robotRow = rowIndex
				wm.robotColumn = 2 * columnIndex
			}
		}

		wm.grid = append(wm.grid, row)
		rowIndex++
	}

	for scanner.Scan() {
		line := scanner.Text()
		for _, b := range line {
			wm.remainingInstructions = append(wm.remainingInstructions, b)
		}
	}
	return wm
}

func printMap(warehouseMap WarehouseMap) {
	for _, row := range warehouseMap.grid {
		for _, b := range row {
			fmt.Print(string(rune(b)))
		}
		fmt.Println("")
	}

	fmt.Println("")
	for _, i := range warehouseMap.remainingInstructions {
		fmt.Print(string(rune(i)))
	}
	fmt.Println("")
}

func canMoveBigBoxInDirection(grid [][]rune, leftColumn int, row int, rowDiff int) bool {
	currentObject := grid[row][leftColumn]

	if currentObject != BoxLeft || grid[row][leftColumn+1] != BoxRight {
		log.Fatal("Invalid box setup")
	}

	// Check up/down from the left side
	switch grid[row+rowDiff][leftColumn] {
	case Wall:
		return false
	case Empty:

	case BoxLeft:
		if !canMoveBigBoxInDirection(grid, leftColumn, row+rowDiff, rowDiff) {
			return false
		}
	case BoxRight:
		if !canMoveBigBoxInDirection(grid, leftColumn-1, row+rowDiff, rowDiff) {
			return false
		}
	case Robot:
		log.Fatal("Pushing into robot?!")
	}

	switch grid[row+rowDiff][leftColumn+1] {
	case Wall:
		return false
	case Empty:

	case BoxLeft:
		if !canMoveBigBoxInDirection(grid, leftColumn+1, row+rowDiff, rowDiff) {
			return false
		}
	case BoxRight:
		if !canMoveBigBoxInDirection(grid, leftColumn, row+rowDiff, rowDiff) {
			return false
		}
	case Robot:
		log.Fatal("Pushing into robot?!")
	}

	return true
}

func moveObjectInDirection(grid [][]rune, column int, row int, columnDiff int, rowDiff int) bool {
	currentObject := grid[row][column]

	if currentObject == Empty {
		return true
	}
	if currentObject != Box && currentObject != Robot && currentObject != BoxLeft && currentObject != BoxRight {
		log.Fatal("Can only move robots and boxes")
	}

	// Part 1, or robots for part 2, or strict left/right
	if currentObject == Box || currentObject == Robot || rowDiff == 0 {
		switch grid[row+rowDiff][column+columnDiff] {
		case Wall:
			return false
		case Empty:
			grid[row+rowDiff][column+columnDiff] = currentObject
			grid[row][column] = Empty
			return true
		case Box, BoxLeft, BoxRight:
			if moveObjectInDirection(grid, column+columnDiff, row+rowDiff, columnDiff, rowDiff) {
				grid[row+rowDiff][column+columnDiff] = currentObject
				grid[row][column] = Empty
				return true
			}
			return false
		case Robot:
			log.Fatal("Pushing into robot?!")
		}

		log.Fatal("Impossible condition")
		return false
	}

	// Part 2, up and down big boxes
	if currentObject == BoxLeft {
		if canMoveBigBoxInDirection(grid, column, row, rowDiff) {
			moveObjectInDirection(grid, column, row+rowDiff, columnDiff, rowDiff)
			moveObjectInDirection(grid, column+1, row+rowDiff, columnDiff, rowDiff)
			grid[row+rowDiff][column] = BoxLeft
			grid[row+rowDiff][column+1] = BoxRight
			grid[row][column] = Empty
			grid[row][column+1] = Empty
			return true
		}
		return false
	}
	if currentObject == BoxRight {
		return moveObjectInDirection(grid, column-1, row, columnDiff, rowDiff)
	}

	log.Fatal("Impossible condition")
	return false
}

func executeOneMove(warehouseMap WarehouseMap) WarehouseMap {
	instruction := warehouseMap.remainingInstructions[0]

	newRow := warehouseMap.robotRow
	newColumn := warehouseMap.robotColumn

	var columnDiff, rowDiff int

	switch instruction {
	case Up:
		rowDiff = -1
	case Right:
		columnDiff = 1
	case Down:
		rowDiff = 1
	case Left:
		columnDiff = -1
	}

	if moveObjectInDirection(warehouseMap.grid, warehouseMap.robotColumn, warehouseMap.robotRow, columnDiff, rowDiff) {
		newRow += rowDiff
		newColumn += columnDiff
	}

	return WarehouseMap{
		grid:                  warehouseMap.grid,
		robotRow:              newRow,
		robotColumn:           newColumn,
		remainingInstructions: warehouseMap.remainingInstructions[1:],
	}
}

func gpsSum(warehouseMap WarehouseMap) int {
	total := 0
	for ri, row := range warehouseMap.grid {
		for ci, b := range row {
			if b == Box || b == BoxLeft {
				total += 100*ri + ci
			}
		}
	}
	return total
}

func part1() {
	wm := parseMap("resources/Day15/input.txt")
	printMap(wm)

	for len(wm.remainingInstructions) > 0 {
		wm = executeOneMove(wm)
		//	printMap(wm)
	}

	fmt.Println("GPS total: ", gpsSum(wm))
}

func part2() {
	wm := parseWideMap("resources/Day15/sample.txt")
	printMap(wm)

	for len(wm.remainingInstructions) > 0 {
		wm = executeOneMove(wm)
		printMap(wm)
	}

	fmt.Println("GPS total: ", gpsSum(wm))
}

func main() {
	//part1()
	part2()
}
