package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func parseLines(fileName string) []string {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines
}

func part1() {
	lines := parseLines("resources/Day22/sample.txt")
	fmt.Println(lines)
}

func part2() {
	lines := parseLines("resources/Day2/sample.txt")
	_ = lines
}

func main() {
	part1()
}
