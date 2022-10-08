package genh

import "sync"

type Once[T any] struct {
	once sync.Once
	v    T
	err  error
}

func (o *Once[T]) Do(fn func() (T, error)) (v T, err error) {
	o.once.Do(func() { o.v, o.err = fn() })
	return o.v, o.err
}
