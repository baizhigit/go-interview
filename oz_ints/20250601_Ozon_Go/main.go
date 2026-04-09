package main

import "fmt"

func main() {
	fmt.Println("main start")
	fmt.Println(isPowerOfTwo(120))
	fmt.Println(isPowerOfTwo(4096))

	fmt.Println(isPowerOfTwo2(120))
	fmt.Println(isPowerOfTwo2(4096))
	fmt.Println("main end")
}

// 8 -> 1000 7 -> 0111 --> 1000 & 0111 = 0
func isPowerOfTwo(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}

// делим число на 2, пока оно делится без остатка.
// Если в итоге получили 1 — это степень двойки
func isPowerOfTwo2(n int) bool {
	if n <= 0 {
		return false
	}
	for n%2 == 0 {
		n /= 2
	}
	return n == 1
}
