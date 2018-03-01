package protocol

import (
	"github.com/yyzybb537/ketty/aop"
	COM "github.com/yyzybb537/ketty/common"
)

type Server interface {
	aop.AopListI

	RegisterMethod(handle COM.ServiceHandle, implement interface{}) error

	Serve() error
}
