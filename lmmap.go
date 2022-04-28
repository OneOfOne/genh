package genh

import (
	"encoding/json"
	"sync"
)

func NewLMultiMap[K1, K2 comparable, V any](sz int) *LMultiMap[K1, K2, V] {
	return &LMultiMap[K1, K2, V]{m: make(map[K1]map[K2]V, sz)}
}

// LMultiMap is a locked multimap
type LMultiMap[K1, K2 comparable, V any] struct {
	mux sync.RWMutex
	m   map[K1]map[K2]V
}

func (lm *LMultiMap[K1, K2, V]) Set(k1 K1, k2 K2, v V) {
	lm.mux.Lock()
	slm := lm.m[k1]
	if slm == nil {
		if lm.m == nil {
			lm.m = make(map[K1]map[K2]V)
		}
		slm = make(map[K2]V)
		lm.m[k1] = slm
	}

	slm[k2] = v
	lm.mux.Unlock()
}

func (lm *LMultiMap[K1, K2, V]) SetChild(k1 K1, v map[K2]V) {
	lm.mux.Lock()
	lm.m[k1] = v
	lm.mux.Unlock()
}

func (lm *LMultiMap[K1, K2, V]) Update(k1 K1, fn func(m map[K2]V) map[K2]V) {
	lm.mux.Lock()
	defer lm.mux.Unlock()

	if m := fn(lm.m[k1]); m != nil {
		if lm.m == nil {
			lm.m = make(map[K1]map[K2]V)
		}
		lm.m[k1] = m
	} else {
		delete(lm.m, k1)
	}
}

func (lm *LMultiMap[K1, K2, V]) DeleteChild(k1 K1, k2 K2) {
	lm.mux.Lock()
	delete(lm.m[k1], k2)
	lm.mux.Unlock()
}

func (lm *LMultiMap[K1, K2, V]) Delete(k1 K1) {
	lm.mux.Lock()
	delete(lm.m, k1)
	lm.mux.Unlock()
}

func (lm *LMultiMap[K1, K2, V]) Keys() (keys []K1) {
	lm.mux.RLock()
	keys = MapKeys(lm.m)
	lm.mux.RUnlock()
	return
}

func (lm *LMultiMap[K1, K2, V]) KeysChild(k1 K1) (keys []K2) {
	lm.mux.RLock()
	keys = MapKeys(lm.m[k1])
	lm.mux.RUnlock()
	return
}

func (lm *LMultiMap[K1, K2, V]) Values(copy bool) (values []map[K2]V) {
	lm.mux.RLock()
	if values = MapValues(lm.m); copy {
		for i, v := range values {
			values[i] = MapClone(v)
		}
	}
	lm.mux.RUnlock()
	return
}

func (lm *LMultiMap[K1, K2, V]) ValuesChild(k1 K1) (values []V) {
	lm.mux.RLock()
	values = MapValues(lm.m[k1])
	lm.mux.RUnlock()
	return
}

func (lm *LMultiMap[K1, K2, V]) Clone() (m map[K1]map[K2]V) {
	lm.mux.RLock()
	m = make(map[K1]map[K2]V, len(lm.m))
	for k, v := range lm.m {
		m[k] = MapClone(v)
	}
	lm.mux.RUnlock()
	return
}

func (lm *LMultiMap[K1, K2, V]) Get(k1 K1, k2 K2) (v V) {
	lm.mux.RLock()
	v = lm.m[k1][k2]
	lm.mux.RUnlock()
	return
}

func (lm *LMultiMap[K1, K2, V]) MustGet(k1 K1, k2 K2, fn func() V) V {
	lm.mux.RLock()
	v, ok := lm.m[k1][k2]
	lm.mux.RUnlock()
	if ok {
		return v
	}

	var nv V
	if fn != nil {
		// create outside lock in case it's heavy, there's a chance it won't be used
		nv = fn()
	}

	lm.mux.Lock()
	defer lm.mux.Unlock()

	if v, ok = lm.m[k1][k2]; ok { // race check
		return v
	}

	if lm.m == nil {
		lm.m = make(map[K1]map[K2]V)
	}

	m := lm.m[k1]
	if m == nil {
		m = make(map[K2]V)
		lm.m[k1] = m
	}

	m[k2] = nv
	return nv
}

func (lm *LMultiMap[K1, K2, V]) GetChild(k1 K1, copy bool) map[K2]V {
	lm.mux.RLock()
	defer lm.mux.RUnlock()
	if !copy {
		return lm.m[k1]
	}

	return MapClone(lm.m[k1])
}

func (lm *LMultiMap[K1, K2, V]) ForEach(fn func(k1 K1, m map[K2]V) bool, rw bool) {
	if rw {
		lm.mux.Lock()
		defer lm.mux.Unlock()
	} else {
		lm.mux.RLock()
		defer lm.mux.RUnlock()
	}

	for k1, m := range lm.m {
		if !fn(k1, m) {
			return
		}
	}
}

func (lm *LMultiMap[K1, K2, V]) ForEachChild(k1 K1, fn func(k2 K2, v V) bool) {
	lm.mux.RLock()
	defer lm.mux.RUnlock()
	for k2, v := range lm.m[k1] {
		if !fn(k2, v) {
			return
		}
	}
}

func (lm *LMultiMap[K1, K2, V]) Clear() {
	lm.mux.Lock()
	MapClear(lm.m)
	lm.mux.Unlock()
}

func (lm *LMultiMap[K1, K2, V]) ClearChild(k1 K1) {
	lm.mux.Lock()
	MapClear(lm.m[k1])
	lm.mux.Unlock()
}

func (lm *LMultiMap[K1, K2, V]) SetMap(m map[K1]map[K2]V) (old map[K1]map[K2]V) {
	lm.mux.Lock()
	old = lm.m
	lm.m = m
	lm.mux.Unlock()
	return
}

func (lm *LMultiMap[K1, K2, V]) Len() (v int) {
	lm.mux.RLock()
	v = len(lm.m)
	lm.mux.RUnlock()
	return
}

func (lm *LMultiMap[K1, K2, V]) LenChild(k K1) (v int) {
	lm.mux.RLock()
	v = len(lm.m[k])
	lm.mux.RUnlock()
	return
}

func (lm *LMultiMap[K1, K2, V]) Raw() map[K1]map[K2]V {
	lm.mux.RLock()
	defer lm.mux.RUnlock()
	return lm.m
}

func (lm *LMultiMap[K1, K2, V]) MarshalJSON() ([]byte, error) {
	lm.mux.RLock()
	defer lm.mux.RUnlock()
	return json.Marshal(lm.m)
}

func (lm *LMultiMap[K1, K2, V]) UnmarshalJSON(p []byte) error {
	lm.mux.Lock()
	defer lm.mux.Unlock()
	return json.Unmarshal(p, &lm.m)
}
