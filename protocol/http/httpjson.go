package http_proto

import (
	U "github.com/yyzybb537/ketty/url"
	P "github.com/yyzybb537/ketty/protocol"
)

type HttpJsonProtocol struct {}

func init() {
	P.RegProtocol("httpjson", new(HttpJsonProtocol))
	U.RegDefaultPort("httpjson", 80)
}

func (this *HttpJsonProtocol) CreateServer(url, driverUrl U.Url) (P.Server, error) {
	return newHttpServer(url, driverUrl, new(P.JsonMarshaler)), nil
}

func (this *HttpJsonProtocol) Dial(url U.Url) (P.Client, error) {
	return newHttpClient(url, new(P.JsonMarshaler))
}
