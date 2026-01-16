package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Response struct {
	Start   int      `json:"start"`
	Size    int      `json:"size"`
	Page    int      `json:"page"`
	Pages   int      `json:"pages"`
	Total   int      `json:"total"`
	Results []Result `json:"results"`
}

type Result struct {
	Number string `json:"number"`
}

const (
	fromPage   = 1001
	toPage     = 1932 // ← лимит страниц (можно 7570:1932)
	limit      = 100
	categoryID = 136 // 0:134	2500:136

	workers    = 5
	reqPerSec  = 5 // ← RATE LIMIT
	maxRetries = 3

	outFile = "numbers.txt"
)

var client = &http.Client{
	Timeout: 20 * time.Second,
}

type RateLimiter struct {
	ticker *time.Ticker
}

func NewRateLimiter(rps int) *RateLimiter {
	return &RateLimiter{
		ticker: time.NewTicker(time.Second / time.Duration(rps)),
	}
}

func (r *RateLimiter) Wait() {
	<-r.ticker.C
}

func fetchOffsetWithRetry(
	ctx context.Context,
	offset int,
	limit int,
	categoryID int,
	rl *RateLimiter,
) (*Response, error) {

	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		rl.Wait()

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			fmt.Sprintf(
				"https://api.izi.me/numbers?limit=%d&offset=%d&categoryId=%d",
				limit,
				offset,
				categoryID,
			),
			nil,
		)
		if err != nil {
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
		} else {
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				var data Response
				if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
					return nil, err
				}
				return &data, nil
			}

			lastErr = fmt.Errorf("status %s", resp.Status)
		}

		// exponential backoff
		sleep := time.Duration(attempt*attempt) * time.Second
		log.Printf("retry %d offset %d after %v", attempt, offset, sleep)
		time.Sleep(sleep)
	}

	return nil, lastErr
}

func worker(
	ctx context.Context,
	wg *sync.WaitGroup,
	offsets <-chan int,
	numbers chan<- string,
	rl *RateLimiter,
	processed *int64,
	total int,
) {
	defer wg.Done()

	for offset := range offsets {
		resp, err := fetchOffsetWithRetry(ctx, offset, limit, categoryID, rl)
		if err != nil {
			log.Printf("offset %d failed: %v", offset, err)
			continue
		}

		for _, r := range resp.Results {
			numbers <- r.Number
			atomic.AddInt64(processed, 1)
		}

		done := atomic.LoadInt64(processed)
		percent := float64(done) / float64(total) * 100

		log.Printf(
			"offset %d done | progress: %d / %d (%.2f%%)",
			offset,
			done,
			total,
			percent,
		)
	}
}

func writer(
	wg *sync.WaitGroup,
	filePath string,
	numbers <-chan string,
) {
	defer wg.Done()

	file, err := os.OpenFile(
		filePath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Fatalf("file error: %v", err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()

	for num := range numbers {
		if _, err := w.WriteString(num + "\n"); err != nil {
			log.Printf("write error: %v", err)
			return
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	offsets := make(chan int, workers)
	numbers := make(chan string, 2000)

	var processed int64
	totalExpected := (toPage - fromPage + 1) * limit

	rl := NewRateLimiter(reqPerSec)

	wgWorkers := sync.WaitGroup{}
	wgWriter := sync.WaitGroup{}

	// writer
	wgWriter.Add(1)
	go writer(&wgWriter, outFile, numbers)

	// workers
	for i := 0; i < workers; i++ {
		wgWorkers.Add(1)
		go worker(
			ctx,
			&wgWorkers,
			offsets,
			numbers,
			rl,
			&processed,
			totalExpected,
		)
	}

	// offsets producer
	for page := fromPage; page <= toPage; page++ {
		offset := (page - 1) * limit
		offsets <- offset
	}
	close(offsets)

	wgWorkers.Wait()
	close(numbers)
	wgWriter.Wait()

	log.Printf("DONE: %d numbers saved", processed)
}
