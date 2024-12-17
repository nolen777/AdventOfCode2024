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
	A                  int64
	B                  int64
	C                  int64
	Program            []byte
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
			a, err := strconv.ParseInt(line[len("Register A: "):], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			state.A = a
		}
		if strings.HasPrefix(line, "Register B: ") {
			b, err := strconv.ParseInt(line[len("Register B: "):], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			state.B = b
		}
		if strings.HasPrefix(line, "Register C: ") {
			c, err := strconv.ParseInt(line[len("Register C: "):], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			state.C = c
		}
		if strings.HasPrefix(line, "Program: ") {
			codes := strings.Split(line[len("Program: "):], ",")
			state.Program = make([]byte, 0, len(codes))
			for _, code := range codes {
				c, err := strconv.ParseInt(code, 10, 8)
				if err != nil {
					log.Fatal(err)
				}
				state.Program = append(state.Program, byte(c))
			}
		}
	}

	state.InstructionPointer = 0
	state.Output = make([]string, 0)
	return state
}

func makeCombo(value byte, s State) int64 {
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

	literal := int64(operand)
	combo := makeCombo(operand, s)

	s.InstructionPointer += 2

	switch inst {
	case adv:
		s.A = s.A >> combo

	case bxl:
		s.B = s.B ^ literal

	case bst:
		s.B = combo % 8

	case jnz:
		if s.A != 0 {
			s.InstructionPointer = int(literal)
		}

	case bxc:
		s.B = s.B ^ s.C

	case out:
		outValue := int(combo % 8)
		s.Output = append(s.Output, strconv.Itoa(outValue))

	case bdv:
		s.B = s.A >> combo

	case cdv:
		s.C = s.A >> combo
	}

	return s
}

func part1() {
	state := parseInitialState("resources/Day17/input.txt")
	fmt.Println(state)

	for state.InstructionPointer < len(state.Program) {
		state = advanceState(state)
		fmt.Println(state)
	}

	fmt.Println(strings.Join(state.Output, ","))
}

func part2() {
	initialState := parseInitialState("resources/Day17/input.txt")

	progStrings := make([]string, 0, len(initialState.Program))
	for _, p := range initialState.Program {
		pI := strconv.Itoa(int(p))
		progStrings = append(progStrings, pI)
	}
	progString := strings.Join(progStrings, ",")
	fmt.Println("Desired:")
	fmt.Println(progString)

	// The strategy here is basically:
	//		For this program, A only ever changes by shifting right by 3, and B and C effectively
	//		don't carry over from one iteration to the next; they're set entirely by the value
	//		of A. And because of the mod 8 in the output, only the last three bits of A matter.
	//		(not sure that's totally true because of the bitshift by 5 on C, but it worked...)
	//		So we can look for the output backwards, one entry at a time, and shift A up by 3
	//		each iteration.
	outputString := ""
	var a int64 = 0
	for outputString != progString {
		// Each pass through shifts A right by 3
		a = a << 3
		for i := int64(0); ; i++ {
			state := initialState
			state.A = a | i

			for state.InstructionPointer < len(state.Program) {
				state = advanceState(state)
			}
			outputString = strings.Join(state.Output, ",")

			if strings.HasSuffix(progString, outputString) {
				fmt.Println("Found it with ", i)
				a = a | i
				fmt.Println(outputString)
				break
			}
		}
	}

	fmt.Println("initial A should be ", a)
}

func main() {
	part1()
	part2()
}
