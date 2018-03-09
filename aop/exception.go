package aop

import (
	"golang.org/x/net/context"
	"github.com/pkg/errors"
	C "github.com/yyzybb537/ketty/context"
	"github.com/yyzybb537/ketty/log"
	"runtime/debug"
)

type ExceptionAop struct {
}

func (this *ExceptionAop) AfterServerInvoke(pCtx *context.Context, req, rsp interface{}) {
	iErr := recover()
	if iErr == nil {
		return
	}

	log.GetLog().Errorf("Ketty.Exception: %+v\nStack: %s", iErr, string(debug.Stack()))

	err, ok := iErr.(error)
	if !ok {
		err = errors.Errorf("%v", iErr)
	}

	err = errors.WithStack(err)
	*pCtx = C.WithError(*pCtx, err)
}

