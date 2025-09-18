package genh

import (
	"bytes"
	"encoding/json"
	"hash/maphash"
	"iter"
	"runtime"
	"sync"
)

func NewSLMap[V any](ln int) *SLMap[V] {
	if ln < 1 {
		ln = runtime.NumCPU()
	}
	var sm SLMap[V]
	sm.init(ln)
	return &sm
}

type SLMap[V any] struct {
	ms []*LMap[string, V]
	s  maphash.Seed
	o  sync.Once
}

func (lm *SLMap[V]) m(k string) *LMap[string, V] {
	lm.initOnce()
	return lm.ms[maphash.String(lm.s, k)%uint64(len(lm.ms))]
}

func (lm *SLMap[V]) initOnce() {
	lm.o.Do(func() {
		if len(lm.ms) == 0 {
			lm.init(runtime.NumCPU())
		}
	})
}

func (lm *SLMap[V]) init(sz int) {
	lm.ms = make([]*LMap[string, V], sz)
	for i := range lm.ms {
		lm.ms[i] = NewLMap[string, V](0)
	}
	lm.s = maphash.MakeSeed()
}

func (lm *SLMap[V]) Set(k string, v V) {
	lm.m(k).Set(k, v)
}

func (lm *SLMap[V]) UpdateKey(k string, fn func(V) V) {
	lm.m(k).UpdateKey(k, fn)
}

func (lm *SLMap[V]) Swap(k string, v V) V {
	return lm.m(k).Swap(k, v)
}

func (lm *SLMap[V]) Delete(k string) {
	lm.m(k).Delete(k)
}

func (lm *SLMap[V]) DeleteGet(k string) V {
	return lm.m(k).DeleteGet(k)
}

func (lm *SLMap[V]) Keys() (keys []string) {
	lm.initOnce()
	ln := 0
	for _, m := range lm.ms {
		ln += m.Len()
	}
	keys = make([]string, 0, ln)
	for _, m := range lm.ms {
		keys = append(keys, m.Keys()...)
	}
	return keys
}

func (lm *SLMap[V]) KeysSeq() iter.Seq[string] {
	lm.initOnce()
	return func(yield func(string) bool) {
		for _, m := range lm.ms {
			for key := range m.KeysSeq() {
				if !yield(key) {
					return
				}
			}
		}
	}
}

func (lm *SLMap[V]) Values() (values []V) {
	lm.initOnce()
	ln := 0
	for _, m := range lm.ms {
		ln += m.Len()
	}
	values = make([]V, 0, ln)
	for _, m := range lm.ms {
		values = append(values, m.Values()...)
	}
	return values
}

func (lm *SLMap[V]) Clone() (out map[string]V) {
	lm.initOnce()
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
	return out
}

func (lm *SLMap[V]) Update(fn func(m map[string]V)) {
	lm.initOnce()
	for _, m := range lm.ms {
		m.Update(fn)
	}
}

func (lm *SLMap[V]) Read(fn func(m map[string]V)) {
	lm.initOnce()
	for _, m := range lm.ms {
		m.Read(fn)
	}
}

func (lm *SLMap[V]) Get(k string) (v V) {
	return lm.m(k).Get(k)
}

func (lm *SLMap[V]) MustGet(k string, fn func() V) V {
	return lm.m(k).MustGet(k, fn)
}

func (lm *SLMap[V]) ForEach(fn func(k string, v V) bool) {
	lm.initOnce()
	for _, m := range lm.ms {
		m.ForEach(func(k string, v V) bool {
			return fn(k, v)
		})
	}
}

func (lm *SLMap[V]) SetMap(m map[string]V) {
	lm.initOnce()
	lm.Clear()
	for k, v := range m {
		lm.Set(k, v)
	}
}

func (lm *SLMap[V]) Clear() {
	lm.initOnce()
	for _, m := range lm.ms {
		m.Clear()
	}
}

func (lm *SLMap[V]) Len() (ln int) {
	lm.initOnce()
	for _, m := range lm.ms {
		ln += m.Len()
	}
	return ln
}

func (lm *SLMap[V]) MarshalJSON() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	first := true
	buf.WriteRune('{')
	lm.ForEach(func(k string, v V) bool {
		if !first {
			buf.WriteRune(',')
		} else {
			first = false
		}
		enc.Encode(k)
		buf.Truncate(buf.Len() - 1)
		buf.WriteRune(':')
		err = enc.Encode(v)
		buf.Truncate(buf.Len() - 1)
		return err == nil
	})
	if err != nil {
		return
	}
	buf.WriteRune('}')

	return buf.Bytes(), nil
}

func (lm *SLMap[V]) UnmarshalJSON(p []byte) (err error) {
	var m map[string]V
	if err = json.Unmarshal(p, &m); err != nil {
		return err
	}
	for k, v := range m {
		lm.Set(k, v)
	}
	return err
}

func (lm *SLMap[V]) MarshalBinary() (_ []byte, err error) {
	var buf bytes.Buffer
	enc := NewMsgpackEncoder(&buf)
	defer PutMsgpackEncoder(enc)
	ln := lm.Len()
	enc.EncodeMapLen(ln)
	lm.ForEach(func(k string, v V) bool {
		if ln--; ln < 0 {
			return false
		}
		enc.EncodeString(k)
		err = enc.Encode(v)
		return err == nil
	})
	if err != nil {
		return
	}
	return buf.Bytes(), nil
}

func (lm *SLMap[V]) UnmarshalBinary(p []byte) error {
	dec := NewMsgpackDecoder(bytes.NewReader(p))
	defer PutMsgpackDecoder(dec)
	ln, err := dec.DecodeMapLen()
	if err != nil {
		return err
	}
	for range ln {
		k, err := dec.DecodeString()
		if err != nil {
			return err
		}
		var v V
		if err = dec.Decode(&v); err != nil {
			return err
		}
		lm.Set(k, v)
	}
	return nil
}
