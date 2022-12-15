package network

type IRequest interface {
	GetMessage() IMessage
	GetConnection() IConnection
}

type Request struct {
	connection IConnection
	message    IMessage
}

func NewRequest(
	connection IConnection,
	message IMessage,
) IRequest {
	return &Request{
		connection: connection,
		message:    message,
	}
}

func (r *Request) GetMessage() IMessage {
	return r.message
}

func (r *Request) GetConnection() IConnection {
	return r.connection
}
