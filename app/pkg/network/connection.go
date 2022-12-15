package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	p_common "gitlab.com/meta-node/meta-node/pkg/common"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
)

var (
	ErrDisconnected         = errors.New("error connection disconnected")
	ErrInvalidMessageLength = errors.New("invalid message length")
	ErrExceedMessageLength  = errors.New("exceed message length")
	ErrNilConnection        = errors.New("nil connection")
)

type IConnection interface {
	SendMessage(message IMessage) error
	Connect() error
	Disconnect() error
	GetIp() string
	GetPort() int
	GetAddress() common.Address
	GetType() string
	Init(common.Address, string)
	ReadRequest() (IRequest, error)
}

func ConnectionFromTcpConnection(tcpConn net.Conn) (IConnection, error) {
	connectionAddress := tcpConn.RemoteAddr().String()
	ip, port, err := p_common.SplitConnectionAddress(connectionAddress)
	if err != nil {
		return nil, err
	}
	return &Connection{
		address: common.Address{},
		ip:      ip,
		port:    port,
		cType:   "",
		tcpConn: tcpConn,
	}, nil
}

func NewConnection(
	address common.Address,
	ip string,
	port int,
	cType string,
) IConnection {
	return &Connection{
		address: address,
		ip:      ip,
		port:    port,
		cType:   cType,
	}
}

type Connection struct {
	mu      sync.Mutex
	address common.Address

	ip   string
	port int

	cType string

	tcpConn net.Conn
}

func (c *Connection) SendMessage(message IMessage) error {
	fmt.Println("message sent vao")
	if c == nil {
		return ErrNilConnection
	}
	b, err := message.Marshal()
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(b)))
	_, err = c.tcpConn.Write(length)
	if err != nil {
		return err
	}
	_, err = c.tcpConn.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) Connect() (err error) {
	connectionAddress := fmt.Sprintf("%v:%v", c.ip, c.port)
	logger.Info("dialing " + connectionAddress)
	c.tcpConn, err = net.Dial("tcp", connectionAddress)
	return err
}

func (c *Connection) Disconnect() error {
	return nil
}

func (c *Connection) GetIp() string {
	return c.ip
}

func (c *Connection) GetPort() int {
	return c.port
}

func (c *Connection) GetAddress() common.Address {
	return c.address
}

func (c *Connection) GetType() string {
	return c.cType
}

func (c *Connection) ReadRequest() (IRequest, error) {
	bLength := make([]byte, 8)
	_, err := io.ReadFull(c.tcpConn, bLength)
	if err != nil {
		switch err {
		case io.EOF:
			return nil, ErrDisconnected
		default:
			return nil, err
		}
	}
	messageLength := binary.LittleEndian.Uint64(bLength)

	maxMsgLength := uint64(1073741824)
	if messageLength > maxMsgLength {
		return nil, ErrExceedMessageLength
	}

	data := make([]byte, messageLength)
	byteRead, err := io.ReadFull(c.tcpConn, data)
	if err != nil {
		switch err {
		case io.EOF:
			return nil, ErrDisconnected
		default:
			return nil, err
		}
	}

	if uint64(byteRead) != messageLength {
		return nil, ErrInvalidMessageLength
	}

	msg := &pb.Message{}
	err = proto.Unmarshal(data[:messageLength], msg)
	if err != nil {
		return nil, err
	}
	return NewRequest(c, NewMessage(msg)), nil
}

func (c *Connection) Init(address common.Address, cType string) {
	c.address = address
	c.cType = cType
}
