package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type User struct {
	Name string
}

func fetch(ctx context.Context, user User) (string, error) {
	// if user.Name == "Ann" {
	// 	return "", errors.New("Ann")
	// }

	select {
	case <-time.After(10 * time.Millisecond):
		return user.Name, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func process(ctx context.Context, users []User) (map[string]int64, error) {
	names := make(map[string]int64, len(users))
	var mu sync.Mutex

	egroup, ectx := errgroup.WithContext(ctx)

	for _, u := range users {
		egroup.Go(
			func() error {
				name, err := fetch(ectx, u)
				if err != nil {
					return err
				}

				mu.Lock()
				names[name] = names[name] + 1
				mu.Unlock()

				return nil
			})
	}

	if err := egroup.Wait(); err != nil {
		return nil, err
	}

	return names, nil
}

func main() {
	fmt.Println("main start")

	names := []User{
		{"Ann"}, {"Bob"}, {"Cindy"}, {"Rob"}, {"Ann"}, {"Bob"}, {"Cindy"}, {"Steve"},
	}
	ctx := context.Background()
	start := time.Now()
	res, err := process(ctx, names)
	if err != nil {
		fmt.Println("an error occured:", err.Error())
	}
	fmt.Println("time:", time.Since(start))
	fmt.Println("main end", res)
}
