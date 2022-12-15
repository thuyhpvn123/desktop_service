package network

import (
 "errors"
 "fmt"

 "github.com/ethereum/go-ethereum/common"
 "gitlab.com/meta-node/meta-node/cmd/node/command"
 "gitlab.com/meta-node/meta-node/cmd/node/processor"
 n_sync "gitlab.com/meta-node/meta-node/cmd/node/sync"
 "gitlab.com/meta-node/meta-node/pkg/block"
 "gitlab.com/meta-node/meta-node/pkg/bls"
 "gitlab.com/meta-node/meta-node/pkg/config"
 "gitlab.com/meta-node/meta-node/pkg/logger"
 "gitlab.com/meta-node/meta-node/pkg/network"
 pb "gitlab.com/meta-node/meta-node/pkg/proto"
)

var (
 ErrorCommandNotFound = errors.New("command not found")
)

type Handler struct {
 keyPair            *bls.KeyPair
 config             config.IConfig
 connectionsManager network.IConnectionsManager
 syncReceiver       n_sync.ISyncReceiver

 consensusProcessor processor.IConsensusProcessor
}

func NewHandler(
 keyPair *bls.KeyPair,
 config config.IConfig,
 connectionsManager network.IConnectionsManager,
 syncReceiver n_sync.ISyncReceiver,
 consensusProcessor processor.IConsensusProcessor,
) *Handler {
 return &Handler{
  keyPair:            keyPair,
  config:             config,
  connectionsManager: connectionsManager,
  syncReceiver:       syncReceiver,
  consensusProcessor: consensusProcessor,
 }
}

func (h *Handler) HandleRequest(request network.IRequest) (err error) {
 cmd := request.GetMessage().GetCommand()
 logger.Info("handling command: " + cmd)
 switch cmd {
 case command.InitConnection:
  return h.handleInitConnection(request)
 // case command.Pack:
 // return h.handlePack(request)
 // sync
 case command.SyncLastBlockData:
  return h.handleSyncLastBlockData(request)
 case command.SyncAccountStatesData:
  return h.handleSyncAccountStatesData(request)
 case command.SyncStakePoolData:
  return h.handleSyncStakePoolData(request)
 case command.SyncBlockData:
  return h.handleSyncBlockData(request)
  // validator messages
  // case command.Entry:
  //  return h.handleEntry(request)
  // case command.ConfirmBlock:
  //  return h.handleConfirmBlock(request)
  // node messages
  // client messages
  // execute miner messages
  // veriy miner messages
 }

 return ErrorCommandNotFound
}

/*
handleInitConnection will receive request from connection
then init that connection with data in request then
add it to connection manager
*/
func (h *Handler) handleInitConnection(request network.IRequest) (err error) {
 conn := request.GetConnection()
 initData := &pb.InitConnection{}
 err = request.GetMessage().Unmarshal(initData)
 if err != nil {
  return err
 }
 address := common.BytesToAddress(initData.Address)
 logger.Info(fmt.Sprintf(
  "init connection from %v type %v", address, initData.Type,
 ))
 conn.Init(address, initData.Type)
 h.connectionsManager.AddConnection(conn)
 return nil
}

/*
handleSyncBlockData will receive request from connection
then get full block data in consensus and send back to connection.
If block not found it will send nil back.
If Get encounters any errors, it will return .
*/
func (h *Handler) handleSyncBlockData(request network.IRequest) (err error) {
 return err
}

func (h *Handler) handleSyncLastBlockData(request network.IRequest) (err error) {
 return h.syncReceiver.SyncLastBlock(request.GetMessage().GetBody())
}

/*
handleSyncStakePoolData will receive StakePoolData from connection
then send it to sync receiver to process
If Get encounters any errors, it will return
*/
func (h *Handler) handleSyncStakePoolData(request network.IRequest) (err error) {
 return h.syncReceiver.SyncStakePool(request.GetMessage().GetBody())
}

/*
handleSyncNexrLeaderScheduleData will receive AccountStates of last block from connection
then send it to sync receiver to process
If Get encounters any errors, it will return
*/
func (h *Handler) handleSyncAccountStatesData(request network.IRequest) (err error) {
 finished, err := h.syncReceiver.SyncAccountStates(request.GetMessage().GetBody())
 if err != nil {
  return
 }
 if finished {
  err = h.syncReceiver.Commit()
 }
 return err
}

/*
handleEntry will receive entry from leader
then unmarshal it and send to consensus to process
If Get encounters any errors, it will return
*/
func (h *Handler) handleEntry(request network.IRequest) (err error) {
 return nil
}


/*
handleConfirmBlock will receive confirm block from leader
then unmarshal it and send to consensus to process
If Get encounters any errors, it will return
*/
func (h *Handler) handleConfirmBlock(request network.IRequest) (err error) {
 confirmBlock, err := block.UnmarshalConfirmBlock(request.GetMessage().GetBody())
 if err != nil {
  return err
 }
 h.consensusProcessor.ProcessConfirmBlock(confirmBlock)
 return nil
}