// TODO

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Warehouse struct {
	items map[string]int
	mu    sync.Mutex
}

func NewWarehouse() *Warehouse {
	return &Warehouse{
		items: map[string]int{
			"phone":  10,
			"laptop": 5,
			"tablet": 7,
		},
	}
}

func (w *Warehouse) ReserveItem(item string, qty int) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	if stock, ok := w.items[item]; ok && stock >= qty {
		w.items[item] -= qty
		return true
	}
	return false
}

type Order struct {
	ID       int
	Item     string
	Quantity int
}

func ProcessOrders(id int, w *Warehouse, orders <-chan Order, wg *sync.WaitGroup) {
	defer wg.Done()
	for order := range orders {
		success := w.ReserveItem(order.Item, order.Quantity)
		status := "FAILED"
		if success {
			status = "SUCCESS"
		}
		fmt.Printf("Worker %d: Processing Order #%d (%d x %s) -> %s\n", id, order.ID, order.Quantity, order.Item, status)

		time.Sleep(time.Second * 3)
	}
}

func main() {
	fmt.Println("main start")
	rand.Seed(time.Now().UnixNano())

	warehouse := NewWarehouse()
	orders := make(chan Order, 10)
	var wg sync.WaitGroup

	numWorkers := 3
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go ProcessOrders(i, warehouse, orders, &wg)
	}

	// Генерируем 10 случайных заказов
	go func() {
		for i := 1; i <= 10; i++ {
			item := []string{"phone", "laptop", "tablet"}[rand.Intn(3)]
			qty := rand.Intn(3) + 1
			orders <- Order{ID: i, Item: item, Quantity: qty}
			time.Sleep(time.Millisecond * 200)
		}
		close(orders)
	}()

	wg.Wait()
	fmt.Println("All orders processed. main end")
}
