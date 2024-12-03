package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func ReadReports(filename string) [][]int {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var reports [][]int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		textLine := scanner.Text()
		reportStrings := strings.Split(textLine, " ")
		report := []int{}

		for _, entryString := range reportStrings {
			entry, err := strconv.Atoi(entryString)
			if err != nil {
				log.Fatal(err)
			}
			report = append(report, entry)
		}

		reports = append(reports, report)
	}

	return reports
}

func isSafe(report []int) bool {
	decreasing := report[1] < report[0]

	for idx := 0; idx < (len(report) - 1); idx++ {
		cur := report[idx]
		next := report[idx+1]

		if cur == next {
			// can't be the same
			return false
		}
		if (next < cur) != decreasing {
			// must be constantly increasing/decreasing
			return false
		}
		if cur-next < -3 {
			// too far apart
			return false
		}
		if cur-next > 3 {
			// too far apart
			return false
		}
	}

	return true
}

func main() {
	reports := ReadReports("resources/Day2/input.txt")

	successCount := 0

	for lineNo, report := range reports {
		isSuccess := isSafe(report)

		if isSuccess {
			fmt.Println(lineNo, " success!")
			successCount++
			continue
		}

		fmt.Println(lineNo, " failure!")
		// Try removing elements to see if it's safe now

		for idx := 0; idx < len(report); idx++ {
			oneRemoved := []int{}
			oneRemoved = append(oneRemoved, report[:idx]...)
			oneRemoved = append(oneRemoved, report[idx+1:]...)

			if isSafe(oneRemoved) {
				fmt.Println("Made it safe by removing index ", idx)
				successCount++
				break
			}
		}
	}

	fmt.Println("Success lines: ", successCount)
}
