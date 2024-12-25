package main

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"os"
	"regexp"
	"strconv"
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

	inputDeps map[string]bool
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

		low := min(i1, i2)
		high := max(i1, i2)

		gates = append(gates, Gate{input1: low, input2: high, gateType: gt, output: o})
	}

	return initialValues, gates
}

func iterateGates(values map[string]bool, gates map[string]Gate) (map[string]Gate, bool) {
	remainingGates := map[string]Gate{}

	changes := false
	for _, gate := range gates {
		i1, ok := values[gate.input1]
		if !ok {
			remainingGates[gate.output] = gate
			continue
		}
		i2, ok := values[gate.input2]
		if !ok {
			remainingGates[gate.output] = gate
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

	gateMap := FillDeps(gates)

	cont := true
	for cont {
		fmt.Println(len(gates), "remaining")
		gateMap, cont = iterateGates(values, gateMap)
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

func TryValues(gates map[string]Gate, position int, x bool, y bool, expectZ0 bool, expectZ1 bool) bool {
	vals := map[string]bool{}
	vals = Set("x", 0, vals, 64)
	vals = Set("y", 0, vals, 64)

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

func oneValidate(gates map[string]Gate, pos int) bool {
	if !TryValues(maps.Clone(gates), pos, false, true, true, false) || !TryValues(maps.Clone(gates), pos, false, false, false, false) {
		return false
	}

	if !TryValues(maps.Clone(gates), pos, true, false, true, false) || !TryValues(maps.Clone(gates), pos, true, true, false, true) {
		return false
	}
	return true
}

func validate(gates map[string]Gate, bits int) (int, bool) {
	for pos := 0; pos <= bits; pos++ {
		if !oneValidate(gates, pos) {
			return pos, false
		}
		//	fmt.Println("Through position ", pos)
	}

	return 0, true
}

func recursiveCheck(gates map[string]Gate, bits int, prev map[string]bool, recurseAfterPos int) (map[string]bool, bool) {
	if len(prev) > 8 {
		return map[string]bool{}, false
	}

	pos, ok := validate(gates, bits)
	if ok {
		return prev, true
	}

	if pos < recurseAfterPos {
		return prev, false
	}

	zN := NameForPosition("z", pos)
	zNP := NameForPosition("z", pos+1)
	if prev[zN] {
		return prev, false
	}
	deps := map[string]bool{}
	for zP, _ := range gates[zN].inputDeps {
		_, ok := gates[zP]
		if ok {
			deps[zP] = true
		}
	}
	for zPD, _ := range gates[zNP].inputDeps {
		_, ok := gates[zPD]
		if ok {
			deps[zPD] = true
		}
	}
	deps[zN] = true
	deps[zNP] = true
	fmt.Println("Checking deps ", deps)
	for depName, _ := range deps {
		if depName[0] == 'x' {
			continue
		}
		if depName[0] == 'y' {
			continue
		}
		if depName != zN && depName[0] == 'z' {
			continue
		}
	swapOuter:
		for swapName, swapper := range gates {
			if swapName == depName {
				continue
			}
			if prev[swapper.output] {
				continue
			}
			for h := pos + 1; h <= bits; h++ {
				if swapper.inputDeps[NameForPosition("x", h)] {
					continue swapOuter
				}
				if swapper.inputDeps[NameForPosition("y", h)] {
					continue swapOuter
				}
			}
			dep := gates[depName]

			newSwapper := swapper
			newDep := dep
			newDep.output = swapName
			newSwapper.output = depName
			gates[swapName] = newDep
			gates[depName] = newSwapper

			if oneValidate(gates, pos) {
				newPrev := maps.Clone(prev)
				newPrev[newSwapper.output] = true
				newPrev[newDep.output] = true

				fmt.Println("Possible success with ", newDep.output, newSwapper.output)
				swaps, ok := recursiveCheck(gates, bits, newPrev, pos+1)
				if ok {
					return swaps, ok
				}
			}

			gates[swapName] = swapper
			gates[depName] = dep
		}
	}

	fmt.Println(zN, zNP)
	return prev, false
}

func FillDeps(gates []Gate) map[string]Gate {
	gateMap := make(map[string]Gate, len(gates))

	for _, gate := range gates {
		gate.inputDeps = map[string]bool{}
		gate.inputDeps[gate.input1] = true
		gate.inputDeps[gate.input2] = true
		gateMap[gate.output] = gate
	}

	changed := true
	for changed {
		changed = false
		for k, v := range gateMap {
			for dep, _ := range v.inputDeps {
				depGate := gateMap[dep]
				for depDep, _ := range depGate.inputDeps {
					if !v.inputDeps[depDep] {
						changed = true
						v.inputDeps[depDep] = true
					}
				}
				gateMap[k] = v
			}
		}
	}
	return gateMap
}

func part2() {
	origValues, origGates := parseGates("resources/Day24/input.txt")
	fmt.Println(origValues)

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

	gateMap := FillDeps(origGates)
	//fmt.Println(gateMap)

	fmt.Println(recursiveCheck(gateMap, bits, map[string]bool{}, 0))

}

func main() {
	part2()
}
