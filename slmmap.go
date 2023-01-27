package genh

import (
	"hash/maphash"
	"runtime"
	"sync"
)

type SLMultiMap[V any] struct {
	ms []*LMultiMap[string, string, V]
	s  maphash.Seed
	o  sync.Once
}

func (lm *SLMultiMap[V]) m(k string) *LMultiMap[string, string, V] {
	lm.o.Do(func() {
		if len(lm.ms) == 0 {
			lm.init(runtime.NumCPU())
		}
	})
	return lm.ms[maphash.String(lm.s, k)%uint64(len(lm.ms))]
}

func (lm *SLMultiMap[V]) init(sz int) {
	lm.ms = make([]*LMultiMap[string, string, V], sz)
	for i := range lm.ms {
		lm.ms[i] = NewLMultiMap[string, string, V](0)
	}
	lm.s = maphash.MakeSeed()
}

func (lm *SLMultiMap[V]) Set(k1, k2 string, v V) {
	lm.m(k1).Set(k1, k2, v)
}

func (lm *SLMultiMap[V]) Get(k1, k2 string) V {
	return lm.m(k1).Get(k1, k2)
}

func (lm *SLMultiMap[V]) MustGet(k1, k2 string, fn func() V) V {
	return lm.m(k1).MustGet(k1, k2, fn)
}
