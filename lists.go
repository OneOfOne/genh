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

type ListNode[T any] struct {
	v          T
	next, prev *ListNode[T]
}

type List[T any] struct {
	head, tail *ListNode[T]
	len        int
}

func (l *List[T]) Append(vs ...T) *List[T] {
	if l == nil {
		l = &List[T]{}
	}
	for _, v := range vs {
		l.Push(v)
	}
	return l
}

func (l *List[T]) Push(v T) {
	l.len++
	n := &ListNode[T]{v: v}
	if l.head == nil {
		l.head, l.tail = n, n
		return
	}

	l.tail.next = n
	n.prev = l.tail
	l.tail = n
}

func (l *List[T]) Pop() (_ T, _ bool) {
	if l.tail == nil {
		return
	}
	l.len--
	n := l.tail
	l.tail = n.prev
	if l.tail != nil {
		l.tail.next = nil
	} else {
		l.head = nil
	}

	return n.v, true
}

func (l *List[T]) Unshift(v T) {
	n := &ListNode[T]{v: v}
	l.len++
	if l.head == nil {
		l.head, l.tail = n, n
		return
	}

	l.head.prev = n
	n.next = l.head
	l.head = n
}

func (l *List[T]) Shift() (_ T, _ bool) {
	if l.head == nil {
		return
	}
	l.len--
	n := l.head
	l.head = n.next
	if l.head != nil {
		l.head.prev = nil
	} else {
		l.tail = nil
	}

	return n.v, true
}

func (l *List[T]) Iter() func() (v T, ok bool) {
	n := l.head
	return func() (v T, ok bool) {
		if ok = n != nil; ok {
			v = n.v
			n = n.next
		}
		return
	}
}

func (l *List[T]) RevIter() func() (v T, ok bool) {
	n := l.tail
	return func() (v T, ok bool) {
		if ok = n != nil; ok {
			v = n.v
			n = n.prev
		}
		return
	}
}

func (l *List[T]) Slice(rev bool) (out []T) {
	if l.head == nil {
		return
	}
	n := l.head
	if rev {
		n = l.tail
	}
	out = make([]T, 0, l.len)
	for n != nil {
		out = append(out, n.v)
		if rev {
			n = n.prev
		} else {
			n = n.next
		}
	}
	return
}

func (l List[T]) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('[')
	enc := json.NewEncoder(&buf)
	for n := l.head; n != nil; n = n.next {
		if buf.Len() > 1 {
			buf.WriteString(", ")
		}
		if err := enc.Encode(n.v); err != nil {
			return nil, err
		}
		// buf.Truncate(buf.Len() - 2)
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

func (l *List[T]) UnmarshalJSON(p []byte) error {
	var v []T
	if err := json.Unmarshal(p, &v); err != nil {
		return err
	}
	l.Append(v...)
	return nil
}

func (l List[T]) MarshalBinary() ([]byte, error) {
	return msgpack.Marshal(&l)
}

func (l *List[T]) UnmarshalBinary(p []byte) error {
	return msgpack.Unmarshal(p, &l)
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

	for i := 0; i < n; i++ {
		var v T
		if err = dec.Decode(&v); err != nil {
			return
		}
		l.Push(v)
	}
	return
}
