// TODO

package main

import (
	"fmt"
	"time"
)

func say(id int, phrase string) {
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Worker %d says: %s\n", id, phrase)
}

func makePool(poolSize int, handler func(int, string)) (func(string), func()) {
	pool := make(chan int, poolSize)

	for i := range poolSize {
		pool <- i
	}

	handle := func(s string) {
		id := <-pool
		go func() {
			defer func() {
				pool <- id
			}()

			handler(id, s)
		}()
	}

	wait := func() {
		for range poolSize {
			<-pool
		}
	}

	return handle, wait
}

func main() {
	fmt.Println("main start")
	phrases := []string{}

	for i := range 100 {
		phrases = append(phrases, fmt.Sprintf("phrase %d", i))
	}

	handle, wait := makePool(5, say)

	for _, phrase := range phrases {
		handle(phrase)
	}

	wait()
	fmt.Println("main end")
}
