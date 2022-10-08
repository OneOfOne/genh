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

type Pool[T any] struct {
	New   func() *T
	Reset func(*T)

	p sync.Pool
}

func (p *Pool[T]) Get() *T {
	v, ok := p.p.Get().(*T)
	if !ok {
		if p.New != nil {
			v = p.New()
		} else {
			v = new(T)
		}
	}
	return v
}

func (p *Pool[T]) Put(v *T) {
	if p.Reset != nil {
		p.Reset(v)
	}
	p.p.Put(v)
}
