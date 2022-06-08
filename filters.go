package genh

func GroupBy[M ~map[MapKey]MapVal, GM map[MapKey][]MapVal, MapKey comparable, MapVal any](in M, fn func(k MapKey, v MapVal) MapKey) (out GM) {
	out = make(GM)
	for k, v := range in {
		gk := fn(k, v)
		out[gk] = append(out[gk], v)
	}
	for k := range out {
		out[k] = Clip(out[k])
	}
	return
}

// Filter filters a slice optionally in place.
func Filter[S ~[]E, E any](in S, f func(E) (keep bool), inplace bool) (out S) {
	if inplace {
		out = in[:0]
	} else {
		out = make(S, 0, len(out))
	}

	for _, v := range in {
		if f(v) {
			out = append(out, v)
		}
	}

	return Clip(out)
}

// SliceMap takes a slice of type E, calls fn on each value of `in` and returns the results as a slice of type `T`
func SliceMap[S ~[]E, E, T any](in S, fn func(E) T) []T {
	out := make([]T, 0, len(in))
	for _, v := range in {
		out = append(out, fn(v))
	}
	return Clip(out)
}

// SliceMap takes a slice of type E, calls fn on each value of `in` and returns the modified in or a copy of it
func SliceMapSameType[S ~[]E, E any](in S, fn func(E) E, inplace bool) (out S) {
	if inplace {
		out = in[:0]
	} else {
		out = make(S, 0, len(out))
	}

	for _, v := range in {
		out = append(out, fn(v))
	}

	return Clip(out)
}

// SliceMapFilter merged SliceMap and Filter
func SliceMapFilter[S ~[]E, E, T any](in S, fn func(E) (val T, ignore bool)) []T {
	out := make([]T, 0, len(in))
	for _, v := range in {
		if nv, ignore := fn(v); !ignore {
			out = append(out, nv)
		}
	}
	return Clip(out)
}

// SliceMapFilter merged SliceMapSameType and Filter
func SliceMapFilterSameType[S ~[]E, E any](in S, fn func(E) (val E, ignore bool), inplace bool) (out S) {
	if inplace {
		out = in[:0]
	} else {
		out = make(S, 0, len(out))
	}

	for _, v := range in {
		if nv, ignore := fn(v); !ignore {
			out = append(out, nv)
		}
	}

	return Clip(out)
}
