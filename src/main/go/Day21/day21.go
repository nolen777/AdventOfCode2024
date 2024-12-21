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
	Panic = byte(255)
)

type Coords struct {
	row    int
	column int
}

type Pad struct {
	rowCount    int
	columnCount int
	keys        [][]byte
}

type NumberPad Pad
type DirectionPad Pad

func (np NumberPad) CoordsOf(b byte) Coords {
	switch b {
	case '7':
		return Coords{0, 0}
	case '8':
		return Coords{0, 1}
	case '9':
		return Coords{0, 2}
	case '4':
		return Coords{1, 0}
	case '5':
		return Coords{1, 1}
	case '6':
		return Coords{1, 2}
	case '1':
		return Coords{2, 0}
	case '2':
		return Coords{2, 1}
	case '3':
		return Coords{2, 2}
	case Panic:
		return Coords{3, 0}
	case '0':
		return Coords{3, 1}
	case A:
		return Coords{3, 2}
	}

	log.Fatal("Bad numberpad byte ", string(rune(b)))
	return Coords{-1, -1}
}

func (dp DirectionPad) CoordsOf(b byte) Coords {
	switch b {
	case Panic:
		return Coords{0, 0}
	case Up:
		return Coords{0, 1}
	case A:
		return Coords{0, 2}
	case Left:
		return Coords{1, 0}
	case Down:
		return Coords{1, 1}
	case Right:
		return Coords{1, 2}
	}

	log.Fatal("Bad directionpad byte ", string(rune(b)))
	return Coords{-1, -1}
}

func (p Pad) CoordsOf(b byte) Coords {
	for rIdx, row := range p.keys {
		for cIdx, entry := range row {
			if entry == b {
				return Coords{rIdx, cIdx}
			}
		}
	}

	log.Fatal(b, "not found in", p)
	return Coords{-1, -1}
}

func CreateNumberPad() NumberPad {
	pad := NumberPad{rowCount: 4, columnCount: 3, keys: [][]byte{}}

	row0 := []byte{'7', '8', '9'}
	row1 := []byte{'4', '5', '6'}
	row2 := []byte{'1', '2', '3'}
	row3 := []byte{Panic, '0', A}

	pad.keys = append(pad.keys, row0, row1, row2, row3)
	return pad
}

func CreateDirectionPad() DirectionPad {
	pad := DirectionPad{rowCount: 2, columnCount: 3, keys: [][]byte{}}

	row0 := []byte{Panic, Up, A}
	row1 := []byte{Left, Down, Right}

	pad.keys = append(pad.keys, row0, row1)

	return pad
}

func (p NumberPad) CoordDirections(from Coords, to Coords, ms string) []string {
	if from == to {
		return []string{ms}
	}

	rDiff := to.row - from.row
	cDiff := to.column - from.column

	directions := []string{}

	if rDiff > 0 && p.keys[from.row+1][from.column] != Panic {
		added := p.CoordDirections(Coords{from.row + 1, from.column}, to, ms+string(Down))
		directions = append(directions, added...)
	}
	if rDiff < 0 && p.keys[from.row-1][from.column] != Panic {
		added := p.CoordDirections(Coords{from.row - 1, from.column}, to, ms+string(Up))
		directions = append(directions, added...)
	}
	if cDiff > 0 && p.keys[from.row][from.column+1] != Panic {
		added := p.CoordDirections(Coords{from.row, from.column + 1}, to, ms+string(Right))
		directions = append(directions, added...)
	}
	if cDiff < 0 && p.keys[from.row][from.column-1] != Panic {
		added := p.CoordDirections(Coords{from.row, from.column - 1}, to, ms+string(Left))
		directions = append(directions, added...)
	}

	return directions
}

func (dp DirectionPad) CoordDirections(from Coords, to Coords, ms string) []string {
	if from == to {
		return []string{ms}
	}

	rDiff := to.row - from.row
	cDiff := to.column - from.column

	directions := []string{}

	if rDiff > 0 && dp.keys[from.row+1][from.column] != Panic {
		added := dp.CoordDirections(Coords{from.row + 1, from.column}, to, ms+string(Down))
		directions = append(directions, added...)
	}
	if rDiff < 0 && dp.keys[from.row-1][from.column] != Panic {
		added := dp.CoordDirections(Coords{from.row - 1, from.column}, to, ms+string(Up))
		directions = append(directions, added...)
	}
	if cDiff > 0 && dp.keys[from.row][from.column+1] != Panic {
		added := dp.CoordDirections(Coords{from.row, from.column + 1}, to, ms+string(Right))
		directions = append(directions, added...)
	}
	if cDiff < 0 && dp.keys[from.row][from.column-1] != Panic {
		added := dp.CoordDirections(Coords{from.row, from.column - 1}, to, ms+string(Left))
		directions = append(directions, added...)
	}

	return directions
}

//func (p Pad) Directions(from byte, to byte, ms string) []string {
//	fC := p.CoordsOf(from)
//	tC := p.CoordsOf(to)
//
//	return p.CoordDirections(fC, tC, ms)
//}

func (np NumberPad) FindSequences(ms string) []string {
	padPosition := np.CoordsOf(A)
	allSeqs := []string{""}

	for _, button := range ms {
		newCoords := np.CoordsOf(byte(button))
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

func (dp DirectionPad) FindSequences(ms string) []string {
	padPosition := dp.CoordsOf(A)
	allSeqs := []string{""}

	for _, button := range ms {
		newCoords := dp.CoordsOf(byte(button))
		nextSeqs := dp.CoordDirections(padPosition, newCoords, "")

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

func (dp DirectionPad) CostFor(seq string, prefix string, remainingDepth int) int {
	subNumPadSeqs := strings.Split(seq, string(A))
	totalCost := 0
	for _, ss := range subNumPadSeqs[:len(subNumPadSeqs)-1] {
		withA := ss + "A"

		var dirPad1Seqs []string
		if dp1s, ok := dirPadCache[withA]; ok {
			dirPad1Seqs = dp1s
			cacheHits++
		} else {
			dirPad1Seqs = dp.FindSequences(withA)
			dirPadCache[withA] = dirPad1Seqs
			cacheMisses++
		}

		minSubstringCost := 999999
		for _, dps := range dirPad1Seqs {
			cost := len(dps)
			if remainingDepth > 0 {
				cost = dp.CostFor(dps, prefix+dps, remainingDepth-1)
			}
			if cost < minSubstringCost {
				minSubstringCost = cost
			}
		}
		totalCost += minSubstringCost
	}
	return totalCost
}

func part1() {
	lines := parseCodes("resources/Day21/sample.txt")
	fmt.Println(lines)

	numPad := CreateNumberPad()
	dirPad := CreateDirectionPad()

	totalComplexity := 0
	for _, line := range lines {
		minCost := 9999999

		numPadSeqs := numPad.FindSequences(string(line))

		_ = dirPad
		for _, numPadSeq := range numPadSeqs {
			cost := dirPad.CostFor(numPadSeq, "", 1)
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

func part2() {
	lines := parseCodes("resources/Day21/sample.txt")
	_ = lines
}

func main() {
	part1()
}
