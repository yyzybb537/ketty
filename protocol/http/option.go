package http_proto

import (
	O "github.com/yyzybb537/ketty/option"
	"github.com/pkg/errors"
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
	} else if o, ok := opt.(*O.Option); ok {
		this.Option = *o
	} else {
		return errors.New("SetOption argument error")
	}
	return nil
}
