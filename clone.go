package genh

import (
	"log"
	"reflect"
)

type Cloner[T any] interface {
	Clone() T
}

func Clone[T any](v T) (cp T) {
	if v, ok := any(v).(Cloner[T]); ok {
		return v.Clone()
	}
	src, dst := reflect.Indirect(reflect.ValueOf(v)), reflect.ValueOf(&cp).Elem()
	ReflectClone(dst, src)
	return
}

func ReflectClone(dst, src reflect.Value) {
	if !src.IsValid() || src.IsZero() {
		return
	}

	styp := src.Type()

	if dst.Kind() == reflect.Ptr && dst.IsNil() {
		dst.Set(reflect.New(styp))
		dst = dst.Elem()
	}

	if styp != dst.Type() {
		log.Panicf("type mismatch %v %v", styp, dst.Type())
	}

	switch src.Kind() {
	case reflect.Slice:
		if src.IsNil() {
			return
		}
		dst.Set(reflect.MakeSlice(styp, src.Len(), src.Cap()))
		fallthrough

	case reflect.Array:
		for i := 0; i < src.Len(); i++ {
			ReflectClone(dst.Index(i), src.Index(i))
		}

	case reflect.Map:
		if src.IsNil() {
			return
		}

		dst.Set(reflect.MakeMapWithSize(styp, src.Len()))
		for it := src.MapRange(); it.Next(); {
			mk, mv := maybeCopy(it.Key()), maybeCopy(it.Value())
			dst.SetMapIndex(mk, mv)
		}

	case reflect.Struct:
		dst.Set(src) // copy private fields
		for i := 0; i < styp.NumField(); i++ {
			f := dst.Field(i)
			if f.CanSet() {
				ReflectClone(dst.Field(i), src.Field(i))
			}
		}

	case reflect.Ptr:
		if src.IsNil() {
			return
		}
		v := reflect.New(styp).Elem()
		ReflectClone(v, src.Elem())
		dst.Set(v)

	case reflect.Interface:
		if src.IsNil() {
			return
		}
		v := reflect.New(src.Elem().Type()).Elem()
		ReflectClone(v, src.Elem())
		dst.Set(v)

	default:
		dst.Set(src)
	}
}

func maybeCopy(src reflect.Value) reflect.Value {
	switch src.Kind() {
	case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Map:
		nv := reflect.New(src.Type()).Elem()
		ReflectClone(nv, src)
		return nv
	case reflect.Interface:
		return maybeCopy(src.Elem())
	default:
		return src
	}
}
