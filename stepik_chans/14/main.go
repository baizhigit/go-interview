// Напишите функцию Run, которая запускает конкурентное выполнение функций fs и дожидается их окончания.
// Если одна или несколько функций из fs завершились с ошибкой, Run возвращает любую из них.

package main

import (
	"errors"
	"fmt"
	"sync"
)

type fn func() error

func main() {
	fmt.Println("main start")

	expErr := errors.New("random error")
	funcs := []fn{
		func() error { return nil },
		func() error { return nil },
		func() error { return expErr },
		func() error { return nil },
	}
	err := Run(funcs...)
	if err != nil && !errors.Is(err, expErr) {
		panic("wrong code")
	}

	fmt.Println("main end", err)
}

func Run(fs ...fn) error {
	errCh := make(chan error, len(fs))
	var wg sync.WaitGroup
	wg.Add(len(fs))
	for _, f := range fs {
		go func(f fn) {
			defer wg.Done()

			errCh <- f()
		}(f)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}
