package grpc_proto

import (
	"github.com/yyzybb537/ketty"
)

type GrpcProtocol struct {}

func init() {
	ketty.RegProtocol("grpc", new(GrpcProtocol))
}

func (this *GrpcProtocol) DefaultPort() int {
	return 0
}

func (this *GrpcProtocol) CreateServer(url, driverUrl ketty.Url) (ketty.Server, error) {
	return newGrpcServer(url, driverUrl), nil
}

func (this *GrpcProtocol) Dial(url ketty.Url) (ketty.Client, error) {
	return newGrpcClient(url)
}
