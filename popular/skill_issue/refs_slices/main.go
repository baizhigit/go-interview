package main

import "fmt"

func changeSlice(s []int) {
	s[1] = 10
}

func getSlice() []int {
	s := []int{1, 2, 3}

	defer changeSlice(s)

	return s
}

func main() {
	fmt.Println("main start")
	p := getSlice()

	fmt.Println("main end", p)
}
