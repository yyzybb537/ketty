package ketty

import (
	"golang.org/x/net/context"
	"github.com/yyzybb537/ketty/aop"
)

// client
type BeforeClientInvokeAop interface {
	BeforeClientInvoke(ctx context.Context, req interface{}) context.Context
}

type AfterClientInvokeAop interface {
	AfterClientInvoke(ctx context.Context, req, rsp interface{}, err error) context.Context
}

type ClientInvokeCleanupAop interface {
	ClientCleanup(ctx context.Context)
}

// server
type BeforeServerInvokeAop interface {
	BeforeServerInvoke(ctx context.Context, req interface{}) context.Context
}

type AfterServerInvokeAop interface {
	AfterServerInvoke(ctx context.Context, req, rsp interface{}, err error) context.Context
}

type ServerInvokeCleanupAop interface {
	ServerCleanup(ctx context.Context)
}

// transport metadata
type ClientTransportMetaDataAop interface {
	ClientSendMetaData(ctx context.Context, metadata map[string]string) context.Context
}

type ServerTransportMetaDataAop interface {
	ServerRecvMetaData(ctx context.Context, metadata map[string]string) context.Context
}

// list
type AopListI interface {
	AddAop(aop ... interface{})

	GetAop() []interface{}
}

type AopList struct {
	aopList []interface{}
}

func (this *AopList) AddAop(aop ... interface{}) {
	this.aopList = append(this.aopList, aop...)
}

func (this *AopList) GetAop() []interface{} {
	return this.aopList
}

var gAopList *AopList = new(AopList)

func DefaultAop() *AopList {
	return gAopList
}

func init() {
	DefaultAop().AddAop(new(aop.TraceAop))
	DefaultAop().AddAop(new(aop.CostAop))
	DefaultAop().AddAop(new(aop.LoggerAop))
}

func GetAop(ctx context.Context) []interface{} {
	aopList, ok := ctx.Value("aop").([]interface{})
	if !ok {
		return nil
	}

	return aopList
}

func SetAop(ctx context.Context, aopList []interface{}) context.Context {
	return context.WithValue(ctx, "aop", aopList)
}

