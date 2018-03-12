package aop

import (
	"golang.org/x/net/context"
	uuid "github.com/satori/go.uuid"
	"github.com/yyzybb537/gls"
)

type TraceAop struct {
}

// client
func (this *TraceAop) BeforeClientInvoke(ctx context.Context, req interface{}) context.Context {
	genTraceID()
	return ctx
}

// server
func (this *TraceAop) BeforeServerInvoke(ctx context.Context, req interface{}) context.Context {
	genTraceID()
	return ctx
}

func (this *TraceAop) ClientCleanup(ctx context.Context) {
	cleanupTraceID(ctx)
}

func (this *TraceAop) ServerCleanup(ctx context.Context) {
	cleanupTraceID(ctx)
}

func (this *TraceAop) ClientSendMetaData(ctx context.Context, metadata map[string]string) context.Context {
	metadata["traceid"] = genTraceID()
	return ctx
}

func (this *TraceAop) ServerRecvMetaData(ctx context.Context, metadata map[string]string) context.Context {
	traceId, exists := metadata["traceid"]
	if exists {
		gls.Set(traceIdKey{}, traceId)
	}
	return ctx
}

type traceIdKey struct {}

func genTraceID() string {
	traceId, ok := gls.Get(traceIdKey{}).(string)
	if !ok || traceId == "" {
		uuid, _ := uuid.NewV4()
		traceId = uuid.String()
		gls.Set(traceIdKey{}, traceId)
	}
	return traceId
}

func GetTraceID() string {
	s, ok := gls.Get(traceIdKey{}).(string)
	if !ok {
		return ""
	}
	return s
}

func cleanupTraceID(ctx context.Context) {
	gls.Del(traceIdKey{})
}

