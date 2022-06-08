package genh

import (
	"encoding/json"
	"sync"
)

// LValue wraps a sync.RWMutex to allow simple and safe operation on the mutex.
type LValue[T any] struct {
	mux sync.RWMutex
	v   T
}

// Update executes fn while the mutex is write-locked and guarantees the mutex is released even in the case of a panic.
func (m *LValue[T]) Update(fn func(old T) T) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.v = fn(m.v)
}

// Read executes fn while the mutex is read-locked and guarantees the mutex is released even in the case of a panic.
func (m *LValue[T]) Read(fn func(v T)) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	fn(m.v)
}

func (m *LValue[T]) Get() T {
	m.mux.RLock()
	v := m.v
	m.mux.RUnlock()
	return v
}

func (m *LValue[T]) Set(v T) {
	m.mux.Lock()
	m.v = v
	m.mux.Unlock()
}

func (m *LValue[T]) Swap(v T) (old T) {
	m.mux.Lock()
	old, m.v = m.v, v
	m.mux.Unlock()
	return
}

func (m *LValue[T]) CompareAndSwap(old, new T, eq func(a, b T) bool) (ok bool) {
	m.mux.Lock()
	defer m.mux.Unlock()
	if ok = eq(m.v, old); ok {
		m.v = new
	}
	return
}

func (m *LValue[T]) MarshalBinary() ([]byte, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return MarshalMsgpack(m.v)
}

func (m *LValue[T]) UnmarshalBinary(b []byte) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	return UnmarshalMsgpack(b, &m.v)
}

func (m *LValue[T]) MarshalJSON() ([]byte, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return json.Marshal(m.v)
}

func (m *LValue[T]) UnmarshalJSON(b []byte) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	return json.Unmarshal(b, &m.v)
}
