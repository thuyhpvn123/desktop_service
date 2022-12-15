package block

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	cm "gitlab.com/meta-node/meta-node/pkg/common"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"gitlab.com/meta-node/meta-node/pkg/receipt"
	"gitlab.com/meta-node/meta-node/pkg/state"
	"google.golang.org/protobuf/proto"
)

type IFullBlock interface {
	GetBlock() IBlock
	GetValidatorSigns() map[cm.PublicKey]cm.Sign
	AddValidatorSign(cm.PublicKey, cm.Sign)
	GetAccountStateChanges() map[common.Address]state.IAccountState
	Marshal() ([]byte, error)
	SetAccountStates(state.IAccountStates)
	GetAccountStates() state.IAccountStates
}

type FullBlock struct {
	block               IBlock
	receipts            []receipt.IReceipt
	validatorSigns      map[cm.PublicKey]cm.Sign
	accountStateChanges map[common.Address]state.IAccountState
	accountStates       state.IAccountStates
}

func UnmarshalFullBlock(b []byte) (IFullBlock, error) {
	fbProto := &pb.FullBlock{}
	err := proto.Unmarshal(b, fbProto)
	if err != nil {
		return nil, err
	}
	return NewFullBlock(fbProto), nil
}

func NewFullBlock(fbProto *pb.FullBlock) IFullBlock {
	block := NewBlock(fbProto.Block)

	receipts := []receipt.IReceipt{}
	for _, v := range fbProto.Receipts {
		receipts = append(receipts, receipt.NewReceipt(v))
	}

	accountStateChanges := map[common.Address]state.IAccountState{}
	for _, v := range fbProto.AccountStateChanges {
		accountStateChanges[common.BytesToAddress(v.Address)] = state.NewAccountState(v)
	}

	validatorSigns := make(map[cm.PublicKey]cm.Sign)
	for i, v := range fbProto.ValidatorSigns {
		validatorSigns[cm.PubkeyFromBytes(common.FromHex(i))] = cm.SignFromBytes(v)
	}

	fullBlock := &FullBlock{
		block:               block,
		receipts:            receipts,
		accountStateChanges: accountStateChanges,
		validatorSigns:      validatorSigns,
	}
	return fullBlock
}

func (fb *FullBlock) GetBlock() IBlock {
	return fb.block
}
func (fb *FullBlock) GetValidatorSigns() map[cm.PublicKey]cm.Sign {
	return fb.validatorSigns
}

func (fb *FullBlock) AddValidatorSign(pk cm.PublicKey, sign cm.Sign) {
	fb.validatorSigns[pk] = sign
}

func (fb *FullBlock) GetAccountStateChanges() map[common.Address]state.IAccountState {
	return fb.accountStateChanges
}

func (fb *FullBlock) SetAccountStates(as state.IAccountStates) {
	fb.accountStates = as.Copy()
}

func (fb *FullBlock) GetAccountStates() state.IAccountStates {
	return fb.accountStates
}

func (fb *FullBlock) Marshal() ([]byte, error) {
	validatorSigns := make(map[string][]byte)
	for i, v := range fb.validatorSigns {
		validatorSigns[hex.EncodeToString(i[:])] = v[:]
	}

	receipts := []*pb.Receipt{}
	for _, v := range fb.receipts {
		receipts = append(receipts, v.GetProto().(*pb.Receipt))
	}

	accountStateChanges := []*pb.AccountState{}
	for _, v := range fb.accountStateChanges {
		accountStateChanges = append(accountStateChanges, v.GetProto().(*pb.AccountState))
	}
	fbProto := &pb.FullBlock{
		Block:               fb.block.GetProto().(*pb.Block),
		ValidatorSigns:      validatorSigns,
		Receipts:            receipts,
		AccountStateChanges: accountStateChanges,
	}
	return proto.Marshal(fbProto)
}
