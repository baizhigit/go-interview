package main

import (
	"fmt"
	"unsafe"
)

func main() {
	println("=== endianness ===")
	var number int32 = 0x12345678
	pointer := unsafe.Pointer(&number)

	fmt.Printf("0x")
	for i := 0; i < 4; i++ {
		byteValue := *(*int8)(unsafe.Add(pointer, i))
		fmt.Printf("%x", byteValue)
	}
	fmt.Println()

	if IsLittleEndian() {
		fmt.Println("Little endian")
	} else {
		fmt.Println("Big endian")
	}
}

func IsLittleEndian() bool {
	var number int16 = 0x0001
	pointer := (*int8)(unsafe.Pointer(&number))
	return *pointer == 1
}

func IsBigEndian() bool {
	return !IsLittleEndian()
}
