package vote

import (
	"github.com/ethereum/go-ethereum/common"
	cm "gitlab.com/meta-node/meta-node/pkg/common"
)

type IVote interface {
	GetValue() interface{}
	GetHash() common.Hash
	GetPublicKey() cm.PublicKey
	GetAddress() common.Address
	GetSign() cm.Sign
}
