package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Machine struct {
	aX     int
	bX     int
	prizeX int
	aY     int
	bY     int
	prizeY int
}

func coordsFrom(line string, startPos int) (int, int) {
	coordsString := line[startPos:]
	nums := strings.Split(coordsString, ", Y")
	x, _ := strconv.Atoi(nums[0])
	y, _ := strconv.Atoi(nums[1][1:])

	return x, y
}

func parseOneMachine(lines [3]string) Machine {
	var m Machine

	m.aX, m.aY = coordsFrom(lines[0], 12)
	m.bX, m.bY = coordsFrom(lines[1], 12)
	m.prizeX, m.prizeY = coordsFrom(lines[2], 9)

	return m
}

func parseMachines(fileName string) []Machine {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ms := []Machine{}

	scanner := bufio.NewScanner(file)

	lines := [3]string{}
	pos := 0
	for scanner.Scan() {
		next := scanner.Text()
		if next == "" {
			continue
		}
		lines[pos] = next
		pos++
		if pos == 3 {
			ms = append(ms, parseOneMachine(lines))
			pos = 0
		}
	}
	return ms
}

func solveOne(m Machine) (int, int, bool) {
	aXNorm := m.aX * m.aY
	bXNorm := m.bX * m.aY
	prizeXNorm := m.prizeX * m.aY

	aYNorm := m.aY * m.aX
	bYNorm := m.bY * m.aX
	prizeYNorm := m.prizeY * m.aX

	if aXNorm != aYNorm {
		log.Fatal("WHAT")
	}

	if (prizeXNorm-prizeYNorm)%(bXNorm-bYNorm) != 0 {
		return 0, 0, false
	}
	bVal := (prizeXNorm - prizeYNorm) / (bXNorm - bYNorm)

	if (prizeXNorm-(bXNorm*bVal))%aXNorm != 0 {
		return 0, 0, false
	}
	aVal := (prizeXNorm - (bXNorm * bVal)) / aXNorm

	return aVal, bVal, true
}

func main() {
	machines := parseMachines("resources/Day13/input.txt")
	fmt.Println(machines)

	const aCost = 3
	const bCost = 1

	//
	// PART 1
	//

	fmt.Println("PART 1")
	totalCost := 0
	for _, m := range machines {
		a, b, ok := solveOne(m)
		if ok {
			fmt.Println("Found solution ", a, ", ", b)
			totalCost += a*aCost + b*bCost
		} else {
			fmt.Println("No solution.")
		}
	}
	fmt.Println("Total cost: ", totalCost)

	//
	// PART 2
	//

	fmt.Println("PART 2")
	const addVal = 10000000000000
	for i := 0; i < len(machines); i++ {
		machines[i].prizeX += addVal
		machines[i].prizeY += addVal
	}

	totalCost = 0
	for _, m := range machines {
		a, b, ok := solveOne(m)
		if ok {
			fmt.Println("Found solution ", a, ", ", b)
			totalCost += a*aCost + b*bCost
		} else {
			fmt.Println("No solution.")
		}
	}
	fmt.Println("Total cost: ", totalCost)

}
