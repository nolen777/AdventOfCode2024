package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func GenerateOne(blockWithoutMetadata []string) []int {
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
func GenerateConfigs(blocks [][]string) ([][]int, [][]int) {
	locks := [][]int{}
	keys := [][]int{}

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

func part1() {
	blocks := parseSchematicBlocks("resources/Day25/sample.txt")
	locks, keys := GenerateConfigs(blocks)
	fmt.Println(locks)
	fmt.Println(keys)
}

func part2() {
	lines := parseSchematicBlocks("resources/Day25/sample.txt")
	_ = lines
}

func main() {
	part1()
}
