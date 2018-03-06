package main

import (
	echo "github.com/yyzybb537/ketty/demo/pb"
	"github.com/yyzybb537/ketty"
	"github.com/yyzybb537/ketty/aop"
	"golang.org/x/net/context"
)

type Auth struct {}
func (this *Auth) CreateAuthorization(ctx context.Context) string {
	return "MyAuthInfo"
}

func main() {
	client, err := ketty.Dial("grpc://127.0.0.1:8090", "")
	if err != nil {
		ketty.GetLog().Errorf("Dial error:%s", err.Error())
		return
	}
	defer client.Close()

	client.AddAop(aop.NewAuthAop(new(Auth), nil))

	req := &echo.Req{ Val : 123 }
	stub := echo.NewKettyEchoServiceClient(client)
	rsp, err := stub.Echo(context.Background(), req)
	if err != nil {
		ketty.GetLog().Errorf("Invoke error:%+v", err)
		return
	}

	ketty.GetLog().Infof("Rsp: %d", rsp.Val)
}

