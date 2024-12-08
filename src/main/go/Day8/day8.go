package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type FrequencyMap struct {
	rowCount    int
	columnCount int
	locations   map[rune][][2]int
}

func parseFrequencyLocations(fileName string) FrequencyMap {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	locationMap := make(map[rune][][2]int, 0)

	row := 0
	columnCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		textLine := scanner.Text()
		columnCount = len(textLine)

		for column, b := range textLine {
			if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') {
				locationMap[b] = append(locationMap[b], [2]int{row, column})
			}
		}

		row++
	}

	return FrequencyMap{rowCount: row, columnCount: columnCount, locations: locationMap}
}

func isValid(frequencyMap FrequencyMap, position [2]int) bool {
	return position[0] >= 0 && position[0] < frequencyMap.rowCount && position[1] >= 0 && position[1] < frequencyMap.columnCount
}

func findAntinodes(frequencyMap FrequencyMap) map[[2]int]bool {
	antinodes := map[[2]int]bool{}
	for _, v := range frequencyMap.locations {
		for i, loc1 := range v {
			for _, loc2 := range v[i+1:] {
				rowDiff := loc2[0] - loc1[0]
				colDiff := loc2[1] - loc1[1]

				pos1 := [2]int{loc2[0] + rowDiff, loc2[1] + colDiff}
				pos2 := [2]int{loc1[0] - rowDiff, loc1[1] - colDiff}

				if isValid(frequencyMap, pos1) {
					antinodes[pos1] = true
				}
				if isValid(frequencyMap, pos2) {
					antinodes[pos2] = true
				}
			}
		}
	}
	return antinodes
}

func main() {
	locationMap := parseFrequencyLocations("resources/Day8/input.txt")

	antinodes := findAntinodes(locationMap)
	fmt.Println(locationMap)
	fmt.Println(antinodes)
	fmt.Println("count: ", len(antinodes))
}
