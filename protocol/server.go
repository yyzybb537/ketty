package protocol

import (
	A "github.com/yyzybb537/ketty/aop"
	COM "github.com/yyzybb537/ketty/common"
	O "github.com/yyzybb537/ketty/option"
)

type Server interface {
	A.AopListI

	SetOption(opt O.OptionI) error

	RegisterMethod(handle COM.ServiceHandle, implement interface{}) error

	Serve() error
}
