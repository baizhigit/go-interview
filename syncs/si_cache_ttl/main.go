// TODO

package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrNotFound = errors.New("key not found")

type ICache interface {
	Get(context.Context, string) (string, error)
	Set(context.Context, string, string) error
	Del(context.Context, string) error
}

type elem struct {
	value    string
	exp_date time.Time
}

type Cache struct {
	storage map[string]elem
	mu      sync.Mutex
	TTL     time.Duration
	done    chan struct{}
}

func New(ttl time.Duration) *Cache {
	cache := &Cache{
		storage: make(map[string]elem),
		TTL:     ttl,
		done:    make(chan struct{}),
	}
	cache.clearByTTL()

	return cache
}

func (c *Cache) clearByTTL() {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.clear()
			case <-c.done:
				return
			}
		}
	}()
}

func (c *Cache) Stop() {
	select {
	case <-c.done:
	default:
		close(c.done)
	}
}

func (c *Cache) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, el := range c.storage {
		if el.exp_date.Before(time.Now()) {
			delete(c.storage, key)
		}
	}
}

func (c *Cache) Get(_ context.Context, key string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	el, ok := c.storage[key]
	if !ok {
		return "", ErrNotFound
	}

	if el.exp_date.Before(time.Now()) {
		c.delete(key)
		return "", ErrNotFound
	}
	return el.value, nil
}

func (c *Cache) Set(_ context.Context, key string, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	el := elem{
		value:    value,
		exp_date: time.Now().Add(c.TTL),
	}
	c.storage[key] = el

	return nil
}

func (c *Cache) delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.storage, key)
}

func (c *Cache) Del(_ context.Context, key string) error {
	c.delete(key)

	return nil
}

func main() {
	fmt.Println("main start")

	fmt.Println("main end")
}
