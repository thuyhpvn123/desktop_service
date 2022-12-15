package bls

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	cm "gitlab.com/meta-node/meta-node/pkg/common"
)

var (
	testSecret1 = "4b1d80795e7404e9bc20fb5d0e783400dbf29ef006ad9518c712a8c60d33c2dd"
	testSecret2 = "6e65d1ec7a396f422d7cce990485d2f2fc7703a3cda1e0da05806249f7e360c9"
	testSecret3 = "5ef8b3caa1c03827c1a2bfa12236450122d9321d8dd74dbe10a643de18e4fd5c"

	testPubkey1 = common.FromHex("a2702ce6bbfb2e013935781bac50a0e168732bd957861e6fbf185d688c82ade34c9f33fead179decb5953b3382b061df")
	testSign1   = common.FromHex("a507c03ab7ebb69a4b3adc22a0347bb2466788e6a3baa174a62bd74cdff60dfd6d6ba9ec6237098f1ceef6013bfeff1d0c8be716266710e1493c422293a676e7f168007324a23435d4590896f97f8e3686cf0c280240b9406800c1cec6bafb5d")
	testHash1   = common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")
)

func TestGenerateKeyPair(t *testing.T) {
	sec, pub, address := GenerateKeyPair()
	fmt.Printf("Secret: %v\nPublic: %v\nAddress: %v\n", common.Bytes2Hex(sec), common.Bytes2Hex(pub), common.Bytes2Hex(address))
	assert.NotNil(t, sec)
	assert.NotNil(t, pub)
	assert.NotNil(t, address)
}

func TestGenerateKeyPairFromPrivateKey(t *testing.T) {
	sec1, pub1, address1 := GenerateKeyPairFromSecretKey(testSecret1)
	sec2, pub2, address2 := GenerateKeyPairFromSecretKey(testSecret2)
	sec3, pub3, address3 := GenerateKeyPairFromSecretKey(testSecret3)
	fmt.Printf("Secret: %v\nPublic: %v\nAddress: %v\n", common.Bytes2Hex(sec1[:]), common.Bytes2Hex(pub1[:]), common.Bytes2Hex(address1[:]))
	fmt.Printf("Secret: %v\nPublic: %v\nAddress: %v\n", common.Bytes2Hex(sec2[:]), common.Bytes2Hex(pub2[:]), common.Bytes2Hex(address2[:]))
	fmt.Printf("Secret: %v\nPublic: %v\nAddress: %v\n", common.Bytes2Hex(sec3[:]), common.Bytes2Hex(pub3[:]), common.Bytes2Hex(address3[:]))
}

func TestSign(t *testing.T) {
	sign1 := Sign(cm.PrivateKeyFromBytes(common.FromHex(testSecret1)), testHash1.Bytes())
	sign2 := Sign(cm.PrivateKeyFromBytes(common.FromHex(testSecret2)), testHash1.Bytes())
	sign3 := Sign(cm.PrivateKeyFromBytes(common.FromHex(testSecret3)), testHash1.Bytes())
	fmt.Printf("Sign1: %v\n", common.Bytes2Hex(sign1[:]))
	fmt.Printf("Sign2: %v\n", common.Bytes2Hex(sign2[:]))
	fmt.Printf("Sign3: %v\n", common.Bytes2Hex(sign3[:]))
}

func TestVerifySign(t *testing.T) {
	VerifySign(cm.PubkeyFromBytes(testPubkey1), cm.SignFromBytes(testSign1), testHash1.Bytes())
}
