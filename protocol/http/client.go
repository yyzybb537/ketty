package http_proto

import (
	C "github.com/yyzybb537/ketty/context"
	COM "github.com/yyzybb537/ketty/common"
	U "github.com/yyzybb537/ketty/url"
	A "github.com/yyzybb537/ketty/aop"
	P "github.com/yyzybb537/ketty/protocol"
	"fmt"
	"strings"
	"net/http"
	"golang.org/x/net/context"
	"github.com/golang/protobuf/proto"
	"bytes"
	"io/ioutil"
	"github.com/pkg/errors"
	"crypto/tls"
	"encoding/json"
)

type HttpClient struct {
	A.AopList

	url U.Url
	tr *http.Transport
	client *http.Client
	m P.Marshaler
}

func newHttpClient(url U.Url, m P.Marshaler) (*HttpClient, error) {
	c := new(HttpClient)
	c.url = url
	c.tr = &http.Transport{
        TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
    }
    c.client = &http.Client{Transport: c.tr}
	c.m = m
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
	ctx = inCtx
	buf, err := this.m.Marshal(req.(proto.Message))
	if err != nil {
		ctx = C.WithError(ctx, errors.WithStack(err))
		return
	}
	
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

	httpRequest, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(buf))
	if err != nil {
		ctx = C.WithError(ctx, errors.WithStack(err))
		return
    }
	httpRequest.Header.Set("Content-Type", "binary/protobuf")
	if len(metadata) > 0 {
		var metadataBuf []byte
		metadataBuf, err = json.Marshal(metadata)
		if err != nil {
			ctx = C.WithError(ctx, errors.WithStack(err))
			return
		}

		httpRequest.Header.Set("KettyMetaData", string(metadataBuf))
	}

	httpResponse, err := this.client.Do(httpRequest)
	if err != nil {
		ctx = C.WithError(ctx, errors.WithStack(err))
		return
    }
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != 200 {
		ctx = C.WithError(ctx, errors.Errorf("status:%d", httpResponse.StatusCode))
		return 
	}

	buf, err = ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		ctx = C.WithError(ctx, errors.WithStack(err))
		return
	}

	err = this.m.Unmarshal(buf, rsp.(proto.Message))
	if err != nil {
		ctx = C.WithError(ctx, errors.WithStack(err))
		return
	}

	return
}
