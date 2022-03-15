package genh

import "sync"

// Locked wraps a sync.RWMutex to allow simple and safe operation on the mutex.
type Locked[T any] struct {
	sync.RWMutex
	v *T
}

// Update executes fn while the mutex is write-locked and guarantees the mutex is released even in the case of a panic.
func (m *Locked[T]) Update(fn func(old *T)) {
	m.Lock()
	defer m.Unlock()
	fn(m.v)
}

// Read executes fn while the mutex is read-locked and guarantees the mutex is released even in the case of a panic.
func (m *Locked[T]) Read(fn func(v *T)) {
	m.RLock()
	defer m.RUnlock()
	fn(m.v)
}

func (m *Locked[T]) Get() *T {
	m.RLock()
	v := m.v
	m.RUnlock()
	return v
}

func (m *Locked[T]) Set(v *T) {
	m.Lock()
	m.v = v
	m.Unlock()
}

func (m *Locked[T]) Swap(v *T) (old *T) {
	m.Lock()
	old, m.v = m.v, v
	m.Unlock()
	return
}
