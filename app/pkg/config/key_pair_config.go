package config

import (
	"github.com/ethereum/go-ethereum/common"
	cm "gitlab.com/meta-node/meta-node/pkg/common"
)

type IKeyPairConfig interface {
	GetPrivateKey() cm.PrivateKey
	GetPublicKey() cm.PublicKey
	GetAddress() common.Address
}
