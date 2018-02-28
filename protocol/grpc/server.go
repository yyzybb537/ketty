package grpc_proto

import (
	"github.com/yyzybb537/ketty"
	kettyContext "github.com/yyzybb537/ketty/context"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	md "google.golang.org/grpc/metadata"
	"net"
)

type GrpcServer struct {
	ketty.AopList

	url       ketty.Url
	driverUrl ketty.Url
	Impl      *grpc.Server
}

func newGrpcServer(url, driverUrl ketty.Url) *GrpcServer {
	s := &GrpcServer{
		url:       url,
		driverUrl: driverUrl,
	}
	s.Impl = grpc.NewServer(grpc.UnaryInterceptor(s.serverIntercept()))
	return s
}

func (this *GrpcServer) RegisterMethod(handle ketty.ServiceHandle, implement interface{}) error {
	this.Impl.RegisterService(handle.Implement().(*grpc.ServiceDesc), implement)
	return nil
}

func (this *GrpcServer) Serve() error {
	addrs := this.url.GetAddrs()
	for _, addr := range addrs {
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			return err
		}

		go func() {
			err := this.Impl.Serve(lis)
			if err != nil {
				ketty.GetLog().Errorf("Serve lis error:%s. addr:%s", err.Error(), addr)
			}
		}()
	}

	if !this.driverUrl.IsEmpty() {
		driver, err := ketty.GetDriver(this.driverUrl.Protocol)
		if err != nil {
			return err
		}

		err = driver.Register(this.driverUrl, this.url)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *GrpcServer) serverIntercept() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (rsp interface{}, err error) {
		rsp, err = this.unaryServerInterceptor(ctx, req, info, handler)
		return
	}
}

func (this *GrpcServer) unaryServerInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (rsp interface{}, err error) {
	rsp, ctx = this.unaryServerInterceptorWithContext(ctx, req, info, handler)
	err = ctx.Err()
	//ketty.GetLog().Infof("unaryServerInterceptor error:%v", err)
	return
}

func (this *GrpcServer) unaryServerInterceptorWithContext(inCtx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (rsp interface{}, ctx context.Context) {
	var err error
	ctx = inCtx
	aopList := this.GetAop()
	if aopList != nil {
		ctx = context.WithValue(ctx, "method", info.FullMethod)
		ctx = context.WithValue(ctx, "remote", "") //TODO
		metadata := map[string]string{}

		grpcMD, hasMetaData := md.FromIncomingContext(ctx)
		//ketty.GetLog().Debugf("server hasMetaData:%v meta:%v", hasMetaData, grpcMD)
		if hasMetaData {
			for k, v := range grpcMD {
				if len(v) > 0 {
					metadata[k] = v[0]
				}
			}
		}

		for _, aop := range aopList {
			caller, ok := aop.(ketty.ServerTransportMetaDataAop)
			if ok {
				ctx = caller.ServerRecvMetaData(ctx, metadata)
				if ctx.Err() != nil {
					return
				}
			}
		}

		for _, aop := range aopList {
			caller, ok := aop.(ketty.BeforeServerInvokeAop)
			if ok {
				ctx = caller.BeforeServerInvoke(ctx, req)
				if ctx.Err() != nil {
					return
				}
			}
		}

		defer func() {
			for _, aop := range aopList {
				caller, ok := aop.(ketty.ServerInvokeCleanupAop)
				if ok {
					caller.ServerCleanup(ctx)
				}
			}
		}()

		for i, _ := range aopList {
			aop := aopList[len(aopList) - i - 1]
			caller, ok := aop.(ketty.AfterServerInvokeAop)
			if ok {
				defer caller.AfterServerInvoke(&ctx, req, rsp)
			}
		}
	}

	rsp, err = handler(ctx, req)
	if err != nil {
		ctx = kettyContext.WithError(ctx, err)
		return
	}

	return
}
