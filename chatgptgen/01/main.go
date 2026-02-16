// TODO

package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {
	fmt.Println("main start")

	ctx := context.Background()

	in := make(chan int)
	go func() {
		defer close(in)
		for i := 1; i <= 10; i++ {
			in <- i
		}
	}()

	out, errs := ParallelMap(ctx, in, 3, func(x int) (int, error) {
		if x == 5 {
			return 0, fmt.Errorf("boom")
		}
		return x * 2, nil
	})

	for {
		select {
		case v, ok := <-out:
			if ok {
				fmt.Println("result:", v)
			}
		case err, ok := <-errs:
			if ok {
				fmt.Println("error:", err)
			}
			return
		}
	}
}

func ParallelMap(
	ctx context.Context,
	in <-chan int,
	workers int,
	fn func(int) (int, error),
) (<-chan int, <-chan error) {
	resCh := make(chan int)
	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(workers)

	ctx, cancel := context.WithCancel(ctx)

	for range workers {
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case i, ok := <-in:
					if !ok {
						return
					}

					res, err := fn(i)
					if err != nil {
						select {
						case errCh <- err:
						default:
						}
						cancel()
						return
					}

					select {
					case <-ctx.Done():
						return
					case resCh <- res:
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resCh)
		close(errCh)
		cancel()
	}()

	return resCh, errCh
}
