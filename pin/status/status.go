package status

import "github.com/ethanvc/pin/pin/status/codes"

type Status struct {
	codeVal     codes.Code
	eventVal    string
	subEventVal string
	msgVal      string
	pcVal       uintptr
}

func NewStatus(code codes.Code, event string) *Status {
	s := &Status{
		codeVal:  code,
		eventVal: event,
	}
	return s
}

func (this *Status) NotOk() bool {
	if this == nil {
		return false
	} else {
		return this.codeVal != codes.OK
	}
}
