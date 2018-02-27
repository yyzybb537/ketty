package main

import (
	echo "github.com/yyzybb537/ketty/test/test_pb"
	"github.com/yyzybb537/ketty"
	kettyLog "github.com/yyzybb537/ketty/log"
	_ "github.com/yyzybb537/ketty/protocol/grpc"
	kettyHttp "github.com/yyzybb537/ketty/protocol/http"
	context "golang.org/x/net/context"
)

type EchoServer struct {
}

func (this *EchoServer) Echo(ctx context.Context, req *echo.Req) (*echo.Rsp, error) {
	return &echo.Rsp{Val:req.Val}, nil
}


func startServer(sUrl string) {
	server, err := ketty.Listen(sUrl, "")
	if err != nil {
		ketty.GetLog().Errorf("Listen error:%s", err.Error())
	}

	server.RegisterMethod(echo.EchoServiceHandle, &EchoServer{})

	err = server.Serve()
	if err != nil {
		ketty.GetLog().Errorf("Serve error:%s", err.Error())
	}
}

var httpUrl = "http://127.0.0.1:8091"
var httpsUrl = "https://127.0.0.1:3050"
var grpcUrl = "grpc://127.0.0.1:8090"
var urls = []string{
	//"http://127.0.0.1",
	httpUrl,
	httpsUrl,
	grpcUrl,
}

func main() {
	kettyHttp.InitTLS("cert.pem", "key.pem")
	for _, sUrl := range urls {
		ketty.GetLog().Infof("Listen url:%s", sUrl)
		startServer(sUrl)
	}
	ketty.SetLog(new(kettyLog.FakeLog))
	q := make(chan int)
	<-q
}

