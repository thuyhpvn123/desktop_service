package network

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/meta-node/meta-node/pkg/bls"
	cm "gitlab.com/meta-node/meta-node/pkg/common"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
)

func getHeaderForCommand(
	pubkey cm.PublicKey,
	command string,
	sign cm.Sign,
	version string,
) *pb.Header {
	return &pb.Header{
		Command: command,
		Pubkey:  pubkey[:],
		Sign:    sign[:],
		Version: version,
	}
}

func SendMessage(
	connection IConnection,
	keyPair *bls.KeyPair,
	command string,
	pbMessage proto.Message,
	sign cm.Sign,
	version string,
) error {
	body := []byte{}
	if pbMessage != nil {
		var err error
		body, err = proto.Marshal(pbMessage)
		if err != nil {
			return err
		}
	}
	if len(sign[:]) == 0 {
		bodyHash := crypto.Keccak256(body)
		sign = bls.Sign(keyPair.GetPrivateKey(), bodyHash)
	}
	messageProto := &pb.Message{
		Header: getHeaderForCommand(keyPair.GetPublicKey(), command, sign, version),
		Body:   body,
	}
	message := NewMessage(messageProto)
	return connection.SendMessage(message)
}

func SendBytes(
	connection IConnection,
	keyPair *bls.KeyPair,
	command string,
	bytes []byte,
	sign cm.Sign,
	version string,
) error {
	fmt.Println("message sent2222222")

	if len(sign[:]) == 0 {
		hash := crypto.Keccak256(bytes)
		sign = bls.Sign(keyPair.GetPrivateKey(), hash)
	}
	messageProto := &pb.Message{
		Header: getHeaderForCommand(keyPair.GetPublicKey(), command, sign, version),
		Body:   bytes,
	}
	message := NewMessage(messageProto)
	return connection.SendMessage(message)
}
