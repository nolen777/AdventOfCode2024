package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func ReadFile(filename string) string {
	b, err := os.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func main() {
	str := ReadFile("resources/Day3/input.txt")

	pattern := `mul\((?P<first>\d{1,4}),(?P<second>\d{1,4})\)|don't\(\)|do\(\)`

	r, _ := regexp.Compile(pattern)

	total := 0
	enabled := true
	for _, m := range r.FindAllStringSubmatch(str, -1) {
		fmt.Println(m)
		switch m[0] {
		case `do()`:
			fmt.Println("Found do()")
			enabled = true
		case `don't()`:
			fmt.Println("Found don't()")
			enabled = false
		default:
			if enabled {
				first, _ := strconv.Atoi(m[1])
				second, _ := strconv.Atoi(m[2])
				total += first * second
			}
		}
	}

	fmt.Println("total ", total)
}
