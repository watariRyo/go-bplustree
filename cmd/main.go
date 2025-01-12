package main

import (
	"fmt"

	"github.com/watariRyo/bptree/tree"
)

func main() {
	tree := tree.NewBPlusTree()

	tree.Insert(10, "A")
	tree.Insert(20, "B")
	tree.Insert(5, "C")
	tree.Insert(15, "D")
	tree.Insert(25, "E")
	tree.Insert(3, "F")
	tree.Insert(35, "G")
	tree.Insert(7, "H")
	tree.Insert(12, "I")
	tree.Insert(9, "J")

	if value, found := tree.Search(20); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}

	if value, found := tree.Search(15); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}

	if value, found := tree.Search(7); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}

	if value, found := tree.Search(30); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}
}
