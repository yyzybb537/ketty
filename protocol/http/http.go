package http_proto

import (
	"github.com/yyzybb537/ketty"
)

type HttpProtocol struct {}

func init() {
	ketty.RegProtocol("http", new(HttpProtocol))
}

func (this *HttpProtocol) DefaultPort() int {
	return 80
}

func (this *HttpProtocol) CreateServer(url, driverUrl ketty.Url) (ketty.Server, error) {
	return newHttpServer(url, driverUrl), nil
}

func (this *HttpProtocol) Dial(url ketty.Url) (ketty.Client, error) {
	return newHttpClient(url)
}
