package grpc_proto

import (
	O "github.com/yyzybb537/ketty/option"
	"github.com/pkg/errors"
	"reflect"
)

type GrpcOption struct {
	O.Option
}

func defaultGrpcOption() *GrpcOption {
	return new(GrpcOption)
}

func (this *GrpcOption) set(opt O.OptionI) error {
	if o, ok := opt.(*GrpcOption); ok {
		*this = *o
	} else if o, ok := opt.(GrpcOption); ok {
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
