package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func parseCodes(fileName string) []string {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines
}

const (
	Up    = '^'
	A     = 'A'
	Left  = '<'
	Down  = 'v'
	Right = '>'
)

type Coords struct {
	row    int
	column int
}

type Pad struct {
	rowCount    int
	columnCount int
	keyCoords   map[byte]Coords
	panicCoords Coords
}

var numberPad = Pad{rowCount: 4, columnCount: 3, keyCoords: map[byte]Coords{
	'7': {0, 0},
	'8': {0, 1},
	'9': {0, 2},
	'4': {1, 0},
	'5': {1, 1},
	'6': {1, 2},
	'1': {2, 0},
	'2': {2, 1},
	'3': {2, 2},
	'0': {3, 1},
	A:   {3, 2},
}, panicCoords: Coords{3, 0}}

var directionPad = Pad{rowCount: 2, columnCount: 3, keyCoords: map[byte]Coords{
	Up:    {0, 1},
	A:     {0, 2},
	Left:  {1, 0},
	Down:  {1, 1},
	Right: {1, 2}}, panicCoords: Coords{0, 0}}

func (p Pad) CoordDirections(from Coords, to Coords, ms string) []string {
	if from == to {
		return []string{ms}
	}
	if from == p.panicCoords {
		log.Fatal("Can't start at Panic")
	}
	if to == p.panicCoords {
		log.Fatal("Can't go to Panic")
	}

	rDiff := to.row - from.row
	cDiff := to.column - from.column

	directions := []string{}

	if rDiff > 0 {
		next := Coords{from.row + 1, from.column}
		if next != p.panicCoords {
			added := p.CoordDirections(next, to, ms+string(Down))
			directions = append(directions, added...)
		}
	}
	if rDiff < 0 {
		next := Coords{from.row - 1, from.column}
		if next != p.panicCoords {
			added := p.CoordDirections(next, to, ms+string(Up))
			directions = append(directions, added...)
		}
	}
	if cDiff > 0 {
		next := Coords{from.row, from.column + 1}
		if next != p.panicCoords {
			added := p.CoordDirections(next, to, ms+string(Right))
			directions = append(directions, added...)
		}
	}
	if cDiff < 0 {
		next := Coords{from.row, from.column - 1}
		if next != p.panicCoords {
			added := p.CoordDirections(next, to, ms+string(Left))
			directions = append(directions, added...)
		}
	}

	return directions
}

func (np Pad) FindSequences(ms string) []string {
	padPosition := np.keyCoords[A]
	allSeqs := []string{""}

	for _, button := range ms {
		newCoords := np.keyCoords[byte(button)]
		nextSeqs := np.CoordDirections(padPosition, newCoords, "")

		extended := make([]string, 0, len(allSeqs)*len(nextSeqs))
		for _, e := range allSeqs {
			for _, ns := range nextSeqs {
				extended = append(extended, e+ns+string(A))
			}
		}
		allSeqs = extended

		padPosition = newCoords
	}

	return allSeqs
}

var dirPadCache = map[string][]string{}
var cacheHits int = 0
var cacheMisses int = 0

type CostCacheKey struct {
	remainingDepth int
	seq            string
}

var costCache = map[CostCacheKey]int{}

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

func (dp Pad) CostFor(seq string, prefix string, remainingDepth int) int {
	superKey := CostCacheKey{remainingDepth, seq}
	if cachedCost, ok := costCache[superKey]; ok {
		return cachedCost
	}

	subNumPadSeqs := strings.Split(seq, string(A))
	totalCost := 0
	for _, ss := range subNumPadSeqs[:len(subNumPadSeqs)-1] {
		withA := ss + "A"
		subCacheKey := CostCacheKey{remainingDepth, ss}

		if subCost, ok := costCache[subCacheKey]; ok {
			totalCost += subCost
			continue
		}

		var dirPad1Seqs []string
		if dp1s, ok := dirPadCache[withA]; ok {
			dirPad1Seqs = dp1s
			cacheHits++
		} else {
			dirPad1Seqs = dp.FindSequences(withA)
			dirPadCache[withA] = dirPad1Seqs
			cacheMisses++
		}

		minSubstringCost := MaxInt
		for _, dps := range dirPad1Seqs {
			cost := len(dps)
			if remainingDepth > 0 {
				cost = dp.CostFor(dps, prefix+dps, remainingDepth-1)
			}
			if cost < minSubstringCost {
				minSubstringCost = cost
			}
		}
		costCache[subCacheKey] = minSubstringCost
		totalCost += minSubstringCost
	}
	costCache[superKey] = totalCost
	return totalCost
}

func calculate(depth int) {
	lines := parseCodes("resources/Day21/input.txt")
	fmt.Println(lines)

	totalComplexity := 0
	for _, line := range lines {
		minCost := MaxInt

		numPadSeqs := numberPad.FindSequences(line)

		for _, numPadSeq := range numPadSeqs {
			cost := directionPad.CostFor(numPadSeq, "", depth)
			if cost < minCost {
				minCost = cost
			}
		}
		fmt.Println(line, " costs ", minCost)

		numPart, err := strconv.Atoi(line[:len(line)-1])
		if err != nil {
			log.Fatal(err)
		}
		totalComplexity += minCost * numPart
	}
	fmt.Println("Total complexity: ", totalComplexity)
}

func part1() {
	calculate(1)
}

func part2() {
	calculate(24)
}

func main() {
	part2()
}
