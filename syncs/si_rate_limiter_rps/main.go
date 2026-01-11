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
	timer := time.NewTimer(2 * time.Second)
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
// Ограничение кол-ва rps

const rps = 100
const burst = 10

func (c client) WithLimiter(ctx context.Context, reqs []Request) {
	ticker := time.NewTicker(time.Second / time.Duration(rps))
	defer ticker.Stop()

	tickets := make(chan struct{}, burst)

	var wg sync.WaitGroup

	for range burst {
		tickets <- struct{}{}
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				select {
				case tickets <- struct{}{}:
				default:
					// tickets уже заполнен, skip
				}

			case <-ctx.Done():
				return
			}
		}
	}()

loop:
	for _, req := range reqs {
		select {
		case <-tickets:
			wg.Add(1)
			go func(r Request) {
				defer wg.Done()

				if err := c.SendRequest(ctx, r); err != nil {
					fmt.Printf("Error sending request %s: %v\n", r.Payload, err)
				}
			}(req)
		case <-ctx.Done():
			break loop
		}
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
