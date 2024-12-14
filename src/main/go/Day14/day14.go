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

func robotPositions(robots []RobotInfo, size CPair) [][]int {
	grid := make([][]int, 0, size.y)
	for i := 0; i < size.y; i++ {
		grid = append(grid, make([]int, size.x))
	}

	for _, robot := range robots {
		grid[robot.position.y][robot.position.x]++
	}

	return grid
}

func printRobots(robots []RobotInfo, size CPair) {
	grid := robotPositions(robots, size)

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
	size := CPair{101, 103}
	robots := parseRobots("resources/day14/input.txt")
	printRobots(robots, size)
	fmt.Println("")
	newRobots := moveRobots(robots, 100, size)
	printRobots(newRobots, size)

	fmt.Println("safety factor: ", safetyFactor(newRobots, size))
}

func part2() {
	size := CPair{101, 103}
	robots := parseRobots("resources/day14/input.txt")
	count := 0
	for {
		count++
		robots = moveRobots(robots, 1, size)
		//printRobots(robots, size)
		grid := robotPositions(robots, size)

		found := false
		for y := 0; y < size.y; y++ {
			failed := false
			for x := size.x/2 - 3; x < size.x/2+3; x++ {
				if grid[y][x] == 0 {
					failed = true
					break
				}
			}
			if !failed {
				found = true
				break
			}
		}

		if found {
			printRobots(robots, size)
			fmt.Println("count is ", count)
			break
		}

		//for y := 0; y < 3; y++ {
		//	if grid[y][5] > 0 && grid[y+1][4] > 0 && grid[y+1][6] > 0 {
		//		printRobots(robots, size)
		//	}
		//}
		//for y := 80; y < 103; y++ {
		//	fail := false
		//	for x := 45; x < 55; x++ {
		//		if grid[y][x] == 0 {
		//			fail = true
		//			break
		//		}
		//	}
		//	if !fail {
		//		printRobots(robots, size)
		//		return
		//	}
		//}
	}
}

func main() {
	//part1()
	part2()
}
