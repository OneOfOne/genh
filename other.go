package genh

// ValuesToPtrs converts a slice of values to a slice of pointers
// optionally copying the values instead of pointing to them in the original slice.
func ValuesToPtrs[T any](vals []T, copy bool) []*T {
	out := make([]*T, 0, len(vals))
	for i := range vals {
		var v *T
		if copy {
			cp := vals[i]
			v = &cp
		} else {
			v = &vals[i]
		}
		out = append(out, v)
	}
	return out
}

func PtrsToValues[T any](vals []*T) []T {
	out := make([]T, 0, len(vals))
	for i := range vals {
		v := vals[i]
		out = append(out, *v)
	}
	return out
}

func Iff[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

func IffFn[T any](cond bool, a, b func() T) T {
	if cond {
		return a()
	}
	return b()
}
