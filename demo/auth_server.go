package main

import (
	echo "github.com/yyzybb537/ketty/demo/pb"
	"github.com/yyzybb537/ketty"
	"github.com/yyzybb537/ketty/aop"
	"golang.org/x/net/context"
	"fmt"
)

type EchoServer struct {}
func (this *EchoServer) Echo(ctx context.Context, req *echo.Req) (*echo.Rsp, error) {
	return &echo.Rsp{Val:req.Val}, nil
}

type Auth struct {}
func (this *Auth) CheckAuthorization(ctx context.Context, authorization string) error {
	if authorization != "MyAuthInfo" {
		return fmt.Errorf("AuthError")
    }
	return nil
}

func main() {
	server, err := ketty.Listen("grpc://127.0.0.1:8090", "")
	if err != nil {
		ketty.GetLog().Errorf("Listen error:%s", err.Error())
	}

	server.AddAop(aop.NewAuthAop(nil, new(Auth)))
	server.RegisterMethod(echo.EchoServiceHandle, &EchoServer{})

	err = server.Serve()
	if err != nil {
		ketty.GetLog().Errorf("Serve error:%s", err.Error())
	}

	q := make(chan int)
	<-q
}

