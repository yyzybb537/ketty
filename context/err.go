package context

import (
	"golang.org/x/net/context"
)

type errorCtx struct {
	context.Context
	err error
}

func WithError(ctx context.Context, err error) context.Context {
	return &errorCtx{ ctx, err }
}

func (this *errorCtx) Err() error {
	return this.err
}
