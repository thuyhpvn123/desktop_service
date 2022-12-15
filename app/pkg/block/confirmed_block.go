package block

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"gitlab.com/meta-node/meta-node/pkg/bls"
	cm "gitlab.com/meta-node/meta-node/pkg/common"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type IConfirmBlock interface {
	GetHash() common.Hash
	GetCount() *uint256.Int
	GetValidatorSigns() map[cm.PublicKey]cm.Sign
	AddValidatorSign(cm.PublicKey, cm.Sign)
	GetProto() protoreflect.ProtoMessage
	Marshal() ([]byte, error)
}

type ConfirmBlock struct {
	proto *pb.ConfirmBlock
}

func NewConfirmBlock(proto *pb.ConfirmBlock) IConfirmBlock {
	return &ConfirmBlock{
		proto: proto,
	}
}

func UnmarshalConfirmBlock(b []byte) (IConfirmBlock, error) {
	pbConfirmBlock := &pb.ConfirmBlock{}
	err := proto.Unmarshal(b, pbConfirmBlock)
	if err != nil {
		return nil, err
	}
	return NewConfirmBlock(pbConfirmBlock), nil
}

func ConfirmBlockFromFullBlock(b IFullBlock) *ConfirmBlock {
	mapPkValidatorSigns := b.GetValidatorSigns()
	mapStringValidatorSigns := make(map[string][]byte)
	for i, v := range mapPkValidatorSigns {
		mapStringValidatorSigns[hex.EncodeToString(i[:])] = v[:]
	}
	block := b.GetBlock()
	return &ConfirmBlock{
		proto: &pb.ConfirmBlock{
			Hash:           block.GetHash().Bytes(),
			Count:          block.GetCount().Bytes(),
			ValidatorSigns: mapStringValidatorSigns,
		},
	}
}
func (cb *ConfirmBlock) GetHash() common.Hash {
	return common.BytesToHash(cb.proto.Hash)
}

func (cb *ConfirmBlock) GetCount() *uint256.Int {
	return uint256.NewInt(0).SetBytes(cb.proto.Count)
}

func (cb *ConfirmBlock) GetValidatorSigns() map[cm.PublicKey]cm.Sign {
	rs := make(map[cm.PublicKey]cm.Sign)
	for i, v := range cb.proto.ValidatorSigns {
		rs[cm.PubkeyFromBytes(common.FromHex(i))] = cm.SignFromBytes(v)
	}
	return rs
}

func (cb *ConfirmBlock) AddValidatorSign(pk cm.PublicKey, sign cm.Sign) {
	cb.proto.ValidatorSigns[hex.EncodeToString(pk[:])] = sign[:]
}

func (cb *ConfirmBlock) GetProto() protoreflect.ProtoMessage {
	return cb.proto
}

func (cb *ConfirmBlock) Marshal() ([]byte, error) {
	return proto.Marshal(cb.proto)
}

func CheckBlockValidatorSigns(block IConfirmBlock) bool {
	validatorSigns := block.GetValidatorSigns()
	for pubkey, sign := range validatorSigns {
		if !bls.VerifySign(pubkey, sign, block.GetHash().Bytes()) {
			return false
		}
	}
	return true
}
