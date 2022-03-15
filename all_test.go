package genh

import "strconv"

// equal is simply ==.
func equal[T comparable](v1, v2 T) bool {
	return v1 == v2
}

// equalNaN is like == except that all NaNs are equal.
func equalNaN[T comparable](v1, v2 T) bool {
	isNaN := func(f T) bool { return f != f }
	return v1 == v2 || (isNaN(v1) && isNaN(v2))
}

// equalStr compares ints and strings.
func equalIntStr(v1 int, v2 string) bool {
	return strconv.Itoa(v1) == v2
}

// offByOne returns true if integers v1 and v2 differ by 1.
func offByOne[Elem Integer](v1, v2 Elem) bool {
	return v1 == v2+1 || v1 == v2-1
}
