package popular

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

type Request struct {
	Payload string
}

type Client interface {
	Send(ctx context.Context, request Request) error
	Limit(ctx context.Context, requests []Request)
}

type client struct{}

func (c *client) Send(ctx context.Context, request Request) error {
	time.Sleep(100 * time.Millisecond)
	fmt.Println("sending req", request.Payload)
	return nil
}

func (c *client) Limit(ctx context.Context, requests []Request) {

}

// огр по коннектам
// var maxConns = 10

// func (c *client) Limit(ctx context.Context, requestsCh chan Request) {
// 	var wg sync.WaitGroup
// 	wg.Add(maxConns)
// 	for range maxConns {
// 		go func() {
// 			defer wg.Done()
// 			for req := range requestsCh {
// 				c.Send(ctx, req)
// 			}
// 		}()
// 	}
// 	wg.Wait()
// }

// огр по горутинам
// var maxGors = 100

// func (c *client) Limit(ctx context.Context, requests []Request) {
// 	tokens := make(chan struct{}, maxGors)

// 	go func() {
// 		for range maxGors {
// 			tokens <- struct{}{}
// 		}
// 	}()

// 	for _, req := range requests {
// 		<-tokens
// 		go func() {
// 			defer func() {
// 				tokens <- struct{}{}
// 			}()
// 			c.Send(ctx, req)
// 		}()
// 	}

// 	for range maxGors {
// 		<-tokens
// 	}
// }

// огр по rps
// var rps = 5
// var burst = 10

// func (c *client) Limit(ctx context.Context, requests []Request) {
// 	ticker := time.NewTicker(time.Second / time.Duration(rps))
// 	defer ticker.Stop()

// 	tickets := make(chan struct{}, burst)
// 	go func() {
// 		for range burst {
// 			tickets <- struct{}{}
// 		}
// 	}()

// 	go func() {
// 		for {
// 			select {
// 			case <-ticker.C:
// 				tickets <- struct{}{}
// 			case <-ctx.Done():
// 				return
// 			}
// 		}
// 	}()

// 	var wg sync.WaitGroup
// 	wg.Add(len(requests))
// 	for _, req := range requests {
// 		<-tickets
// 		go func() {
// 			defer wg.Done()
// 			c.Send(ctx, req)
// 		}()
// 	}

// 	wg.Wait()
// }

func main() {
	fmt.Println("main start")

	client := client{}
	reqs := make([]Request, 1000)
	for i := 0; i < 1000; i++ {
		reqs[i] = Request{Payload: strconv.Itoa(i)}
	}

	// client.Limit(context.Background(), reqs)
	// client.Limit(context.Background(), generate(reqs)) // огр по коннектам
	client.Limit(context.Background(), reqs) // огр по горутинам / rps

	fmt.Println("main end")
}

func generate(reqs []Request) chan Request {
	ch := make(chan Request)

	go func() {
		for _, v := range reqs {
			ch <- v
		}
		close(ch)
	}()

	return ch
}
