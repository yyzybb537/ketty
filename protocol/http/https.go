package http_proto

import (
	"github.com/yyzybb537/ketty"
)

type HttpsProtocol struct {}

func init() {
	ketty.RegProtocol("https", new(HttpsProtocol))
}

var gCertFile, gKeyFile string

func InitTLS(certFile, keyFile string) {
	gCertFile = certFile
	gKeyFile = keyFile
}

func (this *HttpsProtocol) DefaultPort() int {
	return 443
}

func (this *HttpsProtocol) CreateServer(url, driverUrl ketty.Url) (ketty.Server, error) {
	return newHttpServer(url, driverUrl), nil
}

func (this *HttpsProtocol) Dial(url ketty.Url) (ketty.Client, error) {
	return newHttpClient(url)
}
