package grpc_proto

import (
	"github.com/yyzybb537/ketty"
	"fmt"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
)

type GrpcClient struct {
	ketty.AopList

	Impl *grpc.ClientConn
	url ketty.Url
}

func newGrpcClient(url ketty.Url) (*GrpcClient, error) {
	c := new(GrpcClient)
	c.url = url
	err := c.dial(url)
	if err != nil {
		return nil, err
    }
	return c, nil
}

func (this *GrpcClient) dial(url ketty.Url) (err error) {
	this.Impl, err = grpc.Dial(url.SAddr, grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(newGrpcMeta(this)))
	return
}

func (this *GrpcClient) Close() {
	this.Impl.Close()
}

func (this *GrpcClient) Invoke(ctx context.Context, handle ketty.ServiceHandle, method string, req, rsp interface{}) (err error) {
	fullMethodName := fmt.Sprintf("/%s/%s", handle.ServiceName(), method)
	aopList := ketty.GetAop(ctx)
	if aopList != nil {
		ctx = context.WithValue(ctx, "method", fullMethodName)
		ctx = context.WithValue(ctx, "remote", this.url.SAddr)
		metadata := map[string]string{}
		ctx = context.WithValue(ctx, "metadata", metadata)

		for _, aop := range aopList {
			caller, ok := aop.(ketty.ClientTransportMetaDataAop)
			if ok {
				ctx = caller.ClientSendMetaData(ctx, metadata)
			}
		}

		for _, aop := range aopList {
			caller, ok := aop.(ketty.BeforeClientInvokeAop)
			if ok {
				ctx = caller.BeforeClientInvoke(ctx, req)
			}
		}

		defer func() {
			for _, aop := range aopList {
				caller, ok := aop.(ketty.ClientInvokeCleanupAop)
				if ok {
					caller.ClientCleanup(ctx)
				}
			}
        }()

		defer func() {
			for _, aop := range aopList {
				caller, ok := aop.(ketty.AfterClientInvokeAop)
				if ok {
					ctx = caller.AfterClientInvoke(ctx, req, rsp, err)
				}
			}
        }()
	}

	return grpc.Invoke(ctx, fullMethodName, req, rsp, this.Impl)
}

