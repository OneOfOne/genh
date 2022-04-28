package genh

import (
	"encoding/json"
	"sync"
)

type LSlice[T any] struct {
	mux sync.RWMutex
	v   []T
}

func (s *LSlice[T]) Update(fn func(v []T) []T) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.v = fn(s.v)
}

func (s *LSlice[T]) Append(vs ...T) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.v = append(s.v, vs...)
}

func (s *LSlice[T]) Set(i int, v T) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.v[i] = v
}

func (s *LSlice[T]) Insert(i int, vs ...T) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.v = Insert(s.v, i, vs...)
}

func (s *LSlice[T]) Delete(i, j int) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.v = Delete(s.v, i, j)
}

func (s *LSlice[T]) Filter(fn func(T) bool, inplace bool) *LSlice[T] {
	if inplace {
		s.mux.Lock()
		defer s.mux.Unlock()
		s.v = Filter(s.v, fn, true)
		return s
	}
	s.mux.RLock()
	defer s.mux.RLock()
	return &LSlice[T]{v: Filter(s.v, fn, false)}
}

func (s *LSlice[T]) Map(fn func(T) T, inplace bool) *LSlice[T] {
	if inplace {
		s.mux.Lock()
		defer s.mux.Unlock()
		s.v = SliceMapSameType(s.v, fn, true)
		return s
	}
	s.mux.RLock()
	defer s.mux.RLock()
	return &LSlice[T]{v: SliceMapSameType(s.v, fn, false)}
}

func (s *LSlice[T]) Swap(i int, v T) (old T) {
	s.mux.Lock()
	defer s.mux.Unlock()
	old, s.v[i] = s.v[i], v
	return
}

func (s *LSlice[T]) SetSlice(v []T) {
	s.mux.Lock()
	s.v = v
	s.mux.Unlock()
}

func (s *LSlice[T]) Sort(lessFn func(a, b T) bool) {
	s.mux.Lock()
	defer s.mux.Unlock()
	SortFunc(s.v, lessFn)
}

func (s *LSlice[T]) Grow(sz int) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.v = Grow(s.v, sz)
}

func (s *LSlice[T]) ClipTo(len_, cap_ int) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.v = s.v[:len_:cap_]
}

func (s *LSlice[T]) Clip() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.v = Clip(s.v)
}

func (s *LSlice[T]) Len() int {
	s.mux.RLock()
	defer s.mux.RLock()
	return len(s.v)
}

func (s *LSlice[T]) Cap() int {
	s.mux.RLock()
	defer s.mux.RLock()
	return cap(s.v)
}

func (s *LSlice[T]) Get(i int) T {
	s.mux.RLock()
	defer s.mux.RLock()
	return s.v[i]
}

func (s *LSlice[T]) ForEach(fn func(i int, v T) bool) {
	s.mux.RLock()
	defer s.mux.RLock()
	for i, v := range s.v {
		if !fn(i, v) {
			return
		}
	}
}

func (s *LSlice[T]) Search(cmpFn func(v T) int) (v T, found bool) {
	var i int
	s.mux.RLock()
	defer s.mux.RLock()
	if i, found = BinarySearchFunc(s.v, cmpFn); found {
		v = s.v[i]
	}
	return
}

func (s *LSlice[T]) Clone() []T {
	s.mux.RLock()
	defer s.mux.RLock()
	return Clone(s.v)
}

func (s *LSlice[T]) LClone() *LSlice[T] {
	s.mux.RLock()
	defer s.mux.RLock()
	return &LSlice[T]{v: Clone(s.v)}
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
