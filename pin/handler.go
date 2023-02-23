package pin

import (
	"context"
	"github.com/ethanvc/pin/pin/status"
	"reflect"
)

// Handler represent follow function.
// func(context.Context, *ReqType) (*TypeResp, *status.Status)
type Handler struct {
	methodVal reflect.Value
}

func NewHandler(v any) Handler {
	vf := reflect.ValueOf(v)
	return Handler{
		methodVal: vf,
	}
}

func (this Handler) Call(c context.Context, req interface{}) (interface{}, *status.Status) {
	param := [...]reflect.Value{
		reflect.ValueOf(c),
		reflect.ValueOf(req),
	}
	reflectResp := this.methodVal.Call(param[:])
	return reflectResp[0].Interface(), reflectResp[1].Interface().(*status.Status)
}
