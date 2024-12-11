package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func parseNums(fileName string) []int {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	nums := []int{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		textLine := scanner.Text()

		elements := strings.Split(textLine, " ")

		for _, elt := range elements {
			num, err := strconv.Atoi(elt)

			if err != nil {
				fmt.Println(err)
				return nums
			}
			nums = append(nums, num)
		}
		break
	}
	return nums
}

func printStones(stones []int) {
	for _, e := range stones {
		fmt.Print(e, " ")
	}
	fmt.Println("")
}

func digitCount(n int) int {
	return len(strconv.Itoa(n))
}

func split(n int) (int, int) {
	str := strconv.Itoa(n)
	first, err := strconv.Atoi(str[:len(str)/2])
	second, err := strconv.Atoi(str[len(str)/2:])

	if err != nil {
		log.Fatal("Bad atoi")
	}

	return first, second
}

func blink(stones []int) []int {
	idx := 0

	for idx < len(stones) {
		e := stones[idx]

		switch {
		case e == 0:
			stones[idx] = 1
		case digitCount(e)%2 == 0:
			// split
			newStones := make([]int, 0, len(stones)+digitCount(e)/2)
			newStones = append(newStones, stones[:idx]...)
			first, second := split(e)
			newStones = append(newStones, first, second)
			newStones = append(newStones, stones[idx+1:]...)

			stones = newStones
			idx++

		default:
			stones[idx] = e * 2024
		}
		idx++
	}

	return stones
}

func main() {
	stones := parseNums("resources/Day11/sampleinput.txt")

	printStones(stones)

	for i := 0; i < 25; i++ {
		stones = blink(stones)
	}

	fmt.Println(len(stones), " stones")
}
