package main

// import (
// 	"context"
// 	"fmt"
// )

// type ICache interface {
// 	Get(context.Context, string) (string, error)
// 	Set(context.Context, string, string) error
// 	Del(context.Context, string) error
// }

// type Cache struct {
// }

// func New() *Cache {
// }

// func (c *Cache) Get(_ context.Context, key string) (string, error) {
// }

// func (c *Cache) Set(_ context.Context, key string, value string) error {
// }

// func (c *Cache) Del(_ context.Context, key string) error {
// }

// func main() {
// 	fmt.Println("main start")

// 	ctx := context.Background()

// 	cache := New()
// 	cache.Set(ctx, "Club", "Real")
// 	cache.Set(ctx, "City", "Paris")
// 	if city, err := cache.Get(ctx, "City"); err == nil {
// 		fmt.Println("City is", city)
// 	}
// 	cache.Del(ctx, "City")

// 	fmt.Println("main end", cache)
// }
