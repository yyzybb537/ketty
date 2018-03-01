package main

import (
	echo "github.com/yyzybb537/ketty/test/test_pb"
	"github.com/yyzybb537/ketty"
	kettyLog "github.com/yyzybb537/ketty/log"
	kettyHttp "github.com/yyzybb537/ketty/protocol/http"
	context "golang.org/x/net/context"
	"fmt"
	"flag"
	"time"
	"sync/atomic"
)

var clientMap = map[string]ketty.Client{}
var qps int64
var lastQps int64

func bStartClient(sUrl string) {
	client, err := ketty.Dial(sUrl, "")
	if err != nil {
		panic(fmt.Errorf("Dial error:%s", err.Error()))
		return
	}
	defer client.Close()

	for {
		req := &echo.Req{ Val : 123 }
		stub := echo.NewKettyEchoServiceClient(client)
		_, err := stub.Echo(context.Background(), req)
		if err != nil {
			panic(fmt.Errorf("Invoke error:%+v", err))
			return
		}

		atomic.AddInt64(&qps, 1)
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

var urlType int
var nConnections int

func init() {
	flag.IntVar(&urlType, "url", 2, "0:http, 1:https, 2:grpc")
	flag.IntVar(&nConnections, "c", 100, "number of connections")
	flag.Parse()
}

func main() {
	kettyHttp.InitTLS("cert.pem", "key.pem")
	ketty.SetLog(new(kettyLog.FakeLog))
	println("Connect to", urls[urlType])
	for i := 0; i < nConnections; i++ {
		go bStartClient(urls[urlType])
	}

	for {
		time.Sleep(time.Second)
		cur := atomic.LoadInt64(&qps)
		println("QPS:", cur - lastQps)
		lastQps = cur
	}
}

