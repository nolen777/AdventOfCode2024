package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	adv = 0
	bxl = 1
	bst = 2
	jnz = 3
	bxc = 4
	out = 5
	bdv = 6
	cdv = 7
)

type State struct {
	A                  int
	B                  int
	C                  int
	Program            []int
	InstructionPointer int
	Output             []string
}

func parseInitialState(fileName string) State {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	state := State{}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "Register A: ") {
			a, err := strconv.Atoi(line[len("Register A: "):])
			if err != nil {
				log.Fatal(err)
			}
			state.A = a
		}
		if strings.HasPrefix(line, "Register B: ") {
			b, err := strconv.Atoi(line[len("Register B: "):])
			if err != nil {
				log.Fatal(err)
			}
			state.B = b
		}
		if strings.HasPrefix(line, "Register C: ") {
			c, err := strconv.Atoi(line[len("Register C: "):])
			if err != nil {
				log.Fatal(err)
			}
			state.C = c
		}
		if strings.HasPrefix(line, "Program: ") {
			codes := strings.Split(line[len("Program: "):], ",")
			state.Program = make([]int, 0, len(codes))
			for _, code := range codes {
				c, err := strconv.Atoi(code)
				if err != nil {
					log.Fatal(err)
				}
				state.Program = append(state.Program, c)
			}
		}
	}

	state.InstructionPointer = 0
	state.Output = make([]string, 0)
	return state
}

func combo(value int, s State) int {
	switch value {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 2
	case 3:
		return 3
	case 4:
		return s.A
	case 5:
		return s.B
	case 6:
		return s.C
	case 7:
		log.Fatal("Invalid combo operand")
	}

	log.Fatal("Invalid combo operand")
	return 0
}

func advanceState(s State) State {
	inst := s.Program[s.InstructionPointer]
	operand := s.Program[s.InstructionPointer+1]

	s.InstructionPointer += 2

	switch inst {
	case adv:
		numerator := s.A
		denominator := 1 << combo(operand, s)

		s.A = numerator / denominator

	case bxl:
		s.B = s.B ^ operand

	case bst:
		s.B = combo(operand, s) % 8

	case jnz:
		if s.A != 0 {
			s.InstructionPointer = operand
		}

	case bxc:
		s.B = s.B ^ s.C

	case out:
		s.Output = append(s.Output, strconv.Itoa(combo(operand, s)%8))

	case bdv:
		numerator := s.A
		denominator := 2 << combo(operand, s)

		s.B = numerator / denominator

	case cdv:
		numerator := s.A
		denominator := 2 << combo(operand, s)

		s.C = numerator / denominator
	}

	return s
}

func part1() {
	state := parseInitialState("resources/Day17/sample.txt")
	fmt.Println(state)

	for state.InstructionPointer < len(state.Program) {
		state = advanceState(state)
		fmt.Println(state)
	}

	fmt.Println(strings.Join(state.Output, ","))
}

func part2() {

}

func main() {
	part1()
}
