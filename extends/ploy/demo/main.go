package main

import (
	"github.com/yyzybb537/ketty/extends/ploy"
	"github.com/yyzybb537/ketty/extends/ploy/demo/test"
	"github.com/yyzybb537/ketty"
	//"github.com/yyzybb537/ketty/config"
	"golang.org/x/net/context"
	"fmt"
)

// implement grpc service
type TestApi struct {
	ploy.FlowI
}

type Request struct {
	Name		string
}

type Response struct {
	Message		string
}

func (this *TestApi) GetHandle() ketty.ServiceHandle {
	return test.TestHandle
}

func (this *TestApi) Ping(ctx context.Context, req *test.TestRequest) (rsp *test.TestResponse, err error) {
	iReq, err := this.reqTrans(req)
	if err != nil {
		return
	}

	iRsp := &Response{}
	ctx = this.Executor(ctx, iReq, iRsp)
	if ctx != nil {
		err = ctx.Err()
		if err != nil {
			return
		}
	}

	rsp, err = this.rspTrans(iRsp)
	return
}

func (this *TestApi) reqTrans(testReq *test.TestRequest) (req *Request, err error) {
	req = &Request{}
	req.Name = fmt.Sprintf("%d", testReq.Qr.QVal)
	// 协议转换
	return
}
func (this *TestApi) rspTrans(rsp *Response) (testRsp *test.TestResponse, err error) {
	testRsp = &test.TestResponse{}
	testRsp.Message = rsp.Message
	// 协议转换
	return
}

type Config struct {
	Listen string		// http://:8080
	TestRouter string  // /bidding/test
}

var gConfig Config

func main() {
	//err := config.Read(&gConfig, "server.toml")
	//ketty.Assert(err)
	gConfig.Listen = "http.query://:8000"
	gConfig.TestRouter = "/bidding/test"

	server, err := ploy.NewServer(gConfig.Listen, "")
	ketty.Assert(err)

	TA := &TestApi{}
	_, err = server.NewFlow(gConfig.TestRouter, TA)
	ketty.Assert(err)
	initBidding(TA.FlowI)

	err = server.Serve()
	ketty.Assert(err)

	ketty.Hung()
}

type RegionFilling struct {}
func (*RegionFilling) Run(ctx context.Context, req *Request, rsp *Response) context.Context {
	rsp.Message = req.Name
	var err error
	// todo something
	if err != nil {
		return ketty.WithError(ctx, err)
    }

	return ctx
}

type RegionFilling1 struct {}
func (*RegionFilling1) Run(ctx context.Context, req int, rsp *Response) context.Context {
	return ctx
}

type TraceTest struct {}
func (*TraceTest) PloyWillRun(ploy interface{}, ctx context.Context, req *Request, rsp *Response) context.Context {
	fmt.Printf("PloyWillRun %s\n",req.Name)
	return ctx
}
func (*TraceTest) PloyDidRun(ploy interface{}, ctx context.Context, req *Request, rsp *Response) context.Context {
	fmt.Printf("PloyDidRun %s\n",rsp.Message)
	return ctx
}

func initBidding(flow ploy.FlowI) {
	flow.AddTrace(new(TraceTest))
	// 填充地区并且做检索
	flow.AddPloy(new(RegionFilling))
	flow.AddPloy(new(RegionFilling1))
}

