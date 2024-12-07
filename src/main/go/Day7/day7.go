package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Equation struct {
	testValue int
	sequence  []int
}

func parse(filename string) []Equation {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	equations := []Equation{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		textLine := scanner.Text()
		firstSplit := strings.Split(textLine, ": ")

		testValue, _ := strconv.Atoi(firstSplit[0])

		secondSplit := strings.Split(firstSplit[1], " ")
		nums := []int{}

		for _, str := range secondSplit {
			num, _ := strconv.Atoi(str)
			nums = append(nums, num)
		}

		equations = append(equations, Equation{testValue: testValue, sequence: nums})
	}

	return equations
}

var sumCount int
var multCount int
var concatCount int

func isPossible(desiredValue int, sequence []int) bool {
	l := len(sequence)
	rightmostValue := sequence[l-1]

	if l == 1 {
		return desiredValue == rightmostValue
	} else {
		remainingSeq := sequence[:l-1]

		// add
		sumCount++
		if isPossible(desiredValue-rightmostValue, remainingSeq) {
			return true
		}

		// multiply. We only need to do this if the desired total is divisible by the rightmost element.
		if desiredValue%rightmostValue == 0 {
			multCount++
			if isPossible(desiredValue/rightmostValue, remainingSeq) {
				return true
			}
		}

		// concat. We'll make sure that the rightmost digits match, and recurse if so
		for rightmostValue > 0 {
			if (rightmostValue % 10) != (desiredValue % 10) {
				return false
			}
			rightmostValue /= 10
			desiredValue /= 10
		}
		concatCount++
		return isPossible(desiredValue, remainingSeq)
	}
}

func main() {
	equations := parse("resources/Day7/input.txt")

	total := 0
	for _, e := range equations {
		if isPossible(e.testValue, e.sequence) {
			fmt.Println(e, ": Success!")
			total += e.testValue
		}
	}

	fmt.Println("new total:", total)
	fmt.Println("Sum count:    ", sumCount)
	fmt.Println("Mult count:   ", multCount)
	fmt.Println("Concat count: ", concatCount)
}
