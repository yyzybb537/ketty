package grpc_t

import (
	"testing"
	echo "github.com/yyzybb537/ketty/test/test_pb"
	"github.com/yyzybb537/ketty"
	kettyLog "github.com/yyzybb537/ketty/log"
	P "github.com/yyzybb537/ketty/protocol"
	kettyHttp "github.com/yyzybb537/ketty/protocol/http"
	"time"
	context "golang.org/x/net/context"
)

type EchoServer struct {
}

func (this *EchoServer) Echo(ctx context.Context, req *echo.Req) (*echo.Rsp, error) {
	return &echo.Rsp{Val:req.Val}, nil
}

type ExceptionServer struct {
}

func (this *ExceptionServer) Echo(ctx context.Context, req *echo.Req) (*echo.Rsp, error) {
	ketty.GetLog().Infof("panic")
	panic("Echo_Panic")
	return &echo.Rsp{Val:req.Val}, nil
}

func startServer(t *testing.T, sUrl string, driverUrl string) {
	server, err := ketty.Listen(sUrl, driverUrl)
	if err != nil {
		t.Fatalf("Listen error:%s", err.Error())
	}

	server.RegisterMethod(echo.EchoServiceHandle, &EchoServer{})

	err = server.Serve()
	if err != nil {
		t.Fatalf("Serve (%s) error:%s", sUrl, err.Error())
	}
}

func startExceptionServer(t *testing.T, sUrl string) {
	server, err := ketty.Listen(sUrl, "")
	if err != nil {
		t.Fatalf("Listen error:%s", err.Error())
	}

	server.RegisterMethod(echo.EchoServiceHandle, &ExceptionServer{})

	err = server.Serve()
	if err != nil {
		t.Fatalf("Serve error:%s", err.Error())
	}
}
func startClient(t *testing.T, sUrl string, exceptError bool) {
	client, err := ketty.Dial(sUrl, "")
	if err != nil {
		t.Fatalf("Dial error:%s", err.Error())
	}

	req := &echo.Req{ Val : 123 }
	stub := echo.NewKettyEchoServiceClient(client)
	rsp, err := stub.Echo(context.Background(), req)
	if exceptError {
		if err == nil {
			t.Fatalf("Invoke not error")
        }
	} else {
		if err != nil {
			t.Fatalf("Invoke error:%+v", err)
		}

		t.Logf("Echo Val:%d", rsp.Val)
	}
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
var httpjsonUrl = "http.json://127.0.0.1:8092/hj/abc"
var httpsUrl = "https://127.0.0.1:3050"
var httpsjsonUrl = "https.json://127.0.0.1:3051"
var grpcUrl = "grpc://127.0.0.1:8090"
var etcdUrl = "grpc://127.0.0.1:8090"
var urls = []string{
	//"http://127.0.0.1",
	httpUrl,
	httpsUrl,
	grpcUrl,
	httpjsonUrl,
	httpsjsonUrl,
}

func TestGrpc(t *testing.T) {
	kettyHttp.InitTLS("cert.pem", "key.pem")
	ketty.SetLog(new(kettyLog.StdLog))
	ketty.GetLog().Debugf("Protocols:%s", P.DumpProtocols())
	for _, sUrl := range urls {
		t.Logf("Do url:%s", sUrl)
		startServer(t, sUrl, "")
		time.Sleep(time.Millisecond * 100)
		startClient(t, sUrl, false)
	}
}

func TestEtcd(t *testing.T) {
	kettyHttp.InitTLS("cert.pem", "key.pem")
	ketty.SetLog(new(kettyLog.StdLog))
	sUrl := "grpc://127.0.0.1:33009"
	t.Logf("Do url:%s", sUrl)
	driverUrl := "etcd://127.0.0.1:2379/ketty_test"
	startServer(t, sUrl, driverUrl)
	time.Sleep(time.Millisecond * 100)
	startClient(t, driverUrl, false)
}


func TestException(t *testing.T) {
	kettyHttp.InitTLS("cert.pem", "key.pem")
	ketty.SetLog(new(kettyLog.StdLog))
	sUrl := "grpc://127.0.0.1:33008"
	t.Logf("Do url:%s", sUrl)
	startExceptionServer(t, sUrl)
	time.Sleep(time.Millisecond * 100)
	startClient(t, sUrl, true)
}

func TestRouter(t *testing.T) {
	sUrl := "http.json://127.0.0.1:33099"
	startServer(t, sUrl + "/a", "")
	startServer(t, sUrl + "/b", "")
	time.Sleep(time.Millisecond * 100)
	startClient(t, sUrl + "/b", false)
}

func Benchmark_Grpc(b *testing.B) {
	ketty.SetLog(new(kettyLog.FakeLog))
	bStartClient(b, grpcUrl)
}

func Benchmark_Http(b *testing.B) {
	ketty.SetLog(new(kettyLog.FakeLog))
	bStartClient(b, httpUrl)
}

func Benchmark_Https(b *testing.B) {
	ketty.SetLog(new(kettyLog.FakeLog))
	bStartClient(b, httpsUrl)
}

