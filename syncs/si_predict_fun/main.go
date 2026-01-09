package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func unpredictableFunc() int {
	n := rand.Intn(10)
	time.Sleep(time.Duration(n) * time.Second)
	return n
}

var timeout = time.Second * 3

func predictableFunc(ctx context.Context) (int, error) {

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	resCh := make(chan int, 1)
	go func() {
		resCh <- unpredictableFunc()
		close(resCh)
	}()

	select {
	case result := <-resCh:
		return result, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func main() {
	fmt.Println("main start")

	num, err := predictableFunc(context.Background())
	fmt.Println("main end", num, err)
}
