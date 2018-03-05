package http_proto

import (
	U "github.com/yyzybb537/ketty/url"
	P "github.com/yyzybb537/ketty/protocol"
)

type HttpsJsonProtocol struct {}

func init() {
	P.RegProtocol("httpsjson", new(HttpsJsonProtocol))
	U.RegDefaultPort("httpsjson", 443)
}

func (this *HttpsJsonProtocol) CreateServer(url, driverUrl U.Url) (P.Server, error) {
	return newHttpServer(url, driverUrl, new(P.JsonMarshaler)), nil
}

func (this *HttpsJsonProtocol) Dial(url U.Url) (P.Client, error) {
	return newHttpClient(url, new(P.JsonMarshaler))
}
