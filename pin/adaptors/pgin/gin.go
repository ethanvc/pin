package pgin

import (
	"context"
	"github.com/ethanvc/pin/pin"
	"github.com/ethanvc/pin/pin/status"
	"github.com/gin-gonic/gin"
)

func CreateProtocolRequest(ctx *gin.Context) pin.ProtocolRequest {
	return ginRequest{ctx: ctx}
}

type ginRequest struct {
	ctx *gin.Context
}

func (this ginRequest) InitializeRequest(req *pin.Request) *status.Status {
	return nil
}

func (this ginRequest) FinalizeRequest(req *pin.Request) *status.Status {
	return nil
}

func (this ginRequest) ProtocolContext() context.Context {
	return this.ctx.Request.Context()
}

func (this ginRequest) ProtocolObject() any {
	return this.ctx
}
