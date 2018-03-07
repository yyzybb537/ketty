package http_proto

import (
	"net/http"
	"github.com/golang/protobuf/proto"
	P "github.com/yyzybb537/ketty/protocol"
	"crypto/tls"
	"github.com/pkg/errors"
	"io/ioutil"
	"bytes"
)

type HttpInvoker struct {
	tr     *http.Transport
	client *http.Client
	m      P.Marshaler
}

func NewHttpInvoker(m P.Marshaler) *HttpInvoker {
	invoker := new(HttpInvoker)
	invoker.tr = &http.Transport{
        TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
    }
    invoker.client = &http.Client{Transport: invoker.tr}
	invoker.m = m
	return invoker
}

func NewPbHttpInvoker() *HttpInvoker {
	return NewHttpInvoker(new(P.PbMarshaler))
}

func NewJsonHttpInvoker() *HttpInvoker {
	return NewHttpInvoker(new(P.JsonMarshaler))
}

func (this *HttpInvoker) Do(url string, req, rsp proto.Message) error {
	return this.DoWithHeaders(url, req, rsp, nil)
}

func (this *HttpInvoker) DoWithHeaders(url string, req, rsp proto.Message, headers map[string]string) (err error) {
	buf, err := this.m.Marshal(req)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	httpRequest, err := http.NewRequest("POST", url, bytes.NewBuffer(buf))
	if err != nil {
		err = errors.WithStack(err)
		return
    }

	for k, v := range headers {
		httpRequest.Header.Set(k, v)
    }

	httpResponse, err := this.client.Do(httpRequest)
	if err != nil {
		err = errors.WithStack(err)
		return
    }
	defer httpResponse.Body.Close()

	buf, err = ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	if httpResponse.StatusCode == 501 {
		err = errors.Errorf(string(buf))
		return 
	}

	if httpResponse.StatusCode != 200 {
		err = errors.Errorf("error http status:%d", httpResponse.StatusCode)
		return 
	}

	err = this.m.Unmarshal(buf, rsp)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return 
}
