package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	value int
	next  *Node
}

func parseNums(fileName string) *Node {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var head *Node
	var tail *Node

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		textLine := scanner.Text()

		elements := strings.Split(textLine, " ")

		for _, elt := range elements {
			num, err := strconv.Atoi(elt)

			if err != nil {
				fmt.Println(err)
				return head
			}

			newNode := new(Node)
			newNode.value = num
			if head == nil {
				head = newNode
			}
			if tail != nil {
				tail.next = newNode
			}
			tail = newNode
		}
		break
	}
	return head
}

func printStones(head *Node) {
	for iter := head; iter != nil; iter = iter.next {
		fmt.Print(iter.value, " ")
	}
	fmt.Println("")
}

func digitCount(n int) int {
	if n == 0 {
		return 1
	}

	count := 0
	for n != 0 {
		n /= 10
		count++
	}
	return count
}

func split(n int, digitCount int) (int, int) {
	m := 0
	mult := 1
	for i := 0; i < digitCount/2; i++ {
		d := n % 10
		m += d * mult
		n /= 10
		mult *= 10
	}
	return n, m
}

//func blink(stones []int) []int {
//	idx := 0
//
//	for idx < len(stones) {
//		e := stones[idx]
//		dc := digitCount(e)
//
//		switch {
//		case e == 0:
//			stones[idx] = 1
//		case dc%2 == 0:
//			// split
//			newStones := make([]int, 0, len(stones)+dc/2)
//			newStones = append(newStones, stones[:idx]...)
//			first, second := split(e, dc)
//			newStones = append(newStones, first, second)
//			newStones = append(newStones, stones[idx+1:]...)
//
//			stones = newStones
//			idx++
//
//		default:
//			stones[idx] = e * 2024
//		}
//		idx++
//	}
//
//	return stones
//}

func blink(head *Node) *Node {
	idx := 0

	for iter := head; iter != nil; iter = iter.next {
		e := iter.value
		dc := digitCount(e)

		switch {
		case e == 0:
			iter.value = 1
		case dc%2 == 0:
			// split
			first, second := split(e, dc)
			newNode := new(Node)
			newNode.value = second
			newNode.next = iter.next
			iter.value = first
			iter.next = newNode

			iter = newNode

		default:
			iter.value = e * 2024
		}
		idx++
	}

	return head
}

func count(head *Node) int {
	c := 0
	for iter := head; iter != nil; iter = iter.next {
		c++
	}
	return c
}

func main() {
	stones := parseNums("resources/Day11/input.txt")

	printStones(stones)

	for i := 0; i < 25; i++ {
		stones = blink(stones)
		fmt.Println(i)
		//	fmt.Println(stones)
	}

	fmt.Println(count(stones), " stones")
}
