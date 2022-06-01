package genh

func SliceToChan[T any](s []T, cap int) <-chan T {
	if cap == 0 {
		cap = 1
	}
	ch := make(chan T, cap)
	go func() {
		for _, v := range s {
			ch <- v
		}
		close(ch)
	}()
	return ch
}

func ChanToSlice[T any](s <-chan T) []T {
	out := make([]T, cap(s))
	for v := range s {
		out = append(out, v)
	}
	return Clip(out)
}
