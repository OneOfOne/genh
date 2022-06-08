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

func PutMsgpackDecoder(dec *msgpack.Decoder) {
	dec.Reset(nil)
	decPool.Put(dec)
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

func PutMsgpackEncoder(enc *msgpack.Encoder) {
	enc.Reset(nil)
	encPool.Put(enc)
}

func UnmarshalMsgpack(b []byte, v any) error {
	return DecodeMsgpack(bytes.NewReader(b), v)
}

func MarshalMsgpack(v any) ([]byte, error) {
	var buf bytes.Buffer
	err := EncodeMsgpack(&buf, v)
	return buf.Bytes(), err
}

func EncodeMsgpack(w io.Writer, v any) error {
	enc := NewMsgpackEncoder(w)
	err := enc.Encode(v)
	PutMsgpackEncoder(enc)
	return err
}

func DecodeMsgpack(r io.Reader, v any) error {
	dec := NewMsgpackDecoder(r)
	err := dec.Decode(v)
	PutMsgpackDecoder(dec)
	return err
}

// NewMsgpackDecoder returns a new Decoder that writes to w.
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
