package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// 19-40

func main() {
	// 19-40
	err := foo(true)
	if err != nil {
		fmt.Println("fire error 1", err)
	}

	err2 := foo(false)
	if err2 != nil {
		fmt.Println("fire error 2", err2)
	}

	fmt.Println("complete")

	// 23-20
	var urls = []string{
		"https://kaspi.kz",
		"http://ya.ru",
		"https://ya.ru",
		"http://somesite.com",
	}

	processUrls(urls)
}

type CustomError struct {
	Text string
}

func (s *CustomError) Error() string {
	return s.Text
}

func foo(fireError bool) error {
	var c *CustomError

	if fireError {
		return &CustomError{Text: "someError"}
	}

	return c
}

// Указатель на данные равен nil, но указатель на itable не равен nil

// 23-20
var client = http.Client{}

func processUrls(urls []string) {
	// statuses := make(map[int]int)
	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		go func() {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				fmt.Println("request error", err.Error())
				return
			}

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("response error", err.Error())
				return
			}
			defer resp.Body.Close()

			// statuses[resp.StatusCode]++

			if resp.StatusCode == 200 {
				fmt.Printf("%s - ok\n", url)
			} else {
				fmt.Printf("%s - not ok\n", url)
			}

		}()
	}

	wg.Wait()
}
