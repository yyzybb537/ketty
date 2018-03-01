package protocol

import (
	"golang.org/x/net/context"
	A "github.com/yyzybb537/ketty/aop"
	COM "github.com/yyzybb537/ketty/common"
)

type Client interface {
	A.AopListI

	Close()

	Invoke(ctx context.Context, handle COM.ServiceHandle, method string, req, rsp interface{}) error
}

