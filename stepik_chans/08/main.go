// Есть функция executeTask , которая во время исполнения может зависнуть на неопределенно долгое время.
// Реализуйте функцию-обертку executeTaskWithTimeout, которая:
// исполняет executeTask
// принимает аргументом контекст
// завершается либо в результате исполнения executeTask, либо в результате отмены контекста.
// В последнем случае вернуть ошибку контекста.

package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

const timeout = 100 * time.Millisecond

func executeTask() {
	time.Sleep(time.Duration(rand.Intn(3)) * timeout)
}

func executeTaskWithTimeout(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		executeTask()
		// done <- struct{}{}
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	fmt.Println("main start")

	ctx, _ := context.WithTimeout(context.Background(), timeout)
	err := executeTaskWithTimeout(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("main end")
}
