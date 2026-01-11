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
	WithLimiter(ctx context.Context, ch <-chan Request)
}

type client struct {
}

func (c client) SendRequest(ctx context.Context, request Request) error {
	select {
	case <-time.After(250 * time.Millisecond):
		fmt.Println("sending request", request.Payload)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Rate Limiter
// Ограничение кол-ва коннектов
// (worker pool pattern)

const maxConnects = 10

func (c client) WithLimiter(ctx context.Context, ch <-chan Request) {
	var wg sync.WaitGroup
	wg.Add(maxConnects)

	for range maxConnects {
		go func() {
			defer wg.Done()

			for {
				select {
				case req, ok := <-ch:
					if !ok {
						return
					}
					if err := c.SendRequest(ctx, req); err != nil {
						fmt.Printf("Error sending request %s: %v\n", req.Payload, err)
					}
				case <-ctx.Done():
					return
				}
			}
		}()
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
	c.WithLimiter(ctx, generate(requests))
}

func generate(reqs []Request) <-chan Request {
	ch := make(chan Request)

	go func() {
		defer close(ch)

		for _, req := range reqs {
			ch <- req
		}
	}()

	return ch
}
