package genh

import "sync"

func NewLMap[K comparable, V any](sz int) *LMap[K, V] {
	return &LMap[K, V]{m: make(map[K]V, sz)}
}

func LMapOf[K comparable, V any](m map[K]V) *LMap[K, V] {
	return &LMap[K, V]{m: MapClone(m)}
}

type LMap[K comparable, V any] struct {
	mux sync.RWMutex
	m   map[K]V
}

func (lm *LMap[K, V]) Set(k K, v V) {
	lm.mux.Lock()
	if lm.m == nil {
		lm.m = make(map[K]V)
	}
	lm.m[k] = v
	lm.mux.Unlock()
}

func (lm *LMap[K, V]) Swap(k K, v V) V {
	lm.mux.Lock()
	ov := lm.m[k]
	if lm.m == nil {
		lm.m = make(map[K]V)
	}
	lm.m[k] = v
	lm.mux.Unlock()
	return ov
}

func (lm *LMap[K, V]) Delete(k K) {
	lm.mux.Lock()
	delete(lm.m, k)
	lm.mux.Unlock()
}

func (lm *LMap[K1, V]) Keys() (keys []K1) {
	lm.mux.RLock()
	keys = MapKeys(lm.m)
	lm.mux.RUnlock()
	return
}

func (lm *LMap[K1, V]) Values() (values []V) {
	lm.mux.RLock()
	values = MapValues(lm.m)
	lm.mux.RUnlock()
	return
}

func (lm *LMap[K, V]) Get(k K) (v V) {
	lm.mux.RLock()
	v = lm.m[k]
	lm.mux.RUnlock()
	return
}

func (lm *LMap[K, V]) ForEach(fn func(k K, v V) bool) {
	lm.mux.RLock()
	defer lm.mux.RUnlock()
	for k, v := range lm.m {
		if !fn(k, v) {
			return
		}
	}
}

func (lm *LMap[K, V]) Clear() {
	lm.mux.Lock()
	for k := range lm.m {
		delete(lm.m, k)
	}
	lm.mux.Unlock()
}
