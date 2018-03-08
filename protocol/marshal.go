package protocol

import (
	COM "github.com/yyzybb537/ketty/common"
	"github.com/golang/protobuf/proto"
	"encoding/json"
)

type Marshaler interface {
	Marshal(msg proto.Message) ([]byte, error)

	Unmarshal(buf []byte, msg proto.Message) error
}

var MgrMarshaler = COM.NewManager((*Marshaler)(nil))

func init() {
	MgrMarshaler.Register("pb", new(PbMarshaler))
	MgrMarshaler.Register("json", new(JsonMarshaler))
}

// ----------- default protobuf marshaler
type PbMarshaler struct {}

func (this *PbMarshaler) Marshal(msg proto.Message) ([]byte, error) {
	return proto.Marshal(msg)
}

func (this *PbMarshaler) Unmarshal(buf []byte, msg proto.Message) error {
	return proto.Unmarshal(buf, msg)
}

// ----------- json marshaler
type JsonMarshaler struct {}

func (this *JsonMarshaler) Marshal(msg proto.Message) ([]byte, error) {
	return json.Marshal(msg)
}

func (this *JsonMarshaler) Unmarshal(buf []byte, msg proto.Message) error {
	return json.Unmarshal(buf, msg)
}

