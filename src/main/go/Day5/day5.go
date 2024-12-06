package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type rule struct {
	first  int
	second int
}

func parseFile(filename string) ([]rule, [][]int) {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var rules = []rule{}
	var seqs = [][]int{}

	scanner := bufio.NewScanner(file)
	readingRules := true
	for scanner.Scan() {
		textLine := scanner.Text()

		if readingRules {
			if textLine == "" {
				readingRules = false
			} else {
				ruleComponents := strings.Split(textLine, "|")
				first, _ := strconv.Atoi(ruleComponents[0])
				second, _ := strconv.Atoi(ruleComponents[1])
				newRule := rule{first: first, second: second}
				rules = append(rules, newRule)
			}
		} else {
			seqComponents := strings.Split(textLine, ",")
			newSeq := []int{}

			for _, comp := range seqComponents {
				i, _ := strconv.Atoi(comp)
				newSeq = append(newSeq, i)
			}
			seqs = append(seqs, newSeq)
		}
	}

	return rules, seqs
}

func satisfiesRules(rules []rule, seq []int) bool {
	for _, rule := range rules {
		for idx, num := range seq {
			if num != rule.second {
				continue
			}
			for _, later := range seq[idx:] {
				if later == rule.first {
					return false
				}
			}
		}
	}

	return true
}

func middleNumber(seq []int) int {
	return seq[len(seq)/2]
}

func fix(rules []rule, seq []int) []int {
	for _, rule := range rules {
		for idx, num := range seq {
			if num != rule.second {
				continue
			}
			for laterIdx, later := range seq[idx:] {
				if later == rule.first {
					// Swap these two and fix again
					newSeq := make([]int, len(seq))
					copy(newSeq, seq)
					newSeq[idx] = later
					newSeq[laterIdx+idx] = num
					return fix(rules, newSeq)
				}
			}
		}
	}

	return seq
}

func main() {
	rules, seqs := parseFile("resources/Day5/input.txt")

	total := 0
	for _, seq := range seqs {
		if satisfiesRules(rules, seq) {
			fmt.Println("Sequence ", seq, " satisfies!")
		} else {
			newSeq := fix(rules, seq)
			fmt.Println("Fixed sequence is ", seq)
			mid := middleNumber(newSeq)
			fmt.Println("Middle number is ", mid)
			total += mid
		}
	}
	fmt.Println("Total is ", total)
}
