package main

import (
	"errors"
	"math"
)

func main() {
	println("=== int overflow ===")
}

var ErrIntOverflow = errors.New("integer overflow")

func Inc(counter int) (int, error) {
	if counter == math.MaxInt {
		return 0, ErrIntOverflow
	}

	return counter + 1, nil
}

func Add(lhs, rhs int) (int, error) {
	if rhs > 0 {
		if lhs > math.MaxInt-rhs {
			return 0, ErrIntOverflow
		}
	} else if rhs < 0 {
		if lhs < math.MinInt-rhs {
			return 0, ErrIntOverflow
		}
	}

	return lhs + rhs, nil
}

func Mul(lhs, rhs int) (int, error) {
	// Handle zero
	if lhs == 0 || rhs == 0 {
		return 0, nil
	}

	// Handle one
	if lhs == 1 || rhs == 1 {
		return lhs * rhs, nil
	}

	// Special case: -1 * MinInt overflows
	if (lhs == -1 && rhs == math.MinInt) || (rhs == -1 && lhs == math.MinInt) {
		return 0, ErrIntOverflow
	}

	// General overflow check
	if lhs > 0 {
		if rhs > 0 && lhs > math.MaxInt/rhs {
			return 0, ErrIntOverflow
		}
		if rhs < 0 && rhs < math.MinInt/lhs {
			return 0, ErrIntOverflow
		}
	} else {
		if rhs > 0 && lhs < math.MinInt/rhs {
			return 0, ErrIntOverflow
		}
		if rhs < 0 && lhs < math.MaxInt/rhs {
			return 0, ErrIntOverflow
		}
	}

	return lhs * rhs, nil
}
