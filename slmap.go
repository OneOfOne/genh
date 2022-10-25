package genh

import (
	"hash/maphash"
	"runtime"
	"sync"
)

func NewSLMap[V any](ln int) (sm SLMap[V]) {
	if ln < 1 {
		ln = runtime.NumCPU()
	}
	sm.ms = make([]*LMap[string, V], ln)
	for i := range sm.ms {
		sm.ms[i] = NewLMap[string, V](0)
	}
	sm.s = maphash.MakeSeed()
	return
}

type SLMap[V any] struct {
	s    maphash.Seed
	ms   []*LMap[string, V]
	init sync.Once
}

func (lm SLMap[V]) m(k string) *LMap[string, V] {
	lm.init.Do(func() {
		if len(lm.ms) == 0 {
			lm.ms = make([]*LMap[string, V], runtime.NumCPU())
			for i := range lm.ms {
				lm.ms[i] = NewLMap[string, V](0)
			}
			lm.s = maphash.MakeSeed()
		}
	})
	return lm.ms[maphash.String(lm.s, k)%uint64(len(lm.ms))]
}

func (lm SLMap[V]) Set(k string, v V) {
	lm.m(k).Set(k, v)
}

func (lm SLMap[V]) UpdateKey(k string, fn func(V) V) {
	lm.m(k).UpdateKey(k, fn)
}

func (lm SLMap[V]) Swap(k string, v V) V {
	return lm.m(k).Swap(k, v)
}

func (lm SLMap[V]) Delete(k string) {
	lm.m(k).Delete(k)
}

func (lm SLMap[V]) DeleteGet(k string) V {
	return lm.m(k).DeleteGet(k)
}

func (lm SLMap[V]) Keys() (keys []string) {
	ln := 0
	for _, m := range lm.ms {
		ln += m.Len()
	}
	keys = make([]string, 0, ln)
	for _, m := range lm.ms {
		keys = append(keys, m.Keys()...)
	}
	return
}

func (lm SLMap[V]) Values() (values []V) {
	ln := 0
	for _, m := range lm.ms {
		ln += m.Len()
	}
	values = make([]V, 0, ln)
	for _, m := range lm.ms {
		values = append(values, m.Values()...)
	}
	return
}

func (lm SLMap[V]) Clone() (out map[string]V) {
	ln := 0
	for _, m := range lm.ms {
		ln += m.Len()
	}
	out = make(map[string]V, ln)
	for _, m := range lm.ms {
		m.ForEach(func(k string, v V) bool {
			out[k] = v
			return true
		})
	}
	return
}

func (lm SLMap[V]) Update(fn func(m map[string]V)) {
	for _, m := range lm.ms {
		m.Update(fn)
	}
}

func (lm SLMap[V]) Read(fn func(m map[string]V)) {
	for _, m := range lm.ms {
		m.Read(fn)
	}
}

func (lm SLMap[V]) Get(k string) (v V) {
	return lm.m(k).Get(k)
}

func (lm SLMap[V]) MustGet(k string, fn func() V) V {
	return lm.m(k).MustGet(k, fn)
}

func (lm SLMap[V]) ForEach(fn func(k string, v V) bool) {
	for _, m := range lm.ms {
		m.ForEach(func(k string, v V) bool {
			return fn(k, v)
		})
	}
}

func (lm SLMap[V]) Clear() {
	for _, m := range lm.ms {
		m.Clear()
	}
}

func (lm SLMap[V]) Len() (ln int) {
	for _, m := range lm.ms {
		ln += m.Len()
	}
	return
}
