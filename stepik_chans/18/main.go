// Напишите функции produce и main.
// Функция produce:
// на вход получает канал pipe
// бесконечно пишет целые числа в канал pipe, начиная с 0
// по сигналу от main должна завершать работу
// при завершении должна заснуть на 3 секунды, после чего напечатать "produce finished"
// Функция main:
// должна создать канал pipe
// запустить produceCount функций produce и начать читать из канала pipe, печатая каждое число
// при получении числа produceStop из pipe должна перестать читать новые числа из канала и должна отправить
// сигнал в функции produce, заверщающий их работу
// должна дождаться всех сообщений "produce finished" и напечатать "main finished"
// Для реализации требований допускается добавить дополнительные аргументы в функцию produce

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	produceCount = 3
	produceStop  = 10
)

func produce(ctx context.Context, pipe chan<- int, wg *sync.WaitGroup) { // допускается добавить доп. аргументы
	defer wg.Done()
	defer func() {
		time.Sleep(time.Second * 3)
		fmt.Println("produce finished")
	}()

	for i := 0; ; i++ {
		select {
		case pipe <- i:
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	fmt.Println("main start")
	pipe := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	for i := 0; i < produceCount; i++ {
		wg.Add(1)
		go produce(ctx, pipe, &wg)
	}

	for i := range pipe {
		fmt.Println(i)
		if i == produceStop {
			cancel()
			break
		}
	}

	wg.Wait()
	close(pipe)
	fmt.Println("main end")
}
