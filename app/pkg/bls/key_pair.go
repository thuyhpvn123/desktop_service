package bls

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	cm "gitlab.com/meta-node/meta-node/pkg/common"
)

type KeyPair struct {
	publicKey  cm.PublicKey
	privateKey cm.PrivateKey
	address    common.Address
}

func NewKeyPair(privateKey []byte) *KeyPair {
	sec := new(blstSecretKey).Deserialize(privateKey)
	pub := new(blstPublicKey).From(sec).Compress()
	hash := crypto.Keccak256([]byte(pub))
	return &KeyPair{
		privateKey: cm.PrivateKeyFromBytes(sec.Serialize()),
		publicKey:  cm.PubkeyFromBytes(pub),
		address:    common.BytesToAddress(hash[12:]),
	}
}

func (kp *KeyPair) GetPrivateKey() cm.PrivateKey {
	return kp.privateKey
}

func (kp *KeyPair) GetBytesPrivateKey() []byte {
	return kp.privateKey[:]
}

func (kp *KeyPair) GetPublicKey() cm.PublicKey {
	return kp.publicKey
}

func (kp *KeyPair) GetBytesPublicKey() []byte {
	return kp.publicKey[:]
}

func (kp *KeyPair) GetAddress() common.Address {
	return kp.address
}
