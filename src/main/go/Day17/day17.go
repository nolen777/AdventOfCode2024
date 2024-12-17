package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func parseFile(fileName string) []string {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lines := make([]string, 0, 0)
	for scanner.Scan() {
		line := scanner.Text()

		lines = append(lines, line)
	}

	return lines
}

func part1() {
	lines := parseFile("resources/Day17/sample.txt")
	fmt.Println(lines)
}

func part2() {

}

func main() {
	part1()
}
