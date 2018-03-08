package main

import (
	ext_pb "github.com/yyzybb537/ketty/demo/extpb"
	"github.com/yyzybb537/ketty"
	"golang.org/x/net/context"
)

func main() {
	client, err := ketty.Dial("http.json://127.0.0.1:8091", "")
	if err != nil {
		ketty.GetLog().Errorf("Dial error:%s", err.Error())
		return
	}
	defer client.Close()

	req := &ext_pb.Req{
		Qr : &ext_pb.QueryReq{ QVal : 123 },
		Jr : &ext_pb.JsonReq{ JVal : 321 },
	}
	stub := ext_pb.NewKettyEchoServiceClient(client)
	rsp, err := stub.Echo(context.Background(), req)
	if err != nil {
		ketty.GetLog().Errorf("Invoke error:%+v", err)
		return
	}

	ketty.GetLog().Infof("Rsp: %d", rsp.Val)
}

