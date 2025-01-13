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
	tree.Insert(40, "K")
	tree.Insert(38, "L")
	tree.Insert(45, "M")
	tree.Insert(50, "N")
	tree.Insert(55, "O")

	if value, found := tree.Search(10); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}
	if value, found := tree.Search(20); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}
	if value, found := tree.Search(5); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}
	if value, found := tree.Search(15); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}
	if value, found := tree.Search(25); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}
	if value, found := tree.Search(3); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}
	if value, found := tree.Search(35); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}
	if value, found := tree.Search(7); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}
	if value, found := tree.Search(12); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}
	if value, found := tree.Search(9); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}
	if value, found := tree.Search(40); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}
	if value, found := tree.Search(38); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}
	if value, found := tree.Search(45); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}
	if value, found := tree.Search(50); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}
	if value, found := tree.Search(55); found {
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
