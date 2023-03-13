package genh

import (
	"reflect"
	"sync"
)

var cloneCache struct {
	sync.RWMutex
	m map[reflect.Type]int
}

func hasCloner(t reflect.Type) int {
	cloneCache.RLock()
	v, ok := cloneCache.m[t]
	cloneCache.RUnlock()
	if ok {
		return v
	}
	cloneCache.Lock()
	defer cloneCache.Unlock()
	if cloneCache.m == nil {
		cloneCache.m = make(map[reflect.Type]int)
	}

	if hasClone(t) {
		v = 1
	} else if hasClone(reflect.PtrTo(t)) {
		v = 2
	}
	cloneCache.m[t] = v
	return v
}

func hasClone(t reflect.Type) bool {
	m, ok := t.MethodByName("Clone")
	if !ok {
		return false
	}
	if m.Type.NumOut() != 1 {
		return false
	}
	if ot := m.Type.Out(0); ot != m.Type.In(0) {
		return false
	}
	return true
}

type Cloner[T any] interface {
	Clone() T
}

func Clone[T any](v T, keepPrivateFields bool) (cp T) {
	if v, ok := any(v).(Cloner[T]); ok {
		return v.Clone()
	}
	src, dst := reflect.ValueOf(v), reflect.ValueOf(&cp).Elem()
	ReflectClone(dst, src, keepPrivateFields)
	return
}

func ReflectClone(dst, src reflect.Value, keepPrivateFields bool) {
	if !src.IsValid() || src.IsZero() {
		return
	}

	if src.Kind() == reflect.Interface {
		src = src.Elem()
	}

	styp := src.Type()

	if cloneVal(dst, src) {
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
		for i := 0; i < src.Len(); i++ {
			dst, src := dst.Index(i), src.Index(i)

			if dst.Kind() != reflect.Interface {
				ReflectClone(dst, src, keepPrivateFields)
				continue
			}

			if src.Kind() == reflect.Interface {
				src = src.Elem()
			}
			ndst := reflect.New(src.Type()).Elem()
			ReflectClone(ndst, src, keepPrivateFields)
			dst.Set(ndst)

		}

	case reflect.Map:
		if src.IsNil() {
			return
		}

		dst.Set(reflect.MakeMapWithSize(styp, src.Len()))
		for it := src.MapRange(); it.Next(); {
			mk, mv := maybeCopy(it.Key(), keepPrivateFields), maybeCopy(it.Value(), keepPrivateFields)
			dst.SetMapIndex(mk, mv)
		}

	case reflect.Struct:
		if keepPrivateFields {
			dst.Set(src) // copy private fields
		} else {
			dst.Set(reflect.New(styp).Elem())
		}

		for i := 0; i < styp.NumField(); i++ {
			if f := dst.Field(i); f.CanSet() {
				ReflectClone(dst.Field(i), src.Field(i), keepPrivateFields)
			}
		}

	case reflect.Ptr:
		if src.IsNil() {
			return
		}
		ndst := reflect.New(styp.Elem())
		ReflectClone(ndst.Elem(), src.Elem(), keepPrivateFields)
		dst.Set(ndst)

	default:
		dst.Set(src)
	}
}

func maybeCopy(src reflect.Value, copyPrivate bool) reflect.Value {
	switch src.Kind() {
	case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Map:
		nv := reflect.New(src.Type()).Elem()
		ReflectClone(nv, src, copyPrivate)
		return nv
	case reflect.Interface:
		return maybeCopy(src.Elem(), copyPrivate)
	default:
		return src
	}
}

func cloneVal(dst, src reflect.Value) bool {
	var m reflect.Value
	switch hasCloner(src.Type()) {
	case 1:
		m = src.MethodByName("Clone")
	case 2:
		if src.CanAddr() {
			m = src.Addr().MethodByName("Clone")
		}
	default:
		return false
	}

	v := m.Call(nil)[0]
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() == reflect.Ptr && dst.Kind() != reflect.Ptr {
		v = v.Elem()
	}

	dst.Set(v)
	return true
}
