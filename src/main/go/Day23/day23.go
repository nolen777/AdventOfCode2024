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

type Party map[string]bool

func MakeParty(t Triple) Party {
	return Party{t.a: true, t.b: true, t.c: true}
}
func (p Party) AppendedParty(s string) Party {
	pc := maps.Clone(p)
	pc[s] = true
	return pc
}

func (p Party) GetPassword() string {
	elts := make([]string, 0, len(p))
	for k := range p {
		elts = append(elts, k)
	}
	slices.Sort(elts)
	return strings.Join(elts, ",")
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

	parties := map[string]Party{}

	for triple := range triples {
		newParty := MakeParty(triple)
		parties[newParty.GetPassword()] = newParty
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
				newParty := party.AppendedParty(c)
				newParties[newParty.GetPassword()] = newParty
			}
		}
		minLength++
		parties = newParties
	}

	if len(parties) != 1 {
		log.Fatalf("Failed; expected 1 party but got %d\n", len(parties))
	}

	longestPassword := ""
	var longestParty Party
	for k, v := range parties {
		longestPassword = k
		longestParty = v
	}

	fmt.Printf("found of length %d\n", (len(longestParty)))
	fmt.Println(longestPassword)
}

func main() {
	part2()
}
