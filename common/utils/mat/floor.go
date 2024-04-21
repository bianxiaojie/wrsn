package mat

import "math"

func Ceil[T ~float64 | ~int64](value T) T {
	return T(math.Ceil(float64(value)))
}

func Floor[T ~float64 | ~int64](value T) T {
	return T(math.Floor(float64(value)))
}
