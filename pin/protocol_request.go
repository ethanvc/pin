package pin

import (
	"context"
	"github.com/ethanvc/pin/pin/status"
)

type ProtocolRequest interface {
	InitializeRequest(req *Request) *status.Status
	FinalizeRequest(req *Request) *status.Status
	ProtocolContext() context.Context
	ProtocolObject() interface{}
}
