package pack

import (
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type IPack interface {
	GetTimestamp() uint64
	GetHash() common.Hash
	GetProto() protoreflect.Message
	GetBytes() []byte
}

type Pack struct {
	proto *pb.Pack
}

func NewPack(protoPack *pb.Pack) IPack {
	return &Pack{
		proto: protoPack,
	}
}

func Unmarshal(b []byte) (IPack, error) {
	protoPack := &pb.Pack{}
	err := proto.Unmarshal(b, protoPack)
	if err != nil {
		return nil, err
	}
	return NewPack(protoPack), nil
}

func (p *Pack) GetTimestamp() uint64 {
	return uint64(p.proto.TimeStamp)
}

func (p *Pack) GetHash() common.Hash {
	return common.BytesToHash(p.proto.Hash)
}

func (p *Pack) GetProto() protoreflect.Message {
	return p.proto.ProtoReflect()
}

func (p *Pack) GetBytes() []byte {
	b, err := proto.Marshal(p.proto)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return b
}
