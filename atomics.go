package genh

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"runtime"
	"sync/atomic"
)

type (
	AtomicInt8  = AtomicSigned32[int8]
	AtomicInt16 = AtomicSigned32[int16]
	AtomicInt32 = AtomicSigned32[int32]
	AtomicInt64 = AtomicSigned64[int64]
	AtomicInt   = AtomicSigned64[int]

	AtomicUint8  = AtomicUnsigned32[uint8]
	AtomicUint16 = AtomicUnsigned32[uint16]
	AtomicUint32 = AtomicUnsigned32[uint32]
	AtomicUint64 = AtomicUnsigned64[uint64]
	AtomicUint   = AtomicUnsigned64[uint]

	AtomicUintptr = AtomicUnsigned64[uintptr]
)

type AtomicBool struct {
	v AtomicInt32
}

func (v *AtomicBool) Store(val bool) { v.v.Store(Iff[int32](val, 1, 2)) }

func (v *AtomicBool) Load() bool { return v.v.Load() != 0 }

func (v *AtomicBool) Swap(val bool) (old bool) { return v.v.Swap(Iff[int32](val, 1, 2)) != 0 }
func (v *AtomicBool) CompareAndSwap(old, new bool) bool {
	return v.v.CompareAndSwap(Iff[int32](old, 1, 2), Iff[int32](new, 1, 2))
}

func (v *AtomicBool) MarshalJSON() ([]byte, error)   { return json.Marshal(v.Load()) }
func (v *AtomicBool) UnmarshalJSON(b []byte) error   { return unmarshalJSON(b, v.Store) }
func (v *AtomicBool) MarshalBinary() ([]byte, error) { return marshalBinaryI(v.v.Load()) }
func (v *AtomicBool) UnmarshalBinary(b []byte) error { return unmarshalBinaryI(b, v.v.Store) }

type AtomicFloat64 struct {
	v AtomicUint64
}

func (v *AtomicFloat64) Store(val float64) { v.v.Store(math.Float64bits(val)) }
func (v *AtomicFloat64) Load() float64     { return math.Float64frombits(v.v.Load()) }
func (v *AtomicFloat64) Add(val float64) float64 {
	for {
		old := v.Load()
		new := old + val
		if v.CompareAndSwap(old, new) {
			return new
		}
		runtime.Gosched()
	}
}

func (v *AtomicFloat64) Swap(val float64) (old float64) {
	return math.Float64frombits(v.v.Swap(math.Float64bits(val)))
}

func (v *AtomicFloat64) CompareAndSwap(old, new float64) bool {
	return v.v.CompareAndSwap(math.Float64bits(old), math.Float64bits(new))
}

func (v *AtomicFloat64) MarshalJSON() ([]byte, error)   { return json.Marshal(v.Load()) }
func (v *AtomicFloat64) UnmarshalJSON(b []byte) error   { return unmarshalJSON(b, v.Store) }
func (v *AtomicFloat64) MarshalBinary() ([]byte, error) { return marshalBinaryU(v.v.Load()) }
func (v *AtomicFloat64) UnmarshalBinary(b []byte) error { return unmarshalBinaryU(b, v.v.Store) }

type AtomicSigned32[T Signed] struct {
	v atomic.Int32
}

func (v *AtomicSigned32[T]) Store(val T)        { v.v.Store(int32(val)) }
func (v *AtomicSigned32[T]) Load() T            { return T(v.v.Load()) }
func (v *AtomicSigned32[T]) Add(val T) T        { return T(v.v.Add(int32(val))) }
func (v *AtomicSigned32[T]) Swap(val T) (old T) { return T(v.v.Swap(int32(val))) }
func (v *AtomicSigned32[T]) CompareAndSwap(old, new T) bool {
	return v.v.CompareAndSwap(int32(old), int32(new))
}

func (v *AtomicSigned32[T]) MarshalJSON() ([]byte, error)   { return json.Marshal(v.Load()) }
func (v *AtomicSigned32[T]) UnmarshalJSON(b []byte) error   { return unmarshalJSON(b, v.Store) }
func (v *AtomicSigned32[T]) MarshalBinary() ([]byte, error) { return marshalBinaryI(v.Load()) }
func (v *AtomicSigned32[T]) UnmarshalBinary(b []byte) error { return unmarshalBinaryI(b, v.Store) }

type AtomicSigned64[T Signed] struct {
	v atomic.Int64
}

func (v *AtomicSigned64[T]) Store(val T)        { v.v.Store(int64(val)) }
func (v *AtomicSigned64[T]) Load() T            { return T(v.v.Load()) }
func (v *AtomicSigned64[T]) Add(val T) T        { return T(v.v.Add(int64(val))) }
func (v *AtomicSigned64[T]) Swap(val T) (old T) { return T(v.v.Swap(int64(val))) }
func (v *AtomicSigned64[T]) CompareAndSwap(old, new T) bool {
	return v.v.CompareAndSwap(int64(old), int64(new))
}

func (v *AtomicSigned64[T]) MarshalJSON() ([]byte, error)   { return json.Marshal(v.Load()) }
func (v *AtomicSigned64[T]) UnmarshalJSON(b []byte) error   { return unmarshalJSON(b, v.Store) }
func (v *AtomicSigned64[T]) MarshalBinary() ([]byte, error) { return marshalBinaryI(v.Load()) }
func (v *AtomicSigned64[T]) UnmarshalBinary(b []byte) error { return unmarshalBinaryI(b, v.Store) }

type AtomicUnsigned32[T Unsigned] struct {
	v atomic.Uint32
}

func (v *AtomicUnsigned32[T]) Store(val T)        { v.v.Store(uint32(val)) }
func (v *AtomicUnsigned32[T]) Load() T            { return T(v.v.Load()) }
func (v *AtomicUnsigned32[T]) Add(val T) T        { return T(v.v.Add(uint32(val))) }
func (v *AtomicUnsigned32[T]) Swap(val T) (old T) { return T(v.v.Swap(uint32(val))) }
func (v *AtomicUnsigned32[T]) CompareAndSwap(old, new T) bool {
	return v.v.CompareAndSwap(uint32(old), uint32(new))
}

func (v *AtomicUnsigned32[T]) MarshalJSON() ([]byte, error)   { return json.Marshal(v.Load()) }
func (v *AtomicUnsigned32[T]) UnmarshalJSON(b []byte) error   { return unmarshalJSON(b, v.Store) }
func (v *AtomicUnsigned32[T]) MarshalBinary() ([]byte, error) { return marshalBinaryU(v.Load()) }
func (v *AtomicUnsigned32[T]) UnmarshalBinary(b []byte) error { return unmarshalBinaryU(b, v.Store) }

type AtomicUnsigned64[T Unsigned] struct {
	v atomic.Uint64
}

func (v *AtomicUnsigned64[T]) Store(val T)        { v.v.Store(uint64(val)) }
func (v *AtomicUnsigned64[T]) Load() T            { return T(v.v.Load()) }
func (v *AtomicUnsigned64[T]) Add(val T) T        { return T(v.v.Add(uint64(val))) }
func (v *AtomicUnsigned64[T]) Swap(val T) (old T) { return T(v.v.Swap(uint64(val))) }
func (v *AtomicUnsigned64[T]) CompareAndSwap(old, new T) bool {
	return v.v.CompareAndSwap(uint64(old), uint64(new))
}

func (v *AtomicUnsigned64[T]) MarshalJSON() ([]byte, error)   { return json.Marshal(v.Load()) }
func (v *AtomicUnsigned64[T]) UnmarshalJSON(b []byte) error   { return unmarshalJSON(b, v.Store) }
func (v *AtomicUnsigned64[T]) MarshalBinary() ([]byte, error) { return marshalBinaryU(v.Load()) }
func (v *AtomicUnsigned64[T]) UnmarshalBinary(b []byte) error { return unmarshalBinaryU(b, v.Store) }

type noCopy struct{}

// lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) lock()   {}
func (*noCopy) unlock() {}

func unmarshalJSON[T any](b []byte, stFn func(T)) (err error) {
	var val T
	if err = json.Unmarshal(b, val); err != nil {
		return
	}
	stFn(val)
	return
}

func marshalBinaryI[T Signed](val T) ([]byte, error) {
	var b [binary.MaxVarintLen64]byte
	sz := binary.PutVarint(b[:], int64(val))
	return b[:sz:sz], nil
}

func marshalBinaryU[T Unsigned](val T) ([]byte, error) {
	var b [binary.MaxVarintLen64]byte
	sz := binary.PutUvarint(b[:], uint64(val))
	return b[:sz:sz], nil
}

func unmarshalBinaryI[T Signed](b []byte, set func(T)) error {
	v, i := binary.Varint(b)
	if i == 0 && len(b) > 0 {
		return fmt.Errorf("invalid signed varint: %v", b)
	}
	set(T(v))
	return nil
}

func unmarshalBinaryU[T Unsigned](b []byte, set func(T)) error {
	v, i := binary.Uvarint(b)
	if i == 0 && len(b) > 0 {
		return fmt.Errorf("invalid unsigned varint: %v", b)
	}
	set(T(v))
	return nil
}
