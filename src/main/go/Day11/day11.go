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
	if n == 0 {
		return 1
	}

	count := 0
	for n != 0 {
		n /= 10
		count++
	}
	return count
}

func split(n int, digitCount int) (int, int) {
	m := 0
	mult := 1
	for i := 0; i < digitCount/2; i++ {
		d := n % 10
		m += d * mult
		n /= 10
		mult *= 10
	}
	return n, m
}

func blink(stones []int) []int {
	idx := 0
	startLen := len(stones)

	for idx < startLen {
		e := stones[idx]
		dc := digitCount(e)

		switch {
		case e == 0:
			stones[idx] = 1
		case dc%2 == 0:
			// split
			first, second := split(e, dc)
			stones[idx] = first
			stones = append(stones, second)

		default:
			stones[idx] = e * 2024
		}
		idx++
	}

	return stones
}

func main() {
	stones := parseNums("resources/Day11/input.txt")

	printStones(stones)

	for i := 0; i < 75; i++ {
		stones = blink(stones)
		fmt.Println(i)
		//	fmt.Println(stones)
	}

	fmt.Println(len(stones), " stones")
}
