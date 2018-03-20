package http_proto

import (
	"golang.org/x/net/context"
	"net/http"
)

type httpRequestKey struct{}
type httpResponseWriterKey struct{}

func setHttpRequest(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, httpRequestKey{}, req)
}

func GetHttpRequest(ctx context.Context) (*http.Request, bool) {
	req, ok := ctx.Value(httpRequestKey{}).(*http.Request)
	return req, ok
}

func setHttpResponseWriter(ctx context.Context, w http.ResponseWriter) context.Context {
	return context.WithValue(ctx, httpResponseWriterKey{}, w)
}

func GetHttpResponseWriter(ctx context.Context) (http.ResponseWriter, bool) {
	w, ok := ctx.Value(httpResponseWriterKey{}).(http.ResponseWriter)
	return w, ok
}
