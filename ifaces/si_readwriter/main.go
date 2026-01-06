package main

import (
	"fmt"
	"io"
)

func readAll(r []io.Reader) {

}

func convert(a []io.ReadWriter) []io.Reader {
	res := make([]io.Reader, 0, len(a))

	for _, v := range a {
		res = append(res, v)
	}
	return res
}

func main() {
	fmt.Println("main start")

	var a = []io.ReadWriter{}

	// readAll(a)
	readAll(convert(a))
}
