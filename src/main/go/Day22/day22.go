package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
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

type Seq4 struct {
	a, b, c, d int
}

func (s4 Seq4) sum() int {
	return s4.a + s4.b + s4.c + s4.d
}

func (s4 Seq4) slide(newVal int) Seq4 {
	return Seq4{s4.b, s4.c, s4.d, newVal}
}

func part2() {
	secrets := parseSecrets("resources/Day22/input.txt")

	sequenceSums := map[Seq4]int{}
	var expected = Seq4{-2, 1, -1, 3}
	_ = expected

	for _, secret := range secrets {
		newSecret := secret
		lastPrice := newSecret % 10

		thisSeqs := map[Seq4]bool{}

		last4 := Seq4{}

		for i := 0; i < 2000; i++ {
			newSecret = evolve(newSecret)
			newPrice := newSecret % 10
			priceChange := newPrice - lastPrice

			last4 = last4.slide(priceChange)
			if i >= 3 {
				if !thisSeqs[last4] {
					if last4 == expected {
						fmt.Println("here we are")
					}
					thisSeqs[last4] = true
					sequenceSums[last4] += newPrice
				}
			}

			lastPrice = newPrice
		}
	}

	type kv struct {
		k Seq4
		v int
	}
	kvs := make([]kv, 0, len(sequenceSums))
	for k, v := range sequenceSums {
		kvs = append(kvs, kv{k, v})
	}
	slices.SortFunc(kvs, func(kv1 kv, kv2 kv) int { return kv2.v - kv1.v })

	fmt.Println(kvs[0])
}

func main() {
	part2()
}
