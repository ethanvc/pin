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
	structName string
	methodName string
	methodVal  reflect.Value
}

func (this Handler) IsValid() bool {
	return this.methodVal.IsValid()
}

func NewHandlers(v any) []Handler {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Func {
		return []Handler{
			{methodVal: vv},
		}
	} else if vv.Kind() == reflect.Pointer && vv.Elem().Kind() == reflect.Struct {
		var result []Handler
		vt := vv.Type()
		structName := vt.Elem().Name()
		for i := 0; i < vv.NumMethod(); i++ {
			method := vv.Method(i)
			if !IsValidHandler(method) {
				continue
			}
			methodT := vt.Method(i)
			result = append(result, Handler{
				structName: structName,
				methodName: methodT.Name,
				methodVal:  method,
			})
		}
		return result
	} else {
		panic(status.NewStatus(codes.Internal, "ParamInvalid"))
	}
}

var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
var statusPointerType = reflect.TypeOf((*status.Status)(nil))

func IsValidHandler(method reflect.Value) bool {
	t := method.Type()
	if t.NumIn() != 2 || t.NumOut() != 2 {
		return false
	}
	if t.In(0) != contextType {
		return false
	}
	if t.In(1).Kind() != reflect.Pointer {
		return false
	}
	if t.Out(0).Kind() != reflect.Pointer {
		return false
	}
	if t.Out(1) != statusPointerType {
		return false
	}
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
