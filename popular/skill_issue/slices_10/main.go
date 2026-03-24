package main

import "fmt"

const MAX = 5

func main() {
	fmt.Println("main start")
	s := generate()
	res := mutation(s)
	fmt.Println("s", s, res)
	fmt.Println("main end", s[0:MAX])
}

func generate() []int {
	out := make([]int, 0, MAX)
	for i := 1; i < MAX; i++ {
		out = append(out, i)
	}
	return out
}

func mutation(s []int) []int {
	s = append(s, -1)
	return s
}
