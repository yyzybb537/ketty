package http_proto

import (
	"github.com/yyzybb537/ketty"
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
	ketty.AopList

	url ketty.Url
}

func newHttpClient(url ketty.Url) (*HttpClient, error) {
	c := new(HttpClient)
	c.url = url
	return c, nil
}

func (this *HttpClient) Close() {
	
}

func (this *HttpClient) Invoke(ctx context.Context, handle ketty.ServiceHandle, method string, req, rsp interface{}) error {
	buf, err := proto.Marshal(req.(proto.Message))
	if err != nil {
		return errors.WithStack(err)
	}
	
	fullMethodName := fmt.Sprintf("/%s/%s", strings.Replace(handle.ServiceName(), ".", "/", -1), method)
	fullUrl := this.url.ToString() + fullMethodName
	tr := &http.Transport{
        TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}
	metadata := map[string]string{}

	aopList := ketty.GetAop(ctx)
	if aopList != nil {
		ctx = context.WithValue(ctx, "method", fullMethodName)
		ctx = context.WithValue(ctx, "remote", this.url.ToString())

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

	httpRequest, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(buf))
	if err != nil {
		return errors.WithStack(err)
    }
	httpRequest.Header.Set("Content-Type", "binary/protobuf")
	if len(metadata) > 0 {
		var metadataBuf []byte
		metadataBuf, err = json.Marshal(metadata)
		if err != nil {
			return errors.WithStack(err)
		}

		httpRequest.Header.Set("KettyMetaData", string(metadataBuf))
	}

	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		return errors.WithStack(err)
    }
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != 200 {
		return errors.Errorf("status:%d", httpResponse.StatusCode)
	}

	buf, err = ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return errors.WithStack(err)
	}

	err = proto.Unmarshal(buf, rsp.(proto.Message))
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
