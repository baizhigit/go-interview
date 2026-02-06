// Напишите функцию bridge, которая вычитывает очередной канал из канала ins и перенаправляет данные из вычитанного
// канала в возвращаемый канал, пока вычитанный канал открыт и контекст не отменен.
// Функция продолжает работу, пока канал ins открыт и контекст не отменен.

package main

import (
	"context"
	"fmt"
	"reflect"
)

func main() {
	fmt.Println("main start")

	genVals := func() <-chan <-chan interface{} {
		out := make(chan (<-chan interface{}))
		go func() {
			defer close(out)
			for i := 0; i < 3; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				out <- stream
			}
		}()
		return out
	}

	var res []interface{}
	for v := range bridge(context.Background(), genVals()) {
		res = append(res, v)
	}

	if !reflect.DeepEqual(res, []interface{}{0, 1, 2}) {
		panic("wrong code")
	}

	fmt.Println("main end", res)
}

func bridge(ctx context.Context, ins <-chan <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		for {
			var stream <-chan interface{}
			select {
			case <-ctx.Done():
				return
			case s, ok := <-ins:
				if !ok {
					return
				}
				stream = s
			}

			for val := range orDone(ctx, stream) {
				select {
				case out <- val:
				case <-ctx.Done():
				}
			}
		}
	}()

	return out
}

func orDone(ctx context.Context, in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- v:
				case <-ctx.Done():
				}
			}
		}
	}()
	return out
}
