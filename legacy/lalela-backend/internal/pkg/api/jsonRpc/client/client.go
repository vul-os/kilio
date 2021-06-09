package client

type Client interface {
	JSONRPCRequest(method string, request, response interface{}) error
}
