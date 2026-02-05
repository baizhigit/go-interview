// Напишите функцию getFirstResult, которая принимает контекст и запускает конкурентый поиск, возвращая первый
// доступный результат из replicas. Возвращать ошибку контекста, если контекст завершился раньше, чем стал доступен
// какой-то результат из реплики.

// Напишите функцию getResults, которая запускает конкурентный поиск для каждого набора реплик из replicaKinds,
// использую getFirstResult, и возвращает результат для каждого набора реплик.

package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type result struct {
	msg string
	err error
}
type search func() *result
type replicas []search

func fakeSearch(kind string) search {
	return func() *result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return &result{
			msg: fmt.Sprintf("%q result", kind),
		}
	}
}

func getFirstResult(ctx context.Context, replicas replicas) *result {
	if len(replicas) == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := make(chan *result, 1)

	for _, replica := range replicas {
		go func(replica search) {
			select {
			case <-ctx.Done():
			case ch <- replica():
			}
		}(replica)
	}

	select {
	case res := <-ch:
		return res
	case <-ctx.Done():
		return &result{err: ctx.Err()}
	}
}

func getResults(ctx context.Context, replicaKinds []replicas) []*result {
	ch := make(chan *result)
	var wg sync.WaitGroup

	wg.Add(len(replicaKinds))
	for _, reps := range replicaKinds {
		go func(r replicas) {
			defer wg.Done()
			ch <- getFirstResult(ctx, r)
		}(reps)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var res []*result
	for r := range ch {
		res = append(res, r)
	}
	return res
}

func main() {
	fmt.Println("main start")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	replicaKinds := []replicas{
		{fakeSearch("web1"), fakeSearch("web2")},
		{fakeSearch("image1"), fakeSearch("image2")},
		{fakeSearch("video1"), fakeSearch("video2")},
	}

	for _, res := range getResults(ctx, replicaKinds) {
		fmt.Println(res.msg, res.err)
	}

	fmt.Println("main end")
}
