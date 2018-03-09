package http_proto

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	P "github.com/yyzybb537/ketty/protocol"
	//"reflect"
	//"bytes"
	//"strings"
	//"fmt"
	//COM "github.com/yyzybb537/ketty/common"
)

const DefaultMultipartBoundary = "----KettyFromBoundary1234567890123456K"

// ----------- multipart marshaler
type MultipartMarshaler struct{}

func (this *MultipartMarshaler) Marshal(msg proto.Message) ([]byte, error) {
	return nil, errors.Errorf("Unsupport multipart")
}

func (this *MultipartMarshaler) Unmarshal(buf []byte, msg proto.Message) error {
	return errors.Errorf("Unsupport multipart")
}

func init() {
	P.MgrMarshaler.Register("multipart", new(MultipartMarshaler))
}
