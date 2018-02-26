package grpc_proto

import (
	"golang.org/x/net/context"
)

type GrpcMeta struct {
	owner *GrpcClient
}

func newGrpcMeta(owner *GrpcClient) *GrpcMeta {
	return &GrpcMeta{owner : owner}
}

func (this *GrpcMeta) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	metadata, ok := ctx.Value("metadata").(map[string]string)
	if !ok {
		return map[string]string{}, nil
	}

	return metadata, nil
}

func (this *GrpcMeta) RequireTransportSecurity() bool {
	return false
}
