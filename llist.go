package genh

import "sync"

type LList[T any] struct {
	l   List[T]
	mux sync.RWMutex
}

func (l *LList[T]) Append(vs ...T) *LList[T] {
	l.Push(vs...)
	return l
}

func (l *LList[T]) Push(vs ...T) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.l.Push(vs...)
}

func (l *LList[T]) Len() int {
	l.mux.RLock()
	defer l.mux.RUnlock()
	return l.l.Len()
}

func (l *LList[T]) ForEach(fn func(v T) bool) {
	l.mux.RLock()
	defer l.mux.RUnlock()
	l.l.ForEach(fn)
}

func (l *LList[T]) Clear() {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.l.Clear()
}

func (l *LList[T]) Raw() List[T] {
	l.mux.RLock()
	defer l.mux.RUnlock()
	return l.l.Clip()
}

// TODO: rest of the list interface
