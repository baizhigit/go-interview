// Напишите функцию inc, которая принимает на вход канал, читает из него значения
// и записывает эти значения, увеличенные на единицу, в возвращаемый канал.
// Дополните функцию main созданием цепочки каналов, используя inc, так,
// чтобы программа завершалась без паники.

package main

import (
	"fmt"
)

func main() {
	fmt.Println("main start")
	first := make(chan int)
	last := make(<-chan int)
	n := 10
	_ = last

	last = inc(0, first)
	for i := 1; i < n; i++ {
		last = inc(i, last)
	}

	first <- 0
	close(first)
	if n != <-last {
		panic("wrong code")
	}
	fmt.Println("main end")
}

func inc(iter int, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for val := range in {
			fmt.Println("iter val", iter, val+1)
			out <- val + 1
		}
	}()

	return out
}
