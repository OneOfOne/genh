package genh

import (
	"log"
	"reflect"
)

type Cloner[T any] interface {
	Clone() T
}

func TypedClone[T any](v T) (cp T) {
	if v, ok := any(v).(Cloner[T]); ok {
		return v.Clone()
	}
	src, dst := indirect(reflect.ValueOf(v)), reflect.ValueOf(&cp).Elem()
	ReflectCopy(src, dst)
	return
}

func ReflectCopy(src, dst reflect.Value) {
	if !src.IsValid() || src.IsZero() {
		return
	}

	styp := src.Type()

	if dst.Kind() == reflect.Ptr && dst.IsNil() {
		dst.Set(reflect.New(src.Type()))
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
			ReflectCopy(src.Index(i), dst.Index(i))
		}
	case reflect.Map:
		if src.IsNil() {
			return
		}

		dst.Set(reflect.MakeMapWithSize(styp, src.Len()))
		kt, vt := styp.Key(), styp.Elem()
		for it := src.MapRange(); it.Next(); {
			k, v := reflect.New(kt).Elem(), reflect.New(vt).Elem()
			ReflectCopy(it.Key(), k)
			ReflectCopy(it.Value(), v)
			dst.SetMapIndex(k, v)
		}
	case reflect.Struct:
		for i := 0; i < styp.NumField(); i++ {
			ReflectCopy(indirect(src.Field(i)), dst.Field(i))
		}
	case reflect.Ptr, reflect.Interface:
		if src.IsNil() {
			return
		}
		v := reflect.New(src.Elem().Type()).Elem()
		ReflectCopy(src.Elem(), v)
		dst.Set(v)
	default:
		dst.Set(src)
	}
}

func indirect(rv reflect.Value) reflect.Value {
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	return rv
}
