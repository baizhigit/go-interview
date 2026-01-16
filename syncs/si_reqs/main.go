package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	fmt.Println("main start")

	urls := []string{
		"https://google.com",
		"https://amazon.com",
		"https://yandex.ru",
		"https://youtube.com",
		"https://rutracker.net",
		"https://google.com",
		"https://amazon.com",
		"https://yandex.ru",
		"https://youtube.com",
		"https://rutracker.net",
		"https://google.com",
		"https://amazon.com",
		"https://yandex.ru",
		"https://youtube.com",
		"https://rutracker.net",
	}
	fmt.Println(process(urls))

	fmt.Println("main end")
}

var client http.Client

const maxConnects = 4

func process(urls []string) map[int]int {
	statusCodeCounts := make(map[int]int, len(urls))
	var wg sync.WaitGroup
	var mu sync.Mutex

	ch := make(chan string, len(urls))

	go func() {
		defer close(ch)

		for _, url := range urls {
			ch <- url
		}
	}()

	processUrl := func(url string) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			fmt.Println(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()

		mu.Lock()
		statusCodeCounts[resp.StatusCode]++
		mu.Unlock()
	}

	wg.Add(maxConnects)
	for range maxConnects {
		go func() {
			defer wg.Done()

			for url := range ch {
				processUrl(url)
			}
		}()
	}

	wg.Wait()
	return statusCodeCounts
}
