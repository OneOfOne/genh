package genh

import (
	"bytes"
	"encoding"
	"encoding/json"

	"github.com/vmihailenco/msgpack/v5"
)

var (
	_ json.Marshaler             = (*List[any])(nil)
	_ json.Unmarshaler           = (*List[any])(nil)
	_ encoding.BinaryMarshaler   = (*List[any])(nil)
	_ encoding.BinaryUnmarshaler = (*List[any])(nil)
	_ msgpack.CustomEncoder      = (*List[any])(nil)
	_ msgpack.CustomDecoder      = (*List[any])(nil)
)

type (
	StringList  = List[string]
	Int64List   = List[int64]
	Uint64List  = List[uint64]
	Float64List = List[float64]
)

type listNode[T any] struct {
	v    T
	next *listNode[T]
}

func ListOf[T any](vs ...T) (l List[T]) {
	return l.Append(vs...)
}

type List[T any] struct {
	head, tail *listNode[T]
	len        int
}

func (l List[T]) Len() int { return l.len }

func (l List[T]) get(idx int) *listNode[T] {
	if idx > l.len-1 || idx < 0 {
		panic("index out of range")
	}
	n := l.head
	for i := 0; i < idx; i++ {
		n = n.next
	}
	return n
}

func (l *List[T]) Set(idx int, v T) {
	n := l.get(idx)
	n.v = v
}

func (l List[T]) Get(idx int) T {
	return l.get(idx).v
}

func (l List[T]) GetPtr(idx int) *T {
	n := l.get(idx)
	return &n.v
}

func (l *List[T]) Push(v T) {
	l.pushNode(&listNode[T]{v: v})
}

func (l List[T]) Append(vs ...T) List[T] {
	nodes := make([]listNode[T], len(vs))
	for i, v := range vs {
		n := &nodes[i]
		n.v = v
		l.pushNode(n)
	}
	return l
}

func (l *List[T]) pushNode(n *listNode[T]) {
	l.len++
	if l.head == nil {
		l.head, l.tail = n, n
		return
	}

	l.tail.next = n
	l.tail = n
}

func (l *List[T]) Prepend(v T) {
	n := &listNode[T]{v: v}
	l.len++
	if l.head == nil {
		l.head, l.tail = n, n
		return
	}

	n.next = l.head
	l.head = n
}

func (l List[T]) Iter() func() (v T, ok bool) {
	n := l.head

	return func() (v T, ok bool) {
		if ok = n != nil; ok {
			v, n = n.v, n.next
		}
		return
	}
}

func (l List[T]) ForEachPtr(fn func(v *T) bool) {
	for n := l.head; n != nil; n = n.next {
		if !fn(&n.v) {
			break
		}
	}
}

func (l List[T]) IterChan(cap int) <-chan T {
	if cap == 0 {
		cap = 1
	}
	ch := make(chan T, cap)
	go func() {
		defer close(ch)
		for n := l.head; n != nil; n = n.next {
			ch <- n.v
		}
	}()
	return ch
}

func (l List[T]) ForEach(fn func(v T) bool) {
	for n := l.head; n != nil; n = n.next {
		if !fn(n.v) {
			break
		}
	}
}

func (l List[T]) IterPtr() func() (v *T, ok bool) {
	n := l.head
	return func() (v *T, ok bool) {
		if ok = n != nil; ok {
			v, n = &n.v, n.next
		}
		return
	}
}

func (l List[T]) Slice() (out []T) {
	if l.head == nil {
		return
	}

	out = make([]T, 0, l.len)
	for n := l.head; n != nil; n = n.next {
		out = append(out, n.v)
	}
	return
}

func (l List[T]) SlicePtr() (out []*T) {
	if l.head == nil {
		return
	}
	out = make([]*T, 0, l.len)
	for n := l.head; n != nil; n = n.next {
		out = append(out, &n.v)
	}
	return
}

func (l *List[T]) Clear() {
	l.head, l.tail = nil, nil
	l.len = 0
}

func (l List[T]) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('[')
	enc := json.NewEncoder(&buf)
	for n := l.head; n != nil; n = n.next {
		if buf.Len() > 1 {
			buf.WriteString(",")
		}
		if err := enc.Encode(n.v); err != nil {
			return nil, err
		}
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

func (l *List[T]) UnmarshalJSON(p []byte) error {
	var v []T
	if err := json.Unmarshal(p, &v); err != nil {
		return err
	}
	*l = l.Append(v...)
	return nil
}

func (l List[T]) MarshalBinary() ([]byte, error) {
	return MarshalMsgpack(&l)
}

func (l *List[T]) UnmarshalBinary(p []byte) error {
	return UnmarshalMsgpack(p, &l)
}

func (l List[T]) EncodeMsgpack(enc *msgpack.Encoder) (err error) {
	if err = enc.EncodeArrayLen(l.len); err != nil {
		return
	}

	for n := l.head; n != nil; n = n.next {
		if err = enc.Encode(n.v); err != nil {
			return
		}
	}
	return
}

func (l *List[T]) DecodeMsgpack(dec *msgpack.Decoder) (err error) {
	var n int
	if n, err = dec.DecodeArrayLen(); err != nil {
		return
	}
	nodes := make([]listNode[T], n)
	for i := 0; i < n; i++ {
		n := &nodes[i]
		if err = dec.Decode(&n.v); err != nil {
			return
		}
		l.pushNode(n)
	}
	return
}
