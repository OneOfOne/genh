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
	head *listNode[T]
	tail *listNode[T]
	clip *listNode[T]
	len  int
}

func (l List[T]) Len() int {
	return l.len
}

func (l List[T]) Head() (v T) {
	if l.head != nil {
		v = l.head.v
	}
	return
}

func (l List[T]) Tail() (v T) {
	if l.tail != nil {
		v = l.tail.v
	}
	return
}

func (l List[T]) get(idx int) (n *listNode[T]) {
	if idx >= l.len || idx < 0 {
		panic("index out of range")
	}

	if idx == l.len-1 {
		return l.tail
	}

	n = l.head
	for i := 0; i < idx; i++ {
		n = n.next
	}
	return
}

func (l List[T]) ListAt(start, end int) List[T] {
	h := l.get(start)
	if end >= l.len {
		end = l.len - 1
	}

	if end < 0 {
		end = l.len + end
	}

	t := l.get(end)
	nl := List[T]{
		head: h,
		tail: t,
		len:  end - start + 1,
	}
	return nl.Clip()
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

func (l *List[T]) Push(vs ...T) {
	l.len += len(vs)
	if len(vs) == 1 {
		l.pushNode(&listNode[T]{v: vs[0]})
		return
	}

	nodes := make([]listNode[T], len(vs))
	for i, v := range vs {
		n := &nodes[i]
		n.v = v
		l.pushNode(n)
	}
}

// PushSort pushes the item in the order returned by lessFn
// PushSort does *not* work with Clip, use Clone if you really have to.
func (l *List[T]) PushSort(v T, lessFn func(a, b T) bool) {
	l.len++
	nn := &listNode[T]{v: v}
	if l.tail == nil || !lessFn(v, l.tail.v) {
		l.pushNode(nn)
		return
	}

	for n := &l.head; *n != nil; n = l.nextNodePtr(*n) {
		if np := *n; lessFn(v, np.v) {
			nn.next = np
			*n = nn
			return
		}
	}
}

func (l List[T]) Append(vs ...T) List[T] {
	l.Push(vs...)
	return l
}

func (l *List[T]) pushNode(n *listNode[T]) {
	if l.head == nil {
		l.head, l.tail = n, n
		return
	}

	if l.clip == l.tail {
		tn := *l.tail
		l.tail = &tn
	}
	l.tail.next = n
	l.tail = n
}

func (l List[T]) nextNode(n *listNode[T]) *listNode[T] {
	if l.clip == n {
		if l.clip == l.tail {
			return nil
		}
		return l.tail
	}
	return n.next
}

func (l List[T]) nextNodePtr(n *listNode[T]) **listNode[T] {
	if l.clip == n {
		if l.clip == l.tail {
			return nil
		}
		return &l.tail
	}
	return &n.next
}

func (l *List[T]) Prepend(v T) {
	n := &listNode[T]{v: v}
	if l.len++; l.head == nil {
		l.head, l.tail = n, n
		return
	}

	n.next = l.head
	l.head = n
}

func (l List[T]) Slice() (out []T) {
	if l.head == nil {
		return
	}

	out = make([]T, 0, l.len)
	for n := l.head; n != nil; n = l.nextNode(n) {
		out = append(out, n.v)
	}
	return
}

func (l List[T]) Clone() (out List[T]) {
	out.len = l.len
	for n := l.head; n != nil; n = l.nextNode(n) {
		out.pushNode(&listNode[T]{v: n.v})
	}
	return
}

func (l List[T]) count() (out int) {
	if l.head == nil {
		return
	}

	for n := l.head; n != nil; n = l.nextNode(n) {
		out++
	}
	return
}

func (l List[T]) Clip() List[T] {
	l.clip = l.tail
	return l
}

func (l *List[T]) Clear() {
	*l = List[T]{}
}

// Iter is a c++-style iterator:
// it := l.Iter()
// for v := it.Value(); it.Next(); v = it.Value()) {}
func (l *List[T]) Iter() *ListIterator[T] {
	return &ListIterator[T]{l: l, n: l.head}
}

func (l List[T]) IterChan(cap int) <-chan T {
	if cap == 0 {
		cap = 1
	}
	ch := make(chan T, cap)
	go func() {
		defer close(ch)
		for n := l.head; n != nil; n = l.nextNode(n) {
			ch <- n.v
		}
	}()
	return ch
}

func (l List[T]) ForEach(fn func(v T) bool) {
	for n := l.head; n != nil; n = l.nextNode(n) {
		if !fn(n.v) {
			break
		}
	}
}

func (l List[T]) ForEachPtr(fn func(v *T) bool) {
	for n := l.head; n != nil; n = l.nextNode(n) {
		if !fn(&n.v) {
			break
		}
	}
}

func (l List[T]) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('[')
	enc := json.NewEncoder(&buf)
	for n := l.head; n != nil; n = l.nextNode(n) {
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

	for n := l.head; n != nil; n = l.nextNode(n) {
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
		var n listNode[T]
		if err = dec.Decode(&n.v); err != nil {
			return
		}
		l.pushNode(&n)
	}
	return
}

type ListIterator[T any] struct { // i hate how much this feels like java/c++
	l    *List[T]
	n    *listNode[T]
	prev *listNode[T]
}

func (it *ListIterator[T]) Next() (v T, ok bool) {
	if ok = it.n != nil; !ok {
		return
	}
	it.prev = it.n
	v, it.n = it.n.v, it.l.nextNode(it.n)
	return
}

func (it *ListIterator[T]) Set(v T) {
	it.prev.v = v
}

func (it *ListIterator[T]) Delete() {
	it.l.len--
	if it.l.head == it.prev {
		it.l.head = it.n
	} else {
		it.prev.next = it.n.next
		it.prev.v = it.n.v
	}
}
