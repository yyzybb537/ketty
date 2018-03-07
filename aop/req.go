package aop

import (
	"golang.org/x/net/context"
)

type ReqRspAop struct {
}

type requestKey struct {}
type responseKey struct {}

// client
func (this *ReqRspAop) BeforeClientInvoke(ctx context.Context, req interface{}) context.Context {
	return context.WithValue(ctx, requestKey{}, req)
}

func (this *ReqRspAop) AfterClientInvoke(pCtx *context.Context, req, rsp interface{}) {
	*pCtx = context.WithValue(*pCtx, responseKey{}, rsp)
}

// server
func (this *ReqRspAop) BeforeServerInvoke(ctx context.Context, req interface{}) context.Context {
	return context.WithValue(ctx, requestKey{}, req)
}

func (this *ReqRspAop) AfterServerInvoke(pCtx *context.Context, req, rsp interface{}) {
	*pCtx = context.WithValue(*pCtx, responseKey{}, rsp)
}

func GetRequest(ctx context.Context) interface{} {
	return ctx.Value(requestKey{})
}

func GetResponse(ctx context.Context) interface{} {
	return ctx.Value(responseKey{})
}

