package ketty

import (
	"time"
	"golang.org/x/net/context"
	C "github.com/yyzybb537/ketty/context"
	COM "github.com/yyzybb537/ketty/common"
	A "github.com/yyzybb537/ketty/aop"
)

type ServiceHandle COM.ServiceHandle

func Assert(err error) {
	if err != nil {
		panic(err)
    }
}

func Hung() {
	time.Sleep(time.Second)
	GetLog().Infof("ketty service startup")
	for {
		time.Sleep(time.Second * 3600)
    }
}

func WithError(ctx context.Context, err error) context.Context {
	return C.WithError(ctx, err)
}

func GetTraceID() string {
	return A.GetTraceID()
}
