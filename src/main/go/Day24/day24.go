package main

import (
	"bufio"
	"fmt"
	"log"
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

func iterateGates(values map[string]bool, gates []Gate) []Gate {
	remainingGates := []Gate{}

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
	}
	return remainingGates
}

func part1() {
	values, gates := parseGates("resources/Day24/sample.txt")
	fmt.Println(values)
	fmt.Println(gates)

	for len(gates) > 0 {
		fmt.Println(len(gates), "remaining")
		gates = iterateGates(values, gates)
	}

	var result uint64
	for k, v := range values {
		if k[0] != 'z' {
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

	fmt.Println("Result: ", result)
}

func part2() {
	initialValues, gates := parseGates("resources/Day24/small.txt")
	fmt.Println(initialValues)
	_ = gates
}

func main() {
	part1()
}
