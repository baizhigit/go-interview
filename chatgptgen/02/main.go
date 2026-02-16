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

	out, errs := OrderedParallelMap(ctx, in, 3, func(x int) (int, error) {
		if x == 15 {
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

type task struct {
	index int
	val   int
}

func OrderedParallelMap(
	ctx context.Context,
	in <-chan int,
	workers int,
	fn func(int) (int, error),
) (<-chan int, <-chan error) {
	outCh := make(chan int)
	errCh := make(chan error, 1)
	taskCh := make(chan task)

	go func() {
		defer close(taskCh)

		i := 0
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case taskCh <- task{index: i, val: v}:
					i++
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	ctx, cancel := context.WithCancel(ctx)
	resCh := make(chan task)
	var wg sync.WaitGroup
	wg.Add(workers)

	for range workers { // Go 1.22+ синтаксис
		go func() {
			defer wg.Done()

			for t := range taskCh {
				res, err := fn(t.val)
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
				case resCh <- task{index: t.index, val: res}:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	go func() {
		defer close(outCh)
		defer close(errCh)
		defer cancel()

		pending := make(map[int]int)
		nextIndex := 0

		for r := range resCh {
			pending[r.index] = r.val

			for {
				val, ok := pending[nextIndex]
				if !ok {
					break
				}
				delete(pending, nextIndex)

				select {
				case <-ctx.Done():
					return
				case outCh <- val:
					nextIndex++
				}
			}
		}
	}()

	return outCh, errCh
}
