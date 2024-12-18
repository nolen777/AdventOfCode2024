package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Coords struct {
	x int
	y int
}

type Maze [][]rune

const (
	Empty     = '.'
	Corrupted = '#'
	Path      = 'O'
)

func parseCoords(fileName string) []Coords {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	coords := []Coords{}

	for scanner.Scan() {
		line := scanner.Text()
		nums := strings.Split(line, ",")
		x, _ := strconv.Atoi(nums[0])
		y, _ := strconv.Atoi(nums[1])
		coords = append(coords, Coords{x: x, y: y})
	}
	return coords
}

func SetupMaze(size int, corruption []Coords) Maze {
	maze := make([][]rune, 0, size)
	for i := 0; i < size; i++ {
		line := make([]rune, 0, size)
		for j := 0; j < size; j++ {
			line = append(line, Empty)
		}
		maze = append(maze, line)
	}

	for _, c := range corruption {
		maze[c.x][c.y] = Corrupted
	}
	return maze
}

func PrintMaze(m Maze) {
	size := len(m)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			fmt.Print(string(m[x][y]))
		}
		fmt.Println("")
	}
}

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

func CalculateCosts(from Coords, m Maze) [][]int {
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

	costs[from.x][from.y] = 0

	changes := 1
	// Go through every tile repeatedly. On each tile, see if we've identified
	// a cheaper path to it
	for changes > 0 {
		changes = 0
		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				if m[x][y] == Corrupted {
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

//func BestRoute(to Coords, m Maze, costs [][]int) Maze {
//	routeMaze := make([][]rune, 0, len(m))
//	for _, r := range m {
//		routeMaze = append(routeMaze, slices.Clone(r))
//	}
//
//	pos := to
//	for costs[pos.x][pos.y] > 0 {
//		routeMaze[pos.x][pos.y] = Path
//
//	}
//}

func part1() {
	coords := parseCoords("resources/Day18/input.txt")
	corruptionCount := 1024
	size := 71
	fmt.Println(coords)

	maze := SetupMaze(size, coords[:corruptionCount])
	PrintMaze(maze)

	costs := CalculateCosts(Coords{x: 0, y: 0}, maze)
	fmt.Println("Steps: ", costs[size-1][size-1])
}

func part2() {
	coords := parseCoords("resources/Day18/input.txt")
	//corruptionCount := 1024
	size := 71
	fmt.Println(coords)

	low := 0
	high := len(coords)

	for low+1 < high {
		count := (low + high) / 2
		maze := SetupMaze(size, coords[:count])

		costs := CalculateCosts(Coords{x: 0, y: 0}, maze)
		endCost := costs[size-1][size-1]
		if endCost < MaxInt {
			fmt.Println("Success!")
			low = count
		} else {
			fmt.Println("Failed!")
			high = count
		}
	}

	maze := SetupMaze(size, coords[:low])
	PrintMaze(maze)
	costs := CalculateCosts(Coords{x: 0, y: 0}, maze)
	endCost := costs[size-1][size-1]
	_ = endCost

	fmt.Println("Index: ", low)
	fmt.Println("Coords: ", coords[low])

}

func main() {
	//	part1()
	part2()
}
