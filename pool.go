package genh

import "sync"

type BytesPool struct {
	m LMap[uint32, *sync.Pool]
}

func (mp *BytesPool) Get(sz uint32) []byte {
	if sz = sz - (sz % 1024); sz < 1024 { // round to the nearest kb
		sz = 1024
	}
	return *(mp.pool(sz).Get().(*[]byte))
}

func (mp *BytesPool) Put(b []byte) uint32 {
	if cap(b) < 1024 {
		return 0
	}
	sz := uint32(cap(b))
	if sz = sz - (sz % 1024); sz < 1024 { // round to the nearest kb
		sz = 1024
	}
	b = b[:0:sz]
	mp.pool(sz).Put(&b)
	return sz
}

func (mp *BytesPool) pool(sz uint32) *sync.Pool {
	return mp.m.MustGet(sz, func() *sync.Pool {
		p := sync.Pool{
			New: func() any {
				b := make([]byte, 0, sz)
				return &b
			},
		}
		return &p
	})

}
