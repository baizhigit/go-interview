package main

import (
	"fmt"
)

func test(x interface{}) {
	switch x.(type) {
	case int:
		fmt.Println("int", x)
	case string:
		fmt.Println("string", x)
	case nil:
		fmt.Println("nil", x)
	case func() int:
		f := x.(func() int)
		if f == nil {
			fmt.Println("func is nil")
			return
		}
		fmt.Println("func", f())
	default:
		fmt.Println("unknown")
	}
}

func main() {
	fmt.Println("main start")

	var x = func() int { return 1 }
	x = nil
	test(x)
}
