package main

import "fmt"

func main() {
	fmt.Println("main start")

	// Odd numbers > 3 -> 5,7,9... -> cap + 1
	a := []int{}
	a = append(a, []int{1, 2, 3, 4, 5}...)
	fmt.Println("a:", len(a), cap(a))
	b := []int{}
	b = append(b, []int{1, 2, 3, 4, 5, 6}...)
	fmt.Println("b:", len(b), cap(b))
	c := []int{}
	c = append(c, []int{1, 2, 3, 4, 5, 6, 7}...)
	fmt.Println("c:", len(c), cap(c))
	fmt.Println("main end")
}
