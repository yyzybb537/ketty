package grpc_proto

import (
	U "github.com/yyzybb537/ketty/url"
	P "github.com/yyzybb537/ketty/protocol"
	"google.golang.org/grpc"
)

type GrpcProtocol struct {}

func init() {
	P.RegProtocol("grpc", new(GrpcProtocol))
	grpc.EnableTracing = false
}

func (this *GrpcProtocol) CreateServer(url, driverUrl U.Url) (P.Server, error) {
	return newGrpcServer(url, driverUrl), nil
}

func (this *GrpcProtocol) Dial(url U.Url) (P.Client, error) {
	return newGrpcClient(url)
}
