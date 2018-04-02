package main

import (
	ext_pb "github.com/yyzybb537/ketty/demo/extpb"
	"github.com/yyzybb537/ketty"
	http "github.com/yyzybb537/ketty/protocol/http"
	"golang.org/x/net/context"
)

type EchoServer struct {
}

func (this *EchoServer) Echo(ctx context.Context, req *ext_pb.Req) (*ext_pb.Rsp, error) {
	return &ext_pb.Rsp{Val:req.Qr.QVal}, nil
}

func main() {
	server, err := ketty.Listen("http.json://127.0.0.1:8091", "")
	if err != nil {
		ketty.GetLog().Errorf("Listen error:%+v", err)
	}
	
	opt := &http.HttpOption{}
	opt.TimeoutMilliseconds = 100

	err = server.SetOption(opt)
	ketty.Assert(err)

	server.RegisterMethod(ext_pb.EchoServiceHandle, &EchoServer{})

	err = server.Serve()
	if err != nil {
		ketty.GetLog().Errorf("Serve error:%s", err.Error())
	}

	q := make(chan int)
	<-q
}

