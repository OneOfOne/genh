package genh

import (
	"sync"
	"sync/atomic"
	"time"
)

type tmEle[V any] struct {
	v   V
	t   *time.Timer
	la  atomic.Int64
	ttl time.Duration
	sync.RWMutex
}

func (e *tmEle[V]) expired() bool {
	if e.ttl < 1 {
		return false
	}
	return time.Since(time.Unix(0, e.la.Load())) > e.ttl
}

type TimedMap[K comparable, V any] struct {
	m LMap[K, *tmEle[V]]
}

func (tm *TimedMap[K, V]) Set(k K, v V, timeout time.Duration) {
	ele := &tmEle[V]{v: v, ttl: timeout}
	ele.la.Store(time.Now().UnixNano())
	if timeout > 0 {
		ele.t = time.AfterFunc(timeout, func() { tm.deleteEle(k, ele) })
	}
	tm.m.Set(k, ele)
}

func (tm *TimedMap[K, V]) MustGet(k K, vfn func() V, timeout time.Duration) (out V) {
	var ok bool
	if out, ok = tm.GetOk(k); ok {
		return
	}
	out = vfn()
	ele := &tmEle[V]{v: out, ttl: timeout}
	ele.la.Store(time.Now().UnixNano())
	if timeout > 0 {
		ele.t = time.AfterFunc(timeout, func() { tm.deleteEle(k, ele) })
	}
	tm.m.Set(k, ele)
	return
}

func (tm *TimedMap[K, V]) SetUpdateFn(k K, vfn func() V, updateEvery time.Duration) {
	tm.SetUpdateExpireFn(k, vfn, updateEvery, -1)
}

func (tm *TimedMap[K, V]) SetUpdateExpireFn(k K, vfn func() V, updateEvery, expireIfNotAccessedFor time.Duration) {
	if updateEvery < time.Millisecond {
		panic("every must be >= time.Millisecond")
	}

	ele := &tmEle[V]{ttl: expireIfNotAccessedFor}
	ele.la.Store(time.Now().UnixNano())
	var upfn func()
	upfn = func() {
		if ele.expired() {
			tm.deleteEle(k, ele)
			return
		}
		v := vfn()
		ele.Lock()
		if ele.t != nil {
			ele.t.Stop()
		}
		ele.v = v
		ele.t = time.AfterFunc(updateEvery, upfn)
		ele.Unlock()
	}
	upfn()
	tm.m.Set(k, ele)
}

func (tm *TimedMap[K, V]) Get(k K) (v V) {
	v, _ = tm.GetOk(k)
	return
}

func (tm *TimedMap[K, V]) GetOk(k K) (v V, ok bool) {
	ele := tm.m.Get(k)
	if ok = ele != nil && !ele.expired(); ok {
		now := time.Now().UnixNano()
		ele.RLock()
		defer ele.RUnlock()
		v = ele.v
		ele.la.Store(now)
	}
	return
}

func (tm *TimedMap[K, V]) DeleteGet(k K) (v V, ok bool) {
	ele := tm.m.DeleteGet(k)
	if ok = ele != nil && !ele.expired(); ok {
		ele.RLock()
		defer ele.RUnlock()
		v = ele.v
		if ele.t != nil {
			ele.t.Stop()
		}
	}
	return
}

func (tm *TimedMap[K, V]) Delete(k K) {
	if ele := tm.m.DeleteGet(k); ele != nil {
		if ele.t != nil {
			ele.t.Stop()
		}
	}
}

func (tm *TimedMap[K, V]) ForEach(fn func(key K, value V) bool) {
	keys := tm.m.Keys()

	for _, k := range keys {
		if v, ok := tm.GetOk(k); ok {
			if !fn(k, v) {
				return
			}
		}
	}
}
func (tm *TimedMap[K, V]) deleteEle(k K, ele *tmEle[V]) {
	tm.m.Update(func(m map[K]*tmEle[V]) {
		if m[k] == ele {
			delete(m, k)
		}
	})
}
