package grpc_proto

import (
	C "github.com/yyzybb537/ketty/context"
	D "github.com/yyzybb537/ketty/driver"
	O "github.com/yyzybb537/ketty/option"
	COM "github.com/yyzybb537/ketty/common"
	U "github.com/yyzybb537/ketty/url"
	A "github.com/yyzybb537/ketty/aop"
	"github.com/yyzybb537/ketty/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	md "google.golang.org/grpc/metadata"
	"net"
)

type GrpcServer struct {
	A.AopList

	url       U.Url
	driverUrl U.Url
	Impl      *grpc.Server
	opt		  *GrpcOption
}

func newGrpcServer(url, driverUrl U.Url) *GrpcServer {
	s := &GrpcServer{
		url:       url,
		opt : defaultGrpcOption(),
		driverUrl: driverUrl,
	}
	s.Impl = grpc.NewServer(grpc.UnaryInterceptor(s.serverIntercept()))
	return s
}

func (this *GrpcServer) SetOption(opt O.OptionI) error {
	return this.opt.set(opt)
}

func (this *GrpcServer) RegisterMethod(handle COM.ServiceHandle, implement interface{}) error {
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
				log.GetLog().Errorf("Serve lis error:%s. addr:%s", err.Error(), addr)
			}
		}()
	}

	if !this.driverUrl.IsEmpty() {
		driver, err := D.GetDriver(this.driverUrl.Protocol)
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
	//log.GetLog().Infof("unaryServerInterceptor ctx=%+v", ctx)
	rsp, ctx = this.unaryServerInterceptorWithContext(ctx, req, info, handler)
	err = ctx.Err()
	//log.GetLog().Infof("unaryServerInterceptor error:%v", err)
	return
}

func (this *GrpcServer) unaryServerInterceptorWithContext(inCtx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (rsp interface{}, ctx context.Context) {
	var err error
	ctx = inCtx
	aopList := this.GetAop()
	if aopList != nil {
		ctx = context.WithValue(ctx, "method", info.FullMethod)
		p, ok := peer.FromContext(ctx)
		if ok {
			ctx = context.WithValue(ctx, "remote", p.Addr.String())
		} else {
			ctx = context.WithValue(ctx, "remote", "")
		}
		metadata := map[string]string{}

		grpcMD, hasMetaData := md.FromIncomingContext(ctx)
		//log.GetLog().Debugf("server hasMetaData:%v meta:%v", hasMetaData, grpcMD)
		if hasMetaData {
			for k, v := range grpcMD {
				if len(v) > 0 {
					metadata[k] = v[0]
				}
			}
		}

		for _, aop := range aopList {
			caller, ok := aop.(A.ServerTransportMetaDataAop)
			if ok {
				ctx = caller.ServerRecvMetaData(ctx, metadata)
			}
		}

		for _, aop := range aopList {
			caller, ok := aop.(A.BeforeServerInvokeAop)
			if ok {
				ctx = caller.BeforeServerInvoke(ctx, req)
			}
		}

		defer func() {
			for _, aop := range aopList {
				caller, ok := aop.(A.ServerInvokeCleanupAop)
				if ok {
					caller.ServerCleanup(ctx)
				}
			}
		}()

		for i, _ := range aopList {
			aop := aopList[len(aopList) - i - 1]
			caller, ok := aop.(A.AfterServerInvokeAop)
			if ok {
				defer caller.AfterServerInvoke(&ctx, req, rsp)
			}
		}
	}

	if ctx.Err() != nil {
		return
	}

	rsp, err = handler(ctx, req)
	if err != nil {
		ctx = C.WithError(ctx, err)
		return
	}

	return
}
