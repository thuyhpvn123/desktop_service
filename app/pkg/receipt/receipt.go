package receipt

import (
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type IReceipt interface {
	GetProto() protoreflect.ProtoMessage
}

func NewReceipt(proto *pb.Receipt) IReceipt {
	return nil
}
