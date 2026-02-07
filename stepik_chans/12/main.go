// Есть интерфейс waiter который должен:
// 1. параллельно запускать переданные в run функции с указанным контекстом.
// 2. количество параллельных запусков определяется параметром maxParallel при создании waiter через newGroupWait
// 3. возвращать ошибку из wait, если хотя бы одна функция из run вернула ее
// 4. возвращать комбинацию ошибок от вызовов run, если несколько задач завершились с ошибками (можно использовать errors.Join)
// Реализуйте поля и методы структуры waitGroup для интерфейса waiter.

package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type waiter interface {
	wait() error
	run(ctx context.Context, f func(ctx context.Context) error)
}

type pair struct {
	ctx context.Context
	fn  func(ctx context.Context) error
}

type waitGroup struct {
	mu       sync.Mutex
	queue    chan pair
	queueWg  sync.WaitGroup
	workerWg sync.WaitGroup
	err      error
}

func (g *waitGroup) wait() error {
	g.queueWg.Wait()
	close(g.queue)
	g.workerWg.Wait()

	return g.err
}

func (g *waitGroup) run(ctx context.Context, fn func(ctx context.Context) error) {
	g.queueWg.Add(1)

	go func() {
		defer g.queueWg.Done()
		g.queue <- pair{ctx: ctx, fn: fn}
	}()
}

func newGroupWait(maxParallel int) waiter {
	g := &waitGroup{
		queue:    make(chan pair),
		queueWg:  sync.WaitGroup{},
		workerWg: sync.WaitGroup{},
	}

	g.workerWg.Add(maxParallel)
	for i := 0; i < maxParallel; i++ {
		go func() {
			defer g.workerWg.Done()

			for p := range g.queue {
				err := p.fn(p.ctx)
				if err != nil {
					g.mu.Lock()
					g.err = errors.Join(g.err, err)
					g.mu.Unlock()
				}
			}
		}()
	}

	return g
}

func main() {
	fmt.Println("main start")

	g := newGroupWait(2)
	ctx := context.Background()
	expErr1 := errors.New("got error 1")
	expErr2 := errors.New("got error 2")

	g.run(ctx, func(ctx context.Context) error {
		return nil
	})
	g.run(ctx, func(ctx context.Context) error {
		return expErr2
	})
	g.run(ctx, func(ctx context.Context) error {
		return expErr1
	})

	err := g.wait()
	if !errors.Is(err, expErr1) || !errors.Is(err, expErr2) {
		panic("wrong code")
	}

	fmt.Println("main end", err)
}
