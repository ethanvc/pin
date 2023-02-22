package pin

import (
	"context"
	"github.com/ethanvc/pin/pin/status"
	"reflect"
)

type Handler struct {
	DirectFunc func() *status.Status

	obj    reflect.Value
	method reflect.Method
	// only support follow type:
	// func(context.Context) *status.Status
	// func(context.Context) (*TypeResp, *status.Status)
	// func(context.Context, *ReqType) *status.Status
	// func(context.Context, *ReqType) (*TypeResp, *status.Status)
	handlerType HandlerType
}

type HandlerType int

const (
	HandlerTypeDirect HandlerType = iota
	HandlerType11
	HandlerType12
	HandlerType21
	HandlerType22
)

func (this Handler) Call(c context.Context, req interface{}) (interface{}, *status.Status) {
	switch this.handlerType {
	case HandlerTypeDirect:
		return nil, this.DirectFunc()
	case HandlerType11:
		param := [...]reflect.Value{
			this.obj,
			reflect.ValueOf(c),
		}
		reflectResp := this.method.Func.Call(param[:])
		return nil, reflectResp[0].Interface().(*status.Status)
	case HandlerType12:
		param := [...]reflect.Value{
			this.obj,
			reflect.ValueOf(c),
		}
		reflectResp := this.method.Func.Call(param[:])
		return reflectResp[0].Interface(), reflectResp[1].Interface().(*status.Status)
	case HandlerType21:
		param := [...]reflect.Value{
			this.obj,
			reflect.ValueOf(c),
			reflect.ValueOf(req),
		}
		reflectResp := this.method.Func.Call(param[:])
		return nil, reflectResp[0].Interface().(*status.Status)
	case HandlerType22:
		param := [...]reflect.Value{
			this.obj,
			reflect.ValueOf(c),
			reflect.ValueOf(req),
		}
		reflectResp := this.method.Func.Call(param[:])
		return reflectResp[0].Interface(), reflectResp[0].Interface().(*status.Status)
	}
	return nil, nil
}
