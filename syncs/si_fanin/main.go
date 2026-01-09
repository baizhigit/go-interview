package main

import (
	"fmt"
	"sync"
)

func fanin(chans ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for _, ch := range chans {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()

			for val := range ch {
				out <- val
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	fmt.Println("main start")

	res := fanin()
	fmt.Println("main end", res)
}
