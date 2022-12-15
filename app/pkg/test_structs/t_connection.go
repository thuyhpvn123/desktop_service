package test_structs

import (
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/meta-node/meta-node/pkg/network"
)

type TestConnection struct {
	message network.IMessage
	Address common.Address
	Type    string
}

func (c *TestConnection) GetMessage() network.IMessage {
	return c.message
}

func (c *TestConnection) SendMessage(message network.IMessage) error {
	c.message = message
	return nil
}

func (c *TestConnection) Connect() error {
	return nil
}

func (c *TestConnection) Disconnect() error {
	return nil
}

func (c *TestConnection) GetIp() string {
	return ""
}

func (c *TestConnection) GetPort() int {
	return 0
}

func (c *TestConnection) GetAddress() common.Address {
	return c.Address
}

func (c *TestConnection) GetType() string {
	return c.Type
}

func (c *TestConnection) ReadRequest() (network.IRequest, error) {
	return nil, nil
}

func (c *TestConnection) Init(address common.Address, cType string) {
	c.Address = address
	c.Type = cType
}
