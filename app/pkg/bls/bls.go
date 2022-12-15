package bls

import (
	"crypto/rand"
	"encoding/hex"
	"runtime"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	blst "gitlab.com/meta-node/meta-node/pkg/bls/blst/bindings/go"
	cm "gitlab.com/meta-node/meta-node/pkg/common"
)

type blstPublicKey = blst.P1Affine
type blstSignature = blst.P2Affine
type blstAggregateSignature = blst.P2Aggregate
type blstAggregatePublicKey = blst.P1Aggregate
type blstSecretKey = blst.SecretKey

var dstMinPk = []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")

func Init() {
	blst.SetMaxProcs(runtime.GOMAXPROCS(0))
}

func Sign(bPri cm.PrivateKey, bMessage []byte) cm.Sign {
	sk := new(blstSecretKey).Deserialize(bPri[:])
	sign := new(blstSignature).Sign(sk, bMessage, dstMinPk)
	return cm.SignFromBytes(sign.Compress())
}

func GetByteAddress(pubkey []byte) []byte {
	hash := crypto.Keccak256(pubkey)
	address := hash[12:]
	return address
}

func VerifySign(bPub cm.PublicKey, bSig cm.Sign, bMsg []byte) bool {
	return new(blstSignature).VerifyCompressed(bSig[:], true, bPub[:], false, bMsg, dstMinPk)
}

func VerifyAggregateSign(bPubs [][]byte, bSig []byte, bMsgs [][]byte) bool {
	return new(blstSignature).AggregateVerifyCompressed(bSig, true, bPubs, false, bMsgs, dstMinPk)
}

func GenerateKeyPairFromSecretKey(hexSecretKey string) (cm.PrivateKey, cm.PublicKey, common.Address) {
	secByte, _ := hex.DecodeString(hexSecretKey)
	sec := new(blstSecretKey).Deserialize(secByte)
	pub := new(blstPublicKey).From(sec).Compress()
	hash := crypto.Keccak256([]byte(pub))
	return cm.PrivateKeyFromBytes(sec.Serialize()), cm.PubkeyFromBytes(pub), common.BytesToAddress(hash[12:])
}

func randBLSTSecretKey() *blstSecretKey {
	var t [32]byte
	_, _ = rand.Read(t[:])
	secretKey := blst.KeyGen(t[:])
	return secretKey
}

func GenerateKeyPair() ([]byte, []byte, []byte) {
	sec := randBLSTSecretKey()
	pub := new(blstPublicKey).From(sec).Compress()
	hash := crypto.Keccak256([]byte(pub))
	return sec.Serialize(), pub, hash[12:]
}

// func CreateAggregateSignFromTransactions(transactions []*pb.Transaction) []byte {
// 	log.Debugf("CreateAggregateSignFromTransactions total transaction %v", len(transactions))
// 	aggregatedSignature := new(blst.P2Aggregate)
// 	signatures := make([][]byte, len(transactions))
// 	for i, v := range transactions {
// 		signatures[i] = v.Sign
// 	}
// 	aggregatedSignature.AggregateCompressed(signatures, false)
// 	log.Debugf("aggreagtesign %v", aggregatedSignature.ToAffine().Compress())

// 	return aggregatedSignature.ToAffine().Compress()
// }
