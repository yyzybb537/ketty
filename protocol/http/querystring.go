package http_proto

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"reflect"
	"bytes"
	"strings"
	"fmt"
	stdurl "net/url"
	P "github.com/yyzybb537/ketty/protocol"
	COM "github.com/yyzybb537/ketty/common"
)

// ----------- querystring marshaler
type QueryStringMarshaler struct {}

func (this *QueryStringMarshaler) Marshal(msg proto.Message) ([]byte, error) {
	b := bytes.NewBufferString("")
	typ := reflect.TypeOf(msg).Elem()
	val := reflect.ValueOf(msg).Elem()
	numField := typ.NumField()
	for i := 0; i < numField; i++ {
		ftype := typ.Field(i).Type
		fvalue := val.Field(i)
		if ftype.Kind() == reflect.Ptr {
			ftype = ftype.Elem()
			fvalue = fvalue.Elem()
		}

		if ftype.Kind() == reflect.Struct {
			return nil, errors.Errorf("QueryString marshal not support multilayer message")
		}

		svalue, err := COM.V2String(fvalue)
		if err != nil {
			return nil, err
		}
		b.WriteString(fmt.Sprintf("%s=%s", stdurl.QueryEscape(this.getKey(typ.Field(i))), stdurl.QueryEscape(svalue)))

		if i + 1 < numField {
			b.WriteRune('&')
        }
    }

	return b.Bytes(), nil
}

func (this *QueryStringMarshaler) Unmarshal(buf []byte, msg proto.Message) error {
	kvMap, err := stdurl.ParseQuery(string(buf))
	if err != nil {
		return errors.WithStack(err)
	}

	typ := reflect.TypeOf(msg).Elem()
	val := reflect.ValueOf(msg).Elem()
	numField := typ.NumField()
	for i := 0; i < numField; i++ {
		ftype := typ.Field(i).Type
		fvalue := val.Field(i)

		if ftype.Kind() == reflect.Ptr {
			if ftype.Elem().Kind() == reflect.Struct {
				return errors.Errorf("QueryString marshal not support multilayer message")
			}
		}

		ss, exists := kvMap[this.getKey(typ.Field(i))]
		if !exists {
			continue
		}

		if len(ss) == 0 {
			continue
		}

		err = COM.String2V(ss[0], fvalue)
		if err != nil {
			return errors.WithStack(err)
        }
	}

	return nil
}

func (this *QueryStringMarshaler) getKey(sf reflect.StructField) string {
	k := sf.Tag.Get("json")
	if k == "" {
		k = sf.Name
	} else {
		k = strings.SplitN(k, ",", 2)[0]
    }
	return k
}

func init() {
	P.MgrMarshaler.Register("querystring", new(QueryStringMarshaler))
}

