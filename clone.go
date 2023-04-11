package genh

import (
	"reflect"
	"sync"
)

var cloneCache struct {
	m map[reflect.Type]int
	sync.RWMutex
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
	if checkClone {
		if cv := hasCloner(styp); cv != 0 {
			if cloneVal(dst, src, cv) {
				return
			}
		}
	}

	switch styp.Kind() {
	case reflect.Slice:
		if src.IsNil() {
			return
		}
		dst.Set(reflect.MakeSlice(styp, src.Len(), src.Cap()))
		fallthrough

	case reflect.Array:
		hasClone := hasCloner(styp.Elem())
		if hasClone > 0 {
			for i := 0; i < src.Len(); i++ {
				dst, src := dst.Index(i), src.Index(i)
				if !cloneVal(dst, src, hasClone) {
					panic("bad")
				}
			}
			break
		}

		for i := 0; i < src.Len(); i++ {
			dst, src := dst.Index(i), src.Index(i)

			if dst.Kind() != reflect.Interface {
				reflectClone(dst, src, keepPrivateFields, false)
				continue
			}

			if src.Kind() == reflect.Interface {
				src = src.Elem()
			}
			ndst := reflect.New(src.Type()).Elem()
			reflectClone(ndst, src, keepPrivateFields, false)
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
				reflectClone(dst.Field(i), src.Field(i), keepPrivateFields, true)
			}
		}

	case reflect.Ptr:
		if src.IsNil() {
			return
		}
		ndst := reflect.New(styp.Elem())
		reflectClone(ndst.Elem(), src.Elem(), keepPrivateFields, true)
		dst.Set(ndst)

	default:
		dst.Set(src)
	}
}

func maybeCopy(src reflect.Value, copyPrivate bool) reflect.Value {
	switch src.Kind() {
	case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
		nv := reflect.New(src.Type()).Elem()
		reflectClone(nv, src, copyPrivate, hasCloner(nv.Type()) > 0)
		return nv
	case reflect.Interface:
		return maybeCopy(src.Elem(), copyPrivate)
	default:
		return src
	}
}

func cloneVal(dst, src reflect.Value, cv int) bool {
	var m reflect.Value
	switch cv {
	case 1:
		m = src.MethodByName("Clone")
	case 2:
		m = src.Addr().MethodByName("Clone")
	default:
		return false
	}

	v := m.Call(nil)[0]
	if v.Kind() == reflect.Ptr && dst.Kind() != reflect.Ptr {
		v = v.Elem()
	}

	dst.Set(v)
	return true
}
