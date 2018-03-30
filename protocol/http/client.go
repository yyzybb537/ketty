package http_proto

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	A "github.com/yyzybb537/ketty/aop"
	COM "github.com/yyzybb537/ketty/common"
	C "github.com/yyzybb537/ketty/context"
	O "github.com/yyzybb537/ketty/option"
	P "github.com/yyzybb537/ketty/protocol"
	U "github.com/yyzybb537/ketty/url"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

type HttpClient struct {
	A.AopList

	url    U.Url
	tr     *http.Transport
	client *http.Client
	prt    *Proto
	opt    *HttpOption
}

func newHttpClient(url U.Url) (*HttpClient, error) {
	c := new(HttpClient)
	c.url = url
	var err error
	c.prt, err = ParseProto(url.Protocol)
	if err != nil {
		return nil, err
	}
	c.tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c.client = &http.Client{Transport: c.tr}
	return c, nil
}

func (this *HttpClient) SetOption(opt O.OptionI) error {
	return this.opt.set(opt)
}

func (this *HttpClient) Close() {
	this.tr.CloseIdleConnections()
}

func (this *HttpClient) getUrl() string {
	return this.url.ToStringByProtocol(this.url.GetMainProtocol())
}

func (this *HttpClient) Invoke(ctx context.Context, handle COM.ServiceHandle, method string, req, rsp interface{}) error {
	pbReq, ok := req.(proto.Message)
	if !ok {
		return fmt.Errorf("Invoke req is not proto.Message")
	}
	pbRsp, ok := rsp.(proto.Message)
	if !ok {
		return fmt.Errorf("Invoke rsp is not proto.Message")
	}
	ctx = this.invoke(ctx, handle, method, pbReq, pbRsp)
	return ctx.Err()
}

func (this *HttpClient) invoke(inCtx context.Context, handle COM.ServiceHandle, method string, req, rsp proto.Message) (ctx context.Context) {
	var err error
	ctx = inCtx
	fullMethodName := fmt.Sprintf("/%s/%s", strings.Replace(handle.ServiceName(), ".", "/", -1), method)
	metadata := map[string]string{}

	httpRequest, err := http.NewRequest(strings.ToUpper(this.prt.DefaultMethod), this.getUrl(), nil)
	if err != nil {
		ctx = C.WithError(ctx, errors.WithStack(err))
		return
	}

	aopList := A.GetAop(ctx)
	if aopList != nil {
		ctx = context.WithValue(ctx, "method", fullMethodName)
		ctx = context.WithValue(ctx, "remote", this.getUrl())
		ctx = setHttpRequest(ctx, httpRequest)

		for _, aop := range aopList {
			caller, ok := aop.(A.BeforeClientInvokeAop)
			if ok {
				ctx = caller.BeforeClientInvoke(ctx, req)
				if ctx.Err() != nil {
					return
				}
			}
		}

		for _, aop := range aopList {
			caller, ok := aop.(A.ClientTransportMetaDataAop)
			if ok {
				ctx = caller.ClientSendMetaData(ctx, metadata)
				if ctx.Err() != nil {
					return
				}
			}
		}

		defer func() {
			for _, aop := range aopList {
				caller, ok := aop.(A.ClientInvokeCleanupAop)
				if ok {
					caller.ClientCleanup(ctx)
				}
			}
		}()

		for i, _ := range aopList {
			aop := aopList[len(aopList)-i-1]
			caller, ok := aop.(A.AfterClientInvokeAop)
			if ok {
				defer caller.AfterClientInvoke(&ctx, req, rsp)
			}
		}
	}

	headers := map[string]string{
		"KettyMethod": fullMethodName,
	}

	// 鉴权数据用Header发送
	if authorization, exists := metadata[COM.AuthorizationMetaKey]; exists {
		headers["Authorization"] = authorization
		delete(metadata, COM.AuthorizationMetaKey)
	}

	if len(metadata) > 0 {
		var metadataBuf []byte
		metadataBuf, err = json.Marshal(metadata)
		if err != nil {
			ctx = C.WithError(ctx, errors.WithStack(err))
			return
		}

		headers["KettyMetaData"] = string(metadataBuf)
	}

	for k, v := range headers {
		httpRequest.Header.Set(k, v)
	}

	err = this.doHttpRequest(httpRequest, req, rsp)
	if err != nil {
		ctx = C.WithError(ctx, err)
		return
	}

	return
}

func (this *HttpClient) doHttpRequest(httpRequest *http.Request, req, rsp proto.Message) (err error) {
	err = this.writeMessage(httpRequest, req)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	httpResponse, err := this.client.Do(httpRequest)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer httpResponse.Body.Close()

	var buf []byte
	buf, err = ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	if httpResponse.StatusCode == http.StatusBadRequest {
		err = errors.Errorf(string(buf))
		return
	}

	if httpResponse.StatusCode != http.StatusOK {
		err = errors.Errorf("error http status:%d", httpResponse.StatusCode)
		return
	}

	mr, _ := P.MgrMarshaler.Get(this.prt.DefaultMarshaler).(P.Marshaler)
	err = mr.Unmarshal(buf, rsp)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

func (this *HttpClient) writeMessage(httpRequest *http.Request, req proto.Message) error {
	_, isKettyHttpExtend := req.(KettyHttpExtend)
	if !isKettyHttpExtend {
		// Not extend, use default or pb setttings.
		sMr := this.prt.DefaultMarshaler
		if dmr, ok := req.(DefineMarshaler); ok {
			sMr = dmr.KettyMarshal()
		}
		mr, _ := P.MgrMarshaler.Get(sMr).(P.Marshaler)
		buf, err := mr.Marshal(req)
		if err != nil {
			return err
		}

		sTr := this.prt.DefaultTransport
		if dtr, ok := req.(DefineTransport); ok {
			sTr = dtr.KettyTransport()
		}
		tr, _ := MgrTransport.Get(sTr).(DataTransport)
		err = tr.Write(httpRequest, buf)
		if err != nil {
			return err
		}

		return nil
	}

	// use http extend
	typ := reflect.TypeOf(req).Elem()
	val := reflect.ValueOf(req).Elem()
	trMap := map[string]bool{}
	numFields := typ.NumField()
	for i := 0; i < numFields; i++ {
		fvalue := val.Field(i)
		ftype := typ.Field(i).Type
		if !ftype.ConvertibleTo(typeProtoMessage) {
			return fmt.Errorf("Use http extend message, all of fields must be proto.Message! Error message name is %s", typ.Name())
		}

		if fvalue.Interface() == nil {
			// skip nil message
			continue
		}

		sTr := this.prt.DefaultTransport
		if ftype.ConvertibleTo(typeDefineTransport) {
			sTr = fvalue.Convert(typeDefineTransport).Interface().(DefineTransport).KettyTransport()
		}

		// check tr unique
		if _, exists := trMap[sTr]; exists {
			return fmt.Errorf("The message used http extend, transport must be unique. Too many field use transport:%s", sTr)
		}
		trMap[sTr] = true

		sMr := this.prt.DefaultMarshaler
		if ftype.ConvertibleTo(typeDefineMarshaler) {
			sMr = fvalue.Convert(typeDefineMarshaler).Interface().(DefineMarshaler).KettyMarshal()
		}
		if sTr == "query" {
			sMr = "querystring"
		}
		mr, ok := P.MgrMarshaler.Get(sMr).(P.Marshaler)
		if !ok {
			return fmt.Errorf("Unknown marshal(%s) in message:%s", sMr, ftype.Name())
		}
		buf, err := mr.Marshal(fvalue.Interface().(proto.Message))
		if err != nil {
			return err
		}

		tr, ok := MgrTransport.Get(sTr).(DataTransport)
		if !ok {
			return fmt.Errorf("Unknown transport(%s) in message:%s", sTr, ftype.Name())
		}
		err = tr.Write(httpRequest, buf)
		if err != nil {
			return err
		}

		if sTr == "body" {
			var contentType string
			switch sMr {
			case "pb":
				contentType = "application/octet-stream"
			case "querystring":
				contentType = "application/x-www-form-urlencoded"
			case "multipart":
				contentType = "multipart/form-data; boundary=" + DefaultMultipartBoundary
			case "json":
				contentType = "application/json"
			default:
				contentType = "text/plain"
			}
			httpRequest.Header.Set("Content-Type", contentType)
		}
	}

	return nil
}
