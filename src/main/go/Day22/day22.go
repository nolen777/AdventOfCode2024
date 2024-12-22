package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func parseSecrets(fileName string) []int {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	nums := []int{}

	for scanner.Scan() {
		line := scanner.Text()
		i, err := strconv.Atoi(line)
		if err != nil {
			log.Fatal("Invalid input")
		}
		nums = append(nums, i)
	}
	return nums
}

func mix(secret int, value int) int {
	return secret ^ value
}

func prune(secret int) int {
	return secret % 16777216
}

func evolve(secret int) int {
	secret = prune(mix(secret, secret<<6))
	secret = prune(mix(secret, secret>>5))
	secret = prune(mix(secret, secret<<11))

	return secret
}

func part1() {
	secrets := parseSecrets("resources/Day22/sample.txt")
	fmt.Println(secrets)

	sum := 0
	for _, secret := range secrets {
		newSecret := secret
		for i := 0; i < 2000; i++ {
			newSecret = evolve(newSecret)
		}
		fmt.Println(secret, ": ", newSecret)
		sum += newSecret
	}

	fmt.Println("Sum: ", sum)
}

func part2() {
	lines := parseSecrets("resources/Day2/sample.txt")
	_ = lines
}

func main() {
	part1()
}
