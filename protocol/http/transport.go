package http_proto

import (
	COM "github.com/yyzybb537/ketty/common"
	"bytes"
	"net/http"
	"io/ioutil"
	"github.com/yyzybb537/ketty/log"
)

var _ = log.GetLog

type DataTransport interface {
	Write(req *http.Request, buf []byte) error
	Read(req *http.Request) ([]byte, error)
}

var MgrTransport = COM.NewManager((*DataTransport)(nil))

func init() {
	MgrTransport.Register("body", new(BodyTransport))
	MgrTransport.Register("query", new(QueryStringTransport))
	MgrTransport.Register("multipart", new(MultipartTransport))
}

type BodyTransport struct {}
func (this *BodyTransport) Write(req *http.Request, buf []byte) error {
	b := bytes.NewBuffer(buf)
	req.Body = ioutil.NopCloser(b)
	req.ContentLength = int64(b.Len())
	return nil
}
func (this *BodyTransport) Read(req *http.Request) ([]byte, error) {
	buf, err := ioutil.ReadAll(req.Body)
	//println("Body:", string(buf))
	return buf, err
}

type QueryStringTransport struct {}
func (this *QueryStringTransport) Write(req *http.Request, buf []byte) error {
	if req.URL.RawQuery != "" {
		req.URL.RawQuery += "&"
    }
	req.URL.RawQuery += string(buf)
	return nil
}
func (this *QueryStringTransport) Read(req *http.Request) ([]byte, error) {
	return []byte(req.URL.RawQuery), nil
}

type MultipartTransport struct {}
func (this *MultipartTransport) Write(req *http.Request, buf []byte) error {
	return nil
}
func (this *MultipartTransport) Read(req *http.Request) ([]byte, error) {
	return []byte{}, nil
}
