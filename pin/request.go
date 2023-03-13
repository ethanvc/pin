package pin

import (
	"context"
	"github.com/ethanvc/pin/pin/plog"
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
	ReqBytes         []byte
	RespBytes        []byte
	Method           string
	Status           *status.Status
	PatternPath      string
	Handler          Handler
	ProtocolRequest  ProtocolRequest
	interceptorIndex int
	Interceptors     []InterceptorFunc
	Logger           plog.Logger
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
	resp, status := this.Handler.Call(this.C, this.Req)
	this.Resp = resp
	return status
}
