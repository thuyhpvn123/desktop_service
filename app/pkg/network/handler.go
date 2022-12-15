package network

type IHandler interface {
	HandleRequest(IRequest) error
}
