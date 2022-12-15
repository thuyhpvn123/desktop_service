package network

import (
	"context"
	"fmt"
	"net"

	"gitlab.com/meta-node/meta-node/pkg/logger"
)

type ISocketServer interface {
	Listen(string) error
	Stop()

	OnConnect(IConnection)
	OnDisconnect(IConnection)

	HandleConnection(IConnection) error
}

type SocketServer struct {
	connectionsManager IConnectionsManager
	listener           net.Listener
	handler            IHandler

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewSockerServer(
	connectionsManager IConnectionsManager,
	handler IHandler,
) ISocketServer {
	s := &SocketServer{
		connectionsManager: connectionsManager,
		handler:            handler,
	}
	s.ctx, s.cancelFunc = context.WithCancel(context.Background())
	return s
}

func (s *SocketServer) Listen(listenAddress string) error {
	var err error
	s.listener, err = net.Listen("tcp", listenAddress)
	if err != nil {
		return err
	}
	defer func() {
		s.listener.Close()
		s.listener = nil
	}()
	logger.Info(fmt.Sprintf("Listening at %v", listenAddress))
	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			tcpConn, err := s.listener.Accept()
			if err != nil {
				logger.Warn(fmt.Sprintf("Error when accept connection %v\n", err))
				continue
			}
			conn, err := ConnectionFromTcpConnection(tcpConn)
			if err != nil {
				logger.Warn(fmt.Sprintf("error when create connection from tcp connection: %v", err))
				continue
			}
			s.OnConnect(conn)
			go s.HandleConnection(conn)
		}
	}
}

func (s *SocketServer) Stop() {
	s.cancelFunc()
}

func (s *SocketServer) OnConnect(conn IConnection) {
	logger.Info(fmt.Sprintf("On Connect with %v:%v\n", conn.GetIp(), conn.GetPort()))
}

func (s *SocketServer) OnDisconnect(conn IConnection) {
	logger.Info(fmt.Sprintf("On Disconnect with %v:%v - address %v\n", conn.GetIp(), conn.GetPort(), conn.GetAddress()))
	s.connectionsManager.RemoveConnection(conn)
}

func (s *SocketServer) HandleConnection(conn IConnection) error {
	defer func() {
		conn.Disconnect()
		s.OnDisconnect(conn)
	}()
	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			request, err := conn.ReadRequest()
			if err != nil {
				logger.Warn(fmt.Sprintf("error when read request %v", err))
				return err
			}
			s.handler.HandleRequest(request)
		}
	}
}
