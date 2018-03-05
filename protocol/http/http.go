package http_proto

import (
	U "github.com/yyzybb537/ketty/url"
	P "github.com/yyzybb537/ketty/protocol"
)

type HttpProtocol struct {}

func init() {
	P.RegProtocol("http", new(HttpProtocol))
	U.RegDefaultPort("http", 80)
}

func (this *HttpProtocol) CreateServer(url, driverUrl U.Url) (P.Server, error) {
	return newHttpServer(url, driverUrl, new(P.PbMarshaler)), nil
}

func (this *HttpProtocol) Dial(url U.Url) (P.Client, error) {
	return newHttpClient(url, new(P.PbMarshaler))
}
