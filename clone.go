package genh

import (
	"math"
	"reflect"
)

var (
	structCache SLMap[bool]
	cloneCache  SLMap[int]
)

type Cloner[T any] interface {
	Clone() T
}

func Clone[T any](v T, keepPrivateFields bool) (cp T) {
	if v, ok := any(v).(Cloner[T]); ok {
		return v.Clone()
	}
	src, dst := reflect.ValueOf(v), reflect.ValueOf(&cp).Elem()
	reflectClone(dst, src, keepPrivateFields, false)
	return
}

func ReflectClone(dst, src reflect.Value, keepPrivateFields bool) {
	reflectClone(dst, src, keepPrivateFields, true)
}

func reflectClone(dst, src reflect.Value, keepPrivateFields, checkClone bool) {
	if !src.IsValid() || src.IsZero() {
		return
	}

	if src.Kind() == reflect.Interface {
		src = src.Elem()
	}

	styp := src.Type()

	if checkClone && cloneVal(dst, src, isCloner(styp)) {
		return
	}

	switch styp.Kind() {
	case reflect.Slice:
		if src.IsNil() {
			return
		}
		dst.Set(reflect.MakeSlice(styp, src.Len(), src.Cap()))
		fallthrough

	case reflect.Array:
		isIface := styp.Elem().Kind() == reflect.Interface
		simple := isSimple(styp.Elem().Kind())
		for i := 0; i < src.Len(); i++ {
			dst, src := dst.Index(i), src.Index(i)

			if simple {
				dst.Set(src)
				continue
			}

			if !isIface {
				reflectClone(dst, src, keepPrivateFields, true)
				continue
			}

			src = src.Elem()
			ndst := reflect.New(src.Type()).Elem()
			reflectClone(ndst, src, keepPrivateFields, true)
			dst.Set(ndst)

		}

	case reflect.Map:
		if src.IsNil() {
			return
		}

		simpleKey := isSimple(styp.Key().Kind())
		simpleValue := isSimple(styp.Elem().Kind())
		dst.Set(reflect.MakeMapWithSize(styp, src.Len()))
		for it := src.MapRange(); it.Next(); {
			var mk, mv reflect.Value
			if simpleKey {
				mk = it.Key()
			} else {
				mk = maybeCopy(it.Key(), keepPrivateFields, checkClone)
			}
			if simpleValue {
				mv = it.Value()
			} else {
				mv = maybeCopy(it.Value(), keepPrivateFields, checkClone)
			}
			dst.SetMapIndex(mk, mv)
		}

	case reflect.Struct:
		if isSimpleStruct(styp) {
			dst.Set(src)
			return
		}
		if keepPrivateFields {
			dst.Set(src) // copy private fields
		} else {
			dst.Set(reflect.New(styp).Elem())
		}

		for i := 0; i < styp.NumField(); i++ {
			if f := dst.Field(i); f.CanSet() {
				if isSimple(f.Kind()) {
					f.Set(src.Field(i))
					continue
				}
				reflectClone(f, src.Field(i), keepPrivateFields, true)
			}
		}

	case reflect.Ptr:
		if src.IsNil() {
			return
		}
		ndst := reflect.New(styp.Elem())
		if nde := ndst.Elem(); isSimple(nde.Kind()) {
			nde.Set(src.Elem())
		} else {
			reflectClone(nde, src.Elem(), keepPrivateFields, true)
		}
		dst.Set(ndst)

	default:
		dst.Set(src)
	}
}

func isSimpleStruct(t reflect.Type) bool {
	key := t.Name()
	return structCache.MustGet(key, func() bool {
		for i := 0; i < t.NumField(); i++ {
			if !isSimple(t.Field(i).Type.Kind()) {
				return false
			}
		}
		return true
	})
}

func isSimple(k reflect.Kind) bool {
	switch k {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.String:
		return true
	default:
		return false
	}
}

func maybeCopy(src reflect.Value, copyPrivate, checkClone bool) reflect.Value {
	switch src.Kind() {
	case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
		nv := reflect.New(src.Type()).Elem()
		reflectClone(nv, src, copyPrivate, checkClone)
		return nv
	case reflect.Interface:
		return maybeCopy(src.Elem(), copyPrivate, true)
	default:
		return src
	}
}

func isCloner(t reflect.Type) int {
	key := t.Name()
	return cloneCache.MustGet(key, func() int {
		v := math.MaxInt
		if idx := cloneIdx(t); idx != math.MaxInt {
			v = idx + 1
		} else if idx := cloneIdx(reflect.PtrTo(t)); idx != math.MaxInt {
			v = -(idx + 1)
		}
		return v
	})
}

func cloneIdx(t reflect.Type) int {
	m, ok := t.MethodByName("Clone")
	if !ok {
		return math.MaxInt
	}

	if m.Type.NumOut() != 1 {
		return math.MaxInt
	}

	if m.Type.Out(0) != m.Type.In(0) {
		return math.MaxInt
	}

	return m.Index
}

func cloneVal(dst, src reflect.Value, idx int) bool {
	if idx == math.MaxInt {
		return false
	}
	var m reflect.Value
	if idx > 0 {
		m = src.Method(idx - 1)
	} else {
		m = src.Addr().Method(-idx - 1)
	}

	v := m.Call(nil)[0]
	if v.Kind() == reflect.Ptr && dst.Kind() != reflect.Ptr {
		v = v.Elem()
	}

	dst.Set(v)
	return true
}
