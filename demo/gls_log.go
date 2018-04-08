package main

import (
	echo "github.com/yyzybb537/ketty/demo/pb"
	"github.com/yyzybb537/ketty"
	"github.com/yyzybb537/ketty/log"
	"golang.org/x/net/context"
	"github.com/yyzybb537/gls"
)

var _ = gls.Get

type EchoServer struct {}
func (this *EchoServer) Echo(ctx context.Context, req *echo.Req) (*echo.Rsp, error) {
	ketty.GetLog().Debugf("Test Gls Log = %v", log.GetGlsDefaultKey())
	return &echo.Rsp{Val:req.Val}, nil
}

type logKey struct{}

func main() {
	log.SetGlsDefaultKey(logKey{})
	defer log.CleanupGlsDefaultKey()

	opt := log.DefaultLogOption()
	opt.LogCategory = "file"
	opt.OutputFile = "/dev/null"
	lg, err := log.MakeLogger(opt)
	ketty.Assert(err)
	log.SetLog(lg)

	opt = log.DefaultLogOption()
	opt.LogCategory = "std"
	log.BindOption(logKey{}, opt)

	server, err := ketty.Listen("grpc://127.0.0.1:8090", "")
	if err != nil {
		ketty.GetLog().Errorf("Listen error:%s", err.Error())
	}

	server.RegisterMethod(echo.EchoServiceHandle, &EchoServer{})

	err = server.Serve()
	if err != nil {
		ketty.GetLog().Errorf("Serve error:%s", err.Error())
	}

	ketty.Hung()
}
