package aop

import (
	"golang.org/x/net/context"
	"github.com/pkg/errors"
	C "github.com/yyzybb537/ketty/context"
)

type ExceptionAop struct {
}

func (this *ExceptionAop) AfterServerInvoke(pCtx *context.Context, req, rsp interface{}) {
	iErr := recover()
	if iErr == nil {
		return
	}

	err, ok := iErr.(error)
	if !ok {
		err = errors.Errorf("%v", iErr)
	} else {
		err = errors.WithStack(err)
	}
	*pCtx = C.WithError(*pCtx, err)
	//log.GetLog().Infof("exception.AfterServerInvoke err=%v", err)
	//log.GetLog().Infof("exception.AfterServerInvoke err=%v", (*pCtx).Err())
}

