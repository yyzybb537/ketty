package http_proto

import (
	"github.com/yyzybb537/ketty"
	kettyContext "github.com/yyzybb537/ketty/context"
	"net"
	"net/http"
	"reflect"
	"fmt"
	"strings"
	"io/ioutil"
	"google.golang.org/grpc"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"github.com/pkg/errors"
	"encoding/json"
)

type HttpServer struct {
	ketty.AopList

	Impl		*http.Server
	url			ketty.Url
	driverUrl	ketty.Url
	mux         *http.ServeMux
}

func newHttpServer(url, driverUrl ketty.Url) (*HttpServer) {
	s := &HttpServer {
		Impl : &http.Server{},
		url : url,
		driverUrl : driverUrl,
		mux : http.NewServeMux(),
    }
	s.Impl.Handler = s.mux
	return s
}

func (this *HttpServer) doHandler(pattern string, httpRequest *http.Request, requestType reflect.Type, reflectMethod reflect.Value) (rsp interface{}, ctx context.Context) {
	ctx = context.Background()
	var buf []byte
	var err error

	// metadata
	metadata := map[string]string{}
	metadataStr := httpRequest.Header.Get("KettyMetaData")
	if metadataStr != "" {
		err = json.Unmarshal([]byte(metadataStr), &metadata)
		if err != nil {
			ctx = kettyContext.WithError(ctx, errors.WithStack(err))
			return
		}
	}

	var reflectReq reflect.Value
	buf, err = ioutil.ReadAll(httpRequest.Body)
	if err != nil {
		ctx = kettyContext.WithError(ctx, errors.WithStack(err))
		return
	}

	reflectReq = reflect.New(requestType)
	err = proto.Unmarshal(buf, reflectReq.Interface().(proto.Message))
	if err != nil {
		ctx = kettyContext.WithError(ctx, errors.WithStack(err))
		return 
	}
	req := reflectReq.Interface().(proto.Message)

	aopList := this.GetAop()
	if aopList != nil {
		ctx = context.WithValue(ctx, "method", pattern)
		ctx = context.WithValue(ctx, "remote", httpRequest.RemoteAddr)

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

	replies := reflectMethod.Call([]reflect.Value{reflect.ValueOf(context.Background()), reflectReq})
	rsp = replies[0].Interface()
	if replies[1].Interface() != nil {
		err = replies[1].Interface().(error)
		ctx = kettyContext.WithError(ctx, err)
		return 
	}

	return
}

func (this *HttpServer) RegisterMethod(handle ketty.ServiceHandle, implement interface{}) error {
	desc := handle.Implement().(*grpc.ServiceDesc)
	ht := reflect.TypeOf(desc.HandlerType).Elem()
	it := reflect.TypeOf(implement)
	if !it.Implements(ht) {
		return fmt.Errorf("service struct not implements type `%s` interface", desc.ServiceName)
    }

	iv := reflect.ValueOf(implement)
	for _, method := range desc.Methods {
		reflectMethod := iv.MethodByName(method.MethodName)
		requestType := reflectMethod.Type().In(1).Elem()
		pattern := fmt.Sprintf("/%s/%s", strings.Replace(handle.ServiceName(), ".", "/", -1), method.MethodName)
		this.mux.HandleFunc(pattern, func(w http.ResponseWriter, httpRequest *http.Request){
			var err error
			defer func() {
				if err != nil {
					w.WriteHeader(501)
					w.Write([]byte(err.Error()))
				}
			}()

			rsp, ctx := this.doHandler(pattern, httpRequest, requestType, reflectMethod)
			err = ctx.Err()
			if err != nil {
				return 
			}

			// body
			var buf []byte
			buf, err = proto.Marshal(rsp.(proto.Message))
			if err != nil {
				return 
			}

			//ketty.GetLog().Debugf("Write response: %v", buf)

			w.WriteHeader(200)
			w.Write(buf)
		})
	}
	return nil
}

func (this *HttpServer) serve(addr string, proto string) error {
	if proto == "http" {
		lis, err := net.Listen("tcp", ketty.FormatAddr(addr, this.url.Protocol))
		if err != nil {
			return err
		}

		go func() {
			err := this.Impl.Serve(lis)
			if err != nil {
				ketty.GetLog().Errorf("Http.Serve lis error:%s. addr:%s", err.Error(), addr)
			}
        }()
    } else if proto == "https" {
		go func() {
			this.Impl.Addr = ketty.FormatAddr(addr, this.url.Protocol)
			err := this.Impl.ListenAndServeTLS(gCertFile, gKeyFile)
			if err != nil {
				ketty.GetLog().Errorf("Http.ServeTLS lis error:%s. addr:%s", err.Error(), addr)
			}
        }()
    } else {
		return errors.Errorf("Error protocol:%s", proto)
    }

	return nil
}

func (this *HttpServer) Serve() error {
	addrs := this.url.GetAddrs()
	for _, addr := range addrs {
		err := this.serve(addr, this.url.Protocol)
		if err != nil {
			return err
		}
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

