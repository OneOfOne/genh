package genh

import (
	"io"

	"github.com/vmihailenco/msgpack/v5"
	"go.oneofone.dev/genh/internal"
)

type (
	MsgpackEncoder = msgpack.Encoder
	MsgpackDecoder = msgpack.Decoder
)

func UnmarshalMsgpack(b []byte, v any) error {
	return internal.UnmarshalMsgpack(b, v)
}

func MarshalMsgpack(v any) ([]byte, error) {
	return internal.MarshalMsgpack(v)
}

func EncodeMsgpack(w io.Writer, vs ...any) error {
	return internal.EncodeMsgpack(w, vs...)
}

func DecodeMsgpack(r io.Reader, vs ...any) error {
	return internal.DecodeMsgpack(r, vs...)
}

// NewMsgpackDecoder returns a new Decoder that writes to w.
// uses json CustomStructTag, compact floats and ints.
func NewMsgpackEncoder(w io.Writer) *MsgpackEncoder {
	return internal.NewMsgpackEncoder(w)
}

func PutMsgpackEncoder(enc *MsgpackEncoder) {
	internal.PutMsgpackEncoder(enc)
}

// NewMsgpackDecoder returns a new Decoder that reads from r.
// uses json CustomStructTag, and loose interface decoding.
func NewMsgpackDecoder(r io.Reader) *MsgpackDecoder {
	return internal.NewMsgpackDecoder(r)
}

func PutMsgpackDecoder(dec *MsgpackDecoder) {
	internal.PutMsgpackDecoder(dec)
}
