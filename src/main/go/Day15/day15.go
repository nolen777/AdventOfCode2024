package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	Wall  = '#'
	Empty = '.'
	Box   = 'O'
	Robot = '@'
)

const (
	Up    = '^'
	Right = '>'
	Left  = '<'
	Down  = 'v'
)

type WarehouseMap struct {
	grid                  [][]int
	robotRow              int
	robotColumn           int
	remainingInstructions []int
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
		row := []int{}
		for columnIndex, b := range line {
			row = append(row, int(b))
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
			wm.remainingInstructions = append(wm.remainingInstructions, int(b))
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

func moveObjectInDirection(grid [][]int, column int, row int, columnDiff int, rowDiff int) bool {
	currentObject := grid[column][row]

	if currentObject != Box && currentObject != Robot {
		log.Fatal("Can only move robots and boxes")
	}

	switch grid[column+columnDiff][row+rowDiff] {
	case Wall:
		return false
	case Empty:
		grid[column+columnDiff][row+rowDiff] = currentObject
		grid[column][row] = Empty
		return true
	case Box:
		if moveObjectInDirection(grid, column+columnDiff, row+rowDiff, columnDiff, rowDiff) {
			grid[column+columnDiff][row+rowDiff] = currentObject
			grid[column][row] = Empty
			return true
		}
		return false
	case Robot:
		log.Fatal("Pushing into robot?!")
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

	if moveObjectInDirection(warehouseMap.grid, warehouseMap.robotRow, warehouseMap.robotColumn, rowDiff, columnDiff) {
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
			if b == Box {
				total += 100*ri + ci
			}
		}
	}
	return total
}

func part1() {
	wm := parseMap("resources/Day15/sample.txt")
	printMap(wm)

	for len(wm.remainingInstructions) > 0 {
		wm = executeOneMove(wm)
		printMap(wm)
	}

	fmt.Println("GPS total: ", gpsSum(wm))
}

func main() {
	part1()
}
