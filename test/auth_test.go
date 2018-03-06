package grpc_t

import (
	"testing"
	echo "github.com/yyzybb537/ketty/test/test_pb"
	"github.com/yyzybb537/ketty"
	kettyLog "github.com/yyzybb537/ketty/log"
	P "github.com/yyzybb537/ketty/protocol"
	kettyHttp "github.com/yyzybb537/ketty/protocol/http"
	A "github.com/yyzybb537/ketty/aop"
	"time"
	"fmt"
	context "golang.org/x/net/context"
)

type EchoServer struct {}
func (this *EchoServer) Echo(ctx context.Context, req *echo.Req) (*echo.Rsp, error) {
	return &echo.Rsp{Val:req.Val}, nil
}

type Auth struct {}
func (this *Auth) CreateAuthorization(ctx context.Context) string {
	return "MyAuth"
}
func (this *Auth) CheckAuthorization(ctx context.Context, authorization string) error {
	if authorization != "MyAuth" {
		return fmt.Errorf("Auth error")
    }
	return nil
}

var gAop *A.AopList

func startServer(t *testing.T, sUrl string, driverUrl string) {
	server, err := ketty.Listen(sUrl, driverUrl)
	if err != nil {
		t.Fatalf("Listen error:%s", err.Error())
	}

	server.RegisterMethod(echo.EchoServiceHandle, &EchoServer{})
	server.AddAop(gAop.GetAop()...)

	err = server.Serve()
	if err != nil {
		t.Fatalf("Serve (%s) error:%s", sUrl, err.Error())
	}
}

func startClient(t *testing.T, sUrl string, exceptError bool) {
	client, err := ketty.Dial(sUrl, "")
	if err != nil {
		t.Fatalf("Dial error:%s", err.Error())
	}
	client.AddAop(gAop.GetAop()...)

	req := &echo.Req{ Val : 123 }
	stub := echo.NewKettyEchoServiceClient(client)
	rsp, err := stub.Echo(context.Background(), req)
	if exceptError {
		if err == nil {
			t.Fatalf("Invoke not error")
        } else {
			t.Logf("Except error:%v", err)
        }
	} else {
		if err != nil {
			t.Fatalf("Invoke error:%+v", err)
		}

		t.Logf("Echo Val:%d", rsp.Val)
	}
}

var httpUrl = "http://127.0.0.1:8091"
var httpjsonUrl = "httpjson://127.0.0.1:8092"
var httpsUrl = "https://127.0.0.1:3050"
var httpsjsonUrl = "httpsjson://127.0.0.1:3051"
var grpcUrl = "grpc://127.0.0.1:8090"
var etcdUrl = "grpc://127.0.0.1:8090"
var urls = []string{
	httpUrl,
	httpsUrl,
	grpcUrl,
	httpjsonUrl,
	httpsjsonUrl,
}

func TestGrpc(t *testing.T) {
	kettyHttp.InitTLS("cert.pem", "key.pem")
	ketty.SetLog(new(kettyLog.StdLog))
	ketty.GetLog().Debugf("--------------- Protocols:%s", P.DumpProtocols())
	for _, sUrl := range urls {
		testByUrl(t, sUrl)
	}
}

func testByUrl(t *testing.T, sUrl string) {
	t.Logf("Do url:%s", sUrl)
	// 1.all auth
	auth := new(Auth)
	gAop = new(A.AopList)
	gAop.AddAop(A.NewAuthAop(auth, auth))

	startServer(t, sUrl, "")
	time.Sleep(time.Millisecond * 100)
	startClient(t, sUrl, false)

	// 2.Server auth only
	gAop = new(A.AopList)
	time.Sleep(time.Millisecond * 100)
	startClient(t, sUrl, true)
}
