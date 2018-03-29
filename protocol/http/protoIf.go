package http_proto

import (
	"reflect"
	"github.com/golang/protobuf/proto"
)

type KettyHttpExtend interface {
	KettyHttpExtendMessage()
}

type DefineMarshaler interface {
	KettyMarshal() string
}
var typeDefineMarshaler = reflect.TypeOf((*DefineMarshaler)(nil)).Elem()

type DefineTransport interface {
	KettyTransport() string
}
var typeDefineTransport = reflect.TypeOf((*DefineTransport)(nil)).Elem()

var typeProtoMessage = reflect.TypeOf((*proto.Message)(nil)).Elem()
