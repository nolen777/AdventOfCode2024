package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func parseBytes(fileName string) [][]byte {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	grid := [][]byte{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		textLine := scanner.Text()
		byteLine := make([]byte, 0, len(textLine))

		for _, elt := range textLine {
			byteLine = append(byteLine, byte(elt))
		}

		grid = append(grid, byteLine)
	}
	return grid
}

type Coords struct {
	r int
	c int
}

type Plot struct {
	coords   Coords
	value    byte
	regionId int
}

func makePlots(grid [][]byte) [][]Plot {
	plots := make([][]Plot, 0, len(grid))
	for i, row := range grid {
		line := make([]Plot, 0, len(row))
		for j, b := range row {
			line = append(line, Plot{coords: Coords{r: i, c: j}, value: b, regionId: -1})
		}
		plots = append(plots, line)
	}
	return plots
}

func neighbors(plots [][]Plot, coords Coords) []Coords {
	neighbors := []Coords{}
	i := coords.r
	j := coords.c
	if i > 0 {
		neighbors = append(neighbors, Coords{r: i - 1, c: j})
	}
	if j > 0 {
		neighbors = append(neighbors, Coords{r: i, c: j - 1})
	}
	if i < len(plots)-1 {
		neighbors = append(neighbors, Coords{r: i + 1, c: j})
	}
	if j < len(plots[i])-1 {
		neighbors = append(neighbors, Coords{r: i, c: j + 1})
	}

	return neighbors
}

func assignRegion(plots [][]Plot, coords Coords, regionId int) {
	i := coords.r
	j := coords.c

	if plots[i][j].regionId == regionId {
		return
	}
	value := plots[i][j].value
	plots[i][j].regionId = regionId

	for _, neighbor := range neighbors(plots, coords) {
		if plots[neighbor.r][neighbor.c].value == value {
			assignRegion(plots, neighbor, regionId)
		}
	}
}

func findRegions(plots [][]Plot) {
	lastNewRegion := 0

	for i := 0; i < len(plots); i++ {
		for j := 0; j < len(plots[i]); j++ {
			if plots[i][j].regionId == -1 {
				lastNewRegion++
				assignRegion(plots, Coords{r: i, c: j}, lastNewRegion)
			}
		}
	}
}

func perimeterValue(plots [][]Plot, coords Coords) int {
	neighbors := neighbors(plots, coords)
	regionId := plots[coords.r][coords.c].regionId

	perimeterValue := 4 - len(neighbors)
	for _, n := range neighbors {
		if plots[n.r][n.c].regionId != regionId {
			perimeterValue++
		}
	}
	return perimeterValue
}

func calculateValues(plots [][]Plot) []int {
	values := []int{}

	regionId := 1

	for {
		regionArea := 0
		regionPerimeter := 0
		for i, row := range plots {
			for j, plot := range row {
				if plot.regionId == regionId {
					regionArea++
					regionPerimeter += perimeterValue(plots, Coords{r: i, c: j})
				}
			}
		}
		if regionArea == 0 {
			break
		}
		values = append(values, regionArea*regionPerimeter)
		regionId++
	}

	return values
}

func main() {
	grid := parseBytes("resources/Day12/sampleinput.txt")
	plots := makePlots(grid)
	findRegions(plots)
	values := calculateValues(plots)

	fmt.Println(values)
}