package vote

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	cm "gitlab.com/meta-node/meta-node/pkg/common"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type IBlockVote interface {
	GetBlockCount() *uint256.Int
	GetHash() common.Hash
	GetValue() interface{}
	GetPublicKey() cm.PublicKey
	GetAddress() common.Address
	GetSign() cm.Sign
	GetProto() protoreflect.ProtoMessage
}

type BlockVote struct {
	proto *pb.BlockVote
}

func NewBlockVote(
	proto *pb.BlockVote,
) IBlockVote {
	return &BlockVote{
		proto,
	}
}

func UnmarshalBlockVote(
	b []byte,
) (IBlockVote, error) {
	pbBlockVote := &pb.BlockVote{}
	err := proto.Unmarshal(b, pbBlockVote)
	if err != nil {
		return nil, err
	}
	return NewBlockVote(pbBlockVote), nil
}

func (v *BlockVote) GetBlockCount() *uint256.Int {
	return uint256.NewInt(0).SetBytes(v.proto.Count)
}

func (v *BlockVote) GetHash() common.Hash {
	return common.BytesToHash(v.proto.Hash)
}

func (v *BlockVote) GetValue() interface{} {
	if len(v.proto.BlockData) == 0 {
		return nil
	}
	return v.proto.BlockData
}

func (v *BlockVote) GetPublicKey() cm.PublicKey {
	return cm.PubkeyFromBytes(v.proto.Pubkey)
}

func (v *BlockVote) GetAddress() common.Address {
	return cm.GetAddressFromPubkey(v.GetPublicKey())
}

func (v *BlockVote) GetSign() cm.Sign {
	return cm.SignFromBytes(v.proto.Sign)
}

func (v *BlockVote) GetProto() protoreflect.ProtoMessage {
	return v.proto
}
