package pin

import (
	"context"
	"github.com/ethanvc/pin/pin/status"
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
	if this.interceptorIndex == len(this.Interceptors) {
		return this.callHandler()
	} else {
		this.interceptorIndex++
		return this.Interceptors[this.interceptorIndex-1](this)
	}
}

func (this *Request) callHandler() *status.Status {
	return nil
}
