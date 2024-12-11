package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
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

type StoneInfo struct {
	value int
	count int
}

func blink(stones []StoneInfo) []StoneInfo {
	idx := 0
	startLen := len(stones)

	for idx < startLen {
		e := stones[idx]
		dc := digitCount(e.value)

		switch {
		case e.value == 0:
			stones[idx].value = 1
		case dc%2 == 0:
			// split
			first, second := split(e.value, dc)
			stones[idx].value = first
			stones = append(stones, StoneInfo{value: second, count: stones[idx].count})

		default:
			stones[idx].value = e.value * 2024
		}
		idx++
	}

	return stones
}

func consolidate(stones []StoneInfo) []StoneInfo {
	cS := make([]StoneInfo, 0, len(stones))

	sort.Slice(stones, func(i int, j int) bool {
		return stones[i].value < stones[j].value
	})

	idx := 0
	csLen := 0
	for idx < len(stones) {
		newInfo := stones[idx]
		if csLen == 0 || newInfo.value != cS[csLen-1].value {
			cS = append(cS, newInfo)
			csLen++
		} else {
			cS[csLen-1].count += newInfo.count
		}
		idx++
	}
	return cS
}

func countStones(stones []StoneInfo) int {
	totalCount := 0
	for _, stone := range stones {
		totalCount += stone.count
	}
	return totalCount
}

func main() {
	initialNums := parseNums("resources/Day11/input.txt")
	printStones(initialNums)

	stones := consolidate(make([]StoneInfo, 0, len(initialNums)))
	for _, num := range initialNums {
		stones = append(stones, StoneInfo{value: num, count: 1})
	}
	stones = consolidate(stones)

	for i := 0; i < 75; i++ {
		stones = consolidate(blink(stones))
		fmt.Println(i, " unique count: ", len(stones))
		//	fmt.Println(stones)
	}

	fmt.Println(countStones(stones), " stones")
}
