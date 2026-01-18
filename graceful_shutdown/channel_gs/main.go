// Graceful shutdown with channels

package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	log.Println("main start")

	// Shutdown signal
	done := make(chan struct{})

	// Capture OS signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		sig := <-sigCh
		log.Printf("Received signal: %v. Initiating graceful shutdown...", sig)
		close(done)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				log.Println("Worker stopping gracefully...")
				return
			case <-ticker.C:
				log.Println("Doing work...")
			}
		}
	}()

	wg.Wait()
	log.Println("main end")
}
