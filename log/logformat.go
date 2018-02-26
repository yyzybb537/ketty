package log

import (
	"fmt"
	"reflect"
	"strings"
	"bytes"
)

const bytes_limit = 16
const string_limit = 128

// Provides a hook to intercept the execution of logformat.
type LogFormatInterceptor func(name string, v reflect.Value, handler func(reflect.Value) string) string

const KindCount = reflect.UnsafePointer

type LogFormatOptions struct {
	Indent bool
	NoLimit bool
	Interceptors [KindCount]LogFormatInterceptor

	indentNum int
}

type LogFormatOption func(*LogFormatOptions)

var EmptyLogFormatOptions LogFormatOptions = LogFormatOptions{}
var IndentLogFormatOptions LogFormatOptions = LogFormatOptions{ Indent : true }
var Indent LogFormatOptions = IndentLogFormatOptions
var NoLimit LogFormatOptions = LogFormatOptions{ NoLimit : true }

func IndentLogFormatOpt(opt *LogFormatOptions) {
	opt.Indent = true
}

func formatSlice(buf *bytes.Buffer, v reflect.Value, opt *LogFormatOptions) {
	if v.Len() == 0 {
		buf.WriteRune('[')
		buf.WriteRune(']')
		return 
    }

	if v.Index(0).Kind() == reflect.Uint8 {
		if opt.NoLimit {
			formatHex(buf, v.Bytes())
		} else {
			limitBytes(buf, v.Bytes())
		}
		return 
    }

	buf.WriteRune('[')
	for i := 0; i < v.Len(); i++ {
		e := v.Index(i)
		formatV(buf, "", e, opt)
		if i + 1 < v.Len() {
			buf.WriteRune(' ')
		}
	}
	buf.WriteRune(']')
	return
}

func formatArray(buf *bytes.Buffer, v reflect.Value, opt *LogFormatOptions) {
	buf.WriteRune('[')
	for i := 0; i < v.Len(); i++ {
		e := v.Index(i)
		formatV(buf, "", e, opt)
		if i + 1 < v.Len() {
			buf.WriteRune(' ')
		}
	}
	buf.WriteRune(']')
}

func formatStruct(buf *bytes.Buffer, struct_v reflect.Value, opt *LogFormatOptions) {
	typ := struct_v.Type()
	n_field := struct_v.NumField()
	for i := 0; i < n_field; i++ {
		v := struct_v.Field(i)
		field := typ.Field(i)

		if opt.Indent {
			buf.WriteRune('\n')
			buf.WriteString(strings.Repeat("  ", opt.indentNum))
		}

		buf.WriteString(field.Name)
		buf.WriteRune(':')
		opt.indentNum++
		formatV(buf, field.Name, v, opt) 
		opt.indentNum--

		if opt.Indent {
			buf.WriteRune(',')
        } else {
			if i + 1 < n_field {
				buf.WriteRune(',')
				buf.WriteRune(' ')
			}
        }
	}
}

func formatMap(buf *bytes.Buffer, map_v reflect.Value, opt *LogFormatOptions) {
	keys := map_v.MapKeys()
	var i int
	for _, key := range keys {
		value := map_v.MapIndex(key)

		if opt.Indent {
			buf.WriteRune('\n')
			buf.WriteString(strings.Repeat("  ", opt.indentNum))
		}

		formatV(buf, "", key, opt)
		buf.WriteRune(':')
		opt.indentNum++
		formatV(buf, "", value, opt) 
		opt.indentNum--
		if opt.Indent {
			buf.WriteRune(',')
        } else {
			if i + 1 < len(keys) {
				buf.WriteRune(',')
				buf.WriteRune(' ')
			}
        }
		i++
	}
}

func limitString(buf *bytes.Buffer, v string) {
	if len(v) > string_limit {
		buf.WriteRune('"')
		buf.WriteString(v[:string_limit])
		buf.WriteRune('"')
		buf.WriteString(fmt.Sprintf("(%d)", len(v)))
    } else {
		buf.WriteString(v)
    }
}

func formatHex(buf *bytes.Buffer, v []byte) {
	buf.WriteRune('[')
	for i, b := range v {
		buf.WriteString(fmt.Sprintf("%02x", b))
		if i + 1 < len(v) {
			buf.WriteRune(' ')
		}
	}
	buf.WriteRune(']')
	return 
}

func limitBytes(buf *bytes.Buffer, v []byte) {
	if len(v) > bytes_limit {
		formatHex(buf, v[:bytes_limit])
		buf.WriteString(fmt.Sprintf("(%d)", len(v)))
    } else {
		formatHex(buf, v)
    }
}

func callMethod(method reflect.Value) (ret []reflect.Value) {
	defer func(){
		if recover() != nil {
			return 
        }
	}()
	ret = method.Call(nil)
	return
}

func formatV(buf *bytes.Buffer, name string, v reflect.Value, opt *LogFormatOptions) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	handler := func(v reflect.Value) string {
		switch v.Kind() {
		case reflect.Struct:
			if toStringMethod := v.MethodByName("String"); toStringMethod.IsValid() {
				rets := callMethod(toStringMethod)
				if len(rets) > 0 {
					ret_v := rets[0]
					if ret_v.Kind() == reflect.String {
						buf.WriteString(ret_v.String())
						break 
					}
				}
            }

			if opt.Indent {
				buf.WriteRune('\n')
				buf.WriteString(strings.Repeat("  ", opt.indentNum))
				buf.WriteRune('{')
				opt.indentNum++
				formatStruct(buf, v, opt) 
				opt.indentNum--
				buf.WriteRune('\n')
				buf.WriteString(strings.Repeat("  ", opt.indentNum))
				buf.WriteRune('}')
            } else {
				buf.WriteRune('{')
				formatStruct(buf, v, opt) 
				buf.WriteRune('}')
            }
		case reflect.Map:
			if opt.Indent {
				buf.WriteRune('\n')
				buf.WriteString(strings.Repeat("  ", opt.indentNum))
				buf.WriteRune('{')
				opt.indentNum++
				formatMap(buf, v, opt) 
				opt.indentNum--
				buf.WriteRune('\n')
				buf.WriteString(strings.Repeat("  ", opt.indentNum))
				buf.WriteRune('}')
            } else {
				buf.WriteRune('{')
				formatMap(buf, v, opt) 
				buf.WriteRune('}')
            }
		case reflect.Slice:
			formatSlice(buf, v, opt)
		case reflect.String:
			if opt.NoLimit {
				buf.WriteString(v.String())
			} else {
				limitString(buf, v.String())
			}
		case reflect.Array:
			formatArray(buf, v, opt)
		case reflect.Bool:
			buf.WriteString(fmt.Sprintf("%v", v.Bool()))
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			buf.WriteString(fmt.Sprintf("%v", v.Int()))
		case reflect.Uint:
			fallthrough
		case reflect.Uint8:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Uint64:
			buf.WriteString(fmt.Sprintf("%v", v.Uint()))
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			buf.WriteString(fmt.Sprintf("%v", v.Float()))
		}

		return ""
	}

	if opt.Interceptors[v.Kind()] != nil {
		buf.WriteString(opt.Interceptors[v.Kind()](name, v, handler))
    } else {
		handler(v)
    }
}

func LogFormat(req interface{}, opt ... LogFormatOptions) string {
	var o *LogFormatOptions = &EmptyLogFormatOptions
	if len(opt) > 0 {
		o = &opt[0]
	}
	buf := bytes.NewBufferString("")
	formatV(buf, "", reflect.ValueOf(req), o)
	return buf.String()
}

/*
type Base struct {
	x []byte
}

type A struct {
	I int
	I64 int64
	S string
	B []byte
	X Base
}


func main() {
	a := A{ B : []byte(" \\0wsfweifjwoeifjweofjwoefiewjfoiwefweosjidfowoefjweoiofwejfiewjfowejfewfowiejfoiwejofjwoe") }
	a.B[0] = 0
	a.S = "ccwjoeifojweifhowejofwejofwejfowejfooewfjwefjojfiwjowfojfwjwjffejowjoeifojweifhowejofwejofwejfowejfooewfjwefjojfiwjowfojfwjwjffejo"
	fmt.Printf("A:%s", LogFormat(a))
}
*/

