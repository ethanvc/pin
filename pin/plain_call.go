package pin

import (
	"context"
	"github.com/ethanvc/pin/pin/status"
)

type PlainCall[ReqType any, RespType any] struct {
	method string
	f      func(context.Context, *ReqType) (*RespType, *status.Status)
	c      context.Context
	req    *ReqType
	resp   *RespType
	status *status.Status
}

func CreatePlainCall[ReqType any, RespType any](method string,
	f func(context.Context, *ReqType) (*RespType, *status.Status)) *PlainCall[ReqType, RespType] {
	return &PlainCall[ReqType, RespType]{
		method: method,
		f:      f,
	}
}

func (this *PlainCall[ReqType, RespType]) Call(c context.Context, req *ReqType) (*RespType, *status.Status) {
	this.c = c
	this.req = req
	DefaultPlainServer.ProcessRequest(this)
	return this.resp, this.status
}

func (this *PlainCall[ReqType, RespType]) InitializeRequest(req *Request) *status.Status {
	req.Req = this.req
	return nil
}

func (this *PlainCall[ReqType, RespType]) FinalizeRequest(req *Request) *status.Status {
	this.resp = req.Resp.(*RespType)
	this.status = req.Status
	return nil
}

func (this *PlainCall[ReqType, RespType]) ProtocolContext() context.Context {
	return this.c
}

func (this *PlainCall[ReqType, RespType]) ProtocolObject() any {
	return this
}

var DefaultPlainServer Server
