package genh

import (
	"bytes"
	"io"

	"github.com/vmihailenco/msgpack/v5"
)

type (
	MsgpackEncoder = msgpack.Encoder
	MsgpackDecoder = msgpack.Decoder
)

var decPool = Pool[MsgpackDecoder]{
	New: func() *MsgpackDecoder {
		dec := msgpack.NewDecoder(nil)
		dec.SetCustomStructTag("json")
		dec.UseLooseInterfaceDecoding(true)
		return dec
	},
}

func PutMsgpackDecoder(dec *MsgpackDecoder) {
	dec.Reset(nil)
	decPool.Put(dec)
}

var encPool = Pool[MsgpackEncoder]{
	New: func() *MsgpackEncoder {
		enc := msgpack.NewEncoder(nil)
		enc.SetCustomStructTag("json")
		enc.UseCompactFloats(true)
		enc.UseCompactInts(true)
		return enc
	},
}

func PutMsgpackEncoder(enc *MsgpackEncoder) {
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

func EncodeMsgpack(w io.Writer, vs ...any) error {
	enc := NewMsgpackEncoder(w)
	err := enc.EncodeMulti(vs...)
	PutMsgpackEncoder(enc)
	return err
}

func DecodeMsgpack(r io.Reader, vs ...any) error {
	dec := NewMsgpackDecoder(r)
	err := dec.DecodeMulti(vs...)
	PutMsgpackDecoder(dec)
	return err
}

// NewMsgpackDecoder returns a new Decoder that writes to w.
// uses json CustomStructTag, compact floats and ints.
func NewMsgpackEncoder(w io.Writer) *MsgpackEncoder {
	enc := encPool.Get()
	enc.Reset(w)
	return enc
}

// NewMsgpackDecoder returns a new Decoder that reads from r.
// uses json CustomStructTag, and loose interface decoding.
func NewMsgpackDecoder(r io.Reader) *MsgpackDecoder {
	dec := decPool.Get()
	dec.Reset(r)
	return dec
}
