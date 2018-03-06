package http_proto

import (
	C "github.com/yyzybb537/ketty/context"
	COM "github.com/yyzybb537/ketty/common"
	D "github.com/yyzybb537/ketty/driver"
	U "github.com/yyzybb537/ketty/url"
	A "github.com/yyzybb537/ketty/aop"
	P "github.com/yyzybb537/ketty/protocol"
	"github.com/yyzybb537/ketty/log"
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
	A.AopList

	Impl		*http.Server
	url			U.Url
	driverUrl	U.Url
	mux         *http.ServeMux
	m			P.Marshaler
}

func newHttpServer(url, driverUrl U.Url, m P.Marshaler) (*HttpServer) {
	s := &HttpServer {
		Impl : &http.Server{},
		url : url,
		driverUrl : driverUrl,
		mux : http.NewServeMux(),
		m : m,
    }
	s.Impl.Handler = s.mux
	return s
}

func (this *HttpServer) doHandler(pattern string, httpRequest *http.Request, requestType reflect.Type, reflectMethod reflect.Value) (rsp interface{}, ctx context.Context) {
	ctx = context.Background()
	var buf []byte
	var err error

	//log.GetLog().Debugf("HttpServer Request: %s", log.LogFormat(httpRequest, log.Indent))

	// metadata
	metadata := map[string]string{}
	metadataStr := httpRequest.Header.Get("KettyMetaData")
	if metadataStr != "" {
		err = json.Unmarshal([]byte(metadataStr), &metadata)
		if err != nil {
			ctx = C.WithError(ctx, errors.WithStack(err))
			return
		}
	}

	// 鉴权数据从Header提取
	authorization := httpRequest.Header.Get("Authorization")
	if authorization != "" {
		metadata[COM.AuthorizationMetaKey] = authorization
	}

	var reflectReq reflect.Value
	buf, err = ioutil.ReadAll(httpRequest.Body)
	if err != nil {
		ctx = C.WithError(ctx, errors.WithStack(err))
		return
	}

	reflectReq = reflect.New(requestType)
	err = this.m.Unmarshal(buf, reflectReq.Interface().(proto.Message))
	if err != nil {
		ctx = C.WithError(ctx, errors.WithStack(err))
		return 
	}
	req := reflectReq.Interface().(proto.Message)

	aopList := this.GetAop()
	if aopList != nil {
		ctx = context.WithValue(ctx, "method", pattern)
		ctx = context.WithValue(ctx, "remote", httpRequest.RemoteAddr)

		for _, aop := range aopList {
			caller, ok := aop.(A.ServerTransportMetaDataAop)
			if ok {
				ctx = caller.ServerRecvMetaData(ctx, metadata)
				if ctx.Err() != nil {
					return 
				}
			}
		}

		for _, aop := range aopList {
			caller, ok := aop.(A.BeforeServerInvokeAop)
			if ok {
				ctx = caller.BeforeServerInvoke(ctx, req)
				if ctx.Err() != nil {
					return 
				}
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

	replies := reflectMethod.Call([]reflect.Value{reflect.ValueOf(context.Background()), reflectReq})
	rsp = replies[0].Interface()
	if replies[1].Interface() != nil {
		err = replies[1].Interface().(error)
		ctx = C.WithError(ctx, err)
		return 
	}

	return
}

func (this *HttpServer) RegisterMethod(handle COM.ServiceHandle, implement interface{}) error {
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
			buf, err = this.m.Marshal(rsp.(proto.Message))
			if err != nil {
				return 
			}

			//log.GetLog().Debugf("Write response: %v", buf)

			w.WriteHeader(200)
			w.Write(buf)
		})
	}
	return nil
}

func (this *HttpServer) serve(addr string, proto string) error {
	if proto == "http" || proto == "httpjson" {
		lis, err := net.Listen("tcp", U.FormatAddr(addr, this.url.Protocol))
		if err != nil {
			return err
		}

		go func() {
			err := this.Impl.Serve(lis)
			if err != nil {
				log.GetLog().Errorf("Http.Serve lis error:%s. addr:%s", err.Error(), addr)
			}
        }()
    } else if proto == "https" || proto == "httpsjson" {
		go func() {
			this.Impl.Addr = U.FormatAddr(addr, this.url.Protocol)
			err := this.Impl.ListenAndServeTLS(gCertFile, gKeyFile)
			if err != nil {
				log.GetLog().Errorf("Http.ServeTLS lis error:%s. addr:%s", err.Error(), addr)
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

