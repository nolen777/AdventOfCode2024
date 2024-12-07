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

func possibleValues(startValue int, sequence []int) []int {
	allValues := []int{}

	l := len(sequence)
	if l == 0 {
		allValues = []int{startValue}
	} else {
		innerValues := possibleValues(sequence[l-1], sequence[:l-1])
		for _, x := range innerValues {
			allValues = append(allValues, x+startValue)
			if startValue != 0 {
				allValues = append(allValues, x*startValue)
			}
		}
	}

	return allValues
}

func main() {
	equations := parse("resources/Day7/sampleinput.txt")

	total := 0
	for _, e := range equations {
		pv := possibleValues(0, e.sequence)
		for _, v := range pv {
			if v == e.testValue {
				fmt.Println(e, ": Success!")
				total += e.testValue
				break
			}
		}
	}

	fmt.Println("total: ", total)
}
