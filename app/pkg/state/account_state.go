package state

import (
	"github.com/ethereum/go-ethereum/common"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type IAccountState interface {
	GetAddress() common.Address
	Marshal() ([]byte, error)
	GetProto() protoreflect.ProtoMessage
}

func Unmarshal(b []byte) (IAccountState, error) {
	asProto := &pb.AccountState{}
	err := proto.Unmarshal(b, asProto)
	if err != nil {
		return nil, err
	}
	return NewAccountState(asProto), nil
}

type AccountState struct {
	proto *pb.AccountState
}

func NewAccountState(proto *pb.AccountState) IAccountState {
	return &AccountState{
		proto,
	}
}

func (as *AccountState) GetProto() protoreflect.ProtoMessage {
	return as.proto
}

func (as *AccountState) GetAddress() common.Address {
	return common.BytesToAddress(as.proto.Address)
}

func (as *AccountState) Marshal() ([]byte, error) {
	return proto.Marshal(as.proto)
}
