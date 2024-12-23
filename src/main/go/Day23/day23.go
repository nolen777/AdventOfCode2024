package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

type Edge struct {
	a string
	b string
}

// Always make in alpha order
func MakeEdge(a string, b string) Edge {
	if a < b {
		return Edge{a, b}
	}
	return Edge{b, a}
}

type Node struct {
	name  string
	edges map[string]bool
}

func InsertEdge(m map[string]Node, a string, b string) {
	firstNode, ok := m[a]
	if !ok {
		firstNode = Node{name: a, edges: map[string]bool{}}
	}
	firstNode.edges[b] = true
	m[a] = firstNode
}

func ParseNodes(fileName string) map[string]Node {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	nodes := map[string]Node{}

	for scanner.Scan() {
		line := scanner.Text()
		newEdgeNodes := strings.Split(line, "-")
		first := newEdgeNodes[0]
		second := newEdgeNodes[1]

		InsertEdge(nodes, first, second)
		InsertEdge(nodes, second, first)
	}

	return nodes
}

func part1() {
	nodes := ParseNodes("resources/Day23/sample.txt")

	type Triple struct {
		a string
		b string
		c string
	}

	makeTriple := func(a string, b string, c string) Triple {
		elts := []string{a, b, c}
		slices.Sort(elts)
		return Triple{elts[0], elts[1], elts[2]}
	}

	triples := map[Triple]bool{}

	for _, node := range nodes {
		a := node.name

		for alt1 := range node.edges {
			for alt2 := range node.edges {
				if alt1 == alt2 {
					continue
				}
				if nodes[alt1].edges[alt2] {
					if a[0] == 't' || alt1[0] == 't' || alt2[0] == 't' {
						newTriple := makeTriple(a, alt1, alt2)
						triples[newTriple] = true
					}
				}
			}
		}
	}

	for triple := range triples {
		fmt.Printf("%s,%s,%s\n", triple.a, triple.b, triple.c)
	}
	fmt.Printf("%d found!\n", len(triples))

	fmt.Println(nodes)
}

func part2() {
	nodes := ParseNodes("resources/Day23/sample.txt")
	_ = nodes
}

func main() {
	part1()
}
