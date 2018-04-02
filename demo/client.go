package main

import (
	echo "github.com/yyzybb537/ketty/demo/pb"
	"github.com/yyzybb537/ketty"
	grpc "github.com/yyzybb537/ketty/protocol/grpc"
	"golang.org/x/net/context"
)

func main() {
	client, err := ketty.Dial("grpc://127.0.0.1:8090", "")
	if err != nil {
		ketty.GetLog().Errorf("Dial error:%s", err.Error())
		return
	}

	defer client.Close()

	opt := &grpc.GrpcOption{}
	opt.Option.TimeoutMilliseconds = 1000

	err = client.SetOption(opt)
	ketty.Assert(err)

	req := &echo.Req{ Val : 123 }
	stub := echo.NewKettyEchoServiceClient(client)
	rsp, err := stub.Echo(context.Background(), req)
	if err != nil {
		ketty.GetLog().Errorf("Invoke error:%+v", err)
		return
	}

	ketty.GetLog().Infof("Rsp: %d", rsp.Val)
}

