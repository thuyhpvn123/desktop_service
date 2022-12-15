package network

import (
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type IMessage interface {
	Marshal() ([]byte, error)
	Unmarshal(protoStruct protoreflect.ProtoMessage) error
	GetBody() []byte
	GetCommand() string
}

type Message struct {
	proto *pb.Message
}

func NewMessage(pbMessage *pb.Message) IMessage {
	return &Message{
		proto: pbMessage,
	}
}

func (m *Message) Marshal() ([]byte, error) {
	return proto.Marshal(m.proto)
}

func (m *Message) Unmarshal(protoStruct protoreflect.ProtoMessage) error {
	err := proto.Unmarshal(m.proto.Body, protoStruct)
	return err
}

func (m *Message) GetCommand() string {
	return m.proto.Header.Command
}

func (m *Message) GetBody() []byte {
	return m.proto.Body
}
