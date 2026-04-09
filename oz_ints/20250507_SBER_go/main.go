package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	fmt.Println(getRanges([]int{0, 1, 2, 3, 4, 7, 8, 10})) // "0-4,7-8,10"
	fmt.Println(getRanges([]int{4, 7, 10}))                // "4,7,10"
	fmt.Println(getRanges([]int{2, 3, 8, 9}))              // "2-3,8-9"
}

func getRanges(ranges []int) string {
	if len(ranges) == 0 {
		return ""
	}

	var sb strings.Builder
	startIdx := 0

	for i := 1; i <= len(ranges); i++ {
		// End of array or sequence broken
		if i == len(ranges) || ranges[i] != ranges[i-1]+1 {
			if startIdx > 0 {
				sb.WriteString(",")
			}

			// Format: single number or range
			if i-1 == startIdx {
				sb.WriteString(strconv.Itoa(ranges[startIdx]))
			} else {
				sb.WriteString(strconv.Itoa(ranges[startIdx]))
				sb.WriteString("-")
				sb.WriteString(strconv.Itoa(ranges[i-1]))
			}

			startIdx = i
		}
	}

	return sb.String()
}
