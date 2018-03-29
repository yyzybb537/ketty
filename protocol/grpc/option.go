package grpc_proto

import (
	O "github.com/yyzybb537/ketty/option"
	"github.com/pkg/errors"
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
	} else if o, ok := opt.(*O.Option); ok {
		this.Option = *o
	} else {
		return errors.New("SetOption argument error")
	}
	return nil
}
