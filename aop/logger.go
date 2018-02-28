package aop

import (
	"golang.org/x/net/context"
	"github.com/yyzybb537/ketty/log"
)

type LoggerAop struct {
}

// client
func (this *LoggerAop) BeforeClientInvoke(ctx context.Context, req interface{}) context.Context {
	method := ctx.Value("method")
	remote := ctx.Value("remote")
	log.GetLog().Debugf("C-To %s Invoke (%s) TraceID:%s Req %s", remote, method, GetTraceID(), log.LogFormat(req))
	return ctx
}

func (this *LoggerAop) AfterClientInvoke(pCtx *context.Context, req, rsp interface{}) {
	ctx := *pCtx
	method := ctx.Value("method")
	remote := ctx.Value("remote")
	err := ctx.Err()
	if err != nil {
		log.GetLog().Errorf("C-From %s Reply (%s) TraceID:%s Cost:%s Error:%s", remote, method, GetTraceID(), getCostString(ctx), err.Error())
	} else {
		log.GetLog().Debugf("C-From %s Reply (%s) TraceID:%s Rsp:%s Cost:%s", remote, method, GetTraceID(), log.LogFormat(rsp), getCostString(ctx))
    }
}

// server
func (this *LoggerAop) BeforeServerInvoke(ctx context.Context, req interface{}) context.Context {
	method := ctx.Value("method")
	remote := ctx.Value("remote")
	log.GetLog().Debugf("From %s Invoke (%s) TraceID:%s Req:%s", remote, method, GetTraceID(), log.LogFormat(req))
	return ctx
}

func (this *LoggerAop) AfterServerInvoke(pCtx *context.Context, req, rsp interface{}) {
	ctx := *pCtx
	method := ctx.Value("method")
	remote := ctx.Value("remote")
	err := ctx.Err()
	if err != nil {
		log.GetLog().Errorf("To %s Reply (%s) TraceID:%s Error:%s", remote, method, GetTraceID(), err.Error())
	} else {
		log.GetLog().Debugf("To %s Reply (%s) TraceID:%s Rsp:%s Cost:%s", remote, method, GetTraceID(), log.LogFormat(rsp), getCostString(ctx))
    }
}

