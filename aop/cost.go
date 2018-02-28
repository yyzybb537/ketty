package aop

import (
	"golang.org/x/net/context"
	"time"
)

type CostAop struct {
}

// client
func (this *CostAop) BeforeClientInvoke(ctx context.Context, req interface{}) context.Context {
	return context.WithValue(ctx, "begin", time.Now())
}

func (this *CostAop) AfterClientInvoke(pCtx *context.Context, req, rsp interface{}) {
	ctx := *pCtx
	begin := ctx.Value("begin").(time.Time)
	cost := time.Since(begin)
	*pCtx = context.WithValue(ctx, "cost", cost)
}

// server
func (this *CostAop) BeforeServerInvoke(ctx context.Context, req interface{}) context.Context {
	return context.WithValue(ctx, "begin", time.Now())
}

func (this *CostAop) AfterServerInvoke(pCtx *context.Context, req, rsp interface{}) {
	ctx := *pCtx
	begin := ctx.Value("begin").(time.Time)
	cost := time.Since(begin)
	*pCtx = context.WithValue(ctx, "cost", cost)
}

func getCost(ctx context.Context) time.Duration {
	dur, ok := ctx.Value("cost").(time.Duration)
	if !ok {
		return time.Duration(0)
	}
	return dur
}

func getCostString(ctx context.Context) string {
	return getCost(ctx).String()
}
