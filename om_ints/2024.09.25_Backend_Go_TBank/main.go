package main

import (
	"fmt"
	"sync"
	"time"
)

func UpdateProductStock() <-chan map[string]int {
	stockUpdates := make(chan map[string]int)

	go func() {
		defer close(stockUpdates)

		currentStock := map[string]int{
			"Apples":  50,
			"Bananas": 30,
			"Oranges": 20,
			"Grapes":  15,
		}

		for i := 0; i < 5; i++ {
			newStock := make(map[string]int)
			for product, quantity := range currentStock {
				newStock[product] = int(float64(quantity) * 0.95)
			}
			stockUpdates <- newStock
			currentStock = newStock
			time.Sleep(150 * time.Millisecond)
		}
	}()

	return stockUpdates
}

func main() {
	fmt.Println("main start")

	stockStream := UpdateProductStock()

	var stockHistory []map[string]int

	for stock := range stockStream {
		stockHistory = append(stockHistory, stock)
	}

	for i, stock := range stockHistory {
		fmt.Printf("Iteration %d: %v\n", i+1, stock)
	}

	p := RunWriter()
	var prices []map[string]int
	for v := range p {
		prices = append(prices, v)
	}
	for _, price := range prices {
		fmt.Println(price)
	}
	fmt.Println("*******************")

	var wg sync.WaitGroup
	var mu sync.RWMutex
	for i := 0; i < 10; i++ {
		wg.Add(1)
		RunProcessor(&wg, &mu, prices)
	}
	wg.Wait()
}

func RunProcessor(wg *sync.WaitGroup, mu *sync.RWMutex, prices []map[string]int) {
	go func() {
		defer wg.Done()
		for _, price := range prices {
			mu.Lock()
			for key, value := range price {
				price[key] = value + 1
			}
			mu.Unlock()

			mu.RLock()
			fmt.Println(price)
			mu.RUnlock()
		}
	}()
}

func RunWriter() <-chan map[string]int {
	var prices = make(chan map[string]int)
	go func() {
		defer close(prices)
		var currentPrice = map[string]int{
			"AAPL": 163,
			"USD":  117,
			"EUR":  124,
			"NVDA": 234,
		}
		for i := 0; i < 5; i++ {
			newPrice := make(map[string]int)
			for key, value := range currentPrice {
				newPrice[key] = int(float64(value) * 1.3)
			}
			fmt.Println("for", i, currentPrice, newPrice)
			prices <- newPrice
			currentPrice = newPrice
			time.Sleep(time.Millisecond * 100)
		}
	}()
	return prices
}
