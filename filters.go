package genh

// Filter filters a slice optionally in place.
func Filter[E comparable](in []E, f func(E) (keep bool), inplace bool) (out []E) {
	if inplace {
		out = in[:0]
	} else {
		out = make([]E, 0, len(out))
	}

	for _, s := range in {
		if f(s) {
			out = append(out, s)
		}
	}

	return Clip(out)
}
