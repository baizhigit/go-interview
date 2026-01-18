// Graceful shutdown with channels

package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	log.Println("main start")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Printf("Received signal: %v. Initiating graceful shutdown...", ctx.Err())
				return
			case <-ticker.C:
				log.Println("Doing work...")
			}
		}
	}()

	wg.Wait()
	log.Println("main end")
}
