package main

import (
	"fmt"
)

func main() {
	println("=== pointer to pointer ===")
	var value int32 = 100
	pointer := &value

	fmt.Println("value:", *pointer)
	fmt.Println("address:", pointer)

	process(&pointer)

	fmt.Println("value:", *pointer)
	fmt.Println("address:", pointer)
}

func process(temp **int32) {
	fmt.Println("temp:", temp, *temp, &temp)
	var value2 int32 = 200
	*temp = &value2
	// **temp = value2
	fmt.Println("temp:", temp, *temp, &temp)
}
