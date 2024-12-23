package main

import (
	"bufio"
	"fmt"
	"log"
	"maps"
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

type Triple struct {
	a string
	b string
	c string
}

func FindTriples(nodes map[string]Node) map[Triple]bool {
	triples := map[Triple]bool{}

	makeTriple := func(a string, b string, c string) Triple {
		elts := []string{a, b, c}
		slices.Sort(elts)
		return Triple{elts[0], elts[1], elts[2]}
	}

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
	return triples
}

func part1() {
	nodes := ParseNodes("resources/Day23/input.txt")

	triples := FindTriples(nodes)

	for triple := range triples {
		fmt.Printf("%s,%s,%s\n", triple.a, triple.b, triple.c)
	}
	fmt.Printf("%d found!\n", len(triples))
}

func FilterNodes(nodes map[string]Node, condition func(Node) bool) {
	for k, v := range nodes {
		if condition(v) {
			continue
		}
		delete(nodes, k)
	}
}

func part2() {
	nodes := ParseNodes("resources/Day23/input.txt")
	fmt.Println("Start size: ", len(nodes))
	FilterNodes(nodes, func(n Node) bool {
		if n.name[0] == 't' {
			return true
		}
		for k := range n.edges {
			if k[0] == 't' {
				return true
			}
		}
		return false
	})
	fmt.Println("Filtered size: ", len(nodes))

	triples := FindTriples(nodes)

	type Party map[string]bool
	makeParty := func(t Triple) Party {
		return Party{t.a: true, t.b: true, t.c: true}
	}
	appendedParty := func(ps Party, s string) Party {
		pc := maps.Clone(ps)
		pc[s] = true
		return pc
	}

	parties := map[string]Party{}

	getPassword := func(p Party) string {
		elts := make([]string, 0, len(p))
		for k := range p {
			elts = append(elts, k)
		}
		slices.Sort(elts)
		return strings.Join(elts, ",")
	}

	for triple := range triples {
		newParty := makeParty(triple)
		parties[getPassword(newParty)] = newParty
	}

	changes := true
	minLength := 3
	for changes {
		fmt.Println("Party count: ", len(parties))
		if len(parties) == 1 {
			break
		}
		changes = false
		newParties := map[string]Party{}

		// wtf
		for _, party := range parties {
			if len(party) < minLength {
				continue
			}
			candidates := map[string]bool{}

			first := true
			for k := range party {
				if first {
					candidates = maps.Clone(nodes[k].edges)
					for c := range candidates {
						if party[c] {
							delete(candidates, c)
						}
					}
					first = false
					continue
				}
				for c := range candidates {
					if !nodes[k].edges[c] {
						delete(candidates, c)
					}
				}
			}

			for c := range candidates {
				changes = true
				newParty := appendedParty(party, c)
				newParties[getPassword(newParty)] = newParty
			}
		}
		minLength++
		parties = newParties
	}

	longestPassword := ""
	for k := range parties {
		if len(k) > len(longestPassword) {
			longestPassword = k
		}
	}

	fmt.Printf("%d found!\n", len(parties))
	fmt.Println(longestPassword)

	//fmt.Println(password)
}

func main() {
	part2()
}
