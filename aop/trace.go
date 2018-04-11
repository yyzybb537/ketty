package aop

import (
	"golang.org/x/net/context"
	uuid "github.com/satori/go.uuid"
	"github.com/yyzybb537/gls"
//	"github.com/yyzybb537/ketty/log"
)

type TraceAop struct {
}

type createdTraceIdKey struct {}

// client
func (this *TraceAop) BeforeClientInvoke(ctx context.Context, req interface{}) context.Context {
	traceId, created := genTraceID()
	if created {
		ctx = context.WithValue(ctx, createdTraceIdKey{}, traceId)
	}
	return ctx
}

// server
func (this *TraceAop) BeforeServerInvoke(ctx context.Context, req interface{}) context.Context {
	traceId, created := genTraceID()
	if created {
		ctx = context.WithValue(ctx, createdTraceIdKey{}, traceId)
	}
	return ctx
}

func (this *TraceAop) ClientCleanup(ctx context.Context) {
	cleanupTraceID(ctx)
}

func (this *TraceAop) ServerCleanup(ctx context.Context) {
	cleanupTraceID(ctx)
}

func (this *TraceAop) ClientSendMetaData(ctx context.Context, metadata map[string]string) context.Context {
	traceId, created := genTraceID()
	if created {
		ctx = context.WithValue(ctx, createdTraceIdKey{}, traceId)
	}
	metadata["traceid"] = traceId
	return ctx
}

func (this *TraceAop) ServerRecvMetaData(ctx context.Context, metadata map[string]string) context.Context {
	traceId, exists := metadata["traceid"]
	if exists {
//		log.GetLog().Debugf("set trace id:%s", traceId)
		gls.Set(traceIdKey{}, traceId)
		ctx = context.WithValue(ctx, createdTraceIdKey{}, traceId)
	}
	return ctx
}

type traceIdKey struct {}

func genTraceID() (string, bool) {
	traceId, ok := gls.Get(traceIdKey{}).(string)
	if !ok || traceId == "" {
		uuid, _ := uuid.NewV4()
		traceId = uuid.String()
		gls.Set(traceIdKey{}, traceId)
		log.GetLog().Debugf("create trace id:%s", traceId)
		return traceId, true
	}
	return traceId, false
}

func GetTraceID() string {
	s, ok := gls.Get(traceIdKey{}).(string)
	if !ok {
		return ""
	}
	return s
}

func cleanupTraceID(ctx context.Context) {
	ctxTraceId, exists := ctx.Value(createdTraceIdKey{}).(string)
	if exists {
		traceId, exists := gls.Get(traceIdKey{}).(string)
		if exists && ctxTraceId == traceId {
			gls.Del(traceIdKey{})
		}
	}
}

