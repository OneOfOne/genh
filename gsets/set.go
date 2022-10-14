package gsets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"unsafe"

	"go.oneofone.dev/genh/internal"
)

func boolp(b bool) *bool {
	return &b
}

var empty struct{}

func SetOf[T internal.Ordered](keys ...T) Set[T] {
	s := Set[T]{}
	s.Add(keys...)
	return s
}

// Set[T] is a simple set.
type Set[T internal.Ordered] map[T]struct{}

func (s Set[T]) init() Set[T] {
	if s == nil {
		s = Set[T]{}
	}
	return s
}

func (s *Set[T]) Set(keys ...T) Set[T] {
	ss := s.Add(keys...)
	if *s == nil {
		*s = ss
	}
	return ss
}

func (s Set[T]) Add(keys ...T) Set[T] {
	s = s.init()
	for _, k := range keys {
		s[k] = empty
	}
	return s
}

// AddIfNotExists returns true if the key was added, false if it already existed
func (s *Set[T]) AddIfNotExists(key T) bool {
	sm := s.init()
	if *s == nil {
		*s = sm
	}
	if _, ok := sm[key]; ok {
		return false
	}

	sm[key] = empty
	return true
}

func (s Set[T]) Clone() Set[T] {
	ns := make(Set[T], len(s))
	for k, v := range s {
		ns[k] = v
	}
	return ns
}

func (s Set[T]) Merge(os ...Set[T]) Set[T] {
	s = s.init()
	for _, o := range os {
		for k := range o {
			s[k] = empty
		}
	}
	return s
}

func (s Set[T]) Delete(keys ...T) Set[T] {
	for _, k := range keys {
		delete(s, k)
	}
	return s
}

func (s Set[T]) Has(key T) bool {
	_, ok := s[key]
	return ok
}

func (s Set[T]) Equal(os Set[T]) bool {
	if len(os) != len(s) {
		return false
	}

	for k := range s {
		if _, ok := os[k]; !ok {
			return false
		}
	}
	return true
}

func (s Set[T]) Len() int {
	return len(s)
}

func (s Set[T]) Keys() []T {
	if s == nil {
		return nil
	}
	keys := make([]T, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	return keys
}

func (s Set[T]) SortedKeys() []T {
	keys := s.Keys()
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

func (s Set[T]) append(buf *bytes.Buffer, sorted bool) *bytes.Buffer {
	if len(s) == 0 {
		buf.WriteString("[]")
		return buf
	}

	var keys []T
	if sorted {
		keys = s.SortedKeys()
	} else {
		keys = s.Keys()
	}

	_, isString := any(keys[0]).(string)

	verb := "%v"
	if isString {
		verb = "%q"
	}
	buf.WriteByte('[')
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		fmt.Fprintf(buf, verb, k)
	}
	buf.WriteByte(']')

	return buf
}

func (s Set[T]) String() string {
	if len(s) == 0 {
		return "[]"
	}
	buf := s.append(&bytes.Buffer{}, true).Bytes()
	buf = buf[:len(buf):len(buf)]
	return *(*string)(unsafe.Pointer(&buf))
}

func (s Set[T]) MarshalJSON() ([]byte, error) {
	buf := s.append(&bytes.Buffer{}, true).Bytes()
	buf = buf[:len(buf):len(buf)]
	return buf, nil
}

func (s *Set[T]) UnmarshalJSON(data []byte) (err error) {
	var keys []T
	if err = json.Unmarshal(data, &keys); err == nil {
		s.Set(keys...)
	}
	return
}
