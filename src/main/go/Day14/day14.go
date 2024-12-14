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

func mod(a int, b int) int {
	return (a%b + b) % b
}

func moveRobots(robots []RobotInfo, stepCount int, size CPair) []RobotInfo {
	newRobots := make([]RobotInfo, 0, len(robots))

	for _, robot := range robots {
		newRobot := robot
		newRobot.position.x = mod(newRobot.position.x+stepCount*newRobot.velocity.x, size.x)
		newRobot.position.y = mod(newRobot.position.y+stepCount*newRobot.velocity.y, size.y)
		newRobots = append(newRobots, newRobot)
	}

	return newRobots
}

func safetyFactor(robots []RobotInfo, size CPair) int {
	quads := [4]int{}

	for _, robot := range robots {
		if robot.position.y < size.y/2 {
			if robot.position.x < size.x/2 {
				quads[0]++
			} else if robot.position.x > size.x/2 {
				quads[1]++
			}
		} else if robot.position.y > size.y/2 {
			if robot.position.x < size.x/2 {
				quads[2]++
			} else if robot.position.x > size.x/2 {
				quads[3]++
			}
		}
	}

	return quads[0] * quads[1] * quads[2] * quads[3]
}

func part1() {
	size := CPair{11, 7}
	robots := parseRobots("resources/day14/sample.txt")
	printRobots(robots, size)
	fmt.Println("")
	newRobots := moveRobots(robots, 100, size)
	printRobots(newRobots, size)

	fmt.Println("safety factor: ", safetyFactor(newRobots, size))
}

func main() {
	part1()
}
