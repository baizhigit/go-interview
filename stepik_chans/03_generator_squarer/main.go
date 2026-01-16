// Напишите функции generator и squarer .

// Функция generator принимает на вход контекст и слайс целых чисел,
// элементы которого последовательно записываются в возвращаемый канал.

// Функция squarer принимает на вход контекст и канал целых чисел.
// Функция последовательно читает из канал числа,
// возводит их в квадрат и пишет в возвращаемый канал.

// Обе функции должны уметь завершаться по отмене контекста.

package main

import (
	"context"
	"fmt"
)

func main() {
	fmt.Println("main start")
	ctx := context.Background()
	pipeline := squarer(ctx, generator(ctx, 1, 2, 3))
	for x := range pipeline {
		fmt.Println(x)
	}
	fmt.Println("main end")
}

func generator(ctx context.Context, in ...int) <-chan int {
	outChan := make(chan int)

	go func() {
		defer close(outChan)

		for _, v := range in {
			select {
			case outChan <- v:
			case <-ctx.Done():
				return
			}
		}
	}()

	return outChan
}

func squarer(ctx context.Context, inChan <-chan int) <-chan int {
	outChan := make(chan int)

	go func() {
		defer close(outChan)

		for val := range inChan {
			select {
			case outChan <- val * val:
			case <-ctx.Done():
				return
			}
		}
	}()

	return outChan
}
