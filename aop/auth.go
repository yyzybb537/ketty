package aop

import (
	"golang.org/x/net/context"
	C "github.com/yyzybb537/ketty/context"
	COM "github.com/yyzybb537/ketty/common"
)

type AuthorizeClient interface {
	CreateAuthorization(ctx context.Context) string
}

type AuthorizeServer interface {
	CheckAuthorization(ctx context.Context, authorization string) error
}

type AuthAop struct {
	c AuthorizeClient
	s AuthorizeServer
}

func NewAuthAop(c AuthorizeClient, s AuthorizeServer) *AuthAop {
	return &AuthAop{c, s}
}

func (this *AuthAop) ClientSendMetaData(ctx context.Context, metadata map[string]string) context.Context {
	if this.c != nil {
		authorization := this.c.CreateAuthorization(ctx)
		if authorization != "" {
			metadata[COM.AuthorizationMetaKey] = authorization
		}
	}
	return ctx
}

func (this *AuthAop) ServerRecvMetaData(ctx context.Context, metadata map[string]string) context.Context {
	authorization, _ := metadata[COM.AuthorizationMetaKey]
	if this.s != nil {
		err := this.s.CheckAuthorization(ctx, authorization)
		if err != nil {
			return C.WithError(ctx, err)
		}
	}
	return ctx
}

