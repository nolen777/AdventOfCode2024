package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Config []int

func GenerateOne(blockWithoutMetadata []string) Config {
	config := make([]int, len(blockWithoutMetadata[0]), len(blockWithoutMetadata[0]))

	for _, line := range blockWithoutMetadata {
		for idx, c := range line {
			if c == '#' {
				config[idx]++
			}
		}
	}
	return config
}

// returns (locks, keys)
func GenerateConfigs(blocks [][]string) ([]Config, []Config) {
	locks := []Config{}
	keys := []Config{}

	for _, block := range blocks {
		if block[0][0] == '#' {
			// it's a lock
			locks = append(locks, GenerateOne(block[1:]))
		} else {
			keys = append(keys, GenerateOne(block[:len(block)-1]))
		}
	}

	return locks, keys
}

func parseSchematicBlocks(fileName string) [][]string {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	blocks := [][]string{}
	schematicBlock := make([]string, 0, 5)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			blocks = append(blocks, schematicBlock)
			schematicBlock = make([]string, 0, 5)
		} else {
			schematicBlock = append(schematicBlock, line)
		}
	}

	if len(schematicBlock) > 0 {
		blocks = append(blocks, schematicBlock)
	}

	return blocks
}

func GenerateUnderArray(configs []Config) [][]Config {
	configsUnder := make([][]Config, 5, 5)

	for _, config := range configs {
		for idx, ct := range config {
			for i := 0; i < ct; i++ {
				configsUnder[idx][i] = config
			}
		}
	}

	return configsUnder
}

func CountFittingPairs(locks []Config, keys []Config) int {
	//	locksUnder := GenerateUnderArray(locks)
	//	keysUnder := GenerateUnderArray(keys)

	count := 0
	for _, lock := range locks {
	keyLabel:
		for _, key := range keys {
			fmt.Print(lock, key)
			for i := 0; i < 5; i++ {
				if key[i]+lock[i] > 5 {
					fmt.Println("overlap in column ", i+1)
					continue keyLabel
				}
			}
			fmt.Println("fit")
			count++
		}
	}
	return count
}

func part1() {
	blocks := parseSchematicBlocks("resources/Day25/input.txt")
	locks, keys := GenerateConfigs(blocks)

	fmt.Println(CountFittingPairs(locks, keys))
	//fmt.Println(locks)
	//fmt.Println(keys)
}

func part2() {
	lines := parseSchematicBlocks("resources/Day25/sample.txt")
	_ = lines
}

func main() {
	part1()
}
