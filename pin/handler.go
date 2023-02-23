package pin

import (
	"context"
	"github.com/ethanvc/pin/pin/status"
	"reflect"
)

// Handler represent follow function.
// func(context.Context, *ReqType) (*TypeResp, *status.Status)
type Handler struct {
	obj    reflect.Value
	method reflect.Method
}

func (this Handler) Call(c context.Context, req interface{}) (interface{}, *status.Status) {
	param := [...]reflect.Value{
		this.obj,
		reflect.ValueOf(c),
		reflect.ValueOf(req),
	}
	reflectResp := this.method.Func.Call(param[:])
	return reflectResp[0].Interface(), reflectResp[0].Interface().(*status.Status)
}
