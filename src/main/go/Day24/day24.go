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

func InsertHalfAdder(left string, right string, sumName string, carryName string, gates map[string]Gate) map[string]Gate {
	gates[sumName] = Gate{
		input1:   left,
		input2:   right,
		gateType: XOR,
		output:   sumName,
	}

	gates[carryName] = Gate{
		input1:   left,
		input2:   right,
		gateType: AND,
		output:   carryName,
	}

	return gates
}

func IsNormalizedMatch(expected Gate, got Gate, nameToNorm map[string]string) bool {
	if expected.gateType != got.gateType {
		return false
	}

	n1, ok := nameToNorm[got.input1]
	if !ok {
		return false
	}
	n2, ok := nameToNorm[got.input2]
	if !ok {
		return false
	}
	if n1 != expected.input1 && n1 != expected.input2 {
		return false
	}
	if n2 != expected.input2 && n2 != expected.input1 {
		return false
	}
	return true
}

func FindNormalizedMatch(expected Gate, origGateMap map[string]Gate, nameToNorm map[string]string) (Gate, bool) {
	for _, origGate := range origGateMap {
		if IsNormalizedMatch(expected, origGate, nameToNorm) {
			return origGate, true
		}
	}
	return Gate{}, false
}

func expectedGates(bits int, origGateMap map[string]Gate) map[string]Gate {
	nameToNorm := make(map[string]string, len(origGateMap))
	for b := 0; b < bits; b++ {
		nameToNorm[NameForPosition("x", b)] = NameForPosition("x", b)
		nameToNorm[NameForPosition("y", b)] = NameForPosition("y", b)
	}

	gates := make(map[string]Gate, bits*5)

	gates = InsertHalfAdder("x00", "y00", "z00", "_car00", gates)

	for b := 1; b < bits+1; b++ {
		gates = InsertHalfAdder(
			NameForPosition("x", b),
			NameForPosition("y", b),
			NameForPosition("_hSum", b),
			NameForPosition("_hCar", b), gates)

		gates = InsertHalfAdder(
			NameForPosition("_hSum", b),
			NameForPosition("_car", b-1),
			NameForPosition("z", b),
			NameForPosition("_h2Car", b), gates)

		gates[NameForPosition("_car", b)] = Gate{
			input1:   NameForPosition("_hCar", b),
			input2:   NameForPosition("_h2Car", b),
			gateType: OR,
			output:   NameForPosition("_car", b),
		}
	}

	return gates
}

func part2() {
	origValues, origGates := parseGates("resources/Day24/modifiedFile.txt")
	fmt.Println(origValues)

	var bits int
	gateMap := make(map[string]Gate, len(origGates))
	for _, gate := range origGates {
		gateName := gate.output
		gateMap[gateName] = gate

		if gateName[0] == 'z' {
			num, _ := strconv.Atoi(gateName[1:])
			if num > bits {
				bits = num
			}
		}
	}

	expected := expectedGates(bits, gateMap)
	fmt.Println(expected)

	expectedValues := map[string]bool{}
	x := 393842731
	y := 13894701
	Set("x", x, expectedValues, bits)
	Set("y", y, expectedValues, bits)

	cont := true
	remainingGates := expected
	for cont {
		remainingGates, cont = iterateGates(expectedValues, remainingGates)
	}
	result := numFor('z', expectedValues)

	fmt.Printf("Got %d, expected %d, diff %d\n", result, x+y, result-x-y)

	nameToNorm := make(map[string]string, len(origValues))
	normToName := make(map[string]string, len(origValues))
	for b := 0; b < bits; b++ {
		nameToNorm[NameForPosition("x", b)] = NameForPosition("x", b)
		nameToNorm[NameForPosition("y", b)] = NameForPosition("y", b)
		//	nameToNorm[NameForPosition("z", b)] = NameForPosition("z", b)

		normToName[NameForPosition("x", b)] = NameForPosition("x", b)
		normToName[NameForPosition("y", b)] = NameForPosition("y", b)
		//	normToName[NameForPosition("z", b)] = NameForPosition("z", b)
	}

	changes := true
	for changes {
		changes = false
		for givenName, givenGate := range gateMap {
			_, exists := nameToNorm[givenName]
			if exists {
				continue
			}

			normIn1, ok := nameToNorm[givenGate.input1]
			if !ok {
				//	fmt.Printf("Don't have %s yet\n", givenGate.input1)
				continue
			}
			normIn2, ok := nameToNorm[givenGate.input2]
			if !ok {
				//	fmt.Printf("Don't have %s yet\n", givenGate.input2)
				continue
			}

			for normName, normGate := range expected {
				if normGate.gateType != givenGate.gateType {
					continue
				}
				if (normIn1 == normGate.input1 && normIn2 == normGate.input2) ||
					(normIn2 == normGate.input1 && normIn1 == normGate.input2) {
					changes = true
					nameToNorm[givenName] = normName
					normToName[normName] = givenName
					break
				}
			}
		}
	}

	identities := []string{"_hSum", "_hCar", "_h2Car", "_car", "z"}

	changes = true
	for changes {
		changes = false

		for b := 1; b < bits; b++ {
			for _, id := range identities {
				idName := NameForPosition(id, b)

				expectedGate := expected[idName]
				gotGate, ok := gateMap[normToName[idName]]

				isNormalizedMatch := func(expected Gate, got Gate) bool {
					if expected.gateType != got.gateType {
						return false
					}

					n1, ok := nameToNorm[got.input1]
					if !ok {
						return false
					}
					n2, ok := nameToNorm[got.input2]
					if !ok {
						return false
					}
					if n1 != expected.input1 && n1 != expected.input2 {
						return false
					}
					if n2 != expected.input2 && n2 != expected.input1 {
						return false
					}
					return true
				}

				if !ok || !isNormalizedMatch(expectedGate, gotGate) {
					fmt.Printf("Mismatch at %s!\n", idName)

					var shouldBe Gate
					found := false
					for realGateName, realGate := range gateMap {
						if realGate.gateType != expectedGate.gateType {
							continue
						}
						normIn1 := nameToNorm[realGate.input1]
						normIn2 := nameToNorm[realGate.input2]

						if normIn1 != expectedGate.input1 && normIn1 != expectedGate.input2 {
							continue
						}
						if normIn2 != expectedGate.input1 && normIn2 != expectedGate.input2 {
							continue
						}
						fmt.Printf("Found a match for %s\n", realGateName)
						found = true
					}
					if !found {
						fmt.Println("Not good")
					}

					if !ok {
						fmt.Printf("Couldn't find a match for %s\n", idName)
						break
					}
					fmt.Printf("%s maybe?\n", shouldBe)

					fmt.Printf("Swapping %s and %s\n", shouldBe.output, gotGate.output)

					shouldBe.output = gotGate.output
					gateMap[shouldBe.output] = shouldBe
					gotGate.output = normToName[idName]
					gateMap[gotGate.output] = gotGate

					normToName[expectedGate.output] = shouldBe.output
					nameToNorm[shouldBe.output] = expectedGate.output

					changes = true
					break
				}
			}
		}
	}

	//gateMap := FillDeps(origGates)
	////fmt.Println(gateMap)
	//
	//fmt.Println(recursiveCheck(gateMap, bits, map[string]bool{}, 0))

}

func main() {
	part2()
}
