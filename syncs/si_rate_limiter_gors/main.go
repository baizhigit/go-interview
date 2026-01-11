package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Request struct {
	Payload string
}

type Client interface {
	SendRequest(ctx context.Context, request Request) error
	WithLimiter(ctx context.Context, requests []Request)
}

type client struct {
}

func (c client) SendRequest(ctx context.Context, request Request) error {
	timer := time.NewTimer(500 * time.Millisecond)
	defer timer.Stop()

	select {
	case <-timer.C:
		fmt.Println("sending request", request.Payload)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Rate Limiter
// Ограничение кол-ва горутин
// token-based limiter
// Semaphore pattern using buffered channel

const maxGoroutines = 100

func (c client) WithLimiter(ctx context.Context, reqs []Request) {
	tokens := make(chan struct{}, maxGoroutines)

	for range maxGoroutines {
		tokens <- struct{}{}
	}

	var wg sync.WaitGroup

	for _, req := range reqs {
		if ctx.Err() != nil {
			break
		}

		<-tokens
		wg.Add(1)

		go func(r Request) {
			defer func() {
				tokens <- struct{}{}
				wg.Done()
			}()

			if err := c.SendRequest(ctx, r); err != nil {
				fmt.Printf("Error sending request %s: %v\n", r.Payload, err)
			}
		}(req)
	}

	wg.Wait()
}

func main() {
	fmt.Println("main start")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	c := client{}
	requests := make([]Request, 1000)
	for i := 0; i < 1000; i++ {
		requests[i] = Request{Payload: strconv.Itoa(i)}
	}
	c.WithLimiter(ctx, requests)
}
