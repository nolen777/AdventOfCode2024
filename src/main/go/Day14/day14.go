package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type CPair struct {
	x int
	y int
}

type RobotInfo struct {
	position CPair
	velocity CPair
}

func parseRobots(fileName string) []RobotInfo {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	robots := []RobotInfo{}

	scanner := bufio.NewScanner(file)

	pattern := `p\=(?P<x>\-?\d*),(?P<y>\-?\d*) v\=(?P<dx>\-?\d*),(?P<dy>\-?\d*)`
	r, _ := regexp.Compile(pattern)

	for scanner.Scan() {
		line := scanner.Text()

		m := r.FindStringSubmatch(line)

		x, _ := strconv.Atoi(m[1])
		y, _ := strconv.Atoi(m[2])
		dx, _ := strconv.Atoi(m[3])
		dy, _ := strconv.Atoi(m[4])

		robots = append(robots, RobotInfo{position: CPair{x, y}, velocity: CPair{dx, dy}})
	}
	return robots
}

//func moveRobots(robots []RobotInfo, stepCount int, width int, height int) []RobotInfo {
//
//}

func printRobots(robots []RobotInfo, size CPair) {
	grid := make([][]int, 0, size.y)
	for i := 0; i < size.y; i++ {
		grid = append(grid, make([]int, size.x))
	}

	for _, robot := range robots {
		grid[robot.position.y][robot.position.x]++
	}

	for _, row := range grid {
		for _, c := range row {
			if c == 0 {
				fmt.Print(".")
			} else {
				fmt.Print(c)
			}
		}
		fmt.Println("")
	}
}

func part1() {
	size := CPair{11, 7}
	robots := parseRobots("resources/day14/sample.txt")
	printRobots(robots, size)
}

func main() {
	part1()
}
