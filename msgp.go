package genh

import (
	"bytes"
	"io"
	"sync"

	"github.com/vmihailenco/msgpack/v5"
)

var decPool = sync.Pool{
	New: func() any {
		dec := msgpack.NewDecoder(nil)
		dec.SetCustomStructTag("json")
		dec.UseLooseInterfaceDecoding(true)
		return dec
	},
}

var encPool = sync.Pool{
	New: func() any {
		enc := msgpack.NewEncoder(nil)
		enc.SetCustomStructTag("json")
		enc.UseCompactFloats(true)
		enc.UseCompactInts(true)
		return enc
	},
}

func UnmarshalMsgpack(b []byte, v any) error {
	dec := NewMsgpackDecoder(bytes.NewReader(b))
	err := dec.Decode(v)
	dec.Reset(nil)
	decPool.Put(dec)
	return err
}

func MarshalMsgpack(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := NewMsgpackEncoder(&buf)
	err := enc.Encode(v)
	enc.Reset(nil)
	encPool.Put(enc)
	return buf.Bytes(), err
}

// NewMsgpackEncoder returns a new Encoder that writes to w.
// uses json CustomStrucTag, compact floats and ints.
func NewMsgpackEncoder(w io.Writer) *msgpack.Encoder {
	enc := encPool.Get().(*msgpack.Encoder)
	enc.Reset(w)
	return enc
}

// NewMsgpackDecoder returns a new Decoder that reads from r.
// uses json CustomStrucTag, and loose interface decoding.
func NewMsgpackDecoder(r io.Reader) *msgpack.Decoder {
	dec := decPool.Get().(*msgpack.Decoder)
	dec.Reset(r)
	return dec
}
