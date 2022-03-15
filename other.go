package genh

func DecoderOf[T DecoderType](dec T) Decoder[T] {
	return Decoder[T]{d: dec}
}

// Decoder wraps a DecoderType to have a convient Decode() (v T, err error) func
type Decoder[T DecoderType] struct {
	d T
}

func (d *Decoder[T]) Decode() (v T, err error) {
	err = d.d.Decode(&v)
	return
}

// ValuesToPtrs converts a slice of values to a slice of pointers
// optionally copying the values instead of pointing to them in the original slice.
func ValuesToPtrs[T any](vals []T, copy bool) []*T {
	out := make([]*T, 0, len(vals))
	for i := range vals {
		v := &vals[i]
		if copy {
			cp := *v
			v = &cp
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
