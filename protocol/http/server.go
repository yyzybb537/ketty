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
	opts        *Options
}

func newHttpServer(url, driverUrl U.Url) (*HttpServer, error) {
	s := &HttpServer {
		Impl : &http.Server{},
		url : url,
		driverUrl : driverUrl,
		mux : http.NewServeMux(),
    }
	var err error
	s.opts, err = ParseOptions(url.Protocol)
	if err != nil {
		return nil, err
    }
	s.Impl.Handler = s.mux
	return s, nil
}

func (this *HttpServer) parseMessage(httpRequest *http.Request, requestType reflect.Type) (proto.Message, error) {
	req := reflect.New(requestType)
	_, isKettyHttpExtend := req.Interface().(KettyHttpExtend)
	if !isKettyHttpExtend {
		// Not extend, use default
		tr, _ := MgrTransport.Get(this.opts.DefaultTransport).(DataTransport)
		buf, err := tr.Read(httpRequest)
		if err != nil {
			return nil, err
        }

		mr, _ := P.MgrMarshaler.Get(this.opts.DefaultMarshaler).(P.Marshaler)
		err = mr.Unmarshal(buf, req.Interface().(proto.Message))
		if err != nil {
			return nil, err
        }

		return req.Interface().(proto.Message), nil
    }

	// use http extend
	numFields := req.Elem().NumField()
	trMap := map[string]bool{}
	for i := 0; i < numFields; i++ {
		fvalue := req.Elem().Field(i)
		ftype := requestType.Field(i).Type
		if !ftype.ConvertibleTo(typeProtoMessage) {
			return nil, fmt.Errorf("Use http extend message, all of fields must be proto.Message! Error message name is %s", requestType.Name())
        }
		fvalue.Set(reflect.New(ftype.Elem()))

		sTr := this.opts.DefaultTransport
		if ftype.ConvertibleTo(typeDefineTransport) {
			sTr = fvalue.Convert(typeDefineTransport).Interface().(DefineTransport).KettyTransport()
        }

		// check tr unique
		if _, exists := trMap[sTr]; exists {
			return nil, fmt.Errorf("The message used http extend, transport must be unique. Too many field use transport(%s) in message:%s", sTr, requestType.Name())
		}
		trMap[sTr] = true

		tr, ok := MgrTransport.Get(sTr).(DataTransport)
		if !ok {
			return nil, fmt.Errorf("Unknown transport(%s) in message:%s", sTr, ftype.Name())
        }

		buf, err := tr.Read(httpRequest)
		if err != nil {
			return nil, err
        }

		if len(buf) == 0 {
			// skip nil message
			continue
        }

		sMr := this.opts.DefaultMarshaler
		if ftype.ConvertibleTo(typeDefineMarshaler) {
			sMr = fvalue.Convert(typeDefineMarshaler).Interface().(DefineMarshaler).KettyMarshal()
        }

		if sTr == "query" {
			sMr = "querystring"
		}

		mr, ok := P.MgrMarshaler.Get(sMr).(P.Marshaler)
		if !ok {
			return nil, fmt.Errorf("Unknown marshal(%s) in message:%s", sMr, ftype.Name())
        }

		fMessage := fvalue.Interface().(proto.Message)
		err = mr.Unmarshal(buf, fMessage)
		if err != nil {
			return nil, err
        }
	}

	return req.Interface().(proto.Message), nil
}

func (this *HttpServer) doHandler(pattern string, httpRequest *http.Request, requestType reflect.Type, reflectMethod reflect.Value) (rsp interface{}, ctx context.Context) {
	ctx = context.Background()
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

	// 解析Message
	req, err := this.parseMessage(httpRequest, requestType)
	if err != nil {
		ctx = C.WithError(ctx, errors.WithStack(err))
		return 
	}

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

	replies := reflectMethod.Call([]reflect.Value{reflect.ValueOf(context.Background()), reflect.ValueOf(req)})
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
			mr, _ := P.MgrMarshaler.Get(this.opts.DefaultMarshaler).(P.Marshaler)
			buf, err = mr.Marshal(rsp.(proto.Message))
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
	if strings.HasPrefix(proto, "https") {
		go func() {
			this.Impl.Addr = U.FormatAddr(addr, this.url.Protocol)
			err := this.Impl.ListenAndServeTLS(gCertFile, gKeyFile)
			if err != nil {
				log.GetLog().Errorf("Http.ServeTLS lis error:%s. addr:%s", err.Error(), addr)
			}
        }()
    } else if strings.HasPrefix(proto, "http") {
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

