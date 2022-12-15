package network

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type IConnectionsManager interface {
	AddConnection(IConnection)
	RemoveConnection(IConnection)
	GetConnectionByAddress(common.Address) IConnection
	GetConnectionsByType(string) map[common.Address]IConnection
}

type ConnectionsManager struct {
	mu                    sync.RWMutex
	connections           []IConnection
	mapAddressConnections map[common.Address]IConnection
	mapTypeConnections    map[string]map[common.Address]IConnection
}

func NewConnectionsManager(
	connectionTypes []string,
) *ConnectionsManager {
	cm := &ConnectionsManager{}
	cm.mapTypeConnections = make(map[string]map[common.Address]IConnection)
	cm.mapAddressConnections = make(map[common.Address]IConnection)
	for _, v := range connectionTypes {
		cm.mapTypeConnections[v] = make(map[common.Address]IConnection)
	}
	return cm
}

func (cm *ConnectionsManager) AddConnection(conn IConnection) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.connections = append(cm.connections, conn)
	address := conn.GetAddress()
	if (address != common.Address{}) {
		cm.mapAddressConnections[conn.GetAddress()] = conn
		cType := conn.GetType()
		if cType != "" && cm.mapTypeConnections[cType] != nil {
			cm.mapTypeConnections[cType][address] = conn
		}
	}
}

func removeConnectionAtIndex(s []IConnection, index int) []IConnection {
	return append(s[:index], s[index+1:]...)
}

func (cm *ConnectionsManager) RemoveConnection(conn IConnection) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for i, v := range cm.connections {
		if v == conn {
			cm.connections = removeConnectionAtIndex(cm.connections, i)
		}
	}
	address := conn.GetAddress()
	if (address != common.Address{}) {
		delete(cm.mapAddressConnections, conn.GetAddress())
		cType := conn.GetType()
		if cType != "" && cm.mapTypeConnections[cType] != nil {
			delete(cm.mapTypeConnections[cType], conn.GetAddress())
		}
	}
}

func (cm *ConnectionsManager) GetConnectionByAddress(address common.Address) IConnection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.mapAddressConnections[address]
}

func (cm *ConnectionsManager) GetConnectionsByType(cType string) map[common.Address]IConnection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.mapTypeConnections[cType]
}
