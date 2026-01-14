package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("main start")
	reader(doubler(writer()))
}

func writer() <-chan int {
	ch := make(chan int)

	go func() {
		for i := range 10 {
			ch <- i + 1
		}
		close(ch)
	}()

	return ch
}

func doubler(wr <-chan int) <-chan int {
	ch := make(chan int)

	go func() {
		for v := range wr {
			time.Sleep(time.Millisecond * 500)
			ch <- v * 2
		}
		close(ch)
	}()

	return ch
}

func reader(db <-chan int) {
	for v := range db {
		fmt.Println(v)
	}
}
