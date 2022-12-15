package test_structs

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/meta-node/meta-node/pkg/bls"
	cm "gitlab.com/meta-node/meta-node/pkg/common"
)

type TestKeyPairConfig struct {
	PrivateKey cm.PrivateKey
}

func (c *TestKeyPairConfig) GetPrivateKey() cm.PrivateKey {
	return c.PrivateKey
}

func (c *TestKeyPairConfig) GetPublicKey() cm.PublicKey {
	_, pubkey, _ := bls.GenerateKeyPairFromSecretKey(hex.EncodeToString(c.PrivateKey[:]))
	return pubkey
}

func (c *TestKeyPairConfig) GetAddress() common.Address {
	_, _, address := bls.GenerateKeyPairFromSecretKey(hex.EncodeToString(c.PrivateKey[:]))
	return address
}
