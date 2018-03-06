package http_proto

import (
	C "github.com/yyzybb537/ketty/context"
	COM "github.com/yyzybb537/ketty/common"
	U "github.com/yyzybb537/ketty/url"
	A "github.com/yyzybb537/ketty/aop"
	P "github.com/yyzybb537/ketty/protocol"
	"fmt"
	"strings"
	"golang.org/x/net/context"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"encoding/json"
)

type HttpClient struct {
	A.AopList

	url U.Url
	invoker *HttpInvoker
}

func newHttpClient(url U.Url, m P.Marshaler) (*HttpClient, error) {
	c := new(HttpClient)
	c.url = url
	c.invoker = NewHttpInvoker(m)
	return c, nil
}

func (this *HttpClient) Close() {
}

func (this *HttpClient) getUrl() string {
	s := this.url.ToString()
	switch this.url.Protocol {
	case "http":
		fallthrough
	case "https":
		return s

	case "httpjson":
		return s[:4] + s[8:]
	case "httpsjson":
		return s[:5] + s[9:]
    }
	return s
}

func (this *HttpClient) Invoke(ctx context.Context, handle COM.ServiceHandle, method string, req, rsp interface{}) (error) {
	ctx = this.invoke(ctx, handle, method, req, rsp)
	return ctx.Err()
}

func (this *HttpClient) invoke(inCtx context.Context, handle COM.ServiceHandle, method string, req, rsp interface{}) (ctx context.Context) {
	var err error
	ctx = inCtx
	fullMethodName := fmt.Sprintf("/%s/%s", strings.Replace(handle.ServiceName(), ".", "/", -1), method)
	fullUrl := this.getUrl() + fullMethodName
	metadata := map[string]string{}

	aopList := A.GetAop(ctx)
	if aopList != nil {
		ctx = context.WithValue(ctx, "method", fullMethodName)
		ctx = context.WithValue(ctx, "remote", this.getUrl())

		for _, aop := range aopList {
			caller, ok := aop.(A.ClientTransportMetaDataAop)
			if ok {
				ctx = caller.ClientSendMetaData(ctx, metadata)
				if ctx.Err() != nil {
					return 
				}
			}
		}

		for _, aop := range aopList {
			caller, ok := aop.(A.BeforeClientInvokeAop)
			if ok {
				ctx = caller.BeforeClientInvoke(ctx, req)
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
			aop := aopList[len(aopList) - i - 1]
			caller, ok := aop.(A.AfterClientInvokeAop)
			if ok {
				defer caller.AfterClientInvoke(&ctx, req, rsp)
			}
		}
	}

	headers := map[string]string{
		"Content-Type" : "binary/protobuf",
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

	err = this.invoker.DoWithHeaders(fullUrl, req.(proto.Message), rsp.(proto.Message), headers)
	if err != nil {
		ctx = C.WithError(ctx, err)
		return
	}

	return
}
