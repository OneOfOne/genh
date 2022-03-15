package genh

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"sync/atomic"
)

type (
	AtomicInt8  = signedValue32[int8]
	AtomicInt16 = signedValue32[int16]
	AtomicInt32 = signedValue32[int32]
	AtomicInt64 = signedValue64[int64]
	AtomicInt   = signedValue64[int]

	AtomicUint8  = unsignedValue32[uint8]
	AtomicUint16 = unsignedValue32[uint16]
	AtomicUint32 = unsignedValue32[uint32]
	AtomicUint64 = unsignedValue64[uint64]
	AtomicUint   = unsignedValue64[uint]

	AtomicUintptr = unsignedValue64[uintptr]
)

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

type signedValue64[T Signed] struct {
	noCopy noCopy
	v      int64
}

func (v *signedValue64[T]) Store(val T)        { atomic.StoreInt64(&v.v, int64(val)) }
func (v *signedValue64[T]) Load() T            { return T(atomic.LoadInt64(&v.v)) }
func (v *signedValue64[T]) Add(val T) T        { return T(atomic.AddInt64(&v.v, int64(val))) }
func (v *signedValue64[T]) Swap(val T) (old T) { return T(atomic.SwapInt64(&v.v, int64(val))) }
func (v *signedValue64[T]) CompareAndSwap(old, new T) bool {
	return atomic.CompareAndSwapInt64(&v.v, int64(old), int64(new))
}

func (v *signedValue64[T]) MarshalJSON() ([]byte, error)   { return json.Marshal(v.Load()) }
func (v *signedValue64[T]) UnmarshalJSON(b []byte) error   { return unmarshalJSON(b, v.Store) }
func (v *signedValue64[T]) MarshalBinary() ([]byte, error) { return marshalBinaryI(v.Load()) }
func (v *signedValue64[T]) UnmarshalBinary(b []byte) error { return unmarshalBinaryI(b, v.Store) }

type signedValue32[T Signed] struct {
	noCopy noCopy
	v      int32
}

func (v *signedValue32[T]) Store(val T)        { atomic.StoreInt32(&v.v, int32(val)) }
func (v *signedValue32[T]) Load() T            { return T(atomic.LoadInt32(&v.v)) }
func (v *signedValue32[T]) Add(val T) T        { return T(atomic.AddInt32(&v.v, int32(val))) }
func (v *signedValue32[T]) Swap(val T) (old T) { return T(atomic.SwapInt32(&v.v, int32(val))) }
func (v *signedValue32[T]) CompareAndSwap(old, new T) bool {
	return atomic.CompareAndSwapInt32(&v.v, int32(old), int32(new))
}

func (v *signedValue32[T]) MarshalJSON() ([]byte, error)   { return json.Marshal(v.Load()) }
func (v *signedValue32[T]) UnmarshalJSON(b []byte) error   { return unmarshalJSON(b, v.Store) }
func (v *signedValue32[T]) MarshalBinary() ([]byte, error) { return marshalBinaryI(v.Load()) }
func (v *signedValue32[T]) UnmarshalBinary(b []byte) error { return unmarshalBinaryI(b, v.Store) }

type unsignedValue64[T Unsigned] struct {
	noCopy noCopy
	v      uint64
}

func (v *unsignedValue64[T]) Store(val T)        { atomic.StoreUint64(&v.v, uint64(val)) }
func (v *unsignedValue64[T]) Load() T            { return T(atomic.LoadUint64(&v.v)) }
func (v *unsignedValue64[T]) Add(val T) T        { return T(atomic.AddUint64(&v.v, uint64(val))) }
func (v *unsignedValue64[T]) Swap(val T) (old T) { return T(atomic.SwapUint64(&v.v, uint64(val))) }
func (v *unsignedValue64[T]) CompareAndSwap(old, new T) bool {
	return atomic.CompareAndSwapUint64(&v.v, uint64(old), uint64(new))
}

func (v *unsignedValue64[T]) MarshalJSON() ([]byte, error)   { return json.Marshal(v.Load()) }
func (v *unsignedValue64[T]) UnmarshalJSON(b []byte) error   { return unmarshalJSON(b, v.Store) }
func (v *unsignedValue64[T]) MarshalBinary() ([]byte, error) { return marshalBinaryU(v.Load()) }
func (v *unsignedValue64[T]) UnmarshalBinary(b []byte) error { return unmarshalBinaryU(b, v.Store) }

type unsignedValue32[T Unsigned] struct {
	noCopy noCopy
	v      uint32
}

func (v *unsignedValue32[T]) Store(val T)        { atomic.StoreUint32(&v.v, uint32(val)) }
func (v *unsignedValue32[T]) Load() T            { return T(atomic.LoadUint32(&v.v)) }
func (v *unsignedValue32[T]) Add(val T) T        { return T(atomic.AddUint32(&v.v, uint32(val))) }
func (v *unsignedValue32[T]) Swap(val T) (old T) { return T(atomic.SwapUint32(&v.v, uint32(val))) }
func (v *unsignedValue32[T]) CompareAndSwap(old, new T) bool {
	return atomic.CompareAndSwapUint32(&v.v, uint32(old), uint32(new))
}

func (v *unsignedValue32[T]) MarshalJSON() ([]byte, error)   { return json.Marshal(v.Load()) }
func (v *unsignedValue32[T]) UnmarshalJSON(b []byte) error   { return unmarshalJSON(b, v.Store) }
func (v *unsignedValue32[T]) MarshalBinary() ([]byte, error) { return marshalBinaryU(v.Load()) }
func (v *unsignedValue32[T]) UnmarshalBinary(b []byte) error { return unmarshalBinaryU(b, v.Store) }

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
