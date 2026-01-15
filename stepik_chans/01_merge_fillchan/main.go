// Напишите функции merge и fillChan .

// Функция fillChan :
// на вход получает целое число n;
// возвращает канал;
// пишет в этот канал n чисел от 0 до n-1.

// Функция merge :
// получает на вход массив каналов cs ;
// возвращает канал;
// параллельно читает из каждого канала из cs и пишет полученное значение в возвращаемый канал.

package main

import (
	"fmt"
	"sync"
)

// merge - соединяет каналы в один
func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	wg.Add(len(cs))
	out := make(chan int)

	for _, ch := range cs {
		go func(ch <-chan int) {
			defer wg.Done()

			for v := range ch {
				out <- v
			}
		}(ch)
	}

	go func() {
		wg.Wait()

		close(out)
	}()

	return out
}

// fillChan - заполняет канал числами от 0 до n-1
func fillChan(n int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)

		for v := range n {
			ch <- v
		}
	}()
	return ch
}

func main() {
	fmt.Println("main start")
	a := fillChan(2)
	b := fillChan(3)
	c := fillChan(4)
	d := merge(a, b, c)

	for v := range d {
		fmt.Println(v)
	}
	fmt.Println("main end")
}
