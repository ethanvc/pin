package controller

import (
	"context"
	"github.com/ethanvc/pin/dbvisual_backend/service"
	"github.com/ethanvc/pin/pin/status"
)

type Instance struct {
}

type ConnectReq struct {
}

func (this *Instance) Connect(c context.Context, req *ConnectReq) (*service.Instance, *status.Status) {
	return nil, nil
}
