package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type ICache interface {
	Get(context.Context, string) (string, error)
	Set(context.Context, string, string) error
	Del(context.Context, string) error
}

type Cache struct {
	store map[string]string
	mu    sync.Mutex
}

func New() *Cache {
	return &Cache{
		store: make(map[string]string),
	}
}

func (c *Cache) Get(_ context.Context, key string) (string, error) {
	c.mu.Lock()
	val, ok := c.store[key]
	c.mu.Unlock()

	if !ok {
		return "", errors.New("not found")
	}
	return val, nil
}

func (c *Cache) Set(_ context.Context, key string, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = value
	return nil
}

func (c *Cache) Del(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.store, key)
	return nil
}

func main() {
	fmt.Println("main start")

	ctx := context.Background()

	cache := New()
	cache.Set(ctx, "Club", "Real")
	cache.Set(ctx, "City", "Paris")
	if city, err := cache.Get(ctx, "City"); err == nil {
		fmt.Println("City is", city)
	}
	cache.Del(ctx, "City")

	fmt.Println("main end", cache)
}
