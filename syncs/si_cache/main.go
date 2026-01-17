// TODO

package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrNotFound = errors.New("key not found")

type ICache interface {
	Get(context.Context, string) (string, error)
	Set(context.Context, string, string) error
	Del(context.Context, string) error
}

type Cache struct {
	storage map[string]string
	mu      sync.Mutex
}

func New() *Cache {
	return &Cache{
		storage: make(map[string]string),
	}
}

func (c *Cache) Get(_ context.Context, key string) (string, error) {
	c.mu.Lock()
	val, ok := c.storage[key]
	c.mu.Unlock()

	if !ok {
		return "", ErrNotFound
	}
	return val, nil
}

func (c *Cache) Set(_ context.Context, key string, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.storage[key] = value

	return nil
}

func (c *Cache) Del(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.storage, key)

	return nil
}

func main() {
	fmt.Println("main start")

	fmt.Println("main end")
}
