package http_proto

import (
	U "github.com/yyzybb537/ketty/url"
	P "github.com/yyzybb537/ketty/protocol"
)

type HttpProtocol struct {}

func init() {
	P.RegProtocol("http", new(HttpProtocol))
	U.RegDefaultPort("http", 80)

	P.RegProtocol("https", new(HttpProtocol))
	U.RegDefaultPort("https", 443)
}

var gCertFile, gKeyFile string

func InitTLS(certFile, keyFile string) {
	gCertFile = certFile
	gKeyFile = keyFile
}

func (this *HttpProtocol) CreateServer(url, driverUrl U.Url) (P.Server, error) {
	return newHttpServer(url, driverUrl)
}

func (this *HttpProtocol) Dial(url U.Url) (P.Client, error) {
	return newHttpClient(url)
}
