package main

import (
	"fmt"
	"math/rand"
)

func uniqRand(n int, max int) ([]int, error) {
	// 1. Edge case validation
	if n > max {
		return nil, fmt.Errorf("cannot pick %d unique numbers from %d", n, max)
	}

	out := make([]int, 0, n)
	umap := make(map[int]struct{}, n)

	for len(out) < n {
		rand := rand.Intn(max)

		if _, ok := umap[rand]; !ok {
			umap[rand] = struct{}{}
			out = append(out, rand)
		}
	}

	return out, nil
}

func addNum(nums []int) {
	fmt.Printf("addNum len: %d, cap: %d\n", len(nums), cap(nums))
	nums = append(nums, 4)
	fmt.Printf("addNum after len: %d, cap: %d\n", len(nums), cap(nums))
}

func addNums(nums []int) {
	fmt.Printf("addNums len: %d, cap: %d\n", len(nums), cap(nums))
	nums = append(nums, 5, 6)
	fmt.Printf("addNums after len: %d, cap: %d\n", len(nums), cap(nums), nums)
}

func main() {
	fmt.Println(uniqRand(10, 100))

	nums := []int{1, 2, 3}

	addNum(nums[0:2])
	fmt.Println(nums)

	addNums(nums[0:2])
	fmt.Println(nums)
}

// SELECT user_id, count(id) as order_count FROM orders WHERE price_total >= 1000
// GROUP BY user_id ORDER BY order_count DESC
