package genh

import (
	"encoding/json"
	"sync"
)

type LSlice[T any] struct {
	mux sync.RWMutex
	v   []T
}

func (ls *LSlice[T]) Update(fn func(v []T) []T) {
	ls.mux.Lock()
	defer ls.mux.Unlock()
	ls.v = fn(ls.v)
}

func (ls *LSlice[T]) Append(vs ...T) {
	ls.mux.Lock()
	defer ls.mux.Unlock()
	ls.v = append(ls.v, vs...)
}

func (ls *LSlice[T]) Set(i int, v T) {
	ls.mux.Lock()
	defer ls.mux.Unlock()
	ls.v[i] = v
}

func (ls *LSlice[T]) Insert(i int, vs ...T) {
	ls.mux.Lock()
	defer ls.mux.Unlock()
	ls.v = Insert(ls.v, i, vs...)
}

func (ls *LSlice[T]) Delete(i, j int) {
	ls.mux.Lock()
	defer ls.mux.Unlock()
	ls.v = Delete(ls.v, i, j)
}

func (ls *LSlice[T]) Filter(fn func(T) bool, inplace bool) *LSlice[T] {
	if inplace {
		ls.mux.Lock()
		defer ls.mux.Unlock()
		ls.v = Filter(ls.v, fn, true)
		return ls
	}
	ls.mux.RLock()
	defer ls.mux.RUnlock()
	return &LSlice[T]{v: Filter(ls.v, fn, false)}
}

func (ls *LSlice[T]) Map(fn func(T) T, inplace bool) *LSlice[T] {
	if inplace {
		ls.mux.Lock()
		defer ls.mux.Unlock()
		ls.v = SliceMapSameType(ls.v, fn, true)
		return ls
	}
	ls.mux.RLock()
	defer ls.mux.RUnlock()
	return &LSlice[T]{v: SliceMapSameType(ls.v, fn, false)}
}

func (ls *LSlice[T]) Swap(i int, v T) (old T) {
	ls.mux.Lock()
	defer ls.mux.Unlock()
	old, ls.v[i] = ls.v[i], v
	return
}

func (ls *LSlice[T]) SetSlice(v []T) {
	ls.mux.Lock()
	ls.v = v
	ls.mux.Unlock()
}

func (ls *LSlice[T]) Sort(lessFn func(a, b T) bool) {
	ls.mux.Lock()
	defer ls.mux.Unlock()
	SortFunc(ls.v, lessFn)
}

func (ls *LSlice[T]) Grow(sz int) {
	ls.mux.Lock()
	defer ls.mux.Unlock()
	ls.v = Grow(ls.v, sz)
}

func (ls *LSlice[T]) ClipTo(len_, cap_ int) {
	ls.mux.Lock()
	defer ls.mux.Unlock()
	ls.v = ls.v[:len_:cap_]
}

func (ls *LSlice[T]) Clip() {
	ls.mux.Lock()
	defer ls.mux.Unlock()
	ls.v = Clip(ls.v)
}

func (ls *LSlice[T]) Len() int {
	ls.mux.RLock()
	defer ls.mux.RUnlock()
	return len(ls.v)
}

func (ls *LSlice[T]) Cap() int {
	ls.mux.RLock()
	defer ls.mux.RUnlock()
	return cap(ls.v)
}

func (ls *LSlice[T]) Get(i int) T {
	ls.mux.RLock()
	defer ls.mux.RUnlock()
	return ls.v[i]
}

func (ls *LSlice[T]) ForEach(fn func(i int, v T) bool) {
	ls.mux.RLock()
	defer ls.mux.RUnlock()
	for i, v := range ls.v {
		if !fn(i, v) {
			return
		}
	}
}

func (ls *LSlice[T]) Search(cmpFn func(v T) int) (v T, found bool) {
	var i int
	ls.mux.RLock()
	defer ls.mux.RUnlock()
	if i, found = BinarySearchFunc(ls.v, cmpFn); found {
		v = ls.v[i]
	}
	return
}

func (ls *LSlice[T]) Clone() []T {
	ls.mux.RLock()
	defer ls.mux.RUnlock()
	return SliceClone(ls.v)
}

func (ls *LSlice[T]) LClone() *LSlice[T] {
	ls.mux.RLock()
	defer ls.mux.RUnlock()
	return &LSlice[T]{v: SliceClone(ls.v)}
}

func (ls *LSlice[T]) Raw() []T {
	ls.mux.RLock()
	defer ls.mux.RUnlock()
	return ls.v
}

func (lm *LSlice[T]) MarshalJSON() ([]byte, error) {
	lm.mux.RLock()
	defer lm.mux.RUnlock()
	return json.Marshal(lm.v)
}

func (lm *LSlice[T]) UnmarshalJSON(p []byte) error {
	lm.mux.Lock()
	defer lm.mux.Unlock()
	return json.Unmarshal(p, &lm.v)
}
