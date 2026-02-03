// Необходимо написать worker pool: нужно выполнить параллельно numJobs заданий, используя numWorkers горутин,
// которые запущены единожды за время выполнения програмы.
// Для этого напишите функции worker и main.
// Функция worker:
// на вход получает функцию для выполнения f, канал для получения аргументов jobs и канал для записи результатов results
// читает из jobs и записывает результат выполнения f(job) в results.
// Функция main:
// запускает функцию worker в numWorkers горутинах;
// в качестве первого аргумента worker использует функцию multiplier;
// пишет числа от 1 до numJobs в канал jobs;
// читает и выводит полученные значения из канала results, паралелльно работе воркеров

package main

import (
	"fmt"
	"sync"
)

func worker(f func(int) int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		results <- f(job)
	}
}

const numJobs = 5
const numWorkers = 3

func main() {
	fmt.Println("main start")

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	wg := sync.WaitGroup{}

	multiplier := func(x int) int {
		return x * 10
	}

	wg.Add(numWorkers)
	for range numWorkers {
		go func() {
			defer wg.Done()
			worker(multiplier, jobs, results)
		}()
	}

	for i := range numJobs {
		jobs <- i + 1
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for v := range results {
		fmt.Println(v)
	}

	fmt.Println("main end")
}
