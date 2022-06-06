package genh

import (
	"io"
	"os"
)

type TypeEncoder interface {
	Encode(v any) error
}

type TypeDecoder interface {
	Decode(v any) error
}

func EncodeFile[EncT TypeEncoder](fp string, v any, fn func(w io.Writer) EncT) (err error) {
	var f *os.File
	if f, err = os.Create(fp); err != nil {
		return
	}
	defer f.Close()
	return Encode(f, v, fn)
}

func Encode[EncT TypeEncoder](w io.Writer, v any, fn func(w io.Writer) EncT) (err error) {
	enc := fn(w)
	err = enc.Encode(v)
	return
}

func DecodeFile[T any, DecT DecoderType](fp string, fn func(r io.Reader) DecT) (v T, err error) {
	var f *os.File
	if f, err = os.Open(fp); err != nil {
		return
	}
	defer f.Close()
	return Decode[T](f, fn)
}

func Decode[T any, DecT DecoderType](r io.Reader, fn func(r io.Reader) DecT) (v T, err error) {
	dec := fn(r)
	err = dec.Decode(&v)
	return
}
