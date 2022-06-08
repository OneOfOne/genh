package genh

import (
	"encoding/json"
)

func Zero[T any]() (_ T) { return }

type PtrTo[T any] struct {
	v *T
}

func (p *PtrTo[T]) Set(v T) {
	p.v = &v
}

func (p *PtrTo[T]) Unset() {
	p.v = nil
}

func (p *PtrTo[T]) IsSet() bool {
	return p.v != nil
}

func (p *PtrTo[T]) Val() T {
	if p.v != nil {
		return *p.v
	}
	return Zero[T]()
}

func (p PtrTo[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.v)
}

func (p *PtrTo[T]) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &p.v)
}

func (p PtrTo[T]) MarshalBinary() ([]byte, error) {
	return MarshalMsgpack(p.v)
}

func (p *PtrTo[T]) UnmarshalBinary(b []byte) error {
	return UnmarshalMsgpack(b, &p.v)
}
