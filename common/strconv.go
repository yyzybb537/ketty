package common

import (
	"reflect"
	"fmt"
	"strconv"
	"github.com/pkg/errors"
)

func V2String(v reflect.Value) (s string, err error) {
	switch v.Kind() {
	case reflect.Bool:
		s = strconv.FormatBool(v.Bool())

	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		s = strconv.FormatInt(v.Int(), 10)

	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		s = strconv.FormatUint(v.Uint(), 10)

	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		s = fmt.Sprintf("%v", v.Float())

	case reflect.String:
		s = v.String()

	case reflect.Uintptr:
		fallthrough
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		fallthrough
	case reflect.Array:
		fallthrough
	case reflect.Chan:
		fallthrough
	case reflect.Func:
		fallthrough
	case reflect.Interface:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Ptr:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Struct:
		fallthrough
	case reflect.UnsafePointer:
		err = errors.Errorf("V2String unsupport kind: %s", v.Kind().String())
	}
	return
}

func String2V(s string, v reflect.Value) (err error) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Elem().Type()))
		}
		v = v.Elem()
    }

	switch v.Kind() {
	case reflect.Bool:
		var b bool
		b, err = strconv.ParseBool(s)
		if err == nil {
			v.SetBool(b)
		}

	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		var i64 int64
		i64, err = strconv.ParseInt(s, 10, 64)
		if err == nil {
			v.SetInt(i64)
        }

	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		var u64 uint64
		u64, err = strconv.ParseUint(s, 10, 64)
		if err == nil {
			v.SetUint(u64)
        }

	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		var f64 float64
		f64, err = strconv.ParseFloat(s, 64)
		if err == nil {
			v.SetFloat(f64)
        }

	case reflect.String:
		v.SetString(s)

	case reflect.Uintptr:
		fallthrough
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		fallthrough
	case reflect.Array:
		fallthrough
	case reflect.Chan:
		fallthrough
	case reflect.Func:
		fallthrough
	case reflect.Interface:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Ptr:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Struct:
		fallthrough
	case reflect.UnsafePointer:
		err = errors.Errorf("Unsupport kind: %s", v.Kind().String())
	}
	return
}

