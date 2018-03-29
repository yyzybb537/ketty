package protocol

import (
	"golang.org/x/net/context"
	A "github.com/yyzybb537/ketty/aop"
	O "github.com/yyzybb537/ketty/option"
	COM "github.com/yyzybb537/ketty/common"
)

type Client interface {
	A.AopListI

	SetOption(opt O.OptionI) error

	Invoke(ctx context.Context, handle COM.ServiceHandle, method string, req, rsp interface{}) error

	Close()
}

