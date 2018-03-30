package ploy

import (
	"golang.org/x/net/context"
	"github.com/yyzybb537/ketty/extends/ploy/fake_interface"
	"github.com/yyzybb537/ketty"
	"fmt"
)

type FlowI interface{
	Executor(context.Context, interface{}, interface{}) context.Context
	AddPloy(interface{})
	AddTrace(interface{})
}

type ployReqKey struct{}
type ployRespKey struct{}

func GetRequest(ctx context.Context) interface{} {
	return ctx.Value(ployReqKey{})	
}

func GetResponse(ctx context.Context) interface{} {
	return ctx.Value(ployRespKey{})
}

func putRequest(ctx context.Context, request interface{}) context.Context{
	return context.WithValue(ctx, ployReqKey{}, request)	
}

func putResponse(ctx context.Context, response interface{}) context.Context{
	return context.WithValue(ctx, ployRespKey{}, response)	
}
/*
type Ploy interface{
	Run(context.Context, any, any) context.Context
}
type Trace interface{
	TraceBefore(ploy interface{}, ctx context.Context, any, any)
	TraceAfter(ploy interface{}, ctx context.Context, any, any)
}
*/

type BaseFlow struct {
	ployFIs		[]*fake_interface.FakeInterface
	traceFI		*fake_interface.FakeInterface
}

func NewBaseFlow() FlowI {
	baseFlow := &BaseFlow{}
	traceFI := fake_interface.NewFakeInterface()
	traceFI.Add("PloyWillRun", 4)
	traceFI.Add("PloyDidRun", 4)

	baseFlow.traceFI = traceFI

	return baseFlow
}

func (this *BaseFlow) Executor(ctx context.Context, req interface{},resp interface{}) context.Context{
	ctx = putRequest(ctx, req)
	ctx = putResponse(ctx, resp)
	for _, ployFI := range this.ployFIs {
		ctx = this.traceFI.Do("PloyWillRun", ployFI.Interface(), ctx, req, resp)
		ctx = ployFI.Do("Run", ctx, req, resp)
		ctx = this.traceFI.Do("PloyDidRun", ployFI.Interface(), ctx, req, resp)
		if ctx != nil && ctx.Err() != nil {
			return ctx
		}
	}
	return ctx
}

type _init interface {
	Init() 
}

func (this *BaseFlow) AddPloy(imp interface{}) {
	ployFI := fake_interface.NewFakeInterface()
	ployFI.Add("Run", 3)
	err := ployFI.Realize(imp)
	ketty.Assert(err)
	if len(this.ployFIs) != 0 {
		if !ployFI.LookLike(this.ployFIs[0]) {
			ketty.Assert(fmt.Errorf("one new ploy %s don't LookLike the first ploy %s", ployFI.RealizedTypeName(), this.ployFIs[0].RealizedTypeName()))
		}
	}
	i, ok := imp.(_init)
	if ok {
		i.Init()
	}
	this.ployFIs = append(this.ployFIs, ployFI)
	return
}

func (this *BaseFlow) AddTrace(imp interface{}) {
	err := this.traceFI.Realize(imp)
	ketty.Assert(err)
}
