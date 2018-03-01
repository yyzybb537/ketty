package http_proto

import (
	U "github.com/yyzybb537/ketty/url"
	P "github.com/yyzybb537/ketty/protocol"
)

type HttpsProtocol struct {}

func init() {
	P.RegProtocol("https", new(HttpsProtocol))
	U.RegDefaultPort("https", 443)
}

var gCertFile, gKeyFile string

func InitTLS(certFile, keyFile string) {
	gCertFile = certFile
	gKeyFile = keyFile
}

func (this *HttpsProtocol) CreateServer(url, driverUrl U.Url) (P.Server, error) {
	return newHttpServer(url, driverUrl), nil
}

func (this *HttpsProtocol) Dial(url U.Url) (P.Client, error) {
	return newHttpClient(url)
}
