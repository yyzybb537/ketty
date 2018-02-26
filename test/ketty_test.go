package grpc_t

import (
	"testing"
	echo "github.com/yyzybb537/ketty/test/test_pb"
	"github.com/yyzybb537/ketty"
	_ "github.com/yyzybb537/ketty/protocol/grpc"
	kettyHttp "github.com/yyzybb537/ketty/protocol/http"
	"time"
	context "golang.org/x/net/context"
)

type EchoServer struct {
}

func (this *EchoServer) Echo(ctx context.Context, req *echo.Req) (*echo.Rsp, error) {
	return &echo.Rsp{Val:req.Val}, nil
}


func startServer(t *testing.T, sUrl string) {
	server, err := ketty.Listen(sUrl, "")
	if err != nil {
		t.Fatalf("Listen error:%s", err.Error())
	}

	server.RegisterMethod(echo.EchoServiceHandle, &EchoServer{})

	err = server.Serve()
	if err != nil {
		t.Fatalf("Serve error:%s", err.Error())
	}
}

func startClient(t *testing.T, sUrl string) {
	client, err := ketty.Dial(sUrl, "")
	if err != nil {
		t.Fatalf("Dial error:%s", err.Error())
	}

	req := &echo.Req{ Val : 123 }
	stub := echo.NewKettyEchoServiceClient(client)
	rsp, err := stub.Echo(context.Background(), req)
	if err != nil {
		t.Fatalf("Invoke error:%+v", err)
	}

	t.Logf("Echo Val:%d", rsp.Val)
}

func TestGrpc(t *testing.T) {
	urls := []string{
		//"http://127.0.0.1",
		"http://127.0.0.1:8091",
		"https://127.0.0.1:3050",
		"grpc://127.0.0.1:8090",
	}
	kettyHttp.InitTLS("cert.pem", "key.pem")
	for _, sUrl := range urls {
		t.Logf("Do url:%s", sUrl)
		startServer(t, sUrl)
		time.Sleep(time.Millisecond * 100)
		startClient(t, sUrl)
	}
}

