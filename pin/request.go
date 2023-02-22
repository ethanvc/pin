package pin

import (
	"context"
	"github.com/ethanvc/pin/pin/status"
	"github.com/ethanvc/pin/pin/status/codes"
	"time"
)

type InterceptorFunc func(request *Request) *status.Status

type Request struct {
	StartTime        time.Time
	C                context.Context
	Req              interface{}
	Resp             interface{}
	Status           *status.Status
	Handler          Handler
	ProtocolRequest  ProtocolRequest
	interceptorIndex int
	Interceptors     []InterceptorFunc
}

type requestContextKey struct{}

func NewRequest(protocolReq ProtocolRequest) *Request {
	req := &Request{}
	req.StartTime = time.Now()
	req.C = context.WithValue(protocolReq.ProtocolContext(), requestContextKey{}, req)
	req.ProtocolRequest = protocolReq
	return req
}

func RequestFromContext(c context.Context) *Request {
	req, _ := c.Value(requestContextKey{}).(*Request)
	return req
}

func (this *Request) Next() *status.Status {
	intLen := len(this.Interceptors)
	if this.interceptorIndex < intLen {
		this.interceptorIndex++
		return this.Interceptors[this.interceptorIndex-1](this)
	} else if this.interceptorIndex == intLen {
		this.interceptorIndex++
		return this.callHandler()
	} else {
		return status.NewStatus(codes.Internal, "CallNextInWrongPlace")
	}
}

func (this *Request) callHandler() *status.Status {
	if this.Handler.DirectFunc != nil {
		return this.Handler.DirectFunc()
	}
	return nil
}
