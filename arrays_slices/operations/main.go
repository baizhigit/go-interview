package main

import (
	"fmt"
	"unsafe"
)

func main() {
	println("=== arrays and slices operations ===")
}

func accessToArrayElement1() {
	data := [3]int{1, 2, 3}
	idx := 4
	fmt.Println(data[idx]) // panic
	// fmt.Println(data[4])   // compilation error
}

func accessToArrayElement2() {
	data := [3]int{1, 2, 3}
	idx := -1
	fmt.Println(data[idx]) // panic
	// fmt.Println(data[-1])  // compilation error
}

func arrayLen() {
	data := [10]int{}
	fmt.Println(len(data)) // 10
	fmt.Println(cap(data)) // 10
}

func emptyArray() {
	var data [10]byte
	fmt.Println(len(data))           // 10
	fmt.Println(unsafe.Sizeof(data)) // 10
}

func zeroArray() {
	var data [0]byte
	fmt.Println(len(data))           // 0
	fmt.Println(unsafe.Sizeof(data)) // 0
}

func negativeArray() {
	// var data [-1]byte // compilation error
}

func arrayCreation() {
	length1 := 100
	// var data1 [length1]int // compilation error
	_ = length1

	const length2 = 100
	var data2 [length2]int
	_ = data2
}

func accessToSliceElement1() {
	data := make([]int, 3)
	fmt.Println(data[3]) // panic
}

func accessToSliceElement2() {
	data := make([]int, 3, 5)
	fmt.Println(data[3]) // panic
}

func accessToNilSlice() {
	var data []int
	_ = data[0]             // panic
	data = append(data, 10) // ok
	for range data {
	} // ok
}

func increaseCapacity() {
	data := make([]int, 0, 10)
	data = data[:10:100] // panic
}
