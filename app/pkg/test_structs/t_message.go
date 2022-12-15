package test_structs

import (
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type TestMessage struct {
	Message interface{}
	Proto   *pb.Message
}

func (t *TestMessage) Marshal() ([]byte, error) {
	return proto.Marshal(t.Proto)
}

func (t *TestMessage) GetBody() []byte {
	return t.Proto.Body
}

func (t *TestMessage) Unmarshal(protoStruct protoreflect.ProtoMessage) error {
	err := proto.Unmarshal(t.Proto.Body, protoStruct)
	return err
}

func (t *TestMessage) GetCommand() string {
	return t.Proto.Header.Command
}
