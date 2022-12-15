package common

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type PublicKey [48]byte
type PrivateKey [32]byte
type Sign [96]byte

func PubkeyFromBytes(bytes []byte) PublicKey {
	p := PublicKey{}
	copy(p[0:48], bytes)
	return p
}

func SignFromBytes(bytes []byte) Sign {
	s := Sign{}
	copy(s[0:96], bytes)
	return s
}

func PrivateKeyFromBytes(bytes []byte) PrivateKey {
	p := PrivateKey{}
	copy(p[0:32], bytes)
	return p
}

func GetAddressFromPubkey(pk PublicKey) common.Address {
	hash := crypto.Keccak256(pk[:])
	return common.BytesToAddress(hash[12:])
}
