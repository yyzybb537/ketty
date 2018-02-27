package grpc_t

import (
	"testing"
	echo "github.com/yyzybb537/ketty/test/test_pb"
	"github.com/yyzybb537/ketty"
	kettyLog "github.com/yyzybb537/ketty/log"
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

var clientMap = map[string]ketty.Client{}

func bStartClient(b *testing.B, sUrl string) {
	var err error
	client, exists := clientMap[sUrl]
	if !exists {
		client, err = ketty.Dial(sUrl, "")
		if err != nil {
			b.Fatalf("Dial error:%s", err.Error())
		}
		clientMap[sUrl] = client
	}

	for i := 0; i < b.N; i++ {
		req := &echo.Req{ Val : 123 }
		stub := echo.NewKettyEchoServiceClient(client)
		_, err := stub.Echo(context.Background(), req)
		if err != nil {
			b.Fatalf("Invoke error:%+v", err)
		}
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

func TestGrpc(t *testing.T) {
	kettyHttp.InitTLS("cert.pem", "key.pem")
	ketty.SetLog(new(kettyLog.StdLog))
	for _, sUrl := range urls {
		t.Logf("Do url:%s", sUrl)
		startServer(t, sUrl)
		time.Sleep(time.Millisecond * 100)
		startClient(t, sUrl)
	}
}

func Benchmark_Grpc(b *testing.B) {
	ketty.SetLog(new(kettyLog.FakeLog))
	bStartClient(b, grpcUrl)
}

//func Benchmark_Http(b *testing.B) {
	//ketty.SetLog(new(kettyLog.FakeLog))
	//bStartClient(b, httpUrl)
//}
//
//func Benchmark_Https(b *testing.B) {
	//ketty.SetLog(new(kettyLog.FakeLog))
	//bStartClient(b, httpsUrl)
//}
//
