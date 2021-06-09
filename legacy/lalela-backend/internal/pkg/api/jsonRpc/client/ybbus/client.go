package ybbus

import (
	"github.com/rs/zerolog/log"
	"github.com/ybbus/jsonrpc"
	jsonRPCClient "lalela-backend/internal/pkg/api/jsonRpc/client"
	"net/http"
	"time"
)

type client struct {
	preSharedSecret string
	url             string
}

func New(
	url, preSharedSecret string,
) jsonRPCClient.Client {
	return &client{
		preSharedSecret: preSharedSecret,
		url:             url,
	}
}

func (c *client) JSONRPCRequest(method string, request, response interface{}) error {
	rpcResponse, err := jsonrpc.NewClientWithOpts(
		c.url,
		&jsonrpc.RPCClientOpts{
			HTTPClient: &http.Client{Timeout: time.Second * 10},
			CustomHeaders: map[string]string{
				"Pre-Shared-Secret": c.preSharedSecret,
				//"Claims":            marshalledClaimsForHeader,
			},
		},
	).Call(method, request)

	if err != nil {
		return err
	}
	if rpcResponse.Error != nil {
		return rpcResponse.Error
	}

	// parse response
	if err := rpcResponse.GetObject(response); err != nil {
		log.Error().Err(err).Msg("parse response object")
		return err
	}

	return nil
}
