package genh

import "math"

func Min[T Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Abs[T Integer | Float](v T) T {
	return T(math.Abs(float64(v)))
}
