package pin

import (
	"context"
	"github.com/ethanvc/pin/pin/status"
	"github.com/ethanvc/pin/pin/status/codes"
	"reflect"
)

// Handler represent follow function.
// func(context.Context, *ReqType) (*TypeResp, *status.Status)
type Handler struct {
	methodName string
	methodVal  reflect.Value
}

func NewHandlers(v any) []Handler {
	vf := reflect.ValueOf(v)
	if vf.Kind() == reflect.Func {
		return []Handler{
			{methodVal: vf},
		}
	} else if vf.Kind() == reflect.Pointer && vf.Elem().Kind() == reflect.Struct {
		var result []Handler
		for i := 0; i < vf.NumMethod(); i++ {
			method := vf.Method(i)
			result = append(result, Handler{
				methodVal: method,
			})
		}
		return result
	} else {
		panic(status.NewStatus(codes.Internal, "ParamInvalid"))
	}
}

func IsValidHandler(method reflect.Value) bool {
	return true
}

func (this Handler) Call(c context.Context, req interface{}) (interface{}, *status.Status) {
	param := [...]reflect.Value{
		reflect.ValueOf(c),
		reflect.ValueOf(req),
	}
	reflectResp := this.methodVal.Call(param[:])
	return reflectResp[0].Interface(), reflectResp[1].Interface().(*status.Status)
}
