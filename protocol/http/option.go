package http_proto

import (
	O "github.com/yyzybb537/ketty/option"
	"github.com/pkg/errors"
	"reflect"
)

type HttpOption struct {
	O.Option
	ConnectTimeoutMillseconds			int64
	ReadWriteTimeoutMillseconds			int64
	ResponseHeaderTimeoutMillseconds	int64
}

func defaultHttpOption() *HttpOption {
	return new(HttpOption)
}

func (this *HttpOption) set(opt O.OptionI) error {
	if o, ok := opt.(*HttpOption); ok {
		*this = *o
	} else if o, ok := opt.(HttpOption); ok {
		*this = o
	} else if o, ok := opt.(*O.Option); ok {
		this.Option = *o
	} else if o, ok := opt.(O.Option); ok {
		this.Option = o
	} else {
		typ := reflect.TypeOf(opt)
		return errors.Errorf("SetOption argument error. opt={isptr:%t, className:%s}", typ.Kind() == reflect.Ptr, typ.String())
	}
	return nil
}
