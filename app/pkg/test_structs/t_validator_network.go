package test_structs

import (
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/meta-node/meta-node/pkg/network"
)

type TestValidatorNetwork struct {
	ValidatorConnections map[common.Address]network.IConnection
	NodeConnections      map[common.Address]network.IConnection
}

func (t *TestValidatorNetwork) GetValidatorConnections() map[common.Address]network.IConnection {
	return t.ValidatorConnections
}
func (t *TestValidatorNetwork) GetNodeConnections() map[common.Address]network.IConnection {
	return t.NodeConnections
}

func (t *TestValidatorNetwork) GetValidatorAddresses() map[common.Address]interface{} {
	rs := make(map[common.Address]interface{})
	for k := range t.ValidatorConnections {
		rs[k] = true
	}
	return rs
}

func (t *TestValidatorNetwork) GetNodeAddresses() map[common.Address]interface{} {
	rs := make(map[common.Address]interface{})
	for k := range t.NodeConnections {
		rs[k] = true
	}
	return rs
}

func (t *TestValidatorNetwork) GetValidatorConnection(address common.Address) network.IConnection {
	return t.ValidatorConnections[address]
}
