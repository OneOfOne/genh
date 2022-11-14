package gsets

import (
	"encoding/json"
	"sort"
	"sync"

	"go.oneofone.dev/genh/internal"
)

func SafeOf[T internal.Ordered](keys ...T) *SafeSet[T] {
	s := Of(keys...)
	return &SafeSet[T]{s: s}
}

type SafeSet[T internal.Ordered] struct {
	mux sync.RWMutex
	s   Set[T]
}

func (ss *SafeSet[T]) Set(keys ...T) *SafeSet[T] {
	return ss.Add(keys...)
}

func (ss *SafeSet[T]) Add(keys ...T) *SafeSet[T] {
	ss.mux.Lock()
	ss.s.Set(keys...)
	ss.mux.Unlock()
	return ss
}

func (ss *SafeSet[T]) AddIfNotExists(key T) bool {
	ss.mux.Lock()
	added := ss.s.AddIfNotExists(key)
	ss.mux.Unlock()
	return added
}

func (ss *SafeSet[T]) Clone() *SafeSet[T] {
	ss.mux.RLock()
	ns := ss.s.Clone()
	ss.mux.RUnlock()
	return &SafeSet[T]{s: ns}
}

func (ss *SafeSet[T]) MergeSafe(o *SafeSet[T]) *SafeSet[T] {
	ss.mux.Lock()
	o.mux.Lock()
	ss.s.Merge(o.s)
	o.mux.Unlock()
	ss.mux.Unlock()
	return ss
}

func (ss *SafeSet[T]) Merge(o Set[T]) *SafeSet[T] {
	ss.mux.Lock()
	ss.s.Merge(o)
	ss.mux.Unlock()
	return ss
}

func (ss *SafeSet[T]) Delete(keys ...T) *SafeSet[T] {
	ss.mux.Lock()
	ss.s.Delete(keys...)
	ss.mux.Unlock()
	return ss
}

func (ss *SafeSet[T]) Has(key T) bool {
	ss.mux.RLock()
	ok := ss.s.Has(key)
	ss.mux.RUnlock()
	return ok
}

func (ss *SafeSet[T]) Len() int {
	ss.mux.RLock()
	ln := ss.s.Len()
	ss.mux.RUnlock()
	return ln
}

func (ss *SafeSet[T]) Keys() []T {
	ss.mux.RLock()
	keys := ss.s.Keys()
	ss.mux.RUnlock()
	return keys
}

func (ss *SafeSet[T]) SortedKeys() []T {
	keys := ss.Keys()
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

func (ss *SafeSet[T]) MarshalJSON() ([]byte, error) {
	keys := ss.Keys()
	return json.Marshal(keys)
}

func (ss *SafeSet[T]) UnmarshalJSON(data []byte) error {
	var s Set[T]
	if err := s.UnmarshalJSON(data); err != nil {
		return err
	}
	ss.mux.Lock()
	ss.s = s
	ss.mux.Unlock()
	return nil
}

func (ss *SafeSet[T]) MarshalBinary() ([]byte, error) {
	keys := ss.Keys()
	return internal.MarshalMsgpack(keys)
}

func (ss *SafeSet[T]) UnmarshalBinary(data []byte) error {
	var s Set[T]
	if err := s.UnmarshalBinary(data); err != nil {
		return err
	}
	ss.mux.Lock()
	ss.s = s
	ss.mux.Unlock()
	return nil
}
