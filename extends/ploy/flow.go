package ploy

import (
	"golang.org/x/net/context"
	"github.com/yyzybb537/ketty/extends/ploy/fake_interface"
	"github.com/yyzybb537/ketty"
)

type FlowI interface{
	Executor(context.Context, interface{}, interface{}) context.Context
	AddPloy(interface{})
	AddTrace(interface{})
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
	for _, ployFI := range this.ployFIs {
		ctx = this.traceFI.Do("PloyWillRun", ployFI.Interface(), ctx, req, resp)
		ctx = ployFI.Do("Run", ctx, req, resp)
		ctx = this.traceFI.Do("PloyDidRun", ployFI.Interface(), ctx, req, resp)
	}
	return ctx
}

func (this *BaseFlow) AddPloy(imp interface{}) {
	ployFI := fake_interface.NewFakeInterface()
	ployFI.Add("Run", 3)
	err := ployFI.Realize(imp)
	ketty.Assert(err)
	this.ployFIs = append(this.ployFIs, ployFI)
	return
}

func (this *BaseFlow) AddTrace(imp interface{}) {
	err := this.traceFI.Realize(imp)
	ketty.Assert(err)
}
