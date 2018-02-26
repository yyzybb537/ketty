package ketty

import (
	"golang.org/x/net/context"
)

type ClientMethod interface {
	init(client Client)
}

type Client interface {
	AopListI

	Close()

	Invoke(ctx context.Context, handle ServiceHandle, method string, req, rsp interface{}) error
}

