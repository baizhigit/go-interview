package main

import (
	"fmt"
)

type c chan c

func main() {
	fmt.Println("main start")

	var c = make(c, 1)
	c <- c

	for i := 0; i < 1000; i++ {
		select {
		case <-c:
		case <-c:
			c <- c
		default:
			fmt.Println(i)
			return
		}
	}
}
