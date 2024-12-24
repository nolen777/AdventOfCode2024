package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const (
	AND = iota
	OR
	XOR
)

type GateType uint8
type Gate struct {
	input1   string
	input2   string
	gateType GateType
	output   string

	input1DepChains [][]string
	input2DepChains [][]string
}

func parseGates(fileName string) (map[string]bool, []Gate) {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	initialValues := map[string]bool{}

	// first the initial values
	pattern := `(?P<first>.*): (0|1)`
	r, _ := regexp.Compile(pattern)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		match := r.FindStringSubmatch(line)
		switch match[2] {
		case "0":
			initialValues[match[1]] = false
		case "1":
			initialValues[match[1]] = true
		default:
			log.Fatal("Invalid initial value", match[1])
		}
	}

	pattern2 := `(.*) (AND|OR|XOR) (.*) -> (.*)`
	r2, _ := regexp.Compile(pattern2)
	gates := []Gate{}
	for scanner.Scan() {
		line := scanner.Text()

		match := r2.FindStringSubmatch(line)
		i1 := match[1]
		var gt GateType
		switch match[2] {
		case "AND":
			gt = AND
		case "OR":
			gt = OR
		case "XOR":
			gt = XOR
		}
		i2 := match[3]
		o := match[4]

		gates = append(gates, Gate{input1: i1, input2: i2, gateType: gt, output: o})
	}

	return initialValues, gates
}

func iterateGates(values map[string]bool, gates []Gate) ([]Gate, bool) {
	remainingGates := []Gate{}

	changes := false
	for _, gate := range gates {
		i1, ok := values[gate.input1]
		if !ok {
			remainingGates = append(remainingGates, gate)
			continue
		}
		i2, ok := values[gate.input2]
		if !ok {
			remainingGates = append(remainingGates, gate)
			continue
		}
		var o bool
		switch gate.gateType {
		case AND:
			o = i1 && i2
		case OR:
			o = i1 || i2
		case XOR:
			o = i1 != i2
		}
		values[gate.output] = o
		changes = true
	}
	return remainingGates, changes
}

func part1() {
	values, gates := parseGates("resources/Day24/input.txt")
	fmt.Println(values)
	fmt.Println(gates)

	cont := true
	for cont {
		fmt.Println(len(gates), "remaining")
		gates, cont = iterateGates(values, gates)
	}

	//var xNum, yNum, zNum uint64
	//for k, v := range values {
	//	if k[0] != 'z' && k[0] != 'x' && k[0] != 'y' {
	//		continue
	//	}
	//	val, err := strconv.Atoi(k[1:])
	//	if err != nil {
	//		log.Fatal("Invalid identifier ", k)
	//	}
	//	if v {
	//		if k[0] == 'x' {
	//			xNum = xNum | (1 << val)
	//		}
	//		if k[0] == 'y' {
	//			yNum = yNum | (1 << val)
	//		}
	//		if k[0] == 'z' {
	//			zNum = zNum | (1 << val)
	//		}
	//	}
	//}

	fmt.Println("Result: ", numFor('z', values))
}

func numFor(prefix uint8, values map[string]bool) int {
	var result int
	for k, v := range values {
		if k[0] != prefix {
			continue
		}
		val, err := strconv.Atoi(k[1:])
		if err != nil {
			log.Fatal("Invalid identifier ", k)
		}
		if v {
			result = result | (1 << val)
		}
	}
	return result
}

func popcount(v int) int {
	var ct int = 0

	for v != 0 {
		ct += v & 1
		v >>= 1
	}
	return ct
}

func CalcOneDepChain(name string, gateMap map[string]Gate) [][]string {
	if name[0] == 'x' || name[0] == 'y' {
		return [][]string{{name}}
	} else {
		inputGate := gateMap[name]
		if inputGate.input1DepChains != nil && inputGate.input2DepChains != nil {
			newChains := slices.Clone(inputGate.input1DepChains)
			newChains = append(newChains, inputGate.input2DepChains...)

			for idx, chain := range newChains {
				newChains[idx] = append(chain, name)
			}
			return newChains
		}
	}
	return nil
}

func CalculateDependencies(gateMap map[string]Gate) map[string]Gate {
	changed := true
	for changed {
		changed = false
		for name, gate := range gateMap {
			if gate.input1DepChains == nil {
				gate.input1DepChains = CalcOneDepChain(gate.input1, gateMap)
				changed = changed || gate.input1DepChains != nil
			}
			if gate.input2DepChains == nil {
				gate.input2DepChains = CalcOneDepChain(gate.input2, gateMap)
				changed = changed || gate.input2DepChains != nil
			}
			gateMap[name] = gate
		}
	}
	return gateMap
}

func NameForPosition(prefix string, position int) string {
	str := strconv.Itoa(position)
	if position < 10 {
		str = "0" + str
	}
	return prefix + str
}

func Set(prefix string, num int, values map[string]bool, bits int) map[string]bool {
	position := 0
	for position < bits {
		str := NameForPosition(prefix, position)

		values[str] = (num & 1) != 0
		num = num >> 1
		position += 1
	}
	return values
}

func AppendDepIndices(name string, gates []Gate, indices map[int]bool) map[int]bool {
	if name[0] != 'x' && name[0] != 'y' {
		for idx, gate := range gates {
			if indices[idx] {
				continue
			}
			if gate.output != name {
				continue
			}
			indices[idx] = true
			indices = AppendDepIndices(gate.input1, gates, indices)
			indices = AppendDepIndices(gate.input2, gates, indices)
			break
		}
	}
	return indices
}

type IndexPair struct {
	left  int
	right int
}

func recursiveTrySwaps(previousDeps map[string]bool, allSwappedNames map[string]bool, origGates []Gate, startPosition int, bits int) (map[string]bool, bool) {
	if len(allSwappedNames) > 8 {
		return allSwappedNames, false
	}
	if startPosition > bits+1 {
		return allSwappedNames, true
	}
	for pos := 0; pos <= startPosition; pos++ {
		tryVals := func(gates []Gate, position int, x bool, y bool, expectZ0 bool, expectZ1 bool) bool {
			vals := map[string]bool{}
			vals = Set("x", 0, vals, bits)
			vals = Set("y", 0, vals, bits)

			vals[NameForPosition("x", position)] = x
			vals[NameForPosition("y", position)] = y
			cont := true
			for cont {
				gates, cont = iterateGates(vals, gates)
			}
			if vals[NameForPosition("z", position)] != expectZ0 {
				//		fmt.Printf("Failed z0 for x==%t, y==%t at %d\n", x, y, pos)
				return false
			}
			if vals[NameForPosition("z", position+1)] != expectZ1 {
				//		fmt.Printf("Failed z1 for x==%t, y==%t at %d\n", x, y, pos)
				return false
			}
			return true
		}

		var failed bool
		failed = !tryVals(slices.Clone(origGates), pos, false, false, false, false)
		failed = failed || !tryVals(slices.Clone(origGates), pos, false, true, true, false)
		if failed {
			if pos < startPosition {
				return nil, false
			}
			possibleSwapIndices := map[int]bool{}
			possibleSwapIndices = AppendDepIndices(NameForPosition("z", pos), origGates, possibleSwapIndices)

			for swapIdx1, _ := range possibleSwapIndices {
				if allSwappedNames[origGates[swapIdx1].output] {
					continue
				}
				swap1 := origGates[swapIdx1]
				for swapIdx2 := 0; swapIdx2 < len(origGates); swapIdx2++ {
					if possibleSwapIndices[swapIdx2] {
						continue
					}
					swap2 := origGates[swapIdx2]
					if allSwappedNames[swap2.output] {
						continue
					}

					if swap1.input1 == swap2.output || swap1.input2 == swap2.output {
						continue
					}
					if swap2.input1 == swap1.output || swap2.input2 == swap1.output {
						continue
					}

					swap1.output = origGates[swapIdx2].output
					swap2.output = origGates[swapIdx1].output
					swappedGates := slices.Clone(origGates)
					swappedGates[swapIdx2] = swap1
					swappedGates[swapIdx1] = swap2

					swaps, ok := recursiveTrySwaps(previousDeps, allSwappedNames, swappedGates, startPosition+1, bits)
					if ok {
						return swaps, ok
					}
				}
			}

			return nil, false
		}

		failed = !tryVals(slices.Clone(origGates), pos, true, false, true, false)
		failed = failed || !tryVals(slices.Clone(origGates), pos, true, true, false, true)
		if failed {
			if pos < startPosition {
				return nil, false
			}
			possibleSwapIndices := map[int]bool{}
			possibleSwapIndices = AppendDepIndices(NameForPosition("z", pos), origGates, possibleSwapIndices)

			for swapIdx1, _ := range possibleSwapIndices {
				if allSwappedNames[origGates[swapIdx1].output] {
					continue
				}
				swap1 := origGates[swapIdx1]
				for swapIdx2 := 0; swapIdx2 < len(origGates); swapIdx2++ {
					if possibleSwapIndices[swapIdx2] {
						continue
					}
					swap2 := origGates[swapIdx2]
					if allSwappedNames[swap2.output] {
						continue
					}

					if swap1.input1 == swap2.output || swap1.input2 == swap2.output {
						continue
					}
					if swap2.input1 == swap1.output || swap2.input2 == swap1.output {
						continue
					}

					swap1.output = origGates[swapIdx2].output
					swap2.output = origGates[swapIdx1].output
					swappedGates := slices.Clone(origGates)
					swappedGates[swapIdx2] = swap1
					swappedGates[swapIdx1] = swap2

					swaps, ok := recursiveTrySwaps(previousDeps, allSwappedNames, swappedGates, startPosition+1, bits)
					if ok {
						return swaps, ok
					}
				}
			}

			return nil, false
		}
	}

	return recursiveTrySwaps(previousDeps, allSwappedNames, origGates, startPosition+1, bits)
}

type GateTree struct {
	output string

	gateType GateType

	leftName string
	left     *GateTree

	rightName string
	right     *GateTree
}

func FillTrees(trees map[string]*GateTree) map[string]*GateTree {
	allNames := make([]string, 0, len(trees))
	for name, _ := range trees {
		allNames = append(allNames, name)
	}

	for _, name := range allNames {
		root := trees[name]
		if root.left == nil {
			root.left = trees[root.leftName]
		}
		if root.right == nil {
			root.right = trees[root.rightName]
		}
	}

	return trees
}

func part2() {
	origValues, origGates := parseGates("resources/Day24/input.txt")
	fmt.Println(origValues)

	//var z08i int
	//var z08 Gate
	//var qtwi int
	//var qtw Gate
	//for idx, gate := range origGates {
	//	if gate.output == "z08" {
	//		z08i = idx
	//		z08 = gate
	//	}
	//	if gate.output == "qtw" {
	//		qtwi = idx
	//		qtw = gate
	//	}
	//}
	//
	//z08.output = "qtw"
	//qtw.output = "z08"
	//origGates[z08i] = qtw
	//origGates[qtwi] = z08

	//trees := make(map[string]*GateTree, len(origGates))
	//for len(origGates) > 0 {
	//	gate := origGates[len(origGates)-1]
	//	if gate.input1[0] == 'z' || gate.input2[0] == 'z' {
	//		log.Fatal("whelp")
	//	}
	//	tree := GateTree{output: gate.output, gateType: gate.gateType, leftName: gate.input1, left: nil, rightName: gate.input2, right: nil}
	//	trees[gate.output] = &tree
	//	origGates = origGates[:len(origGates)-1]
	//}
	//
	//FillTrees(trees)
	//
	//z := trees["z08"]
	//var printTree func(tree *GateTree)
	//printTree = func(tree *GateTree) {
	//	fmt.Println(tree.output)
	//	if tree.leftName[0] == 'x' || tree.leftName[0] == 'y' {
	//		fmt.Println(tree.leftName)
	//	} else {
	//		printTree(tree.left)
	//	}
	//	if tree.rightName[0] == 'x' || tree.rightName[0] == 'y' {
	//		fmt.Println(tree.rightName)
	//	} else {
	//		printTree(tree.right)
	//	}
	//}
	//printTree(z)

	var bits int
	for name, _ := range origValues {
		if name[0] != 'x' {
			continue
		}
		num, _ := strconv.Atoi(name[1:])
		if num > bits {
			bits = num
		}
	}
	bits += 1

	allSwappedNames, _ := recursiveTrySwaps(map[string]bool{}, map[string]bool{}, origGates, 0, bits)
	swappedNameSlice := make([]string, 0, len(allSwappedNames))
	for k, _ := range allSwappedNames {
		swappedNameSlice = append(swappedNameSlice, k)
	}

	slices.Sort(swappedNameSlice)
	fmt.Println(strings.Join(swappedNameSlice, ","))
}

func main() {
	part2()
}
